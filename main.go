package main

import (
	"github.com/diegofalk/go-video-packager/database"
	"github.com/op/go-logging"
	"os"
	"runtime"
)

// TODO: move to config file
const maxQueuedPackagingJobs = 10	// max supported queued packaging jobs
const logFormat = `%{color}%{time:2006-01-02 15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`

// shared instances (all thread safe)
var db *database.Mongodb
var queue chan string
var log *logging.Logger

func main() {
	// use as much goroutines as CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// init log
	var format = logging.MustStringFormatter(logFormat)
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	log = logging.MustGetLogger("log")
	log.SetBackend(logging.AddModuleLevel(backendFormatter))

	// init queue
	queue = make(chan string, maxQueuedPackagingJobs)

	// init mongo
	db = database.NewMongodb()
	err := db.Init()
	if err != nil {
		log.Critical("Mongo init error: %s", err.Error())
		panic(err)
	}

	// start packager (in a new goroutine)
	go packagerRun()

	// start API
	apiRun(":8081")
}
