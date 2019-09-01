package tables

import (
	"context"
	"log"

	// "github.com/aquasecurity/kube-query/utils"
	"github.com/kolide/osquery-go/plugin/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

// NodesTable implements the Table interface,
// Uses kubeclient to extract information about pods
type NodesTable struct {
	columns        []table.ColumnDefinition
	name           string
	client  	   kubernetes.Interface
	metricsClient  *metrics.Clientset
}

// NewNodesTable creates a new NodesTable
// saves given initialized kubernetes client
func NewNodesTable(kubeclient kubernetes.Interface, mc *metrics.Clientset) *NodesTable {
	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("role"),
		table.TextColumn("external_ip"),
		table.TextColumn("kernel_version"),
		table.TextColumn("kubelet_version"),
		table.TextColumn("cpu_usage"),
		table.TextColumn("memory_usage"),
	}
	return &NodesTable{
		name:    "kubernetes_nodes",
		columns: columns,
		client:  kubeclient,
		metricsClient: mc,
	}
}

// Name Returns name of table
func (t *NodesTable) Name() string {
	return t.name
}

// Columns Retrieves the initialized columns
func (t *NodesTable) Columns() []table.ColumnDefinition {
	return t.columns
}

// Generate uses the api to retrieve information on all pods
func (t *NodesTable) Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	nodes, err := t.client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Could not get nodes from Api")
		return nil, err
	}

	rows := make([]map[string]string, len(nodes.Items))
	for i, node := range nodes.Items {
		currRow := map[string]string{
			"name":	node.Status.NodeInfo.BootID,
			"kernel_version": node.Status.NodeInfo.KernelVersion,
			"kubelet_version": node.Status.NodeInfo.KubeletVersion,			
			"role": "slave", // default to slave, unless decided master 
		}

		// checking if it's a master node
		if _, hasMasterRoleLabel := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; hasMasterRoleLabel {
			currRow["role"] = "master"
		}

		// setting addresses
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeHostName {
				currRow["name"] = address.Address 
			} else if address.Type == corev1.NodeExternalIP {
				currRow["external_ip"] = address.Address 
			}
		}

		// if nodename exists, we extract metrics using the name
		if nodename := currRow["name"]; nodename != "" {
			metrics, err := t.metricsClient.MetricsV1beta1().NodeMetricses().Get(nodename, metav1.GetOptions{})
			if err == nil {
				// metrics.Usage is of type *resource.Quantity. from:
				// https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go
				if cpu, ok := metrics.Usage[corev1.ResourceCPU]; ok {
					currRow["cpu_usage"] = cpu.String()
				}
				if memory, ok := metrics.Usage[corev1.ResourceMemory]; ok {
					currRow["memory_usage"] = memory.String()
				}
			}
		}

		rows[i] = currRow
	}
	return rows, nil
}
