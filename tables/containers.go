package tables

import (
	"context"
	"log"

	"github.com/aquasecurity/kube-query/utils"
	"github.com/kolide/osquery-go/plugin/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ContainersTable implements the Table interface,
// Uses kubeclient to extract information about pods
type ContainersTable struct {
	columns []table.ColumnDefinition
	name    string
	client  kubernetes.Interface
}

// NewContainersTable creates a new ContainersTable
// saves given initialized kubernetes client
func NewContainersTable(kubeclient kubernetes.Interface) *ContainersTable {
	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("pod_uid"),
		table.TextColumn("image"),
		table.TextColumn("privileged"), // TODO: should we make this an IntegerColumn?
	}
	return &ContainersTable{
		name:    "kubernetes_containers",
		columns: columns,
		client:  kubeclient,
	}
}

// Name Returns name of table
func (t *ContainersTable) Name() string {
	return t.name
}

// Columns Retrieves the initialized columns
func (t *ContainersTable) Columns() []table.ColumnDefinition {
	return t.columns
}

// Generate uses the api to retrieve information on all pods
func (t *ContainersTable) Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	pods, err := t.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Could not get pods from Api")
		return nil, err
	}

	// TODO: think of an efficient way to create the slice without reallocating
	rows := make([]map[string]string, 0, len(pods.Items))

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			currRow := map[string]string{
				"name":       container.Name,
				"pod_uid":    string(pod.UID),
				"image":      container.Image,
				"privileged": "False",
			}
			if container.SecurityContext != nil {
				if container.SecurityContext.Privileged != nil {
					currRow["privileged"] = utils.Bool2str(*container.SecurityContext.Privileged)
				}
			}
			rows = append(rows, currRow)
		}
	}
	return rows, nil
}
