package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	ylib "github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	kubeapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/storage"
)

func TestPodTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pod/input/pod.yaml")
	assert.Nil(t, err)
	pod := kubeapi.Pod{}
	err = yaml.Unmarshal(yamlFile, &pod)
	assert.Nil(t, err)
	template, values := podTemplate(pod)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pod/output/pod_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/pod/output/pod_value.yaml", values.value)
	assert.Nil(t, err)
}

func TestReplicationControllerTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/rc/input/rc.yaml")
	assert.Nil(t, err)
	rc := kubeapi.ReplicationController{}
	err = yaml.Unmarshal(yamlFile, &rc)
	assert.Nil(t, err)
	template, values := replicationControllerTemplate(rc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/rc/output/rc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/rc/output/rc_value.yaml", values.value)
}

func TestReplicaSetTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/replicaset/input/replicaset.yaml")
	assert.Nil(t, err)
	rcSet := ext.ReplicaSet{}
	err = yaml.Unmarshal(yamlFile, &rcSet)
	assert.Nil(t, err)
	template, values := replicaSetTemplate(rcSet)
	expectedTemplate, err := ioutil.ReadFile("../testdata/replicaset/output/replicaset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/replicaset/output/replicaset_value.yaml", values.value)
}

func TestJobTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/job/input/job.yaml")
	assert.Nil(t, err)
	job := batch.Job{}
	err = yaml.Unmarshal(yamlFile, &job)
	assert.Nil(t, err)
	template, values := jobTemplate(job)
	expectedTemplate, err := ioutil.ReadFile("../testdata/job/output/job_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/job/output/job_value.yaml", values.value)
}

func TestChartForConfigMap(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/configmap/input/configmap.yaml")
	assert.Nil(t, err)
	configMap := kubeapi.ConfigMap{}
	err = yaml.Unmarshal(yamlFile, &configMap)
	assert.Nil(t, err)
	template, values := configMapTemplate(configMap)
	expectedTemplate, err := ioutil.ReadFile("../testdata/configmap/output/configmap_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/configmap/output/configmap_value.yaml", values.value)
}

func TestDaemonsetTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/daemon/input/daemon.yaml")
	assert.Nil(t, err)
	daemonset := ext.DaemonSet{}
	err = yaml.Unmarshal(yamlFile, &daemonset)
	assert.Nil(t, err)
	template, values := daemonsetTemplate(daemonset)
	expectedTemplate, err := ioutil.ReadFile("../testdata/daemon/output/daemon_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/daemon/output/daemon_value.yaml", values.value)
}

func TestSecretTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/secret/input/secret.yaml")
	assert.Nil(t, err)
	secret := kubeapi.Secret{}
	err = yaml.Unmarshal(yamlFile, &secret)
	assert.Nil(t, err)
	template, values := secretTemplate(secret)
	expectedTemplate, err := ioutil.ReadFile("../testdata/secret/output/secret_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/secret/output/secret_value.yaml", values.value)
}

func TestPVTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pv/input/pv.yaml")
	assert.Nil(t, err)
	pv := kubeapi.PersistentVolume{}
	err = yaml.Unmarshal(yamlFile, &pv)
	assert.Nil(t, err)
	template, values := pvTemplate(pv)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pv/output/pv_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/pv/output/pv_value.yaml", values.value)
}

func TestServiceTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/service/input/service.yaml")
	assert.Nil(t, err)
	svc := kubeapi.Service{}
	err = yaml.Unmarshal(yamlFile, &svc)
	assert.Nil(t, err)
	template, values := serviceTemplate(svc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/service/output/service_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/service/output/service_value.yaml", values.value)
}

func TestPVCTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pvc/input/pvc.yaml")
	assert.Nil(t, err)
	pvc := kubeapi.PersistentVolumeClaim{}
	err = yaml.Unmarshal(yamlFile, &pvc)
	assert.Nil(t, err)
	template, values := pvcTemplate(pvc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pvc/output/pvc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/pvc/output/pvc_value.yaml", values.persistence)
}

func TestDeploymentTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/deployment/input/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	template, values := deploymentTemplate(deployment)
	expectedTemplate, err := ioutil.ReadFile("../testdata/deployment/output/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/deployment/output/deployment_value.yaml", values.value)
}

func TestStorageClassTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/storageclass/input/storageclass.yaml")
	assert.Nil(t, err)
	storageclass := storage.StorageClass{}
	err = yaml.Unmarshal(yamlFile, &storageclass)
	assert.Nil(t, err)
	template, values := storageClassTemplate(storageclass)
	expectedTemplate, err := ioutil.ReadFile("../testdata/storageclass/output/storageclass_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/storageclass/output/storageclass_value.yaml", values.value)
}

func TestStatefulsetTemplate(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/statefulset/input/statefulset.yaml")
	assert.Nil(t, err)
	statefulset := apps.StatefulSet{}
	err = yaml.Unmarshal(yamlFile, &statefulset)
	assert.Nil(t, err)
	template, values := statefulsetTemplate(statefulset)
	expectedTemplate, err := ioutil.ReadFile("../testdata/statefulset/output/statefulset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/statefulset/output/statefulset_value.yaml", values.value)
}

func TestServiceTemplateWithClusterIP(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/service_clusterIP/input/service.yaml")
	assert.Nil(t, err)
	svc := kubeapi.Service{}
	err = yaml.Unmarshal(yamlFile, &svc)
	assert.Nil(t, err)
	template, values := serviceTemplate(svc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/service_clusterIP/output/service_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/service_clusterIP/output/service_value.yaml", values.value)
}

func TestChartForVolume(t *testing.T) {
	yamlFiles := ReadLocalFiles("../testdata/mix_objects/check_volume/input")
	tmp, err := ioutil.TempDir(os.TempDir(), "test")
	defer os.Remove(tmp)
	assert.Nil(t, err)
	g := Generator{
		ChartName: "test",
		YamlFiles: yamlFiles,
		Location:  tmp,
	}
	chdir, err := g.Create()
	assert.Nil(t, err)
	files, err := ioutil.ReadDir("../testdata/mix_objects/check_volume/output")
	assert.Nil(t, err)
	for _, v := range files {
		acturalData, err := ioutil.ReadFile(filepath.Join(chdir, "templates", v.Name()))
		assert.Nil(t, err)
		expectedData, err := ioutil.ReadFile(filepath.Join("../testdata/mix_objects/check_volume/output", v.Name()))
		assert.Equal(t, string(expectedData), string(acturalData))
	}
	os.Remove(chdir)
}

func valueChecker(t *testing.T, expectedPath string, value map[string]interface{}) {
	valuesInfo, err := ylib.Marshal(value)
	assert.Nil(t, err)
	expectedValues, err := ioutil.ReadFile(expectedPath)
	assert.Nil(t, err)
	assert.Equal(t, string(expectedValues), string(valuesInfo))
}
