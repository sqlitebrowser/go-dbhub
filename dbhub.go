package dbhub

// An library for working with databases on DBHub.io

// TODO:
//   * Add tests for each function
//   * Create function(s) for listing indexes in the remote database
//   * Create function to list columns in a table or view
//   * Create function for returning a list of available databases
//   * Create function for downloading complete database
//   * Create function for uploading complete database
//   * Create function for retrieving database details (size, branch, commit list, whatever else is useful)
//   * Make a reasonable example application written in Go

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"

	com "github.com/sqlitebrowser/dbhub.io/common"
)

const (
	LibraryVersion = "0.0.1"
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

// ChangeServer changes the address all queries will be sent to.  Useful for testing and development.
func (c *Connection) ChangeServer(s string) {
	c.Server = s
}

// Columns returns the column information for a given table or view
func (c Connection) Columns(dbowner, dbname, table string) (columns []com.APIJSONColumn, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)
	data.Set("table", table)

	// Fetch the list of columns
	queryUrl := c.Server + "/v1/columns"
	err = sendRequest(queryUrl, data, &columns)
	if err != nil {
		log.Printf(err.Error())
	}
	return
}

// Indexes returns the list of indexes present in the database, along with the table they belong to
func (c Connection) Indexes(dbowner, dbname string) (idx map[string]string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of indexes
	queryUrl := c.Server + "/v1/indexes"
	err = sendRequest(queryUrl, data, &idx)
	if err != nil {
		log.Printf(err.Error())
	}
	return
}

// Query runs a SQL query (SELECT only) on the chosen database, returning the results.
// The "blobBase64" boolean specifies whether BLOB data fields should be base64 encoded in the output, or just skipped
// using an empty string as a placeholder.
func (c Connection) Query(dbowner, dbname string, blobBase64 bool, sql string) (out Results, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)
	data.Set("sql", base64.StdEncoding.EncodeToString([]byte(sql)))

	// Run the query on the remote database
	var returnedData []com.DataRow
	queryUrl := c.Server + "/v1/query"
	err = sendRequest(queryUrl, data, &returnedData)
	if err != nil {
		log.Printf(err.Error())
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
					if _, ok := l.Value.(string); ok {
						oneRow.Fields = append(oneRow.Fields, base64.StdEncoding.EncodeToString([]byte(l.Value.(string))))
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
func (c Connection) Tables(dbowner, dbname string) (tbl []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of tables
	queryUrl := c.Server + "/v1/tables"
	err = sendRequest(queryUrl, data, &tbl)
	if err != nil {
		log.Printf(err.Error())
	}
	return
}

// Views returns the list of views in the database
func (c Connection) Views(dbowner, dbname string) (views []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of views
	queryUrl := c.Server + "/v1/views"
	err = sendRequest(queryUrl, data, &views)
	if err != nil {
		log.Printf(err.Error())
	}
	return
}
