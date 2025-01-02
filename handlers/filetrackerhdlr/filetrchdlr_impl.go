package filetrchdlr

import (
	"net/http"

	"github.com/fabrikiot/goutils/apiserver"
	"github.com/fabrikiot/goutils/faberr"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
	"github.com/varasheb/fileconfig_api.git/handlers"
	"github.com/varasheb/fileconfig_api.git/services/filetrcsvc"
)

func (o *FileTrcHdlr) listmyfiles(w http.ResponseWriter, r *http.Request) {
	res, err := o.fileTrcSvcI.Listmyfiles()
	if err != nil {
		apiserver.APIResponseInternalServerError(w, r, "INTERNAL_ERROR", "", err.Error())
		return
	}
	apiserver.APIResponseOK(w, r, res, "sucess")
}

func (o *FileTrcHdlr) getbyfilename(w http.ResponseWriter, r *http.Request) {
	var req filetrcsvc.GetFileByNameReq
	handlers.CallService(w, r, &req, func(reqparsed interface{}) (interface{}, *faberr.FabErr) {
		reqmapped, ok := reqparsed.(*filetrcsvc.GetFileByNameReq)
		if !ok {
			return nil, apperr.ErrInvalidReqBody
		}
		return o.fileTrcSvcI.GetByFilename(reqmapped.FileName)
	}, "success")

}
func (o *FileTrcHdlr) createconfigfiles(w http.ResponseWriter, r *http.Request) {
	var req filetrcsvc.FileTrcReq
	handlers.CallService(w, r, &req, func(reqparsed interface{}) (interface{}, *faberr.FabErr) {
		reqmapped, ok := reqparsed.(*filetrcsvc.FileTrcReq)
		if !ok {
			return nil, apperr.ErrInvalidReqBody
		}
		return o.fileTrcSvcI.CreateConfigFiles(reqmapped)
	}, "success")
}
