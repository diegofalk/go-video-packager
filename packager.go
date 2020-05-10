package main

import (
	"fmt"
	"github.com/diegofalk/go-video-packager/database"
	"os/exec"
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
		err = db.UpdateStreamStatus(streamID, "DONE")
		if err != nil {
			panic(err)
		}
		fmt.Printf("processed %s\n", streamID)
		time.Sleep(5 * time.Second)
	}
}

func doPackage(stream database.Stream, content database.Content) error {
	// define all paths
	contentPath := "content/" + content.Name
	streamFolder := "stream/" + stream.ID.Hex() + "/"
	mpdPath := streamFolder + stream.ID.Hex() + ".mpd"

	// command attribs
	app := "./packager-osx"
	audioSegments := "in=" + contentPath + ",stream=audio,init_segment=" + streamFolder + "audio/init.mp4,segment_template=" + streamFolder + "audio/$Number$.m4s"
	videoSegments := "in=" + contentPath + ",stream=video,init_segment=" + streamFolder + "video/init.mp4,segment_template=" + streamFolder + "video/$Number$.m4s"
	opt1 := "--mpd_output"

	// execute command
	cmd := exec.Command(app, audioSegments, videoSegments, opt1, mpdPath)
	fmt.Print(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
