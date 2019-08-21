package tables

import (
	"context"
	"log"

	"github.com/kolide/osquery-go/plugin/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodsTable implements the Table interface,
// Uses kubeclient to extract information about pods
type PodsTable struct {
	columns []table.ColumnDefinition
	name    string
	client  kubernetes.Interface
}

// NewPodsTable creates a new PodsTable
// saves given initialized kubernetes client
func NewPodsTable(kubeclient kubernetes.Interface) *PodsTable {
	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("namespace"),
		table.TextColumn("ip"),
		table.TextColumn("service_account"),
		table.TextColumn("node_name"),
	}
	return &PodsTable{
		name:    "kubernetes_pods",
		columns: columns,
		client:  kubeclient,
	}
}

// Name Returns name of table
func (t *PodsTable) Name() string {
	return t.name
}

// Columns Retrieves the initialized columns
func (t *PodsTable) Columns() []table.ColumnDefinition {
	return t.columns
}

// Generate uses the api to retrieve information on all pods
func (t *PodsTable) Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	pods, err := t.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Could not get pods from Api")
		return nil, err
	}
	rows := make([]map[string]string, len(pods.Items))
	for i, pod := range pods.Items {
		rows[i] = map[string]string{
			"name":            pod.Name,
			"namespace":       pod.Namespace,
			"ip":              pod.Status.PodIP,
			"service_account": pod.Spec.ServiceAccountName,
			"node_name":       pod.Spec.NodeName,
		}
	}
	return rows, nil
}
