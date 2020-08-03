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

	// Retrieve the list of views in the remote database
	views, err := db.Views("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"})
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of views
	fmt.Println("Views:")
	for _, j := range views {
		fmt.Printf("  * %s\n", j)
	}
	fmt.Println()
}
