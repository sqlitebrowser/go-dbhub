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
	indexes, err := db.Indexes("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"})
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of indexes
	fmt.Println("Indexes:")
	for _, j := range indexes {
		fmt.Printf("  * '%s' on table '%s'\n", j.Name, j.Table)
		for _, l := range j.Columns {
			fmt.Printf("      Column name: %v\n", l.Name)
			fmt.Printf("      Column ID: %v\n", l.CID)
		}
	}
}
