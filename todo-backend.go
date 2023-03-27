package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

func noteCallHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		//handle Get, get All Notes
		getAllNotes()

	case "POST":
		//handle POST
		//createUpdateNote(nil)

	case "DELETE":
		//Handle deleting note
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
	getAllNotesByPoster(3)

	http.HandleFunc("/note", noteCallHandler)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}

type Note struct {
	Note   string
	Id     string
	Poster int
}

// globaali declaraatio?
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"note": {
			Name: "note",
			Indexes: map[string]*memdb.IndexSchema{
				"note": {
					Name:    "note",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Note"},
				},
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Id"},
				},
				"poster": {
					Name:    "poster",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Poster"},
				},
			},
		},
	},
}
var db, err = memdb.NewMemDB(schema)

func testAddNotes() {
	// Create a write transaction
	txn := db.Txn(true)

	// Insert some people
	note := []*Note{
		{"This is test note number 1", "f543dd77-826f-4809-a42c-174ac3e5b140", 30},
		{"This is test note number 2", "d4ad9db6-5635-44f9-ad75-cae291bb21b4", 35},
		{"This is test note number 3", "da85f4f2-cdd2-4190-a530-be12e425cd44", 21},
		{"This is test note number 4", "7169a2fe-cc09-11ed-afa1-0242ac120002", 53},
	}
	for _, p := range note {
		if err := txn.Insert("note", p); err != nil {
			panic(err)
		}
	}
	// Commit the transaction
	txn.Commit()
}

func getAllNotes() []Note {
	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// query all the people
	it, err := txn.Get("note", "note")
	if err != nil {
		panic(err)
	}

	fmt.Println("All the notes:")
	AllNotesArray := make([]Note, 0, 1)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Note)
		AllNotesArray = append(AllNotesArray, *p)
		fmt.Printf("  %s note with poster %d\n", p.Note, p.Poster)
	}
	return AllNotesArray
}

func getAllNotesByPoster(posterId int) []Note {
	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// query all the people
	it, err := txn.Get("note", "poster", posterId)
	if err != nil {
		panic(err)
	}

	fmt.Println("All the notes by poster:")
	AllNotesArray := make([]Note, 0, 1)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Note)
		AllNotesArray = append(AllNotesArray, *p)
		fmt.Printf("  %s note with poster %d\n", p.Note, p.Poster)
	}
	return AllNotesArray
}

func createUpdateNote(note Note) {

	txn := db.Txn(true)
	if note.Id != "" {
		note.Id = uuid.New().String()
	}
	if err := txn.Insert("note", note); err != nil {
		panic(err)
	}
}

func deleteNote(uuid string) {
	txn := db.Txn(true)
	txn.Delete("note", Note{Id: uuid})
}

func testGetNotes() {

	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// Lookup by id
	raw, err := txn.First("note", "id", "f543dd77-826f-4809-a42c-174ac3e5b140")
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %s!\n", raw.(*Note).Note)

	deleteNote(raw.(*Note).Id)

	// List all the notes
	it, err := txn.Get("note", "poster", 30)
	if err != nil {
		panic(err)
	}

	fmt.Println("All the notes:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Note)
		fmt.Printf("  %s\n", p.Note)
	}
}
