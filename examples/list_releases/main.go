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

	// Retrieve the release info for the remote database
	rels, err := db.Releases("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Display the release info
	fmt.Println("Releases:")
	for i, j := range rels {
		fmt.Printf("  * %s\n", i)
		fmt.Printf("      Commit: %v\n", j.Commit)
		fmt.Printf("      Date: %v\n", j.Date.Format(time.RFC1123))
		fmt.Printf("      Size: %v bytes\n", j.Size)
		fmt.Printf("      Releaser Name: %v\n", j.ReleaserName)
		fmt.Printf("      Releaser Email: %v\n", j.ReleaserEmail)
		fmt.Printf("      Description: %v\n", j.Description)
	}
}
