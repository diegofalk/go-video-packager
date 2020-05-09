package main

import (
	"github.com/diegofalk/go-video-packager/database"
	"runtime"
)

const maxQueuedPackagingJobs = 10

var db *database.Mongodb
var data = make(chan string, maxQueuedPackagingJobs)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init mongo
	db = database.NewMongodb()
	err := db.Init()
	if err != nil {
		panic(err)
	}

	// create comunication queue
	go packagerRun()
	apiRun(":8081")
}
