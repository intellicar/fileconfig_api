package filetrcsvc

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/fabrikiot/goutils/faberr"
	"github.com/fabrikiot/goutils/fabpgsql"
)

type FileTrcSvc struct {
	pgsqlI *fabpgsql.PGSqlDBService
	logger *log.Logger

	isstopped    *atomic.Bool
	activecalls  *atomic.Int64
	activeThread *sync.WaitGroup
}

func NewfileTrcSvc(pgsqlI *fabpgsql.PGSqlDBService, logger *log.Logger) *FileTrcSvc {
	return &FileTrcSvc{
		pgsqlI: pgsqlI,
		logger: logger,

		isstopped:    &atomic.Bool{},
		activecalls:  &atomic.Int64{},
		activeThread: &sync.WaitGroup{},
	}
}
func (o *FileTrcSvc) Start() {

}

func (o *FileTrcSvc) Listmyfiles() ([]*FileTrc, *faberr.FabErr) {
	return o.listmyfiles()
}
