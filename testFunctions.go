package main

import "fmt"

func testAddNotes() {
	// Create a write transaction
	txn := db.Txn(true)

	// Insert some people
	note := []*Note{
		{"This is test note number 1", "f543dd77-826f-4809-a42c-174ac3e5b140", 30},
		{"This is test note number 2", "d4ad9db6-5635-44f9-ad75-cae291bb21b4", 35},
		{"This is test note number 3", "da85f4f2-cdd2-4190-a530-be12e425cd44", 21},
		{"This is test note number 4", "7169a2fe-cc09-11ed-afa1-0242ac120002", 53},
		{"This is test note number 65", "7169a2fe-cc09-11ed-afa1-0242ac120003", 30},
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
