package main

import (
	"github.com/diegofalk/go-video-packager/database"
	"github.com/op/go-logging"
	"os"
	"runtime"
)

const maxQueuedPackagingJobs = 10

var db *database.Mongodb
var data = make(chan string, maxQueuedPackagingJobs)
var log = logging.MustGetLogger("log")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	// init log
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	log.SetBackend(logging.AddModuleLevel(backendFormatter))

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init mongo
	db = database.NewMongodb()
	err := db.Init()
	if err != nil {
		log.Critical("Mongo init error: %s", err.Error())
		panic(err)
	}

	// create comunication queue
	go packagerRun()
	apiRun(":8081")
}
