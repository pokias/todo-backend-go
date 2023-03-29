package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func noteCallHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		//handle Get, get Notes
		id := strings.TrimPrefix(req.URL.Path, "/note/")
		if id != "" {
			//get posters notes
		} else {
			//all notes
			getAllNotes()
		}

	case "POST":
		//Create or update note
		var n Note
		err := json.NewDecoder(req.Body).Decode(&n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createUpdateNote(n)
		req.Response.StatusCode = 200

	case "DELETE":
		//Handle deleting note
		deleteNote(strings.TrimPrefix(req.URL.Path, "/note/"))
	}
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	testAddNotes()
	testGetNotes()
	getAllNotes()
	getAllNotesByPoster(30)

	http.HandleFunc("/note/", noteCallHandler)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}
