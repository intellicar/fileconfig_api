package releaseconfigsvc

type ReleasedFile struct {
	ConfigID      int64   `json:"configid"`
	Group         *string `json:"group"`
	BoardVersion  *string `json:"boardversion"`
	ReleaseDate   int64   `json:"releaseDate"`
	Filename      string  `json:"filename"`
	SIM           string  `json:"sim"`
	NRFBootloader *string `json:"nrfbootloader"`
	ReleaseNote   *string `json:"releasenote"`
	IsLatest      bool    `json:"islatest"`
	IsValid       bool    `json:"isvalid"`
	CreatedBy     string  `json:"createdby"`
	GroupType     *string `json:"grouptype"`
	CreatedAt     int64   `json:"createdat"`
}

type ReleasedFileReq struct {
	Group         string `json:"group" validate:"required"`
	GroupType     string `json:"grouptype" validate:"required"`
	BoardVersion  string `json:"boardversion" validate:"required"`
	ReleaseDate   int64  `json:"releaseDate" validate:"required"`
	Filename      string `json:"filename" validate:"required"`
	SIM           string `json:"sim" validate:"required"`
	NRFBootloader string `json:"nrfbootloader" validate:"required"`
	ReleaseNote   string `json:"releasenote" validate:"required"`
	IsLatest      bool   `json:"islatest" validate:"required"`
	IsValid       bool   `json:"isvalid" validate:"required"`
	CreatedBy     string `json:"createdby" validate:"required"`
}

type ReleasedFileUpdReq struct {
	Group         string `json:"group,omitempty" validate:"omitempty"`
	GroupType     string `json:"grouptype,omitempty" validate:"omitempty"`
	BoardVersion  string `json:"boardversion,omitempty" validate:"omitempty"`
	ReleaseDate   int64  `json:"releaseDate,omitempty" validate:"omitempty"`
	Filename      string `json:"filename,omitempty" validate:"required"`
	NewFilename   string `json:"newfilename,omitempty" validate:"omitempty"`
	SIM           string `json:"sim,omitempty" validate:"omitempty"`
	NRFBootloader string `json:"nrfbootloader,omitempty" validate:"omitempty"`
	ReleaseNote   string `json:"releasenote,omitempty" validate:"omitempty"`
	IsLatest      *bool  `json:"islatest,omitempty" validate:"omitempty"`
	IsValid       *bool  `json:"isvalid,omitempty" validate:"omitempty"`
	UpdatedBy     string `json:"updatedby,omitempty" validate:"required"`
}

type ReleasedFileDelReq struct {
	Filename  string `json:"filename" validate:"required"`
	UpdatedBy string `json:"updatedby" validate:"required"`
}
