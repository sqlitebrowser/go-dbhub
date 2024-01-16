package dbhub

import (
	"time"
)

type APIJSONColumn struct {
	Cid       int    `json:"column_id"`
	Name      string `json:"name"`
	DataType  string `json:"data_type"`
	NotNull   bool   `json:"not_null"`
	DfltValue string `json:"default_value"`
	Pk        int    `json:"primary_key"`
}

type APIJSONIndexColumn struct {
	CID  int    `json:"id"`
	Name string `json:"name"`
}

type APIJSONIndex struct {
	Name    string               `json:"name"`
	Table   string               `json:"table"`
	Columns []APIJSONIndexColumn `json:"columns"`
}

type BranchEntry struct {
	Commit      string `json:"commit"`
	CommitCount int    `json:"commit_count"`
	Description string `json:"description"`
}

type BranchListResponseContainer struct {
	Branches      map[string]BranchEntry `json:"branches"`
	DefaultBranch string                 `json:"default_branch"`
}

type CommitEntry struct {
	AuthorEmail    string    `json:"author_email"`
	AuthorName     string    `json:"author_name"`
	CommitterEmail string    `json:"committer_email"`
	CommitterName  string    `json:"committer_name"`
	ID             string    `json:"id"`
	Message        string    `json:"message"`
	OtherParents   []string  `json:"other_parents"`
	Parent         string    `json:"parent"`
	Timestamp      time.Time `json:"timestamp"`
	Tree           DBTree    `json:"tree"`
}

type ValType int

const (
	Binary ValType = iota
	Image
	Null
	Text
	Integer
	Float
)

type DataValue struct {
	Name  string
	Type  ValType
	Value interface{}
}

type DataRow []DataValue

type DBTree struct {
	ID      string        `json:"id"`
	Entries []DBTreeEntry `json:"entries"`
}

type DBTreeEntryType string

const (
	TREE     DBTreeEntryType = "tree"
	DATABASE                 = "db"
	LICENCE                  = "licence"
)

type DBTreeEntry struct {
	EntryType    DBTreeEntryType `json:"entry_type"`
	LastModified time.Time       `json:"last_modified"`
	LicenceSHA   string          `json:"licence"`
	Name         string          `json:"name"`
	Sha256       string          `json:"sha256"`
	Size         int64           `json:"size"`
}

type ExecuteResponseContainer struct {
	RowsChanged int    `json:"rows_changed"`
	Status      string `json:"status"`
}

type MetadataResponseContainer struct {
	Branches  map[string]BranchEntry          `json:"branches"`
	Commits   map[string]CommitEntry          `json:"commits"`
	DefBranch string                          `json:"default_branch"`
	Releases  map[string]ReleaseEntry         `json:"releases"`
	Tags      map[string]TagEntry             `json:"tags"`
	WebPage   string                          `json:"web_page"`
}

type ReleaseEntry struct {
	Commit        string    `json:"commit"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	ReleaserEmail string    `json:"email"`
	ReleaserName  string    `json:"name"`
	Size          int64     `json:"size"`
}

type TagEntry struct {
	Commit      string    `json:"commit"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	TaggerEmail string    `json:"email"`
	TaggerName  string    `json:"name"`
}

type WebpageResponseContainer struct {
	WebPage string `json:"web_page"`
}

type DiffType string

const (
	ActionAdd DiffType = "add"
	ActionDelete DiffType = "delete"
	ActionModify DiffType = "modify"
)

type SchemaDiff struct {
	ActionType DiffType `json:"action_type"`
	Sql        string   `json:"sql,omitempty"`
	Before     string   `json:"before"`
	After      string   `json:"after"`
}

type DataDiff struct {
	ActionType DiffType      `json:"action_type"`
	Sql        string        `json:"sql,omitempty"`
	Pk         []DataValue   `json:"pk"`
	DataBefore []interface{} `json:"data_before,omitempty"`
	DataAfter  []interface{} `json:"data_after,omitempty"`
}

type DiffObjectChangeset struct {
	ObjectName string      `json:"object_name"`
	ObjectType string      `json:"object_type"`
	Schema     *SchemaDiff `json:"schema,omitempty"`
	Data       []DataDiff  `json:"data,omitempty"`
}

type Diffs struct {
	Diff []DiffObjectChangeset `json:"diff"`
}
