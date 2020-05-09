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

func apiRun(httpListenAddress string) {
	router := mux.NewRouter().StrictSlash(true)
	apiRegister(router)
	log.Fatal(http.ListenAndServe(httpListenAddress, router))
}

func apiRegister(router *mux.Router) {
	router.HandleFunc("/", apiHomeHandler)
	router.HandleFunc("/publish/{path:.*\\.mp4}", apiPublishHandler).Methods("POST")
	router.HandleFunc("/package", apiPackageHandler).Methods("POST")
}

func apiHomeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi, I'm go-video-packager")
}

func apiPublishHandler(w http.ResponseWriter, r *http.Request) {
	// get arrival time
	now := time.Now()
	arrivalTime := strconv.FormatInt(now.Unix(), 10)

	// get uploaded file
	var uploadedContent storage.UploadedContent
	uploadedContent.Name = arrivalTime + "_" + mux.Vars(r)["path"] // add timestamp to avoid duplicates
	r.Body.Read(uploadedContent.Data)
	defer r.Body.Close()

	// save it locally
	err := uploadedContent.Save()
	if err != nil {
		panic(err)
	}

	// Save content model on DB
	content := database.Content{
		Name: uploadedContent.Name,
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

	fmt.Println("File " + uploadedContent.Name + " Uploaded successfully")
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
	}
	id, err := db.InsertStream(stream)
	if err != nil {
		panic(err)
	}

	fmt.Printf("stream: %v. ID: %s", stream, id)
}
