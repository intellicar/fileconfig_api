package filetrcsvc

type FileTrc struct {
	FileID      string `json:"fileid"`
	FileName    string `json:"filename"`
	FileType    string `json:"filetype"`
	Filehash    string `json:"filehash"`
	Signature   string `json:"signature"`
	Plsign      string `json:"plsign"`
	Description string `json:"description"`
	CreatedAt   int    `json:"createdat"`
	CreatedBy   string `json:"createdby"`
}

type GetFileByNameReq struct {
	FileName string `json:"filename"`
}

type FileTrcReq struct {
	FileID      string `json:"fileid" validate:"required"`
	FileName    string `json:"filename" validate:"required"`
	FileType    string `json:"filetype" validate:"required"`
	Filehash    string `json:"filehash" validate:"required"`
	Signature   string `json:"signature" validate:"required"`
	Plsign      string `json:"plsign" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	CreatedBy   string `json:"createdby" validate:"required"`
}
