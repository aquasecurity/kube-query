package tables

import (
	"context"
	"log"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/aquasecurity/kube-query/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeploymentsTable implements the Table interface,
// Uses kubeclient to extract information about pods
type DeploymentsTable struct {
	columns []table.ColumnDefinition
	name    string
	client  kubernetes.Interface
}

// NewDeploymentsTable creates a new DeploymentsTable
// saves given initialized kubernetes client
func NewDeploymentsTable(kubeclient kubernetes.Interface) *DeploymentsTable {
	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("namespace"),
		table.TextColumn("selector"),
	}
	return &DeploymentsTable{
		name:    "kubernetes_deployments",
		columns: columns,
		client:  kubeclient,
	}
}

// Name Returns name of table
func (t *DeploymentsTable) Name() string {
	return t.name
}

// Columns Retrieves the initialized columns
func (t *DeploymentsTable) Columns() []table.ColumnDefinition {
	return t.columns
}

// Generate uses the api to retrieve information on all pods
func (t *DeploymentsTable) Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	depls, err := t.client.ExtensionsV1beta1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Could not get deployments from Api")
		return nil, err
	}
	rows := make([]map[string]string, len(depls.Items))
	for i, depl := range depls.Items {
		rows[i] = map[string]string{
			"name":            depl.Name,
			"namespace":       depl.Namespace,
			// TODO: check whether Selector is always using LabelSelector
			"selector":        utils.Map2Str(depl.Spec.Selector.MatchLabels),
		}
	}
	return rows, nil
}
