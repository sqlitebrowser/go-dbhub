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

	// Retrieve the list of indexes in the remote database
	indexes, err := db.Indexes("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of indexes
	fmt.Println("Indexes:")
	for i, j := range indexes {
		fmt.Printf("  * '%s' on table '%s'\n", i, j)
	}
	fmt.Println()
}
