package dbhub

// A Go library for working with databases on DBHub.io

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

const (
	version = "0.1.0"
)

// New creates a new DBHub.io connection object.  It doesn't connect to DBHub.io to do this.  Connection only occurs
// when subsequent functions (eg Query()) are called.
func New(key string) (Connection, error) {
	c := Connection{
		APIKey:           key,
		Server:           "https://api.dbhub.io",
		VerifyServerCert: true,
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

// ChangeVerifyServerCert changes whether to verify the server provided https certificate.  Useful for testing and development.
func (c *Connection) ChangeVerifyServerCert(b bool) {
	c.VerifyServerCert = b
}

// Branches returns a list of all available branches of a database along with the name of the default branch
func (c Connection) Branches(dbOwner, dbName string) (branches map[string]BranchEntry, defaultBranch string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the list of branches and the default branch
	var response BranchListResponseContainer
	queryUrl := c.Server + "/v1/branches"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &response)

	// Extract information for return values
	branches = response.Branches
	defaultBranch = response.DefaultBranch
	return
}

// Columns returns the column information for a given table or view
func (c Connection) Columns(dbOwner, dbName string, ident Identifier, table string) (columns []APIJSONColumn, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)
	data.Set("table", table)

	// Fetch the list of columns
	queryUrl := c.Server + "/v1/columns"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &columns)
	return
}

// Commits returns the details of all commits for a database
func (c Connection) Commits(dbOwner, dbName string) (commits map[string]CommitEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the commits
	queryUrl := c.Server + "/v1/commits"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &commits)
	return
}

// Databases returns the list of standard databases in your account
func (c Connection) Databases() (databases []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)

	// Fetch the list of databases
	queryUrl := c.Server + "/v1/databases"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &databases)
	return
}

// DatabasesLive returns the list of Live databases in your account
func (c Connection) DatabasesLive() (databases []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)
	data.Set("live", "true")

	// Fetch the list of databases
	queryUrl := c.Server + "/v1/databases"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &databases)
	return
}

// Delete deletes a database in your account
func (c Connection) Delete(dbName string) (err error) {
	// Prepare the API parameters
	data := c.PrepareVals("", dbName, Identifier{})

	// Delete the database
	queryUrl := c.Server + "/v1/delete"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, nil)
	if err != nil && err.Error() == "no rows in result set" { // Feels like a dodgy workaround
		err = fmt.Errorf("Unknown database\n")
	}
	return
}

// Diff returns the differences between two commits of two databases, or if the details on the second database are left empty,
// between two commits of the same database. You can also specify the merge strategy used for the generated SQL statements.
func (c Connection) Diff(dbOwnerA, dbNameA string, identA Identifier, dbOwnerB, dbNameB string, identB Identifier, merge MergeStrategy) (diffs Diffs, err error) {
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
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &diffs)
	return
}

// Download returns the database file
func (c Connection) Download(dbOwner, dbName string, ident Identifier) (db io.ReadCloser, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the database file
	queryUrl := c.Server + "/v1/download"
	db, err = sendRequest(queryUrl, c.VerifyServerCert, data)
	if err != nil {
		return
	}
	return
}

// Execute executes a SQL statement (INSERT, UPDATE, DELETE) on the chosen database.
func (c Connection) Execute(dbOwner, dbName string, sql string) (rowsChanged int, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})
	data.Set("sql", base64.StdEncoding.EncodeToString([]byte(sql)))

	// Run the query on the remote database
	var execResponse ExecuteResponseContainer
	queryUrl := c.Server + "/v1/execute"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &execResponse)
	if err != nil {
		return
	}
	rowsChanged = execResponse.RowsChanged
	return
}

// Indexes returns the list of indexes present in the database, along with the table they belong to
func (c Connection) Indexes(dbOwner, dbName string, ident Identifier) (idx []APIJSONIndex, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of indexes
	queryUrl := c.Server + "/v1/indexes"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &idx)
	return
}

// Metadata returns the metadata (branches, releases, tags, commits, etc) for the database
func (c Connection) Metadata(dbOwner, dbName string) (meta MetadataResponseContainer, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the list of databases
	queryUrl := c.Server + "/v1/metadata"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &meta)
	return
}

// PrepareVals creates an url.Values container holding the API key, database owner, name, and database identifier.  The
// url.Values container is then used for the requests to DBHub.io.
func (c Connection) PrepareVals(dbOwner, dbName string, ident Identifier) (data url.Values) {
	// Prepare the API parameters
	data = url.Values{}
	if c.APIKey != "" {
		data.Set("apikey", c.APIKey)
	}
	if dbOwner != "" {
		data.Set("dbowner", dbOwner)
	}
	if dbName != "" {
		data.Set("dbname", dbName)
	}
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
	var returnedData []DataRow
	queryUrl := c.Server + "/v1/query"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &returnedData)
	if err != nil {
		return
	}

	// Loop through the results, converting it to a more concise output format
	for _, j := range returnedData {

		// Construct a single row
		var oneRow ResultRow
		for _, l := range j {
			switch l.Type {
			case Float, Integer, Text:
				// Float, integer, and text fields are added to the output
				oneRow.Fields = append(oneRow.Fields, fmt.Sprint(l.Value))
			case Binary:
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

// Releases returns the details of all releases for a database
func (c Connection) Releases(dbOwner, dbName string) (releases map[string]ReleaseEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the releases
	queryUrl := c.Server + "/v1/releases"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &releases)
	return
}

// Tables returns the list of tables in the database
func (c Connection) Tables(dbOwner, dbName string, ident Identifier) (tbl []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of tables
	queryUrl := c.Server + "/v1/tables"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &tbl)
	return
}

// Tags returns the details of all tags for a database
func (c Connection) Tags(dbOwner, dbName string) (tags map[string]TagEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the tags
	queryUrl := c.Server + "/v1/tags"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &tags)
	return
}

// Views returns the list of views in the database
func (c Connection) Views(dbOwner, dbName string, ident Identifier) (views []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of views
	queryUrl := c.Server + "/v1/views"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &views)
	return
}

// Upload uploads a new standard database, or a new revision of a database
func (c Connection) Upload(dbName string, info UploadInformation, dbBytes *[]byte) (err error) {
	// Prepare the API parameters
	data := c.PrepareVals("", dbName, info.Ident)
	data.Del("dbowner") // The upload function always stores the database in the account of the API key user
	if info.CommitMsg != "" {
		data.Set("commitmsg", info.CommitMsg)
	}
	if info.SourceURL != "" {
		data.Set("sourceurl", info.SourceURL)
	}
	if !info.LastModified.IsZero() {
		data.Set("lastmodified", info.LastModified.Format(time.RFC3339))
	}
	if info.Licence != "" {
		data.Set("licence", info.Licence)
	}
	if info.Public != "" {
		data.Set("public", info.Public)
	}
	if info.Force {
		data.Set("force", "true")
	}
	if !info.CommitTimestamp.IsZero() {
		data.Set("committimestamp", info.CommitTimestamp.Format(time.RFC3339))
	}
	if info.AuthorName != "" {
		data.Set("authorname", info.AuthorName)
	}
	if info.AuthorEmail != "" {
		data.Set("authoremail", info.AuthorEmail)
	}
	if info.CommitterName != "" {
		data.Set("committername", info.CommitterName)
	}
	if info.CommitterEmail != "" {
		data.Set("committeremail", info.CommitterEmail)
	}
	if info.OtherParents != "" {
		data.Set("otherparents", info.OtherParents)
	}
	if info.ShaSum != "" {
		data.Set("dbshasum", info.ShaSum)
	}

	// Upload the database
	var body io.ReadCloser
	queryUrl := c.Server + "/v1/upload"
	body, err = sendUpload(queryUrl, c.VerifyServerCert, &data, dbBytes)
	if body != nil {
		defer body.Close()
	}
	if err != nil {
		if body != nil {
			// If there's useful error info in the returned JSON, return that as the error message
			var z JSONError
			err = json.NewDecoder(body).Decode(&z)
			if err != nil {
				return
			}
			err = fmt.Errorf("%s", z.Msg)
		}
	}
	return
}

// UploadLive uploads a new Live database
func (c Connection) UploadLive(dbName string, dbBytes *[]byte) (err error) {
	// Prepare the API parameters
	data := c.PrepareVals("", dbName, Identifier{})
	data.Del("dbowner") // The upload function always stores the database in the account of the API key user
	data.Set("live", "true")

	// Upload the database
	var body io.ReadCloser
	queryUrl := c.Server + "/v1/upload"
	body, err = sendUpload(queryUrl, c.VerifyServerCert, &data, dbBytes)
	if body != nil {
		defer body.Close()
	}
	if err != nil {
		if body != nil {
			// If there's useful error info in the returned JSON, return that as the error message
			var z JSONError
			err = json.NewDecoder(body).Decode(&z)
			if err != nil {
				return
			}
			err = fmt.Errorf("%s", z.Msg)
		}
	}
	return
}

// Webpage returns the URL of the database file in the webUI.  eg. for web browsers
func (c Connection) Webpage(dbOwner, dbName string) (webPage WebpageResponseContainer, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the releases
	queryUrl := c.Server + "/v1/webpage"
	err = sendRequestJSON(queryUrl, c.VerifyServerCert, data, &webPage)
	return
}
