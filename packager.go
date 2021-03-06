package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/diegofalk/go-video-packager/database"
	"os/exec"
)

// TODO: move to config file
const localContentPath string = "content/"
const localStreamsPath string = "stream/"
const shakaBinary string = "bin/packager-linux" // shaka packager
const baseUrl string = "http://localhost:8081/"

func packagerRun() {
	for {
		// get stream id
		streamID := <-queue // blocked, waiting for new streams

		// run packaging task on a new goroutine
		go packageTask(streamID)
	}
	log.Infof("packager ended")
}

func packageTask(streamID string) {
	// get stream
	stream, err := db.GetStream(streamID)
	if err != nil {
		log.Errorf("Stream %s not found: %s", streamID, err.Error())
	}

	// get content
	content, err := db.GetContent(stream.ContentID)
	if err != nil {
		log.Errorf("Content %s not found: %s", stream.ContentID, err.Error())
		db.UpdateStreamStatus(streamID, "FAILED")
		return
	}

	// do package
	err = doPackage(stream, content)
	if err != nil {
		log.Errorf("Packaging error: %s", err.Error())
		db.UpdateStreamStatus(streamID, "FAILED")
		return
	}

	// TODO: use config
	url := baseUrl + localStreamsPath + streamID + "/" + streamID + ".mpd"
	err = db.UpdateStreamUrl(streamID, url)
	if err != nil {
		log.Errorf("Error updating URL: %s", err.Error())
		db.UpdateStreamStatus(streamID, "FAILED")
		return
	}
	// update status
	err = db.UpdateStreamStatus(streamID, "DONE")
	if err != nil {
		log.Errorf("Error updating status: %s", err.Error())
	}

	log.Infof("packaging job ended for streamID: %s", streamID)
}

func doPackage(stream database.Stream, content database.Content) error {
	// define all paths
	contentPath := localContentPath + content.Name
	streamFolder := localStreamsPath + stream.ID.Hex() + "/"
	mpdPath := streamFolder + stream.ID.Hex() + ".mpd"

	// segment templates
	audioSegments := "in=" + contentPath + ",stream=audio,init_segment=" + streamFolder +
		"audio/init.mp4,segment_template=" + streamFolder + "audio/$Number$.m4s,drm_label=ALL"
	videoSegments := "in=" + contentPath + ",stream=video,init_segment=" + streamFolder +
		"video/init.mp4,segment_template=" + streamFolder + "video/$Number$.m4s,drm_label=ALL"

	var args []string

	// args definition
	args = append(args, audioSegments, videoSegments)
	args = append(args,"--mpd_output", mpdPath)
	args = append(args,"--base_urls", "")
	args = append(args,"--generate_static_live_mpd")

	// add extra args if it is encrypted
	if len(stream.Key) != 0 {

		// decode keys
		kidHex := base64toHexString(stream.Kid)
		keyHex := base64toHexString(stream.Key)

		encryptionKeys := "label=ALL:key_id=" + kidHex + ":key=" + keyHex

		args = append(args,"--protection_scheme", "cbcs")
		args = append(args,"--enable_raw_key_encryption")
		args = append(args,"--keys", encryptionKeys)
	}

	// run command
	output, err  := exec.Command(shakaBinary, args...).CombinedOutput()
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
