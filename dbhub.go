package dbhub

// A Go library for working with databases on DBHub.io

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

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
func (c Connection) Branches(dbOwner, dbName string) (branches map[string]com.BranchEntry, defaultBranch string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the list of branches and the default branch
	var response com.BranchListResponseContainer
	queryUrl := c.Server + "/v1/branches"
	err = sendRequestJSON(queryUrl, data, &response)

	// Extract information for return values
	branches = response.Branches
	defaultBranch = response.DefaultBranch
	return
}

// Columns returns the column information for a given table or view
func (c Connection) Columns(dbOwner, dbName string, ident Identifier, table string) (columns []com.APIJSONColumn, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)
	data.Set("table", table)

	// Fetch the list of columns
	queryUrl := c.Server + "/v1/columns"
	err = sendRequestJSON(queryUrl, data, &columns)
	return
}

// Commits returns the details of all commits for a database
func (c Connection) Commits(dbOwner, dbName string) (commits map[string]com.CommitEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the commits
	queryUrl := c.Server + "/v1/commits"
	err = sendRequestJSON(queryUrl, data, &commits)
	return
}

// Databases returns the list of databases in your account
func (c Connection) Databases() (databases []string, err error) {
	// Prepare the API parameters
	data := url.Values{}
	data.Set("apikey", c.APIKey)

	// Fetch the list of databases
	queryUrl := c.Server + "/v1/databases"
	err = sendRequestJSON(queryUrl, data, &databases)
	return
}

// Delete deletes a database in your account
func (c Connection) Delete(dbName string) (err error) {
	// Prepare the API parameters
	data := c.PrepareVals("", dbName, Identifier{})

	// Delete the database
	queryUrl := c.Server + "/v1/delete"
	err = sendRequestJSON(queryUrl, data, nil)
	if err != nil && err.Error() == "no rows in result set" { // Feels like a dodgy workaround
		err = fmt.Errorf("Unknown database\n")
	}
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
	err = sendRequestJSON(queryUrl, data, &diffs)
	return
}

// Download returns the database file
func (c Connection) Download(dbOwner, dbName string, ident Identifier) (db io.ReadCloser, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the database file
	queryUrl := c.Server + "/v1/download"
	db, err = sendRequest(queryUrl, data)
	if err != nil {
		return
	}
	return
}

// Indexes returns the list of indexes present in the database, along with the table they belong to
func (c Connection) Indexes(dbOwner, dbName string, ident Identifier) (idx []com.APIJSONIndex, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of indexes
	queryUrl := c.Server + "/v1/indexes"
	err = sendRequestJSON(queryUrl, data, &idx)
	return
}

// Metadata returns the metadata (branches, releases, tags, commits, etc) for the database
func (c Connection) Metadata(dbOwner, dbName string) (meta com.MetadataResponseContainer, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the list of databases
	queryUrl := c.Server + "/v1/metadata"
	err = sendRequestJSON(queryUrl, data, &meta)
	return
}

// PrepareVals creates a url.Values container holding the API key, database owner, name, and database identifier.  The
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
	var returnedData []com.DataRow
	queryUrl := c.Server + "/v1/query"
	err = sendRequestJSON(queryUrl, data, &returnedData)
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

// Releases returns the details of all releases for a database
func (c Connection) Releases(dbOwner, dbName string) (releases map[string]com.ReleaseEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the releases
	queryUrl := c.Server + "/v1/releases"
	err = sendRequestJSON(queryUrl, data, &releases)
	return
}

// Tables returns the list of tables in the database
func (c Connection) Tables(dbOwner, dbName string, ident Identifier) (tbl []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of tables
	queryUrl := c.Server + "/v1/tables"
	err = sendRequestJSON(queryUrl, data, &tbl)
	return
}

// Tags returns the details of all tags for a database
func (c Connection) Tags(dbOwner, dbName string) (tags map[string]com.TagEntry, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the tags
	queryUrl := c.Server + "/v1/tags"
	err = sendRequestJSON(queryUrl, data, &tags)
	return
}

// Views returns the list of views in the database
func (c Connection) Views(dbOwner, dbName string, ident Identifier) (views []string, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, ident)

	// Fetch the list of views
	queryUrl := c.Server + "/v1/views"
	err = sendRequestJSON(queryUrl, data, &views)
	return
}

// Upload uploads a new database, or a new revision of a database
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
	body, err = sendUpload(queryUrl, &data, dbBytes)
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
func (c Connection) Webpage(dbOwner, dbName string) (webPage com.WebpageResponseContainer, err error) {
	// Prepare the API parameters
	data := c.PrepareVals(dbOwner, dbName, Identifier{})

	// Fetch the releases
	queryUrl := c.Server + "/v1/webpage"
	err = sendRequestJSON(queryUrl, data, &webPage)
	return
}
