package dbhub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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

// ChangeServer changes the address all Queries will be sent to.  Useful for testing and development.
func (c *Connection) ChangeServer(s string) {
	c.Server = s
}

// Indexes returns the list of indexes present in the database
func (c Connection) Indexes(dbowner, dbname string) (idx map[string]string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of indexes
	var resp *http.Response
	resp, err = http.PostForm(c.Server+"/v1/indexes", data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Basic error handling, depending on the status code received from the server
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		return
	}

	if resp.StatusCode != 200 {
		// TODO: Figure out what should be returned for other 2** status messages
		return
	}

	// Convert the response into the list of indexes
	err = json.NewDecoder(resp.Body).Decode(&idx)
	if err != nil {
		log.Fatal(err)
	}
	return
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
	var resp *http.Response
	resp, err = http.PostForm(c.Server+"/v1/query", data)
	if err != nil {
		return
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
	err = json.NewDecoder(resp.Body).Decode(&returnedData)
	if err != nil {
		log.Fatal(err)
	}

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

// Tables returns the list of tables present in the database
func (c Connection) Tables(dbowner, dbname string) (tbl []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of tables
	var resp *http.Response
	resp, err = http.PostForm(c.Server+"/v1/tables", data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Basic error handling, depending on the status code received from the server
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		return
	}

	if resp.StatusCode != 200 {
		// TODO: Figure out what should be returned for other 2** status messages
		return
	}

	// Convert the response into the list of tables
	err = json.NewDecoder(resp.Body).Decode(&tbl)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Views returns the list of views present in the database
func (c Connection) Views(dbowner, dbname string) (vws []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)

	// Fetch the list of views
	var resp *http.Response
	resp, err = http.PostForm(c.Server+"/v1/views", data)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Basic error handling, depending on the status code received from the server
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		return
	}

	if resp.StatusCode != 200 {
		// TODO: Figure out what should be returned for other 2** status messages
		return
	}

	// Convert the response into the list of views
	err = json.NewDecoder(resp.Body).Decode(&vws)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// TODO: Create function(s) for listing indexes in the remote database

// TODO: Create function to list columns in a table (or view?)

// TODO: Create function for returning a list of available databases

// TODO: Create function for downloading complete database

// TODO: Create function for uploading complete database

// TODO: Create function for retrieving database details (size, branch, commit list, whatever else is useful)

// TODO: Make a reasonable example application written in Go
