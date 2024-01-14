package dbhub

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/fileutils"
	sqlite "github.com/gwenn/gosqlite"
	"github.com/sqlitebrowser/dbhub.io/common"
	"github.com/stretchr/testify/assert"
)

// For now, the tests require the DBHub.io dev docker container be running on its standard
// ports (that means the API server is listening on https://localhost:9444)

func TestMain(m *testing.M) {
	log.Println("Seeding the database...")

	// Disable https cert validation for our tests
	insecureTLS := tls.Config{InsecureSkipVerify: true}
	insecureTransport := http.Transport{TLSClientConfig: &insecureTLS}
	client := http.Client{Transport: &insecureTransport}

	// Seed the database
	resp, err := client.Get("https://localhost:9443/x/test/seed")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Database seed request returned http code '%d'.  Aborting tests.", resp.StatusCode)
	}
	log.Println("Database seeding completed ok.")

	// Run the tests
	log.Println("Running the tests...")
	m.Run()
}

// TestBranches verifies retrieving the branch and default branch information using the API
func TestBranches(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve branch details for the database
	branches, defaultBranch, err := conn.Branches("default", "Assembly Election 2017 with view.sqlite")
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the returned branch information matches what we're expecting
	assert.Len(t, branches, 1)
	assert.Contains(t, branches, "main")
	assert.Equal(t, "main", defaultBranch)
}

// TestColumns verifies retrieving the list of column names for a database using the API
func TestColumns(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the column info for a table or view in the remote database
	table := "Candidate_Names"
	columns, err := conn.Columns("default", "Assembly Election 2017 with view.sqlite", Identifier{Branch: "main"}, table)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the returned column information matches what we're expecting
	assert.Len(t, columns, 2)
	assert.Contains(t, columns[0].Name, "Firstname")
	assert.Equal(t, "Surname", columns[1].Name)
}

// TestCommits verifies retrieving commit information for standard databases using the API
func TestCommits(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the commit info for a remote database
	commits, err := conn.Commits("default", "Assembly Election 2017 with view.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Verify the returned commit information matches what we're expecting
	assert.Len(t, commits, 1)

	// Abort early if the returned length of commits isn't what we're expecting
	if len(commits) != 1 {
		return
	}

	// Retrieve the the first commit id
	var ids []string
	for id := range commits {
		ids = append(ids, id)
	}
	firstID := ids[0]

	// Verify the commit information is what we're expecting
	assert.Equal(t, "default@docker-dev.dbhub.io", commits[firstID].AuthorEmail)
	assert.Equal(t, "Default system user", commits[firstID].AuthorName)
	assert.Equal(t, "Initial commit", commits[firstID].Message)
	assert.Equal(t, "Assembly Election 2017 with view.sqlite", commits[firstID].Tree.Entries[0].Name)
	assert.Equal(t, int64(73728), commits[firstID].Tree.Entries[0].Size)
	assert.Equal(t, "9cb18719bddb949043abc1ba089dd7c4845ab024ddbe4ad19e9334da4e5b8cdc", commits[firstID].Tree.Entries[0].Sha256)
	assert.Equal(t, "9348ddfd44da5a127c59141981954746a860ec8e03e0412cf3af7134af0f97e2", commits[firstID].Tree.Entries[0].LicenceSHA)
}

// TestDatabases verifies retrieving the list of standard databases using the API
func TestDatabases(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the list of databases for the user
	databases, err := conn.Databases()
	if err != nil {
		t.Errorf("Connecting to the API server failed: %v", err)
		return
	}

	// If no databases were found, the test failed
	if len(databases) == 0 {
		t.Error("No databases found")
		return
	}

	// Verify the expected database names were returned, and only them
	assert.Contains(t, databases, "Assembly Election 2017 with view.sqlite")
	assert.Contains(t, databases, "Assembly Election 2017.sqlite")
	assert.Len(t, databases, 2)
	return
}

// TestDatabases verifies retrieving the list of live databases using the API
func TestDatabasesLive(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the list of live databases for the user
	databases, err := conn.DatabasesLive()
	if err != nil {
		t.Error(err)
		return
	}

	// If no databases were found, the test failed
	if len(databases) == 0 {
		t.Error("No databases found")
		return
	}

	// Verify the expected database name was returned
	assert.Contains(t, databases, "Join Testing with index.sqlite")
	assert.Len(t, databases, 1)
	return
}

// TestDiff verifies the Diff API call
func TestDiff(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Read the example database file into memory
	dbFile := filepath.Join("examples", "upload", "example.db")
	z, err := os.ReadFile(dbFile)
	if err != nil {
		t.Error(err)
		return
	}

	// Upload the example database
	dbName := "uploadtest1.sqlite"
	err = conn.Upload(dbName, UploadInformation{}, &z)
	if err != nil {
		t.Error(err)
		return
	}
	t.Cleanup(func() {
		// Delete the uploaded database when the test exits
		err = conn.Delete(dbName)
		if err != nil {
			t.Error(err)
			return
		}
	})

	// Copy the database file to a temp location so we can make some changes
	newFile := filepath.Join(t.TempDir(), "diff-"+common.RandomString(8)+".sqlite")
	_, err = fileutils.CopyFile(dbFile, newFile)
	if err != nil {
		t.Error(err)
		return
	}

	// Make some changes to the copied database file
	sdb, err := sqlite.Open(newFile, sqlite.OpenReadWrite|sqlite.OpenFullMutex)
	if err != nil {
		t.Error()
		return
	}
	dbQuery := `
		CREATE TABLE foo (first integer);
		INSERT INTO foo (first) values (10);
		INSERT INTO foo (first) values (20);`
	err = sdb.Exec(dbQuery)
	if err != nil {
		t.Error()
		sdb.Close()
		return
	}
	sdb.Close()

	// Retrieve the initial commit id for the database
	dbOwner := "default"
	commitMap, err := conn.Commits(dbOwner, dbName)
	if err != nil {
		t.Error(err)
		return
	}
	var commits []string
	for idx := range commitMap {
		commits = append(commits, idx)
	}
	firstCommit := commits[0]

	// Upload the copied file as a new commit
	z, err = os.ReadFile(newFile)
	if err != nil {
		t.Error(err)
		return
	}
	uploadCommit := Identifier{CommitID: firstCommit}
	err = conn.Upload(dbName, UploadInformation{Ident: uploadCommit}, &z)
	if err != nil {
		t.Error(err)
		return
	}

	// Retrieve the new commit id for the database
	commitMap, err = conn.Commits(dbOwner, dbName)
	if err != nil {
		t.Error(err)
		return
	}
	var secondCommit string
	for idx := range commitMap {
		if idx != firstCommit {
			secondCommit = idx
		}
	}

	// Do the diff using NewPkMerge
	commit1 := Identifier{CommitID: firstCommit}
	commit2 := Identifier{CommitID: secondCommit}
	diffs, err := conn.Diff(dbOwner, dbName, commit1, "", "", commit2, NewPkMerge)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the changes
	assert.Len(t, diffs.Diff, 1)
	assert.Equal(t, "foo", diffs.Diff[0].ObjectName)
	assert.Equal(t, "table", diffs.Diff[0].ObjectType)
	assert.Equal(t, common.DiffType("add"), diffs.Diff[0].Schema.ActionType)
	assert.Equal(t, "CREATE TABLE foo (first integer);", diffs.Diff[0].Schema.Sql)
	assert.Equal(t, "", diffs.Diff[0].Schema.Before)
	assert.Equal(t, "CREATE TABLE foo (first integer)", diffs.Diff[0].Schema.After)
	assert.Equal(t, `INSERT INTO "foo"("first") VALUES(10);`, diffs.Diff[0].Data[0].Sql)
	assert.Equal(t, `INSERT INTO "foo"("first") VALUES(20);`, diffs.Diff[0].Data[1].Sql)

	// Diff with PreservePkMerge
	diffs, err = conn.Diff(dbOwner, dbName, commit1, "", "", commit2, PreservePkMerge)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the changes
	assert.Len(t, diffs.Diff, 1)
	assert.Equal(t, "foo", diffs.Diff[0].ObjectName)
	assert.Equal(t, "table", diffs.Diff[0].ObjectType)
	assert.Equal(t, common.DiffType("add"), diffs.Diff[0].Schema.ActionType)
	assert.Equal(t, "CREATE TABLE foo (first integer);", diffs.Diff[0].Schema.Sql)
	assert.Equal(t, "", diffs.Diff[0].Schema.Before)
	assert.Equal(t, "CREATE TABLE foo (first integer)", diffs.Diff[0].Schema.After)
	assert.Equal(t, `INSERT INTO "foo"("first") VALUES(10);`, diffs.Diff[0].Data[0].Sql)
	assert.Equal(t, `INSERT INTO "foo"("first") VALUES(20);`, diffs.Diff[0].Data[1].Sql)

	// Diff with NoMerge
	diffs, err = conn.Diff(dbOwner, dbName, commit1, "", "", commit2, NoMerge)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the changes
	assert.Len(t, diffs.Diff, 1)
	assert.Equal(t, "foo", diffs.Diff[0].ObjectName)
	assert.Equal(t, "table", diffs.Diff[0].ObjectType)
	assert.Equal(t, common.DiffType("add"), diffs.Diff[0].Schema.ActionType)
	assert.Equal(t, "", diffs.Diff[0].Schema.Sql)
	assert.Equal(t, "", diffs.Diff[0].Schema.Before)
	assert.Equal(t, "CREATE TABLE foo (first integer)", diffs.Diff[0].Schema.After)
	assert.Equal(t, "", diffs.Diff[0].Data[0].Sql)
	assert.Equal(t, "", diffs.Diff[0].Data[1].Sql)
}

// TestExecute verifies the Execute API call
func TestExecute(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Execute a SQL statement
	dbQuery := `INSERT INTO table1 (id, Name) VALUES (7, "Stuff")`
	rowsChanged, err := conn.Execute("default", "Join Testing with index.sqlite", dbQuery)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the result
	assert.Equal(t, 1, rowsChanged)

	// Execute another SQL statement
	dbQuery = `UPDATE table1 SET Name = "New Stuff" WHERE id = 1 OR id = 7`
	rowsChanged, err = conn.Execute("default", "Join Testing with index.sqlite", dbQuery)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the result
	assert.Equal(t, 2, rowsChanged)
}

// TestIndexes verifies the Indexes API call
func TestIndexes(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the index information for the database
	indexes, err := conn.Indexes("default", "Join Testing with index.sqlite", Identifier{})
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the index information
	assert.Len(t, indexes, 1)
	assert.Equal(t, "table1", indexes[0].Table)
	assert.Equal(t, "stuff", indexes[0].Name)
	assert.Empty(t, indexes[0].Columns[0].CID)
	assert.Equal(t, "id", indexes[0].Columns[0].Name)
}

// TestMetadata verifies the metadata API call
func TestMetadata(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the metadata information for the database
	meta, err := conn.Metadata("default", "Assembly Election 2017 with view.sqlite")
	if err != nil {
		t.Error(err)
		return
	}

	// Get the commit id of the first commit
	var firstCommit string
	for _, c := range meta.Commits {
		firstCommit = c.ID
	}

	// Verify the metadata info
	assert.Equal(t, "https://docker-dev.dbhub.io:9443/default/Assembly Election 2017 with view.sqlite", meta.WebPage)
	assert.Equal(t, "main", meta.DefBranch)
	assert.Equal(t, "", meta.Branches["main"].Description)
	assert.Equal(t, 1, meta.Branches["main"].CommitCount)
	assert.Empty(t, meta.Releases)
	assert.Empty(t, meta.Tags)
	assert.Equal(t, firstCommit, meta.Commits[firstCommit].ID)
	assert.Equal(t, "Initial commit", meta.Commits[firstCommit].Message)
	assert.Equal(t, "Default system user", meta.Commits[firstCommit].AuthorName)
	assert.Equal(t, "default@docker-dev.dbhub.io", meta.Commits[firstCommit].AuthorEmail)
	assert.Equal(t, int64(73728), meta.Commits[firstCommit].Tree.Entries[0].Size)
}

// TestQuery verifies the Query API call
func TestQuery(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Query the database
	dbQuery := `
		SELECT id, Name
		FROM table1
		ORDER BY Name DESC`
	result, err := conn.Query("default", "Join Testing with index.sqlite", Identifier{}, false, dbQuery)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the result
	assert.Len(t, result.Rows, 7)
	assert.Contains(t, result.Rows, ResultRow{Fields: []string{"2", "Bar"}})
	assert.Contains(t, result.Rows, ResultRow{Fields: []string{"3", "Baz"}})
	assert.Contains(t, result.Rows, ResultRow{Fields: []string{"4", "Blumph"}})
	assert.Contains(t, result.Rows, ResultRow{Fields: []string{"5", "Blargo"}})
	assert.Contains(t, result.Rows, ResultRow{Fields: []string{"6", "Batty"}})
}

// TestReleases verifies the Releases API call
func TestReleases(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the releases for a database
	releases, err := conn.Releases("default", "Assembly Election 2017.sqlite")
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the retrieved information
	assert.Len(t, releases, 2)
	assert.Equal(t, "First release", releases["first"].Description)
	assert.Equal(t, "Example Releaser", releases["first"].ReleaserName)
	assert.Equal(t, "example@example.org", releases["first"].ReleaserEmail)
	assert.Equal(t, "Second release", releases["second"].Description)
	assert.Equal(t, "Example Releaser", releases["second"].ReleaserName)
	assert.Equal(t, "example@example.org", releases["second"].ReleaserEmail)
}

// TestTables verifies the Tables API call
func TestTables(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve table information for a database
	tbls, err := conn.Tables("default", "Assembly Election 2017.sqlite", Identifier{})
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the returned information
	assert.Len(t, tbls, 3)
	assert.Contains(t, tbls, "Candidate_Information")
	assert.Contains(t, tbls, "Constituency_Turnout_Information")
	assert.Contains(t, tbls, "Elected_Candidates")
}

// TestTags verifies the Tags API call
func TestTags(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Retrieve the tags for a database
	tags, err := conn.Tags("default", "Assembly Election 2017.sqlite")
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the retrieved information
	assert.Len(t, tags, 2)
	assert.Equal(t, "First tag", tags["first"].Description)
	assert.Equal(t, "Example Tagger", tags["first"].TaggerName)
	assert.Equal(t, "example@example.org", tags["first"].TaggerEmail)
	assert.Equal(t, "Second tag", tags["second"].Description)
	assert.Equal(t, "Example Tagger", tags["second"].TaggerName)
	assert.Equal(t, "example@example.org", tags["second"].TaggerEmail)
}

// TestUpload verifies uploading a standard database via the API
func TestUpload(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Read the example database file into memory
	dbFile := filepath.Join("examples", "upload", "example.db")
	z, err := os.ReadFile(dbFile)
	if err != nil {
		t.Error(err)
		return
	}

	// Upload the example database
	dbName := "testupload.sqlite"
	err = conn.Upload(dbName, UploadInformation{}, &z)
	if err != nil {
		t.Error(err)
		return
	}

	// Verify the file contents.  This is done by downloading the database and doing a byte comparison to ensure its
	// identical to the upload
	downloaded, err := conn.Download("default", dbName, Identifier{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Cleanup(func() {
		// Delete the uploaded file when the function exits
		err = conn.Delete(dbName)
		if err != nil {
			t.Error(err)
			return
		}
	})
	data, err := io.ReadAll(downloaded)
	if err != nil {
		t.Error(err)
		return
	}
	result := bytes.Compare(z, data)
	if result != 0 {
		t.Errorf("Standard database upload succeeded, but failed verification when downloading it and comparing to the original")
		return
	}
}

// TestUploadLive verifies uploading a live database via the API
func TestUploadLive(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Read the example database file into memory
	dbA := filepath.Join("examples", "upload", "example.db")
	z, err := os.ReadFile(dbA)
	if err != nil {
		t.Error(err)
		return
	}

	// Upload the database
	dbName := "testuploadlive.sqlite"
	err = conn.UploadLive(dbName, &z)
	if err != nil {
		t.Error(err)
		return
	}

	// *** Verify the database ***

	// This is done by downloading and comparing (database diff) the database contents with the database file that was uploaded
	downloaded, err := conn.Download("default", dbName, Identifier{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Cleanup(func() {
		// Delete the uploaded database when the test exits
		err = conn.Delete(dbName)
		if err != nil {
			t.Error(err)
			return
		}
	})
	data, err := io.ReadAll(downloaded)
	if err != nil {
		t.Error(err)
		return
	}
	dbB := filepath.Join(t.TempDir(), "diff-"+common.RandomString(8)+".sqlite")
	err = os.WriteFile(dbB, data, 0750)
	if err != nil {
		t.Error(err)
		return
	}

	// Do the comparison
	diffs, err := common.DBDiff(dbA, dbB, common.NoMerge, false)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, diffs.Diff)
}

// TestViews verifies the Views API call
func TestViews(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Get the list of views in the database
	views, err := conn.Views("default", "Assembly Election 2017 with view.sqlite", Identifier{})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, views, 1)
	assert.Equal(t, "Candidate_Names", views[0])
}

// TestWebpage verifies the Webpage API call
func TestWebpage(t *testing.T) {
	// Create the local test server connection
	conn := serverConnection("Rh3fPl6cl84XEw2FeWtj-FlUsn9OrxKz9oSJfe6kho7jT_1l5hizqw")

	// Gather the data then test the result
	pageData, err := conn.Webpage("default", "Assembly Election 2017.sqlite")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "https://docker-dev.dbhub.io:9443/default/Assembly Election 2017.sqlite", pageData.WebPage)
}

// serverConnection is a utility function that sets up the API connection object to the test server, ready for use
func serverConnection(apiKey string) Connection {
	// Create a new DBHub.io API object
	db, err := New(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	db.ChangeServer("https://localhost:9444")
	db.ChangeVerifyServerCert(false)
	return db
}
