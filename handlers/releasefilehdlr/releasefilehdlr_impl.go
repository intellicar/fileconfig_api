package releasefilehdlr

import (
	"net/http"

	"github.com/fabrikiot/goutils/apiserver"
	"github.com/fabrikiot/goutils/faberr"
	"github.com/varasheb/fileconfig_api.git/common/apperr"
	"github.com/varasheb/fileconfig_api.git/handlers"
	"github.com/varasheb/fileconfig_api.git/services/releaseconfigsvc"
)

func (o *ReleaseFileHdlr) listmyfiles(w http.ResponseWriter, r *http.Request) {
	res, err := o.releaseConfigSvcI.Listmyfiles()
	if err != nil {
		apiserver.APIResponseInternalServerError(w, r, "INTERNAL_ERROR", "", err.Error())
		return
	}
	apiserver.APIResponseOK(w, r, res, "success")
}

func (o *ReleaseFileHdlr) createreleaseconfig(w http.ResponseWriter, r *http.Request) {
	var req releaseconfigsvc.ReleasedFileReq
	handlers.CallService(w, r, &req, func(reqparsed interface{}) (interface{}, *faberr.FabErr) {
		reqmapped, ok := reqparsed.(*releaseconfigsvc.ReleasedFileReq)
		if !ok {
			return nil, apperr.ErrInvalidReqBody
		}
		return o.releaseConfigSvcI.CreateReleaseConfig(reqmapped)
	}, "success")

}

func (o *ReleaseFileHdlr) updatereleaseconfig(w http.ResponseWriter, r *http.Request) {
	var req releaseconfigsvc.ReleasedFileUpdReq
	handlers.CallService(w, r, &req, func(reqparsed interface{}) (interface{}, *faberr.FabErr) {
		reqmapped, ok := reqparsed.(*releaseconfigsvc.ReleasedFileUpdReq)
		if !ok {
			return nil, apperr.ErrInvalidReqBody
		}
		return o.releaseConfigSvcI.UpdateReleaseConfig(reqmapped)
	}, "success")
}

func (o *ReleaseFileHdlr) deletereleaseconfig(w http.ResponseWriter, r *http.Request) {
	var req releaseconfigsvc.ReleasedFileDelReq
	handlers.CallService(w, r, &req, func(reqparsed interface{}) (interface{}, *faberr.FabErr) {
		reqmapped, ok := reqparsed.(*releaseconfigsvc.ReleasedFileDelReq)
		if !ok {
			return nil, apperr.ErrInvalidReqBody
		}
		return o.releaseConfigSvcI.DeleteReleaseConfig(reqmapped.Filename, reqmapped.UpdatedBy)
	}, "success")
}
