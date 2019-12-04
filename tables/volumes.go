package tables

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kolide/osquery-go/plugin/table"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// VolumesTable implements the Table interface,
// Uses kubeclient to extract information about pods
type VolumesTable struct {
	columns []table.ColumnDefinition
	name    string
	client  kubernetes.Interface
}

type VolumeRow struct {
	volume  *corev1.Volume
	fromPod string
}

// NewVolumesTable creates a new VolumesTable
// saves given initialized kubernetes client
func NewVolumesTable(kubeclient kubernetes.Interface) *VolumesTable {
	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("type"),
		table.TextColumn("source"),
		table.TextColumn("from_pod"),
	}
	return &VolumesTable{
		name:    "kubernetes_volumes",
		columns: columns,
		client:  kubeclient,
	}
}

// Name Returns name of table
func (t *VolumesTable) Name() string {
	return t.name
}

// Columns Retrieves the initialized columns
func (t *VolumesTable) Columns() []table.ColumnDefinition {
	return t.columns
}

// Generate uses the api to retrieve information on all pods
func (t *VolumesTable) Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	allVolumes := t.getVolumesFromAllPods()
	rows := make([]map[string]string, len(allVolumes))

	for i, volumeRow := range allVolumes {
		rows[i] = map[string]string{
			"name":     volumeRow.volume.Name,
			"from_pod": volumeRow.fromPod,
		}
		// if the volume source is not with zero value
		if (corev1.VolumeSource{}) != volumeRow.volume.VolumeSource {
			rows[i]["type"], rows[i]["source"] = t.getPathAndTypeFromVolume(&volumeRow.volume.VolumeSource)
		}
	}
	return rows, nil
}

// getPathAndTypeFromVolume gets the
func (t *VolumesTable) getPathAndTypeFromVolume(volume *corev1.VolumeSource) (string, string) {
	var typ, source string
	// Because the VolumeSource struct contains alot of optional fields,
	// We use the marshal unmarshal to filter the zero values, and get the
	// json name representation of the only non zero type in the struct
	if bytes, err := json.Marshal(*volume); err == nil {
		output := make(map[string]map[string]interface{})
		_ = json.Unmarshal(bytes, &output)
		for k, v := range output {
			typ = k
			strRepr, _ := json.Marshal(v)
			source = string(strRepr)
		}
	}
	return typ, source
}

func (t *VolumesTable) getVolumesFromAllPods() []*VolumeRow {
	pods, err := t.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Println("could not get pods from k8s api")
		return nil
	}
	volumes := make([]*VolumeRow, 0)
	for _, pod := range pods.Items {
		if pod.Spec.Volumes != nil {
			for _, volume := range pod.Spec.Volumes {
				volumes = append(volumes, &VolumeRow{
					volume:  volume.DeepCopy(),
					fromPod: pod.Name,
				})
			}
		}
	}
	return volumes
}
