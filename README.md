[![GoDoc](https://godoc.org/github.com/sqlitebrowser/go-dbhub?status.svg)](https://godoc.org/github.com/sqlitebrowser/go-dbhub)
[![Go Report Card](https://goreportcard.com/badge/github.com/sqlitebrowser/go-dbhub)](https://goreportcard.com/report/github.com/sqlitebrowser/go-dbhub)

A Go library for accessing and using SQLite databases stored remotely on DBHub.io

*This is an early stage work in progress*

What works now:

* Run read-only queries (eg SELECT statements) on databases, returning the results as JSON
* List the names of tables, views, and indexes present in a database
* List the columns present in a table or view, along with their details

Still to do:

* Tests for each function
* Retrieve index details for a database
* Return the list of available databases
* Download a complete database
* Upload a complete database
* Retrieve database commit history details (size, branch, commit list, whatever else is useful)

Example code:

```
// Create a new DBHub.io API object
db, err := dbhub.New("YOUR_API_KEY_HERE")
if err != nil {
    log.Fatal(err)
}

// Retrieve the list of tables in the remote database
tables, err := db.Tables("justinclift", "Join Testing.sqlite")
if err != nil {
    log.Fatal(err)
}

// Display the retrieved list of tables
fmt.Println("Tables:")
for _, j := range tables {
    fmt.Printf("  * %s\n", j)
}

// Run a SQL query on the remote database
r, err := db.Query("justinclift", "Join Testing.sqlite", false,
    `SELECT table1.Name, table2.value
        FROM table1 JOIN table2
        USING (id)
        ORDER BY table1.id`)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Query results:\n\t%v\n", r)
fmt.Println()
```

Output:

```
Tables:
  * table1
  * table2

Query results:
        {[{[Foo 5]} {[Bar 10]} {[Baz 15]} {[Blumph 12.5000]} {[Blargo 8]} {[Batty 3]}]}
```
  
Please try it out, submits PRs, and report any weirdness or bugs you encounter. :smile:
