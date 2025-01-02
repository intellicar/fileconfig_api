package filetrcsvc

import (
	"context"
	"errors"
	"time"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/jackc/pgx/v5"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
	"github.com/varasheb/fileconfig_api.git/common/apperr/utils"
)

var cachelist []*FileTrc

func (o *FileTrcSvc) listmyfiles() ([]*FileTrc, *faberr.FabErr) {
	if len(cachelist) > 0 {
		return cachelist, nil
	}

	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}
	query := `SELECT fileid, filename, filetype, filehash, signature, plsign, description, createdat, createdby FROM "FileTracker".configfiles ORDER BY createdat`
	ctx, ctxcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxcancelf()
	rows, scanner := dbI.Query(ctx, query)
	if scanner != nil {
		return nil, apperr.ErrInternal.NewWData(scanner.Error())
	}
	defer rows.Close()
	var filelist []*FileTrc
	for rows.Next() {
		file := &FileTrc{}
		scanner := rows.Scan(&file.FileID, &file.FileName, &file.FileType, &file.Filehash, &file.Signature, &file.Plsign, &file.Description, &file.CreatedAt, &file.CreatedBy) //
		if scanner != nil {
			return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
		}
		filelist = append(filelist, file)
	}
	cachelist = filelist //cache
	return filelist, nil
}

func (o *FileTrcSvc) getbyfilename(filename string) (*FileTrc, *faberr.FabErr) {
	if len(cachelist) > 0 {
		for _, file := range cachelist {
			if file.FileName == filename {
				return file, nil
			}
		}
	}
	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}

	query := `SELECT fileid, filename, filetype, filehash, signature, plsign, description, createdat, createdby 
	          FROM "FileTracker".configfiles 
	          WHERE filename = $1`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	row := dbI.QueryRow(ctx, query, filename)
	file := &FileTrc{}
	err := row.Scan(
		&file.FileID, &file.FileName, &file.FileType, &file.Filehash, &file.Signature, &file.Plsign, &file.Description, &file.CreatedAt, &file.CreatedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFileNameNotExists
		}
		return nil, apperr.ErrDBQueryRow.NewWData(err.Error())
	}

	return file, nil
}

func (o *FileTrcSvc) createconfigfiles(file *FileTrcReq) (*FileTrc, *faberr.FabErr) {

	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return nil, apperr.ErrDBGetConn
	}

	tx, txErr := dbI.BeginTx(context.Background(), pgx.TxOptions{})
	if txErr != nil {
		return nil, apperr.ErrDBTxn.NewWData(txErr.Error())
	}
	defer func() {
		if txErr != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := `INSERT INTO "FileTracker".configfiles (fileid, filename, filetype, filehash, signature, plsign, description, createdby)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING fileid, filename, filetype, filehash, signature, plsign, description, createdat, createdby`
	qctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var configFile = &FileTrc{}

	scanner := tx.QueryRow(qctx, query,
		file.FileID, file.FileName, file.FileType, file.Filehash, file.Signature, file.Plsign, file.Description, file.CreatedBy,
	).Scan(
		&configFile.FileID, &configFile.FileName, &configFile.FileType, &configFile.Filehash,
		&configFile.Signature, &configFile.Plsign, &configFile.Description, &configFile.CreatedAt, &configFile.CreatedBy,
	)

	if scanner != nil {
		if utils.IsPsqlDuplicateKeyErrorMessage(scanner.Error()) {
			tx.Rollback(context.Background())
			return nil, ErrFileNameAlreadyExists
		}
		tx.Rollback(context.Background())
		return nil, apperr.ErrDBQueryRow.NewWData(scanner.Error())
	}

	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		return nil, apperr.ErrDBInsert.NewWData(commitErr.Error())
	}
	o.clearCache()

	return configFile, nil
}

func (o *FileTrcSvc) clearCache() {
	cachelist = []*FileTrc{}
}
