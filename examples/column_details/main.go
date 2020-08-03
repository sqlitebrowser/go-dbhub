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

	// Retrieve the column info for a table or view in the remote database
	table := "table1"
	columns, err := db.Columns("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"}, table)
	if err != nil {
		log.Fatal(err)
	}

	// Display the retrieved column details
	fmt.Printf("Columns on table or view '%s':\n", table)
	for _, j := range columns {
		fmt.Printf("  * '%v':\n", j.Name)
		fmt.Printf("      Autoincrement: %v\n", j.Autoinc)
		fmt.Printf("      Cid: %v\n", j.Cid)
		fmt.Printf("      Collation Sequence: %v\n", j.CollSeq)
		fmt.Printf("      Data Type: %v\n", j.DataType)
		fmt.Printf("      Default Value: %v\n", j.DfltValue)
		fmt.Printf("      Not Null: %v\n", j.NotNull)
		fmt.Printf("      Primary Key: %v\n", j.Pk)
	}
	fmt.Println()
}
