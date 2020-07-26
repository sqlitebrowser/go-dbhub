package dbhub

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
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
func (c Connection) Query(dbowner, dbname, sql string) (Results, error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("dbowner", dbowner)
	data.Set("dbname", dbname)
	data.Set("sql", base64.StdEncoding.EncodeToString([]byte(sql)))

	// Run the query on the remote database
	res, err := http.PostForm(c.Server+"/v1/query", data)
	if err != nil {
		return Results{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// The returned status code indicates something went wrong
		return Results{}, fmt.Errorf(res.Status)
	}
	if res.StatusCode == 200 {
		// The query ran successfully, so prepare and return the results
		// TODO: TBD
		fmt.Printf("Results: %v\n", res)
	}

	// TODO: Figure out what should be returned here
	return Results{}, nil
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
