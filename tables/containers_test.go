package tables

import (
	"context"
	"errors"
	"testing"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewContainersTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	ct := NewContainersTable(tc)

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("pod_uid"),
		table.TextColumn("image"),
		table.TextColumn("privileged"),
	}

	assert.Equal(t, &ContainersTable{
		columns: expectedColumns,
		name:    "kubernetes_containers",
		client:  tc,
	}, ct)

	assert.Equal(t, "kubernetes_containers", ct.Name())
	assert.Equal(t, expectedColumns, ct.Columns())
}

func isPrivileged(b bool) *bool {
	return &b
}

func TestContainersTable_Generate(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	_, _ = tc.CoreV1().Pods("testing-namespace").Create(&v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "testing-namespace",
			UID:       "test-uid",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "test-container-1",
					Image: "test-container-image",
					SecurityContext: &v1.SecurityContext{
						Privileged: isPrivileged(true),
					},
				}, {
					Name:  "test-container-2",
					Image: "test-container-image",
				},
			},
		},
		Status: v1.PodStatus{},
	})

	ct := NewContainersTable(tc)

	m, err := ct.Generate(context.TODO(), table.QueryContext{})
	assert.NoError(t, err)
	assert.Equal(t, []map[string]string{
		{"image": "test-container-image", "name": "test-container-1", "pod_uid": "test-uid", "privileged": "True"},
		{"image": "test-container-image", "name": "test-container-2", "pod_uid": "test-uid", "privileged": "False"},
	}, m)

	t.Run("sad path, list pods returns an error", func(t *testing.T) {
		tc := testclient.NewSimpleClientset()
		tc.CoreV1().(*fake.FakeCoreV1).Fake.PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("unable to list pods")
		})
		ct := NewContainersTable(tc)
		genTable, err := ct.Generate(context.TODO(), table.QueryContext{})
		assert.Equal(t, errors.New("unable to list pods"), err)
		assert.Nil(t, genTable)
	})
}
