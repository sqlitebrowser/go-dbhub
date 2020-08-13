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

	// Delete a remote database
	dbName := "Join Testing.sqlite"
	err = db.Delete(dbName)
	if err != nil {
		log.Fatal(err)
	}

	// Display a success message
	fmt.Printf("Database '%s' deleted\n", dbName)
}
