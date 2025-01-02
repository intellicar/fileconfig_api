package filetrcsvc

import "github.com/fabrikiot/goutils/faberr"

var (
	ErrFileNameAlreadyExists = faberr.NewFabErr("INVALID_REQ_FILENAME_ALREADY_EXISTS", nil, "Filename already exists")
	ErrFileNameNotExists     = faberr.NewFabErr("INVALID_REQ_FILENAME_NOT_EXISTS", nil, "Filename not matching")
)
