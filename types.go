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
