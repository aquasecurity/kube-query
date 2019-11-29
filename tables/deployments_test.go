package tables

import (
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewDeploymentsTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	dt := NewDeploymentsTable(tc)
	assert.Equal(t, &DeploymentsTable{
		name: "kubernetes_deployments",
		columns: []table.ColumnDefinition{
			table.TextColumn("name"),
			table.TextColumn("namespace"),
			table.TextColumn("selector"),
		},
		client: tc,
	}, dt)
}
