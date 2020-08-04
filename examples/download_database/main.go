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

	// Retrieve the remote database file
	dbName := "Join Testing.sqlite"
	dbStream, err := db.Download("justinclift", dbName, dbhub.Identifier{})
	if err != nil {
		log.Fatal(err)
	}

	// Save the database file in the current directory
	buf, err := ioutil.ReadAll(dbStream)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dbName, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved database file as '%s'\n", dbName)
}
