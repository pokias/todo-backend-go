package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-memdb"
)

func noteCallHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		//handle Gett

	case "POST":
		//handle POST

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

	//initialize db
	initialize()
	testAddNotes()
	testGetNotes()

	http.HandleFunc("/note", noteCallHandler)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}

type Note struct {
	Note   string
	Poster int
	Id     int
}

var db memdb.MemDB

func initialize() {

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"note": &memdb.TableSchema{
				Name: "note",
				Indexes: map[string]*memdb.IndexSchema{
					"note": &memdb.IndexSchema{
						Name:    "note",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Note"},
					},
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"poster": &memdb.IndexSchema{
						Name:    "poster",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Poster"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	_ = db
}

func testAddNotes() {
	// Create a write transaction
	txn := db.Txn(true)

	// Insert some people
	note := []*Note{
		{"This is test note number 1", 30, 30},
		{"This is test note number 2", 31, 35},
		{"This is test note number 3", 32, 21},
		{"This is test note number 4", 33, 53},
	}
	for _, p := range note {
		if err := txn.Insert("note", p); err != nil {
			panic(err)
		}
	}
	// Commit the transaction
	txn.Commit()
}

func testGetNotes() {

	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// Lookup by id
	raw, err := txn.First("note", "id", 53)
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %s!\n", raw.(*Note).Note)

	// List all the notes
	it, err := txn.Get("note", "poster", "id")
	if err != nil {
		panic(err)
	}

	fmt.Println("All the notes:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Note)
		fmt.Printf("  %s\n", p.Note)
	}

	// Range scan over people with ages between 25 and 35 inclusive
	it, err = txn.LowerBound("note", "poster", 25)
	if err != nil {
		panic(err)
	}

	fmt.Println("People aged 25 - 35:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Note)
		if p.Id > 35 {
			break
		}
		fmt.Printf("  %s is aged %d\n", p.Note, p.Poster)
	}
	// Output:
	// Hello Joe!
	// All the people:
	//   Dorothy
	//   Joe
	//   Lucy
	//   Tariq
	// People aged 25 - 35:
	//   Joe is aged 30
	//   Lucy is aged 35
}
