package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type Note struct {
	Note   string `json:"note"`
	Id     string `json:"id"`
	Poster int    `json:"poster"`
}

// global declaration?
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

func getAllNotes() []Note {

	if err != nil {
		//Failed to create db with schema
		panic(err)
	}

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
