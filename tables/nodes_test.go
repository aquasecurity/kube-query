package tables

import (
	"context"
	"errors"
	"testing"

	"k8s.io/client-go/kubernetes/typed/core/v1/fake"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
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

func TestNodesTable_Generate(t *testing.T) {
	kc := testclient.NewSimpleClientset()
	mc := metricsclient.NewSimpleClientset()

	_, _ = kc.CoreV1().Nodes().Create(&v1.Node{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-slave",
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			NodeInfo: v1.NodeSystemInfo{
				BootID:         "node-slave",
				KernelVersion:  "1.2.3",
				KubeletVersion: "4.5.6",
			},
			Addresses: []v1.NodeAddress{
				{
					Type:    v1.NodeExternalIP,
					Address: "1.2.3.4",
				},
				{
					Type:    v1.NodeHostName,
					Address: "node-slave-hostname",
				},
			},
		},
	})

	_, _ = kc.CoreV1().Nodes().Create(&v1.Node{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node-master",
			Labels: map[string]string{"node-role.kubernetes.io/master": "masterlabel"},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			NodeInfo: v1.NodeSystemInfo{
				BootID:         "node-master",
				KernelVersion:  "1.2.3",
				KubeletVersion: "4.5.6",
			},
		},
	})

	mc.Fake.PrependReactor("get", "nodes", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1beta1.NodeMetrics{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Timestamp:  metav1.Time{},
			Window:     metav1.Duration{},
			Usage: v1.ResourceList{
				"cpu":    resource.NewQuantity(16, resource.BinarySI).DeepCopy(),
				"memory": resource.NewQuantity(10*1024*1024*1024, resource.BinarySI).DeepCopy(),
			},
		}, nil
	})

	nt := NewNodesTable(kc, mc)
	m, err := nt.Generate(context.TODO(), table.QueryContext{})
	assert.NoError(t, err)

	assert.Equal(t, []map[string]string{
		{
			"cpu_usage": "16", "external_ip": "1.2.3.4", "kernel_version": "1.2.3", "kubelet_version": "4.5.6",
			"memory_usage": "10Gi", "name": "node-slave-hostname", "role": "slave",
		},
		{
			"cpu_usage": "16", "kernel_version": "1.2.3", "kubelet_version": "4.5.6",
			"memory_usage": "10Gi", "name": "node-master", "role": "master",
		},
	}, m)

	t.Run("sad path, node list returns an error", func(t *testing.T) {
		kc := testclient.NewSimpleClientset()
		kc.CoreV1().(*fake.FakeCoreV1).Fake.PrependReactor("list", "nodes", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("unable to get nodes")
		})
		mc := metricsclient.NewSimpleClientset()
		ct := NewNodesTable(kc, mc)

		genTable, err := ct.Generate(context.TODO(), table.QueryContext{})
		assert.Equal(t, errors.New("unable to get nodes"), err)
		assert.Nil(t, genTable)
	})
}
