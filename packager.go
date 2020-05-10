package main

import (
	"encoding/base64"
	"encoding/hex"
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

	// decode base64 keys
	//kidHex := base64toHexString(stream.Kid)
	//fmt.Println(kidHex)
	//keyHex := base64toHexString(stream.Key)
	//fmt.Println(keyHex)

	// command attribs
	app := "./packager-osx"
	audioSegments := "in=" + contentPath + ",stream=audio,init_segment=" + streamFolder +
		"audio/init.mp4,segment_template=" + streamFolder + "audio/$Number$.m4s,drm_label=ALL"
	videoSegments := "in=" + contentPath + ",stream=video,init_segment=" + streamFolder +
		"video/init.mp4,segment_template=" + streamFolder + "video/$Number$.m4s,drm_label=ALL"
	//encryptionKeys := "label=ALL:key_id=" + kidHex + ":key=" + keyHex

	//opt1 := "--enable_raw_key_encryption"
	opt2 := "--mpd_output"
	//opt3 := "--keys"

	// execute command
	cmd := exec.Command(app, audioSegments, videoSegments, opt2, mpdPath)
	fmt.Print(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
func base64toHexString(base64Key string) string {
	fmt.Println(base64Key)
	fmt.Println(len(base64Key))
	binaryKey, _ := base64.RawStdEncoding.DecodeString(base64Key[:len(base64Key)])
	return hex.EncodeToString(binaryKey)
}
