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

	// Retrieve the list of tables in the remote database
	tables, err := db.Tables("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"})
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of tables
	fmt.Println("Tables:")
	for _, j := range tables {
		fmt.Printf("  * %s\n", j)
	}
}
