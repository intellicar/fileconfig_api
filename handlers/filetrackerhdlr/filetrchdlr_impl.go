package filetrchdlr

import (
	"net/http"

	"github.com/fabrikiot/goutils/apiserver"
)

func (o *FileTrcHdlr) listmyfiles(w http.ResponseWriter, r *http.Request) {
	res, err := o.fileTrcSvcI.Listmyfiles()
	if err != nil {
		apiserver.APIResponseInternalServerError(w, r, "INTERNAL_ERROR", "", err.Error())
		return
	}
	apiserver.APIResponseOK(w, r, res, "sucess")
}
