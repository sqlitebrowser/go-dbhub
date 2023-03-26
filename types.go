package dbhub

import "time"

// Connection is a simple container holding the API key and address of the DBHub.io server
type Connection struct {
	APIKey           string `json:"api_key"`
	Server           string `json:"server"`
	VerifyServerCert bool   `json:"verify_certificate"`
}

// Identifier holds information used to identify a specific commit, tag, release, or the head of a specific branch
type Identifier struct {
	Branch   string `json:"branch"`
	CommitID string `json:"commit_id"`
	Release  string `json:"release"`
	Tag      string `json:"tag"`
}

// JSONError holds information about an error condition, in a useful JSON format
type JSONError struct {
	Msg string `json:"error"`
}

// MergeStrategy specifies the type of SQL statements included in the diff results.
// The SQL statements can be used for merging databases and depending on whether and
// how you want to merge you should choose your merge strategy.
type MergeStrategy int

const (
	// NoMerge removes any SQL statements for merging from the diff results
	NoMerge MergeStrategy = iota

	// PreservePkMerge produces SQL statements which preserve the values of the primary key columns.
	// Executing these statements on the first database produces a database similar to the second.
	PreservePkMerge

	// NewPkMerge produces SQL statements which generate new values for the primary key columns when
	// executed. This avoids a couple of possible conflicts and allows merging more distant databases.
	NewPkMerge
)

// ResultRow is used for returning the results of a SQL query as a slice of strings
type ResultRow struct {
	Fields []string
}

// Results is used for returning the results of a SQL query as a slice of strings
type Results struct {
	Rows []ResultRow
}

// UploadInformation holds information used when uploading
type UploadInformation struct {
	Ident           Identifier `json:"identifier"`
	CommitMsg       string     `json:"commitmsg"`
	SourceURL       string     `json:"sourceurl"`
	LastModified    time.Time  `json:"lastmodified"`
	Licence         string     `json:"licence"`
	Public          string     `json:"public"`
	Force           bool       `json:"force"`
	CommitTimestamp time.Time  `json:"committimestamp"`
	AuthorName      string     `json:"authorname"`
	AuthorEmail     string     `json:"authoremail"`
	CommitterName   string     `json:"committername"`
	CommitterEmail  string     `json:"committeremail"`
	OtherParents    string     `json:"otherparents"`
	ShaSum          string     `json:"dbshasum"`
}
