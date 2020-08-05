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

	// Retrieve the databases in your account
	databases, err := db.Databases()
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of databases
	fmt.Println("Databases:")
	for _, j := range databases {
		fmt.Printf("  * %s\n", j)
	}
}
