package tables

import (
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewContainersTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	ct := NewContainersTable(tc)

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("pod_uid"),
		table.TextColumn("image"),
		table.TextColumn("privileged"),
	}

	assert.Equal(t, &ContainersTable{
		columns: expectedColumns,
		name:    "kubernetes_containers",
		client:  tc,
	}, ct)

	assert.Equal(t, "kubernetes_containers", ct.Name())
	assert.Equal(t, expectedColumns, ct.Columns())
}
