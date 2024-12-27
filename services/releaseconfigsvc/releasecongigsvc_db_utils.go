package releaseconfigsvc

import (
	"context"
	"errors"
	"time"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/jackc/pgx/v5"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
)

func (o *ReleaseFileConfigSvc) isFileNameExists(filename string) (bool, *faberr.FabErr) {
	dbI, getdberr := o.pgsqlI.GetDBInstance()
	if getdberr != nil {
		return false, apperr.ErrDBGetConn
	}
	querry := `SELECT filename FROM "FileRelease"."ReleaseConfig" WHERE filename = $1`
	ctx, ctxcancelf := context.WithTimeout(context.Background(), time.Second*30)
	defer ctxcancelf()
	var qfilename string
	scanner := dbI.QueryRow(ctx, querry, filename).Scan(&qfilename)
	if scanner != nil {
		if errors.Is(scanner, pgx.ErrNoRows) {
			return false, nil
		}
		return false, apperr.ErrDBQueryRow.NewWData(scanner.Error())
	}
	return true, nil
}
