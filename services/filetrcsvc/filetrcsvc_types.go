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
