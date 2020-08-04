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

	// Retrieve the list of branches of the remote database
	user := "justinclift"
	database := "Marine Litter Survey (Keep Northern Ireland Beautiful).sqlite"
	branches, defaultBranch, err := db.Branches(user, database)
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved list of branches
	fmt.Println("Branches:")
	for branchName, branchDetails := range branches {
		var defaultBranchText string
		if branchName == defaultBranch {
			defaultBranchText = ", default branch"
		}
		fmt.Printf("  * %s (commits: %d%s)\n", branchName, branchDetails.CommitCount, defaultBranchText)
	}
}
