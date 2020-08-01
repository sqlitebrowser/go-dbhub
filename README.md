A Go library for accessing and using SQLite databases on DBHub.io

*This is an early stage work in progress*

What works now:

* Running any read-only query (eg SELECT statements) on databases, returning the results
* Listing the tables, views, indexes, and columns present in a database

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
```

Output:

```
Tables:
  * table1
  * table2
```
  
Please try it out, and report any weirdness or bugs you encounter. :smile:
