package tables

import (
	"context"
	"errors"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/stretchr/testify/assert"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewVolumesTable(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	dt := NewVolumesTable(tc)

	expectedColumns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("type"),
		table.TextColumn("source"),
		table.TextColumn("from_pod"),
	}

	assert.Equal(t, &VolumesTable{
		name:    "kubernetes_volumes",
		columns: expectedColumns,
		client:  tc,
	}, dt)

	assert.Equal(t, "kubernetes_volumes", dt.Name())
	assert.Equal(t, expectedColumns, dt.Columns())
}

func TestVolumesTable_Generate(t *testing.T) {
	tc := testclient.NewSimpleClientset()
	_, _ = tc.CoreV1().Pods("testing-namespace").Create(&v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo-pod-with-two-volumes",
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "volume-1",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/foo1/bar/baz",
							Type: nil,
						},
					},
				},
				//{	// This is not supported as desired: https://github.com/aquasecurity/kube-query/issues/10
				//	Name: "volume-2",
				//	VolumeSource: v1.VolumeSource{
				//		HostPath: &v1.HostPathVolumeSource{
				//			Path: "/foo2/bar/baz",
				//			Type: nil,
				//		},
				//	},
				//},
			},
		},
		Status: v1.PodStatus{},
	})

	dt := NewVolumesTable(tc)
	genTable, err := dt.Generate(context.TODO(), table.QueryContext{})
	assert.NoError(t, err)
	assert.Equal(t, []map[string]string{
		{
			"from_pod": "foo-pod-with-two-volumes", "name": "volume-1", "source": `{"path":"/foo1/bar/baz"}`, "type": "hostPath",
		},
		//{
		//	"from_pod": "foo-pod-with-two-volumes", "name": "volume-2", "source": `{"path":"/foo2/bar/baz"}`, "type": "hostPath",
		//},
	}, genTable)

	t.Run("sad path, list pod returns an error", func(t *testing.T) {
		tc := testclient.NewSimpleClientset()
		tc.CoreV1().(*fake.FakeCoreV1).Fake.PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("unable to list pods")
		})
		ct := NewVolumesTable(tc)
		genTable, err := ct.Generate(context.TODO(), table.QueryContext{})
		assert.NoError(t, err)
		assert.Empty(t, genTable)
	})
}
