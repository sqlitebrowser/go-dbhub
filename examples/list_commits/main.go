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

	// Retrieve the commits for the remote database
	commits, err := db.Commits("justinclift", "Join Testing.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Display the commits
	fmt.Println("Commits:")
	for i, j := range commits {
		fmt.Printf("  * %s\n", i)
		if j.CommitterName != "" {
			fmt.Printf("      Committer Name: %v\n", j.CommitterName)
		}
		if j.CommitterEmail != "" {
			fmt.Printf("      Committer Email: %v\n", j.CommitterEmail)
		}
		fmt.Printf("      Timestamp: %v\n", j.Timestamp.Format(time.RFC1123))
		fmt.Printf("      Author Name: %v\n", j.AuthorName)
		fmt.Printf("      Author Email: %v\n", j.AuthorEmail)
		if j.Message != "" {
			fmt.Printf("      Message: %v\n", j.Message)
		}
		if j.Parent == "" {
			fmt.Println("      Parent: NONE")
		} else {
			fmt.Printf("      Parent: %v\n", j.Parent)
		}
	}
}
