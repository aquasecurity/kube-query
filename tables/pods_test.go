package tables

import (
	"context"
	"errors"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"

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

func TestPodsTable_Generate(t *testing.T) {
	tc := testclient.NewSimpleClientset()

	_, _ = tc.CoreV1().Pods("testing-namespace").Create(&v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foopod",
			Namespace: "testing-namespace",
			UID:       "foo-123-bar-456-baz",
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "testing-service-account-name",
			NodeName:           "testing-node-name",
		},
		Status: v1.PodStatus{
			PodIP: "1.2.3.4",
			Phase: "running",
		},
	})

	_, _ = tc.CoreV1().Pods("testing-namespace").Create(&v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bazpod",
			Namespace: "testing-namespace",
			UID:       "baz-789-foo-123-bar",
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "testing-service-account-name",
			NodeName:           "testing-node-name",
		},
		Status: v1.PodStatus{
			PodIP: "4.5.6.7",
			Phase: "finished",
		},
	})

	pt := NewPodsTable(tc)
	m, err := pt.Generate(context.TODO(), table.QueryContext{})
	assert.NoError(t, err)

	assert.Equal(t, []map[string]string{
		{"ip": "1.2.3.4", "name": "foopod", "namespace": "testing-namespace", "node_name": "testing-node-name",
			"phase": "running", "service_account": "testing-service-account-name", "uid": "foo-123-bar-456-baz"},
		{"ip": "4.5.6.7", "name": "bazpod", "namespace": "testing-namespace", "node_name": "testing-node-name",
			"phase": "finished", "service_account": "testing-service-account-name", "uid": "baz-789-foo-123-bar"},
	}, m)

	t.Run("sad path, list pod returns an error", func(t *testing.T) {
		tc := testclient.NewSimpleClientset()
		tc.CoreV1().(*fake.FakeCoreV1).Fake.PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("unable to list pods")
		})
		ct := NewPodsTable(tc)
		genTable, err := ct.Generate(context.TODO(), table.QueryContext{})
		assert.Equal(t, errors.New("unable to list pods"), err)
		assert.Nil(t, genTable)
	})
}
