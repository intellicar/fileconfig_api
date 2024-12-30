package filetrcsvc

import (
	"context"
	"time"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
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
	cachelist = filelist

	return filelist, nil
}
