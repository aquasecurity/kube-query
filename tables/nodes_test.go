package tables

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kolide/osquery-go/plugin/table"
	testclient "k8s.io/client-go/kubernetes/fake"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestNewNodesTable(t *testing.T) {
	kc := testclient.NewSimpleClientset()
	mc := metricsclient.NewSimpleClientset()

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("role"),
		table.TextColumn("external_ip"),
		table.TextColumn("kernel_version"),
		table.TextColumn("kubelet_version"),
		table.TextColumn("cpu_usage"),
		table.TextColumn("memory_usage"),
	}

	nt := NewNodesTable(kc, mc)
	assert.Equal(t, &NodesTable{
		name:          "kubernetes_nodes",
		columns:       expectedColumns,
		client:        kc,
		metricsClient: mc,
	}, nt)

	assert.Equal(t, "kubernetes_nodes", nt.Name())
	assert.Equal(t, expectedColumns, nt.Columns())

}
