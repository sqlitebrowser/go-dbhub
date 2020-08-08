package main

import (
	"fmt"
	"log"

	"github.com/sqlitebrowser/go-dbhub"
)

func main() {
	// Create a new DBHub.io API object
	db, err := dbhub.New("YOUR_API_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the metadata for the remote database
	wp, err := db.Webpage("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Display the web page for the database
	fmt.Printf("Web page: '%s'\n", wp.WebPage)
}
