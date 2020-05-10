package main

import (
	"encoding/json"
	"fmt"
	"github.com/diegofalk/go-video-packager/database"
	"github.com/diegofalk/go-video-packager/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type publishResponse struct {
	ContentID string `json:"content_id"`
}

type packageRequest struct {
	ContentID string `json:"content_id"`
	Key       string `json:"key"`
	Kid       string `json:"kid"`
}

type packageResponse struct {
	StreamID string `json:"stream_id"`
}

type streamResponse struct {
	Url string `json:"url"`
	Key string `json:"key"`
	Kid string `json:"kid"`
}

func apiRun(httpListenAddress string) {
	router := mux.NewRouter().StrictSlash(true)
	apiRegister(router)
	log.Fatal(http.ListenAndServe(httpListenAddress, router))
}

func apiRegister(router *mux.Router) {
	router.HandleFunc("/", apiHomeHandler)
	router.HandleFunc("/publish/{path:.*\\.mp4}", apiPublishHandler).Methods("POST")
	router.HandleFunc("/package", apiPackageHandler).Methods("POST")
	router.HandleFunc("/streaminfo/{stream_id}", apiStreaminfoHandler).Methods("GET")
	router.HandleFunc("/stream/{stream_id}/{path:.*\\.mpd}", apiStreamMpdHandler).Methods("GET")
	router.HandleFunc("/stream/{stream_id}/{folder}/{chunk}", apiStreamChunkHandler).Methods("GET")
}

func apiHomeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi, I'm go-video-packager")
}

func apiPublishHandler(w http.ResponseWriter, r *http.Request) {
	// get arrival time
	now := time.Now()
	arrivalTime := strconv.FormatInt(now.Unix(), 10)

	// save uploaded file localy
	fileName := arrivalTime + "_" + mux.Vars(r)["path"] // add timestamp to avoid duplicates
	err := storage.SaveContentFile(r.Body, fileName)
	if err != nil {
		panic(err)
	}
	// Save content model on DB
	content := database.Content{
		Name: fileName,
	}
	id, err := db.InsertContent(content)
	if err != nil {
		panic(err)
	}

	// Write response
	err = json.NewEncoder(w).Encode(publishResponse{id})
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("File " + fileName + " Uploaded successfully")
}

func apiPackageHandler(w http.ResponseWriter, r *http.Request) {
	var requestData packageRequest

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the content by ID exists
	_, err = db.GetContent(requestData.ContentID)
	if err != nil {
		panic(err)
	}

	// create the stream model
	stream := database.Stream{
		ContentID: requestData.ContentID,
		Key:       requestData.Key,
		Kid:       requestData.Kid,
		Status:    "PACKAGING",
		Url:       "",
	}
	streamID, err := db.InsertStream(stream)
	if err != nil {
		panic(err)
	}

	// Write response
	err = json.NewEncoder(w).Encode(packageResponse{streamID})
	if err != nil {
		panic(err)
	}

	select {
	case data <- streamID:
		fmt.Printf("inserted packaging job ID: %s\n", streamID)
	default:
		panic(fmt.Errorf("packaging queue full"))
	}

	w.Header().Set("Content-Type", "application/json")
}

func apiStreaminfoHandler(w http.ResponseWriter, r *http.Request) {
	// get stream id
	streamID := mux.Vars(r)["stream_id"]

	// get stream
	stream, err := db.GetStream(streamID)
	if err != nil {
		panic(err)
	}

	switch stream.Status {
	case "FAILED":
		http.Error(w, "Packaging failed", http.StatusBadRequest)
		return
	case "PACKAGING":
		http.Error(w, "Packaging in progress", http.StatusAccepted)
		return
	case "DONE":
		break
	default:
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}
	// after this point, we only got DONE cases

	// create the stream model
	response := streamResponse{
		Url: stream.Url,
		Key: stream.Key,
		Kid: stream.Kid,
	}

	// Write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
}

func apiStreamMpdHandler(w http.ResponseWriter, r *http.Request) {
	// get uploaded file
	path := mux.Vars(r)["path"]
	streamID := mux.Vars(r)["stream_id"]
	fileName := streamID + "/" + path

	err := storage.LoadStreamFile(w, fileName)
	if err != nil {
		panic(err)
	}
}

func apiStreamChunkHandler(w http.ResponseWriter, r *http.Request) {
	// get uploaded file
	chunk := mux.Vars(r)["chunk"]
	streamID := mux.Vars(r)["stream_id"]
	folder := mux.Vars(r)["folder"]
	fileName := streamID + "/" + folder + "/" + chunk

	err := storage.LoadStreamFile(w, fileName)
	if err != nil {
		panic(err)
	}
}
