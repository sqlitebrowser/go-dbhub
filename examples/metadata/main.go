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

	// Retrieve the metadata for the remote database
	meta, err := db.Metadata("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// * Display the retrieved metadata *

	// Display the database branches
	fmt.Println("Branches:")
	for i := range meta.Branches {
		fmt.Printf("  * %s\n", i)
	}
	fmt.Printf("Default branch: %s\n", meta.DefBranch)

	// Display the database releases
	if len(meta.Releases) != 0 {
		fmt.Println("Releases:")
		for i := range meta.Releases {
			fmt.Printf("  * %s\n", i)
		}
	}

	// Display the database tags
	if len(meta.Tags) != 0 {
		fmt.Println("Tags:")
		for i := range meta.Tags {
			fmt.Printf("  * %s\n", i)
		}
	}

	// Display the database commits
	fmt.Println("Commits:")
	for _, j := range meta.Commits {
		fmt.Printf("  * %s, %v\n", j.ID, j.Timestamp.Format(time.RFC1123))
	}

	// Display the web page for the database
	fmt.Printf("Web page: %s\n", meta.WebPage)
}
