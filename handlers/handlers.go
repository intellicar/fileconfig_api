package handlers

import (
	"net/http"
	"strings"

	"github.com/fabrikiot/goutils/apiserver"
	"github.com/fabrikiot/goutils/faberr"
)

type callServiceCallbackFn func(interface{}) (interface{}, *faberr.FabErr)

func CallService(w http.ResponseWriter, r *http.Request, reqObj interface{}, svcf callServiceCallbackFn, succmsg string) {
	// 1. Parse the body according to the request..
	parseerrdata, parseerr := apiserver.JsonBodyParser(r.Body, reqObj)
	if parseerr != nil {
		apiserver.APIResponseBadRequest(w, r, "MALFORMED_REQUEST", parseerrdata, parseerr.Error())
		return
	}

	// 2. Now that the body is available call the service fn callback..
	svcresp, svcerr := svcf(reqObj)
	if svcerr != nil {
		if strings.HasPrefix(svcerr.ErrCode, "ERR_SERVER") || strings.HasPrefix(svcerr.ErrCode, "ERR_DB") {
			apiserver.APIResponseInternalServerError(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		} else if strings.HasPrefix(svcerr.ErrCode, "ERR_JSON") {
			apiserver.APIResponseInternalServerError(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		} else if strings.HasPrefix(svcerr.ErrCode, "ERR_HTTP") {
			apiserver.APIFailedInternalAPICall(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg, http.StatusInternalServerError)
			return
		} else if strings.HasPrefix(svcerr.ErrCode, "ERR_REQ") {
			apiserver.APIResponseBadRequest(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		} else if strings.Contains(svcerr.ErrCode, "INVALID") {
			apiserver.APIResponseBadRequest(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		} else if strings.HasPrefix(svcerr.ErrCode, "ERR_TOKEN") {
			apiserver.APIResponseUnauthorized(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		} else if strings.HasPrefix(svcerr.ErrCode, "UNAUTHORIZED") {
			apiserver.APIResponseUnauthorized(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
			return
		}

		apiserver.APIResponseUnprocessableEntity(w, r, svcerr.ErrCode, svcerr.ErrData, svcerr.ErrMsg)
		return
	}

	// 3. Service err is nil, svc response with 200 should be returned...
	apiserver.APIResponseOK(w, r, svcresp, succmsg)
}
