package tables

import (
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewPodsTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("uid"),
		table.TextColumn("name"),
		table.TextColumn("namespace"),
		table.TextColumn("ip"),
		table.TextColumn("service_account"),
		table.TextColumn("node_name"),
		table.TextColumn("phase"),
	}

	pt := NewPodsTable(tc)
	assert.Equal(t, &PodsTable{
		name:    "kubernetes_pods",
		client:  tc,
		columns: expectedColumns,
	}, pt)

	assert.Equal(t, "kubernetes_pods", pt.Name())
	assert.Equal(t, expectedColumns, pt.Columns())
}
