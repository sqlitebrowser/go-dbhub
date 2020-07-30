package dbhub

type Connection struct {
	APIKey string `json:"api_key"`
	Server string `json:"server"`
}

type ResultRow struct {
	Fields []string
}

type Results struct {
	Rows []ResultRow
}
