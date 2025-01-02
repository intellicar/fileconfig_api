package filetrchdlr

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/varasheb/fileconfig_api.git/services/filetrcsvc"
)

type FileTrcHdlr struct {
	fileTrcSvcI *filetrcsvc.FileTrcSvc
	logger      *log.Logger
}

func NewReleaseFileHdlr(fileTrcSvcI *filetrcsvc.FileTrcSvc, logger *log.Logger) *FileTrcHdlr {
	return &FileTrcHdlr{
		fileTrcSvcI: fileTrcSvcI,
		logger:      logger,
	}
}

func (o *FileTrcHdlr) RegisterRoutes(router chi.Router) {
	router.Get("/", o.listmyfiles)
	router.Post("/getbyfilename", o.getbyfilename)
	router.Post("/", o.createconfigfiles)

}
