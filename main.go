package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/fabrikiot/goutils/apiserver"
	"github.com/fabrikiot/goutils/configparser"
	"github.com/fabrikiot/goutils/fabpgsql"
	"github.com/natefinch/lumberjack"
	"github.com/varasheb/fileconfig_api.git/config"
	filetrchdlr "github.com/varasheb/fileconfig_api.git/handlers/filetrackerhdlr"
	"github.com/varasheb/fileconfig_api.git/handlers/releasefilehdlr"
	"github.com/varasheb/fileconfig_api.git/services/filetrcsvc"
	"github.com/varasheb/fileconfig_api.git/services/releaseconfigsvc"
)

var ConfigFileName string

func init() {
	flag.StringVar(&ConfigFileName, "configfile", "", "Please pass the config file")

	flag.Parse()
	if ConfigFileName == "" {
		log.Fatal("Please pass the config file as the argument")
	}
	log.Print("Config File to be used:" + ConfigFileName)
}

func addLoggerFile(logDir string, logFile string, logger *log.Logger) {
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(path.Join(logDir, logFile)),
		MaxSize:    5, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	logger.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	logger.SetOutput(lumberjackLogger)
}

func setupLogger(logDirPath string) {
	logFilePath := path.Join(logDirPath, "general.log")
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(logFilePath),
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	log.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	log.SetOutput(lumberjackLogger)
}

func main() {
	fmt.Println("start")
	var serverConfigI config.ConfigServer

	parserr := configparser.ParseConfigFromFile(ConfigFileName, &serverConfigI)
	if parserr != nil {
		log.Fatal("Please provide a valid config file")
	}
	// 1. logpath
	logDirPath := path.Join(*serverConfigI.ScratchDir, "logs")
	setupLogger(logDirPath)

	// 2. Create the pgsql service required by the namespacesvc...
	pgsqlI := fabpgsql.NewPGSqlDBService(*serverConfigI.PgSQLConfig.PgURL)

	// 3. Services...
	serviceLoggerI := log.New(os.Stdout, "SERVICE", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "serives.log", serviceLoggerI)

	// 1. Services...
	newfilesvcI := releaseconfigsvc.NewfileConfigSvc(pgsqlI, serviceLoggerI)
	newfilesvcI.Start()

	filetrcsvcI := filetrcsvc.NewfileTrcSvc(pgsqlI, serviceLoggerI)
	filetrcsvcI.Start()

	// 2. Handlers..
	handlerLoggerI := log.New(os.Stdout, "HANDLER:", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "handlers.log", handlerLoggerI)

	releasefilehdlrI := releasefilehdlr.NewReleaseFileHdlr(newfilesvcI, handlerLoggerI)
	filetrackerhdlrI := filetrchdlr.NewReleaseFileHdlr(filetrcsvcI, handlerLoggerI)

	// 3. Start the server...
	activeThreads := &sync.WaitGroup{}
	apiservercallback := apiserver.ApiServerStateCallback{
		Started: func() {
			log.Println("Api server started")
		},
		Stopped: func() {
			log.Println("Api server stopped")
			activeThreads.Done()
		},
	}
	listenaddr := fmt.Sprintf(":%d", serverConfigI.APIServerConfig.Port)

	accessLoggerI := log.New(os.Stdout, "ACCESS:", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "access.log", accessLoggerI)
	handlerMap := apiserver.NewApiServerHandlerMap(accessLoggerI)
	handlerMap.AddHandler("/api/v1/releaseconfigs", releasefilehdlrI)
	handlerMap.AddHandler("/api/v1/filetracker", filetrackerhdlrI)
	apiserver := apiserver.NewApiServer(listenaddr, handlerMap.ToRouter(), accessLoggerI, apiservercallback)
	activeThreads.Add(1)
	starterr := apiserver.Start()
	if starterr != nil {
		log.Fatal("Api server start failed", starterr)
	}
	activeThreads.Wait()

	// 4. The server seems to have stopped...
	stoperr := apiserver.Stop()
	if stoperr != nil {
		log.Fatal("Api server stop failed", stoperr)
	}
	// 5. Stop all the services one-by-one..
	pgsqlI.ClearDBInstance()
	log.Print("Api Server shutdown successfully...")

}
