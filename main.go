package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

type todo struct {
	Title  string
	Phrase string
	List   []string
	Table  []string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	lst := []string{"itqm", "asd", "qwe"}
	tbl := []string{"line1", "liune 2 is basdmasd", "iulxcvlkqwe lkjzs", "qwezxcqawe"}
	context := todo{"This is a title", "A paragraph is such a line of text", lst, tbl}
	template, err := template.ParseFiles("test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
