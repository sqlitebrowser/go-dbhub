package dbhub

// A Go library for working with databases on DBHub.io

import (
	"encoding/base64"
	"fmt"
	"net/url"

	com "github.com/sqlitebrowser/dbhub.io/common"
)

const (
	version = "0.0.2"
)

// New creates a new DBHub.io connection object.  It doesn't connect to DBHub.io to do this.  Connection only occurs
// when subsequent functions (eg Query()) are called.
func New(key string) (Connection, error) {
	c := Connection{
		APIKey: key,
		Server: "https://api.dbhub.io",
	}
	return c, nil
}

// ChangeAPIKey updates the API key used for authenticating with DBHub.io.
func (c *Connection) ChangeAPIKey(k string) {
	c.APIKey = k
}

// ChangeServer changes the address for communicating with DBHub.io.  Useful for testing and development.
func (c *Connection) ChangeServer(s string) {
	c.Server = s
}

// Branches returns a list of all available branches of a database along with the name of the default branch
func (c Connection) Branches(dbowner, dbname string) (branches map[string]com.BranchEntry, defaultBranch string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of branches and the default branch
	var response com.BranchListResponseContainer
	queryUrl := c.Server + "/v1/branches"
	err = sendRequest(queryUrl, data, &response)

	// Extract information for return values
	branches = response.Entries
	defaultBranch = response.Default
	return
}

// Columns returns the column information for a given table or view
func (c Connection) Columns(dbOwner, dbName string, ident Identifier, table string) (columns []com.APIJSONColumn, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)
	data.Set("table", table)

	// Fetch the list of columns
	queryUrl := c.Server + "/v1/columns"
	err = sendRequest(queryUrl, data, &columns)
	return
}

// Diff returns the differences between two commits of two databases, or if the details on the second database are left empty,
// between two commits of the same database. You can also specify the merge strategy used for the generated SQL statements.
func (c Connection) Diff(dbOwnerA, dbNameA string, identA Identifier, dbOwnerB, dbNameB string, identB Identifier, merge MergeStrategy) (diffs com.Diffs, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner_a", dbOwnerA)
	data.Set("dbname_a", dbNameA)
	if identA.Branch != "" {
		data.Set("branch_a", identA.Branch)
	}
	if identA.CommitID != "" {
		data.Set("commit_a", identA.CommitID)
	}
	if identA.Release != "" {
		data.Set("release_a", identA.Release)
	}
	if identA.Tag != "" {
		data.Set("tag_a", identA.Tag)
	}
	data.Set("dbowner_b", dbOwnerB)
	data.Set("dbname_b", dbNameB)
	if identB.Branch != "" {
		data.Set("branch_b", identB.Branch)
	}
	if identB.CommitID != "" {
		data.Set("commit_b", identB.CommitID)
	}
	if identB.Release != "" {
		data.Set("release_b", identB.Release)
	}
	if identB.Tag != "" {
		data.Set("tag_b", identB.Tag)
	}
	if merge == PreservePkMerge {
		data.Set("merge", "preserve_pk")
	} else if merge == NewPkMerge {
		data.Set("merge", "new_pk")
	} else {
		data.Set("merge", "none")
	}

	// Fetch the diffs
	queryUrl := c.Server + "/v1/diff"
	err = sendRequest(queryUrl, data, &diffs)
	return
}

// Indexes returns the list of indexes present in the database, along with the table they belong to
func (c Connection) Indexes(dbOwner, dbName string, ident Identifier) (idx map[string]string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of indexes
	queryUrl := c.Server + "/v1/indexes"
	err = sendRequest(queryUrl, data, &idx)
	return
}

// PrepareVals creates a url.Values container holding the API key, database owner, name, and database identifier.  The
// url.Values container is then used for the requests to DBHub.io.
func (c Connection) PrepareVals(dbOwner, dbName string, ident Identifier) (data url.Values) {
	// Prepare the API parameters
	data = url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbOwner)
	data.Set("dbname", dbName)
	if ident.Branch != "" {
		data.Set("branch", ident.Branch)
	}
	if ident.CommitID != "" {
		data.Set("commit", ident.CommitID)
	}
	if ident.Release != "" {
		data.Set("release", ident.Release)
	}
	if ident.Tag != "" {
		data.Set("tag", ident.Tag)
	}
	return
}

// Query runs a SQL query (SELECT only) on the chosen database, returning the results.
// The "blobBase64" boolean specifies whether BLOB data fields should be base64 encoded in the output, or just skipped
// using an empty string as a placeholder.
func (c Connection) Query(dbOwner, dbName string, ident Identifier, blobBase64 bool, sql string) (out Results, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)
	data.Set("sql", base64.StdEncoding.EncodeToString([]byte(sql)))

	// Run the query on the remote database
	var returnedData []com.DataRow
	queryUrl := c.Server + "/v1/query"
	err = sendRequest(queryUrl, data, &returnedData)
	if err != nil {
		return
	}

	// Loop through the results, converting it to a more concise output format
	for _, j := range returnedData {

		// Construct a single row
		var oneRow ResultRow
		for _, l := range j {
			switch l.Type {
			case com.Float, com.Integer, com.Text:
				// Float, integer, and text fields are added to the output
				oneRow.Fields = append(oneRow.Fields, fmt.Sprint(l.Value))
			case com.Binary:
				// BLOB data is optionally Base64 encoded, or just skipped (using an empty string as placeholder)
				if blobBase64 {
					// Safety check. Make sure we've received a string
					if s, ok := l.Value.(string); ok {
						oneRow.Fields = append(oneRow.Fields, base64.StdEncoding.EncodeToString([]byte(s)))
					} else {
						oneRow.Fields = append(oneRow.Fields, fmt.Sprintf("unexpected data type '%T' for returned BLOB", l.Value))
					}
				} else {
					oneRow.Fields = append(oneRow.Fields, "")
				}
			default:
				// All other value types are just output as an empty string (for now)
				oneRow.Fields = append(oneRow.Fields, "")
			}
		}

		// Add the row to the output list
		out.Rows = append(out.Rows, oneRow)
	}
	return
}

// Tables returns the list of tables in the database
func (c Connection) Tables(dbOwner, dbName string, ident Identifier) (tbl []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of tables
	queryUrl := c.Server + "/v1/tables"
	err = sendRequest(queryUrl, data, &tbl)
	return
}

// Views returns the list of views in the database
func (c Connection) Views(dbOwner, dbName string, ident Identifier) (views []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of views
	queryUrl := c.Server + "/v1/views"
	err = sendRequest(queryUrl, data, &views)
	return
}
