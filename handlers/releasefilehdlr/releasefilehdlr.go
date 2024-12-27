package releasefilehdlr

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/varasheb/fileconfig_api.git/services/releaseconfigsvc"
)

type ReleaseFileHdlr struct {
	releaseConfigSvcI *releaseconfigsvc.ReleaseFileConfigSvc
	logger            *log.Logger
}

func NewReleaseFileHdlr(releasecConfigSvcI *releaseconfigsvc.ReleaseFileConfigSvc, logger *log.Logger) *ReleaseFileHdlr {
	return &ReleaseFileHdlr{
		releaseConfigSvcI: releasecConfigSvcI,
		logger:            logger,
	}
}

func (o *ReleaseFileHdlr) RegisterRoutes(router chi.Router) {
	router.Get("/", o.listmyfiles)
	router.Post("/", o.createreleaseconfig)
	router.Put("/", o.updatereleaseconfig)
	router.Delete("/", o.deletereleaseconfig)

}
