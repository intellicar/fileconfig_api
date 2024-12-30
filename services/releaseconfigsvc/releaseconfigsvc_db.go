package releaseconfigsvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/jackc/pgx/v5"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
	"github.com/varasheb/fileconfig_api.git/common/apperr/utils"
)

var cachelist []*ReleasedFile

func (o *ReleaseFileConfigSvc) listmyfiles() ([]*ReleasedFile, *faberr.FabErr) {

	if len(cachelist) > 0 {
		return cachelist, nil
	}

	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}
	query := `
	SELECT configid, "group", boardversion, "releasedate", filename, sim, 
		   nrfbootloader, releasenote, islatest, isvalid, createdby, grouptype, createdAt
	FROM "FileRelease"."ReleaseConfig"
	WHERE isdelete = false
	ORDER BY "releasedate"`
	ctx, ctxcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxcancelf()
	rows, scanner := dbI.Query(ctx, query)
	if scanner != nil {
		return nil, apperr.ErrInternal.NewWData(scanner.Error())
	}
	defer rows.Close()
	var filelist []*ReleasedFile
	for rows.Next() {
		file := &ReleasedFile{}
		scanner := rows.Scan(&file.ConfigID, &file.Group, &file.BoardVersion, &file.ReleaseDate, &file.Filename,
			&file.SIM, &file.NRFBootloader, &file.ReleaseNote, &file.IsLatest, &file.IsValid, &file.CreatedBy,
			&file.GroupType, &file.CreatedAt)
		if scanner != nil {
			return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
		}
		filelist = append(filelist, file)
	}

	cachelist = filelist

	return filelist, nil
}

func (o *ReleaseFileConfigSvc) createreleaseconfig(file *ReleasedFileReq) (*ReleasedFile, *faberr.FabErr) {
	isexists, err := o.isFileNameExists(file.Filename)
	if err != nil {
		return nil, err
	}
	if isexists {
		return nil, ErrFileNameAlreadyExists
	}
	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}
	tx, err1 := dbI.BeginTx(context.Background(), pgx.TxOptions{})
	if err1 != nil {
		return nil, apperr.ErrDBTxn.NewWData(err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	querry := `INSERT INTO "FileRelease"."ReleaseConfig" ("group", boardversion, "releasedate", filename, sim, nrfbootloader, releasenote, islatest, isvalid, createdby, grouptype)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
               RETURNING configid, "group", boardversion, "releasedate", filename, sim, nrfbootloader, islatest, isvalid, createdby, createdat, grouptype;`
	qctx, qcexcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer qcexcancelf()

	var releaseFile = &ReleasedFile{}
	scanner := dbI.QueryRow(qctx, querry,
		file.Group, file.BoardVersion, file.ReleaseDate, file.Filename,
		file.SIM, file.NRFBootloader, file.ReleaseNote, file.IsLatest,
		file.IsValid, file.CreatedBy, file.GroupType).Scan(
		&releaseFile.ConfigID, &releaseFile.Group, &releaseFile.BoardVersion, &releaseFile.ReleaseDate,
		&releaseFile.Filename, &releaseFile.SIM, &releaseFile.NRFBootloader,
		&releaseFile.IsLatest, &releaseFile.IsValid, &releaseFile.CreatedBy, &releaseFile.CreatedAt, &releaseFile.GroupType)

	if scanner != nil {
		if utils.IsPsqlDuplicateKeyErrorMessage(scanner.Error()) {
			return nil, ErrFileNameAlreadyExists
		}
		o.logger.Printf("ERROR: %v", scanner.Error())
		return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
	}
	err1 = tx.Commit(context.Background())
	if err1 != nil {
		return nil, apperr.ErrDBInsert.NewWData(err.Error())
	}
	o.clearCache()
	return releaseFile, nil
}

func (o *ReleaseFileConfigSvc) updatereleaseconfig(req *ReleasedFileUpdReq) (*ReleasedFile, *faberr.FabErr) {
	isexists, err := o.isFileNameExists(req.Filename)
	if err != nil {
		return nil, err
	}
	if !isexists {
		return nil, ErrFileNameNotExists
	}
	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}
	tx, err1 := dbI.BeginTx(context.Background(), pgx.TxOptions{})
	if err1 != nil {
		return nil, apperr.ErrDBTxn.NewWData(err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := `UPDATE "FileRelease"."ReleaseConfig" SET `
	var updates []string
	var args []interface{}
	argCounter := 1

	if req.IsLatest != nil {
		updates = append(updates, fmt.Sprintf("islatest = $%d", argCounter))
		args = append(args, *req.IsLatest)
		argCounter++
	}
	if req.IsValid != nil {
		updates = append(updates, fmt.Sprintf("isvalid = $%d", argCounter))
		args = append(args, *req.IsValid)
		argCounter++
	}
	if req.Group != "" {
		updates = append(updates, fmt.Sprintf("\"group\" =$%d", argCounter)) // "group" instead of group
		args = append(args, req.Group)
		argCounter++
	}
	if req.GroupType != "" {
		updates = append(updates, fmt.Sprintf("grouptype = $%d", argCounter))
		args = append(args, req.GroupType)
		argCounter++
	}
	if req.BoardVersion != "" {
		updates = append(updates, fmt.Sprintf("boardversion = $%d", argCounter))
		args = append(args, req.BoardVersion)
		argCounter++
	}
	if req.ReleaseDate != 0 {
		updates = append(updates, fmt.Sprintf("releasedate = $%d", argCounter))
		args = append(args, req.ReleaseDate)
		argCounter++
	}
	if req.NewFilename != "" {
		updates = append(updates, fmt.Sprintf("filename = $%d", argCounter))
		args = append(args, req.NewFilename)
		argCounter++
	}
	if req.SIM != "" {
		updates = append(updates, fmt.Sprintf("sim = $%d", argCounter))
		args = append(args, req.SIM)
		argCounter++
	}
	if req.NRFBootloader != "" {
		updates = append(updates, fmt.Sprintf("nrfbootloader = $%d", argCounter))
		args = append(args, req.NRFBootloader)
		argCounter++
	}
	if req.ReleaseNote != "" {
		updates = append(updates, fmt.Sprintf("releasenote = $%d", argCounter))
		args = append(args, req.ReleaseNote)
		argCounter++
	}
	if req.UpdatedBy != "" {
		updates = append(updates, fmt.Sprintf("updatedby = $%d", argCounter))
		args = append(args, req.UpdatedBy)
		argCounter++
	}

	if len(updates) == 0 {
		return nil, faberr.NewFabErr("ERR_NO_UPDATES", nil, "no fields to update")
	}
	updates = append(updates, fmt.Sprintf("updatedat = $%d", argCounter))
	args = append(args, time.Now().UnixMilli())
	argCounter++
	whereclause := fmt.Sprintf(" WHERE filename = $%d ", argCounter)
	args = append(args, req.Filename)

	query += strings.Join(updates, ", ") + whereclause +
		`RETURNING configid, "group", boardversion, "releasedate", filename, sim, nrfbootloader, islatest, isvalid, createdby, createdat, grouptype;`

	qctx, qcexcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer qcexcancelf()
	var releaseFile = &ReleasedFile{}
	scanner := dbI.QueryRow(qctx, query, args...).Scan(
		&releaseFile.ConfigID, &releaseFile.Group, &releaseFile.BoardVersion, &releaseFile.ReleaseDate,
		&releaseFile.Filename, &releaseFile.SIM, &releaseFile.NRFBootloader,
		&releaseFile.IsLatest, &releaseFile.IsValid, &releaseFile.CreatedBy, &releaseFile.CreatedAt, &releaseFile.GroupType)

	if scanner != nil {
		o.logger.Printf("ERROR: %v", scanner.Error())
		return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
	}
	err1 = tx.Commit(context.Background())
	if err1 != nil {
		return nil, apperr.ErrDBInsert.NewWData(err.Error())
	}
	o.clearCache()
	return releaseFile, nil
}

func (o *ReleaseFileConfigSvc) deletereleaseconfig(filename string, updatedby string) (*ReleasedFile, *faberr.FabErr) {
	isexists, err := o.isFileNameExists(filename)
	if err != nil {
		return nil, err
	}
	if !isexists {
		return nil, ErrFileNameNotExists
	}
	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}
	tx, err1 := dbI.BeginTx(context.Background(), pgx.TxOptions{})
	if err1 != nil {
		return nil, apperr.ErrDBTxn.NewWData(err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	querry := `UPDATE "FileRelease"."ReleaseConfig" SET isdelete = true, updatedat = $1, updatedby = $2 WHERE filename = $3 RETURNING configid, "group", boardversion, "releasedate", filename, sim, nrfbootloader, islatest, isvalid, createdby, createdat, grouptype;`
	qctx, qcexcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer qcexcancelf()

	var releaseFile = &ReleasedFile{}
	updatedat := time.Now().UnixMilli()
	scanner := dbI.QueryRow(qctx, querry, updatedat, updatedby, filename).Scan(
		&releaseFile.ConfigID, &releaseFile.Group, &releaseFile.BoardVersion, &releaseFile.ReleaseDate,
		&releaseFile.Filename, &releaseFile.SIM, &releaseFile.NRFBootloader,
		&releaseFile.IsLatest, &releaseFile.IsValid, &releaseFile.CreatedBy, &releaseFile.CreatedAt, &releaseFile.GroupType)

	if scanner != nil {
		if utils.IsPsqlDuplicateKeyErrorMessage(scanner.Error()) {
			return nil, ErrFileNameAlreadyExists
		}
		o.logger.Printf("ERROR: %v", scanner.Error())
		return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
	}
	err1 = tx.Commit(context.Background())
	if err1 != nil {
		return nil, apperr.ErrDBInsert.NewWData(err.Error())
	}
	o.clearCache()

	return releaseFile, nil

}

func (o *ReleaseFileConfigSvc) clearCache() {
	cachelist = []*ReleasedFile{}
}
