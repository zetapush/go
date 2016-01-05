package zpservice

import (
	"encoding/json"
)

const (
	NO_TABLE       = "Happens when the specified table does not exist (tables are created by a createTable call)."
	NO_SUCH_COLUMN = "Happens when the specified column does not exist (columns are defined by a createTable or addColumns call)."
	BAD_COLUMN     = "Happens when the input data is incorrect for the specified column"
)

/*get*/
type GdaGet struct {
	Table string `json:"table"`
	Key   string `json:"key"`
	Owner string `json:"owner,omitempty"`
}

type GdaGetResult struct {
	Request GdaGet           `json:"request"`
	Result  *json.RawMessage `json:"result"`
}

/*getCells*/
type GdaCellsGet struct {
	Table  string   `json:"table"`
	Column string   `json:"column"`
	Key    string   `json:"key"`
	Key2   []string `json:"key2,omitempty"`
	Owner  string   `json:"owner,omitempty"`
}

type GdaCellsResult struct {
	Request GdaCellsGet      `json:"request"`
	Result  *json.RawMessage `json:"result"`
}

/*put*/
type GdaPut struct {
	Table  string           `json:"table"`
	Column string           `json:"column"`
	Key    string           `json:"key"`
	Key2   []string         `json:"key2,omitempty"`
	Owner  string           `json:"owner,omitempty"`
	Data   *json.RawMessage `json:"data"`
}

/*puts*/
type GdaPuts struct {
	Table string       `json:"table"`
	Rows  []GdaPutsRow `json:"rows,omitempty"`
	Owner string       `json:"owner,omitempty"`
}

type GdaPutsRow struct {
	Data *json.RawMessage `json:"data"`
	Key  string           `json:"key"`
}

type GdaPutsResult struct {
	Inserted int64  `json:"inserted"`
	Table    string `json:"table"`
	Owner    string `json:"owner"`
}

/*range*/
type GdaRange struct {
	Table   string          `json:"table"`
	Start   string          `json:"start"`
	Stop    string          `json:"stop"`
	Columns []GdaColumnSpec `json:"columns"`
	Page    Pagination      `json:"page"`
	Owner   string          `json:"owner,omitempty"`
}

type GdaColumnSpec struct {
	Column string   `json:"column"`
	Key2   []string `json:"key2,omitempty"`
}

type GdaRangeResult struct {
	Request GdaRange         `json:"request"`
	Result  *json.RawMessage `json:"result"`
}

/*removeCell*/
type GdaCellRequest struct {
	Table  string   `json:"table"`
	Column string   `json:"column"`
	Key    string   `json:"key"`
	Key2   []string `json:"key2,omitempty"`
	Owner  string   `json:"owner,omitempty"`
}

/*removeColumn*/
type GdaColumnRequest struct {
	Table  string `json:"table"`
	Column string `json:"column"`
	Key    string `json:"key"`
	Owner  string `json:"owner,omitempty"`
}

/*removeRow*/
type GdaRowRequest struct {
	Table string `json:"table"`
	Key   string `json:"key"`
	Owner string `json:"owner,omitempty"`
}
