package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sqlitebrowser/go-dbhub"
)

func main() {
	// Create a new DBHub.io API object
	db, err := dbhub.New("YOUR_API_KEY_HERE")
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the tags for the remote database
	tags, err := db.Tags("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Display the tags
	fmt.Println("Tags:")
	for i, j := range tags {
		fmt.Printf("  * %s\n", i)
		fmt.Printf("      Commit: %v\n", j.Commit)
		fmt.Printf("      Date: %v\n", j.Date.Format(time.RFC1123))
		fmt.Printf("      Tagger Name: %v\n", j.TaggerName)
		fmt.Printf("      Tagger Email: %v\n", j.TaggerEmail)
		fmt.Printf("      Description: %v\n", j.Description)
	}
}
