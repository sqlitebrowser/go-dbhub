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

	// Run a query on the remote database
	showBlobs := false
	r, err := db.Query("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"},
		showBlobs, `SELECT table1.Name, table2.value
			FROM table1 JOIN table2
			USING (id)
			ORDER BY table1.id`)
	if err != nil {
		log.Fatal(err)
	}

	// Display the query result (without unmarshalling)
	fmt.Printf("Query results:\n\t%v\n", r)
	fmt.Println()
}
