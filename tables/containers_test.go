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
	assert.Equal(t, &ContainersTable{
		columns: []table.ColumnDefinition{
			table.TextColumn("name"),
			table.TextColumn("pod_uid"),
			table.TextColumn("image"),
			table.TextColumn("privileged"),
		},
		name:   "kubernetes_containers",
		client: tc,
	}, ct)
}
