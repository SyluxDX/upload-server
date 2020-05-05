package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type simpleMsg struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload Hit")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(simpleMsg{Code: 404, Msg: "File not found"})
		return
	}
	defer file.Close()
	// Create a temporary file
	tempFile, err := ioutil.TempFile("temp", "*.upload")
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// send respose
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(simpleMsg{Code: 200, Msg: fmt.Sprintf("%v uploaded to %v", handler.Filename, tempFile.Name())})
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(234)
	json.NewEncoder(w).Encode(simpleMsg{Code: 404, Msg: "random text"})
	// fmt.Fprintf(w, "Welcome to the HomePage!\n")
	log.Println("Endpoint Hit: homePage")
}

func staticUploadFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/testupload", staticUploadFile)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func setupStuff() {
	os.Mkdir("upload", os.ModePerm)
	os.Mkdir("temp", os.ModePerm)
}

func main() {
	setupStuff()
	log.Println("starting server")
	setupRoutes()
}
