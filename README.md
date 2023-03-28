[![GoDoc](https://godoc.org/github.com/sqlitebrowser/go-dbhub?status.svg)](https://godoc.org/github.com/sqlitebrowser/go-dbhub)
[![Go Report Card](https://goreportcard.com/badge/github.com/sqlitebrowser/go-dbhub)](https://goreportcard.com/report/github.com/sqlitebrowser/go-dbhub)

A Go library for accessing and using SQLite databases stored remotely on DBHub.io

### What works now

* (Experimental) Upload, delete, and list your "Live" databases
* (Experimental) Execute INSERT/UPDATE/DELETE statements on your "Live" databases
* (Experimental) List the tables, views, indexes, and columns in your "Live" databases
* Run read-only queries (eg SELECT statements) on databases, returning the results as JSON
* Upload and download your databases
* List the databases in your account
* List the tables, views, and indexes present in a database
* List the columns in a table, view or index, along with their details
* List the branches, releases, tags, and commits for a database
* Generate diffs between two databases, or database revisions
* Download the database metadata (size, branches, commit list, etc.)
* Retrieve the web page URL of a database

### Still to do

* Have the backend server correctly use the incoming branch, release, and tag information
* Tests for each function
* Investigate what would be needed for this to work through the Go SQL API
  * Probably need to improve the Query approach, to at least support placeholders and argument parameters
* Anything else people suggest and seems like a good idea :smile:

### Requirements

* [Go](https://golang.org/dl/) version 1.17 or above
* A DBHub.io API key
  * These can be generated in your [Settings](https://dbhub.io/pref) page, when logged in.

### Example code

#### Create a new DBHub.io API object

```
db, err := dbhub.New("YOUR_API_KEY_HERE")
if err != nil {
    log.Fatal(err)
}
```

#### Retrieve the list of tables in a remote database
```
// Run the `Tables()` function on the new API object
tables, err := db.Tables("justinclift", "Join Testing.sqlite", dbhub.Identifier{Branch: "master"})
if err != nil {
    log.Fatal(err)
}

// Display the retrieved list of tables
fmt.Println("Tables:")
for _, j := range tables {
    fmt.Printf("  * %s\n", j)
}
```

##### Output
```
Tables:
  * table1
  * table2
```

#### Run a SQL query on a remote database
```
// Do we want to display BLOBs as base64?
showBlobs := false

// Run the query
result, err := db.Query("justinclift", "Join Testing.sqlite",
    dbhub.Identifier{ Branch: "master" }, showBlobs,
    `SELECT table1.Name, table2.value
    FROM table1 JOIN table2
    USING (id)
    ORDER BY table1.id`)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Query results (JSON):\n\t%v\n", result)
fmt.Println()
```

##### Output
```
Query results (JSON):
        {[{[Foo 5]} {[Bar 10]} {[Baz 15]} {[Blumph 12.5000]} {[Blargo 8]} {[Batty 3]}]}
```

#### Generate and display the difference between two commits of a remote database
```
// The databases we want to see differences for
db1Owner := "justinclift"
db1Name := "Join Testing.sqlite"
db1Commit := dbhub.Identifier{
    CommitID: "c82ba65add364427e9af3f540be8bf98e8cd6bdb825b07c334858e816c983db0" }
db2Owner := ""
db2Name := ""
db2Commit := dbhub.Identifier{
    CommitID: "adf78104254ece17ff40dab80ae800574fa5d429a4869792a64dcf2027cd9cd9" }

// Create the diff
diffs, err := db.Diff(db1Owner, db1Name, db1Commit, db2Owner, db2Name, db2Commit,
    dbhub.PreservePkMerge)
if err != nil {
    log.Fatal(err)
}

// Display the diff
fmt.Printf("SQL statements for turning the first commit into the second:\n")
for _, i := range diffs.Diff {
    if i.Schema != nil {
        fmt.Printf("%s\n", i.Schema.Sql)
    }
    for _, j := range i.Data {
        fmt.Printf("%s\n", j.Sql)
    }
}
```

##### Output
```
SQL statements for turning the first commit into the second:
CREATE VIEW joinedView AS
SELECT table1.Name, table2.value
FROM table1 JOIN table2
USING (id)
ORDER BY table1.id;
```

### Further examples

* [SQL Query](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/sql_query/main.go) - Run a SQL query, return the results as JSON
* [List databases](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_databases/main.go) - List the databases present in your account
* [List tables](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_tables/main.go) - List the tables present in a database
* [List views](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_views/main.go) - List the views present in a database
* [List indexes](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_indexes/main.go) - List the indexes present in a database
* [Retrieve column details](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/column_details/main.go) - Retrieve the details of columns in a table
* [List branches](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_branches/main.go) - List all branches of a database
* [List releases](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_releases/main.go) - Display the releases for a database
* [List tags](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_tags/main.go) - Display the tags for a database
* [List commits](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/list_commits/main.go) - Display the commits for a database
* [Generate diff between two revisions](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/diff_commits/main.go) - Figure out the differences between two databases or two versions of one database
* [Upload database](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/upload/main.go) - Upload a new database file
* [Download database](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/download_database/main.go) - Download the complete database file
* [Delete database](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/delete_database/main.go) - Delete a database
* [Retrieve metadata](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/metadata/main.go) - Download the database metadata (size, branches, commit list, etc)
* [Web page](https://github.com/sqlitebrowser/go-dbhub/blob/master/examples/webpage/main.go) - Get the URL of the database file in the webUI.  eg. for web browsers

Please try it out, submits PRs to extend or fix things, and report any weirdness or bugs you encounter. :smile:
