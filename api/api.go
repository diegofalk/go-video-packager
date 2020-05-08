package api

import (
	"fmt"
	"github.com/diegofalk/go-video-packager/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func Run(httpListenAddress string) {
	router := mux.NewRouter().StrictSlash(true)
	register(router)
	log.Fatal(http.ListenAndServe(httpListenAddress, router))
}

func register(router *mux.Router) {
	router.HandleFunc("/", apiHomeHandler)
	router.HandleFunc("/publish/{path:.*\\.mp4}", apiPublishHandler).Methods("POST")
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

	io.WriteString(w, "File " + uploadedContent.Name + " Uploaded successfully")
	fmt.Println("File " + uploadedContent.Name + " Uploaded successfully")
}

func apiHomeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi, I'm go-video-packager")
}