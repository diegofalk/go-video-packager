package main

import (
	"encoding/json"
	"fmt"
	"github.com/diegofalk/go-video-packager/database"
	"github.com/diegofalk/go-video-packager/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	// get uploaded file
	var uploadedContent storage.UploadedContent
	uploadedContent.Name = mux.Vars(r)["path"]
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
	err = db.SaveContent(content)
	if err != nil {
		panic(err)
	}

	// Get contentID
	id, err := db.GetContentID(content)
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

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("request: %v", requestData)
}
