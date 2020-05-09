package main

import (
	"fmt"
	"github.com/diegofalk/go-video-packager/database"
	"time"
)

func packagerRun() {
	for {

		// get stream id
		streamID := <-data

		// get stream
		stream, err := db.GetStream(streamID)
		if err != nil {
			panic(err)
		}

		// get content
		content, err := db.GetContent(stream.ContentID)
		if err != nil {
			panic(err)
		}

		// do package
		err = doPackage(stream, content)
		if err != nil {
			panic(err)
		}

		// update stream
		err = doPackage(stream, content)
		if err != nil {
			panic(err)
		}

		// update stream
		err = db.UpdateStreamStatus(streamID, "DONE")
		if err != nil {
			panic(err)
		}
		fmt.Printf("processed %s\n", streamID)
		time.Sleep(5 * time.Second)
	}
}

func doPackage(stream database.Stream, content database.Content) error {
	return nil
}
