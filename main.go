package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type simpleMsg struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}
type confimation struct {
	Name  string
	Lines []string
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("myFile")
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

	st := string(fileBytes)
	fmt.Println(st)

	// request Confirmation page
	confirm := confimation{tempFile.Name(), []string{}}
	// check uploaded file
	confirm.Lines = parseTesters(fileBytes)

	template, err := template.ParseFiles("confirm.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, confirm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processFile(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("filename")
	filename := filepath.Base(path)
	// move file to upload folder
	os.Rename(path, filepath.Join("upload", filename))
	// Serve final page
	http.ServeFile(w, r, "thankyou.html")
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "upload.html")
	})
	http.HandleFunc("/confirm", uploadFile)
	http.HandleFunc("/thankyou", processFile)

	url := "127.0.0.1"
	port := 8080
	log.Printf("Serving on %[1]s port %[2]d (http://%[1]s:%[2]d/)\n", url, port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", url, port), nil)
}

func setupStuff() {
	os.Mkdir("upload", os.ModePerm)
	os.Mkdir("temp", os.ModePerm)
}

func main() {
	log.Println("Preping server")
	setupStuff()
	setupRoutes()
}
