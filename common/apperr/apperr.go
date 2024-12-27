package apperr

import "github.com/fabrikiot/goutils/faberr"

var (
	// Server related
	ErrInvalidState  = faberr.NewFabErr("ERR_SERVER_INVALID_STATE", nil, "invalid state to do this operation")
	ErrIsBusy        = faberr.NewFabErr("ERR_SERVER_BUSY", nil, "the processor is currently busy")
	ErrCacheNotReady = faberr.NewFabErr("ERR_SERVER_CACHE_NOT_READY", nil, "cache is still not ready, we are still loading")
	ErrInternal      = faberr.NewFabErr("ERR_SERVER_INTERNAL", nil, "internal server error")
	ErrStopped       = faberr.NewFabErr("ERR_SERVER_SERVICE_STOPPING", nil, "service is stopping")

	// DB related
	ErrDBGetConn     = faberr.NewFabErr("ERR_DB_CONN_GET_FAILED", nil, "db connection unsuccessful")
	ErrDBQuery       = faberr.NewFabErr("ERR_DB_QUERY_FAILED", nil, "failed to query")
	ErrDBQueryRow    = faberr.NewFabErr("ERR_DB_QUERY_ROW", nil, "failed to query row")
	ErrDBTxn         = faberr.NewFabErr("ERR_DB_TRANSACTION_FAILED", nil, "db transaction failed")
	ErrDBInsert      = faberr.NewFabErr("ERR_DB_INSERT_FAILED", nil, "db insert has failed")
	ErrDBScan        = faberr.NewFabErr("ERR_DB_SCAN", nil, "error while scanning the response")
	ErrDBAuditInsert = faberr.NewFabErr("ERR_DB_AUDIT_INSERT_FAILED", nil, "audit insert failed")

	// Request related
	ErrUnknownReq         = faberr.NewFabErr("ERR_REQ_UNKNOWN", nil, "unknown request")
	ErrInvalidReqBody     = faberr.NewFabErr("ERR_REQ_INVALID", nil, "invalid request")
	ErrInvalidEmailString = faberr.NewFabErr("ERR_INVALID_EMAIL_STRING", nil, "invalid email string")
	ErrInvalidParam       = faberr.NewFabErr("ERR_PARAM_INVALID", nil, "invalid parameter request")

	// Token related
	ErrAuthHeaderInvalid  = faberr.NewFabErr("ERR_TOKEN_AUTH_HEADER_INVALID", nil, "auth header invalid")
	ErrBearerTokenMissing = faberr.NewFabErr("ERR_TOKEN_BEARER_TOKEN_MISSING", nil, "bearer token missing")
	ErrTokenInvalid       = faberr.NewFabErr("ERR_TOKEN_INVALID", nil, "token invalid")
	ErrVerifyTokenFailed  = faberr.NewFabErr("ERR_TOKEN_VERIFY_FAILED", nil, "token verification failed")
	ErrTokenGen           = faberr.NewFabErr("ERR_TOKEN_GEN_FAILED", nil, "token gen failed")

	// HTTP Request related
	ErrCreatingHttpReq           = faberr.NewFabErr("ERR_HTTP_CREATING_REQ", nil, "error creating http request")
	ErrMakingHttpReq             = faberr.NewFabErr("ERR_HTTP_MAKING_REQ", nil, "error making http request")
	ErrReadingHttpRsp            = faberr.NewFabErr("ERR_HTTP_READING_RSP", nil, "error reading http response")
	ErrUpstreamHttpRsp           = faberr.NewFabErr("ERR_HTTP_UPSTREAM_RSP", nil, "error upstream http response")
	ErrUpstreamUnexpectedHttpRsp = faberr.NewFabErr("ERR_HTTP_UPSTREAM_UNEXPECTED_RSP", nil, "error upstream unexpected http response")

	// JSON errors
	ErrMarshalJson   = faberr.NewFabErr("ERR_JSON_MARSHAL", nil, "error marshalling json")
	ErrUnmarshalJson = faberr.NewFabErr("ERR_JSON_UNMARSHAL", nil, "error unmarshalling json")
	ErrParsingJson   = faberr.NewFabErr("ERR_JSON_PARSING", nil, "error parsing json")

	// GCS Upload Errors
	ErrGcsStorageWrite      = faberr.NewFabErr("ERR_GCS_WRITE", nil, "error while writing to the cloud storage")
	ErrGcsStorageWriteClose = faberr.NewFabErr("ERR_GCS_WRITE_CLOSE", nil, "error while closing the client after writing")
	ErrGcsFileDelete        = faberr.NewFabErr("ERR_GCS_FILE_DELETE_FAILED", nil, "failed to delete file in gcs")
	ErrGcsFileReadFail      = faberr.NewFabErr("ERR_GCS_FILE_READ_FAILED", nil, "failed to read file from gcs")
	ErrGcsFileReader        = faberr.NewFabErr("ERR_GCS_FILE_READER", nil, "error getting gcs file reader")

	// Other common errors
	ErrInvalidNameFormat     = faberr.NewFabErr("INVALID_NAME_FORMAT", nil, "name can only have lowercase alphabets, hyphen and numbers")
	ErrInvalidDispNameFormat = faberr.NewFabErr("INVALID_DISP_NAME_FORMAT", nil, "disp name can only have alphabets, hyphen, whitespace and numbers")
	ErrJSONUnMarshal         = faberr.NewFabErr("ERR_JSON_UNMARSHAL", nil, "error unmarshalling")
	ErrMaxReadSize           = faberr.NewFabErr("ERR_MAX_READ_SIZE", nil, "error reached max read size")

	ErrUnauthorizedAction = faberr.NewFabErr("UNAUTHORIZED_ACTION", nil, "user not authorized to do this action")

	ErrUnauthorizedLoginDomain = faberr.NewFabErr("UNAUTHORIZED_LOGIN_DOMAIN", nil, "user not authorized to login from this domain")
)
