package tables

import (
	"context"

	"github.com/kolide/osquery-go/plugin/table"
)

// Table inteface defines the basic Table implementation mechanism for os-query
type Table interface {
	// TODO: change the return value of Genreate() to be map to inteface{} ?
	Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error)
	Columns() []table.ColumnDefinition
	Name() string
}
