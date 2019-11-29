package tables

import (
	"context"
	"errors"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewDeploymentsTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	dt := NewDeploymentsTable(tc)

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("namespace"),
		table.TextColumn("selector"),
	}

	assert.Equal(t, &DeploymentsTable{
		name:    "kubernetes_deployments",
		columns: expectedColumns,
		client:  tc,
	}, dt)

	assert.Equal(t, "kubernetes_deployments", dt.Name())
	assert.Equal(t, expectedColumns, dt.Columns())
}

func TestDeploymentsTable_Generate(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	_, _ = tc.ExtensionsV1beta1().Deployments("testing-namespace").Create(&v1beta1.Deployment{
		TypeMeta: v1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "testing-deployment",
			Namespace: "testing-namespace",
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"foo": "bar"},
			},
		},
		Status: v1beta1.DeploymentStatus{},
	})

	ct := NewDeploymentsTable(tc)

	m, err := ct.Generate(context.TODO(), table.QueryContext{})
	assert.NoError(t, err)
	assert.Equal(t, []map[string]string{
		{"name": "testing-deployment", "namespace": "testing-namespace", "selector": "foo=bar,"},
	}, m)

	t.Run("sad path, list deployments returns an error", func(t *testing.T) {
		tc := testclient.NewSimpleClientset()
		tc.CoreV1().(*fake.FakeCoreV1).Fake.PrependReactor("list", "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("unable to list deployments")
		})
		ct := NewDeploymentsTable(tc)
		genTable, err := ct.Generate(context.TODO(), table.QueryContext{})
		assert.Equal(t, errors.New("unable to list deployments"), err)
		assert.Nil(t, genTable)
	})
}
