package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sqlitebrowser/go-dbhub"
)

func main() {
	// Create a new DBHub.io API object
	db, err := dbhub.New("YOUR_API_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}

	// Read the database file into memory
	var myDB []byte
	myDB, err = ioutil.ReadFile("example.db")
	if err != nil {
		log.Fatal(err)
	}

	// Prepare any information you want to include with the upload (eg a commit message, etc)
	info := dbhub.UploadInformation{
		CommitMsg: "An example upload",
	}

	// Upload the database
	err = db.Upload("somedb.sqlite", info, &myDB)
	if err != nil {
		log.Fatal(err)
	}

	// Display a success message
	fmt.Println("Database uploaded")
}
