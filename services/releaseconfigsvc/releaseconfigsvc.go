package releaseconfigsvc

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/fabrikiot/goutils/fabpgsql"
)

type ReleaseFileConfigSvc struct {
	pgsqlI *fabpgsql.PGSqlDBService
	logger *log.Logger

	isstopped    *atomic.Bool
	activecalls  *atomic.Int64
	activeThread *sync.WaitGroup
}

func NewfileConfigSvc(pgsqlI *fabpgsql.PGSqlDBService, logger *log.Logger) *ReleaseFileConfigSvc {
	return &ReleaseFileConfigSvc{
		pgsqlI: pgsqlI,
		logger: logger,

		isstopped:    &atomic.Bool{},
		activecalls:  &atomic.Int64{},
		activeThread: &sync.WaitGroup{},
	}
}

func (o *ReleaseFileConfigSvc) Start() {

}

func (o *ReleaseFileConfigSvc) Listmyfiles() ([]*ReleasedFile, *faberr.FabErr) {
	return o.listmyfiles()
}

func (o *ReleaseFileConfigSvc) CreateReleaseConfig(req *ReleasedFileReq) (*ReleasedFile, *faberr.FabErr) {
	return o.createreleaseconfig(req)
}

func (o *ReleaseFileConfigSvc) UpdateReleaseConfig(req *ReleasedFileUpdReq) (*ReleasedFile, *faberr.FabErr) {
	return o.updatereleaseconfig(req)
}

func (o *ReleaseFileConfigSvc) DeleteReleaseConfig(filename string, updatedby string) (*ReleasedFile, *faberr.FabErr) {
	return o.deletereleaseconfig(filename, updatedby)
}
