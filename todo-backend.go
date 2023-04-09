package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func noteCallHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		//handle Get, get Notes
		id := strings.TrimPrefix(req.URL.Path, "/note/")
		if id != "" {
			intId, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("error durin conversion to int")
			}
			//get posters notes
			getAllNotesByPoster(intId)

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

		w.Header().Set("Content-Type", "application/json")
		note := Note{
			Note:   n.Note,
			Id:     n.Id,
			Poster: n.Poster,
		}
		json.NewEncoder(w).Encode(note)

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
