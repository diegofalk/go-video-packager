package main

import (
	"github.com/diegofalk/go-video-packager/database"
)

var db *database.Mongodb

func main() {
	db = database.NewMongodb()
	err := db.Init()
	if err != nil {
		panic(err)
	}
	apiRun(":8081")
}
