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

		// TODO: use config
		url := "http://localhost:8081/stream/" + stream.ID.Hex() + "/" + stream.ID.Hex() + ".mpd"
		err = db.UpdateStreamUrl(streamID, url)
		if err != nil {
			panic(err)
		}

		fmt.Printf("processed %s\n", streamID)
		time.Sleep(5 * time.Second)
	}
	fmt.Printf("packager ended")
}

func doPackage(stream database.Stream, content database.Content) error {
	// define all paths
	contentPath := "content/" + content.Name
	streamFolder := "stream/" + stream.ID.Hex() + "/"
	mpdPath := streamFolder + stream.ID.Hex() + ".mpd"

	//// decode base64 keys
	//kidHex := base64toHexString(stream.Kid)
	//keyHex := base64toHexString(stream.Key)

	// command attribs
	app := "./packager-osx"
	audioSegments := "in=" + contentPath + ",stream=audio,init_segment=" + streamFolder +
		"audio/init.mp4,segment_template=" + streamFolder + "audio/$Number$.m4s,drm_label=ALL"
	videoSegments := "in=" + contentPath + ",stream=video,init_segment=" + streamFolder +
		"video/init.mp4,segment_template=" + streamFolder + "video/$Number$.m4s,drm_label=ALL"
	//encryptionKeys := "label=ALL:key_id=" + kidHex + ":key=" + keyHex
	//scheme := "cbcs"
	baseUrl := ""

	//opt0 := "--protection_scheme"
	//opt1 := "--enable_raw_key_encryption"
	opt2 := "--mpd_output"
	//opt3 := "--keys"
	opt4 := "--base_urls"

	// execute command
	output, err := exec.Command(app, audioSegments, videoSegments /*opt0, scheme, opt1, opt3, encryptionKeys,*/, opt2, mpdPath, opt4, baseUrl).CombinedOutput()
	if err != nil {
		fmt.Printf("Command error: %s\n", output)
		return err
	}
	return nil
}

func base64toHexString(base64Key string) string {
	binaryKey, _ := base64.RawStdEncoding.DecodeString(base64Key[:len(base64Key)])
	return hex.EncodeToString(binaryKey)
}
