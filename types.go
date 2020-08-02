package dbhub

// Connection is a simple container holding the API key and address of the DBHub.io server
type Connection struct {
	APIKey string `json:"api_key"`
	Server string `json:"server"`
}

// ResultRow is used for returning the results of a SQL query as a slice of strings
type ResultRow struct {
	Fields []string
}

// Results is used for returning the results of a SQL query as a slice of strings
type Results struct {
	Rows []ResultRow
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
