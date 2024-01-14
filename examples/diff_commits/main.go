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

	// Retrieve the differences between two commits of the same database
	user := "justinclift"
	database := "DB4S download stats.sqlite"
	commit1 := dbhub.Identifier{CommitID: "34cbeebfc347a09406707f4220cd40f60778692523d2e7d227ccd92f4125c9ea"}
	commit2 := dbhub.Identifier{CommitID: "bc6a07955811d86db79e9b4f7fdc3cb2360d40da793066510d792588a8bf8de2"}
	mergeMode := dbhub.PreservePkMerge
	diffs, err := db.Diff(user, database, commit1, "", "", commit2, mergeMode)
	if err != nil {
		log.Fatal(err)
	}

	// Display the SQL statements needed to turn the first version of the database into the second.
	// This should produce a similar output to the sqldiff utility.
	fmt.Printf("SQL statements for turning the first version into the second:\n")
	for _, i := range diffs.Diff { // There is one item for each modified database object
		// Print schema changes to this object if there are any
		if i.Schema != nil {
			fmt.Printf("%s\n", i.Schema.Sql)
		}

		// Loop over all data changes in this object if there are any
		for _, j := range i.Data {
			fmt.Printf("%s\n", j.Sql)
		}
	}
	fmt.Println()
}
