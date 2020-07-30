package dbhub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	com "github.com/sqlitebrowser/dbhub.io/common"
)

// New creates a new DBHub.io connection object.  It doesn't connect to DBHub.io to do this.
func New(key string) (Connection, error) {
	c := Connection{
		APIKey: key,
		Server: "https://api.dbhub.io",
	}
	return c, nil
}

// Query runs a SQL query (SELECT only) on the chosen database, returning the results
func (c Connection) Query(dbowner, dbname, sql string) (out Results, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)
	data.Set("sql", base64.StdEncoding.EncodeToString([]byte(sql)))

	// Run the query on the remote database
	resp, err := http.PostForm(c.Server+"/v1/query", data)
	if err != nil {
		return Results{}, err
	}
	defer resp.Body.Close()

	// Basic error handling, depending on the status code received from the server
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The returned status code indicates something went wrong
		return Results{}, fmt.Errorf(resp.Status)
	}

	if resp.StatusCode != 200 {
		// TODO: Figure out what should be returned for other 2** status messages
		return
	}

	// The query ran successfully, so prepare and return the results
	var returnedData []com.DataRow
	json.NewDecoder(resp.Body).Decode(&returnedData)

	// Construct the result list
	for _, j := range returnedData {

		// Construct a single row
		var oneRow ResultRow
		for _, l := range j {
			// Float, integer, and text fields are added to the output
			switch l.Type {
			case com.Float, com.Integer, com.Text:
				oneRow.Fields = append(oneRow.Fields, fmt.Sprint(l.Value))
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

// ChangeServer changes the address all Queries will be sent to.  Useful for testing and development.
func (c *Connection) ChangeServer(s string) {
	c.Server = s
}

// TODO: Create function(s) for listing tables and indexes in the remote database

// TODO: Create function for returning a list of available databases

// TODO: Create function for downloading complete database

// TODO: Create function for uploading complete database

// TODO: Create function for retrieving database details (size, branch, commit list, whatever else is useful)

// TODO: Make a reasonable example application written in Go
