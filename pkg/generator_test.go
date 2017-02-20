package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func TestChartForPod(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/pod/input/pod.yaml")
	assert.Nil(t, err)
	pod := kubeapi.Pod{}
	err = yaml.Unmarshal(yamlFile, &pod)
	assert.Nil(t, err)
	name := pod.Name
	cleanUpObjectMeta(&pod.ObjectMeta)
	cleanUpPodSpec(&pod.Spec)
	template, values := podTemplate(pod)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/pod/output/pod_chart.yaml")
	assert.Nil(t, err)
	expectedValues, err := ioutil.ReadFile("../testdata/pod/output/pod_value.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueFileData, err := ylib.Marshal(valueFile)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(expectedValues)), strings.TrimSpace(string(valueFileData)))
}

func TestChartForRc(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/rc/input/rc.yaml")
	assert.Nil(t, err)
	rc := kubeapi.ReplicationController{}
	err = yaml.Unmarshal(yamlFile, &rc)
	assert.Nil(t, err)
	cleanUpObjectMeta(&rc.ObjectMeta)
	cleanUpPodSpec(&rc.Spec.Template.Spec)
	name := rc.Name
	template, values := replicationControllerTemplate(rc)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/rc/output/rc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	expectedValues, err := ioutil.ReadFile("../testdata/rc/output/rc_value.yaml")
	assert.Nil(t, err)
	valueFileData, err := ylib.Marshal(valueFile)
	assert.Nil(t, err)
	assert.Equal(t, string(expectedValues), string(valueFileData))
}

func TestChartForReplicaSet(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/replicaset/input/replicaset.yaml")
	assert.Nil(t, err)
	rcSet := ext.ReplicaSet{}
	err = yaml.Unmarshal(yamlFile, &rcSet)
	assert.Nil(t, err)
	cleanupForReplicaSets(&rcSet)
	name := rcSet.Name
	template, values := replicaSetTemplate(rcSet)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/replicaset/output/replicaset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/replicaset/output/replicaset_value.yaml", valueFile)
}

func TestChartForJob(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/job/input/job.yaml")
	assert.Nil(t, err)
	job := batch.Job{}
	err = yaml.Unmarshal(yamlFile, &job)
	assert.Nil(t, err)
	cleanUpObjectMeta(&job.ObjectMeta)
	cleanUpPodSpec(&job.Spec.Template.Spec)
	template, values := jobTemplate(job)
	name := job.Name
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/job/output/job_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/job/output/job_value.yaml", valueFile)
}

func TestChartForConfigMap(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/configmap/input/configmap.yaml")
	assert.Nil(t, err)
	configMap := kubeapi.ConfigMap{}
	err = yaml.Unmarshal(yamlFile, &configMap)
	assert.Nil(t, err)
	cleanUpObjectMeta(&configMap.ObjectMeta)
	template, _ := configMapTemplate(configMap)
	expectedTemplate, err := ioutil.ReadFile("../testdata/configmap/output/configmap_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForDaemonsets(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/daemon/input/daemon.yaml")
	assert.Nil(t, err)
	daemonset := ext.DaemonSet{}
	err = yaml.Unmarshal(yamlFile, &daemonset)
	assert.Nil(t, err)
	cleanUpObjectMeta(&daemonset.ObjectMeta)
	cleanUpPodSpec(&daemonset.Spec.Template.Spec)
	template, values := daemonsetTemplate(daemonset)
	name := daemonset.Name
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/daemon/output/daemon_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/daemon/output/daemon_value.yaml", valueFile)
}

func TestChartForSecret(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/secret/input/secret.yaml")
	assert.Nil(t, err)
	secret := kubeapi.Secret{}
	err = yaml.Unmarshal(yamlFile, &secret)
	name := secret.Name
	assert.Nil(t, err)
	cleanUpObjectMeta(&secret.ObjectMeta)
	template, values := secretTemplate(secret)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/secret/output/secret_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/secret/output/secret_value.yaml", valueFile)
}

func TestChartForPv(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/pv/input/pv.yaml")
	assert.Nil(t, err)
	pv := kubeapi.PersistentVolume{}
	err = yaml.Unmarshal(yamlFile, &pv)
	assert.Nil(t, err)
	name := pv.Name
	cleanUpObjectMeta(&pv.ObjectMeta)
	template, values := pvTemplate(pv)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/pv/output/pv_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/pv/output/pv_value.yaml", valueFile)
}

func TestChartForService(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/service/input/service.yaml")
	assert.Nil(t, err)
	svc := kubeapi.Service{}
	err = yaml.Unmarshal(yamlFile, &svc)
	assert.Nil(t, err)
	name := svc.Name
	cleanUpObjectMeta(&svc.ObjectMeta)
	template, values := serviceTemplate(svc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/service/output/service_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	values.MergeInto(valueFile, generateSafeKey(name))
	valueChecker(t, "../testdata/service/output/service_value.yaml", valueFile)
}

func TestChartForPvc(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/pvc/input/pvc.yaml")
	assert.Nil(t, err)
	pvc := kubeapi.PersistentVolumeClaim{}
	err = yaml.Unmarshal(yamlFile, &pvc)
	assert.Nil(t, err)
	cleanUpObjectMeta(&pvc.ObjectMeta)
	cleanUpDecorators(pvc.ObjectMeta.Annotations)
	template, values := pvcTemplate(pvc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pvc/output/pvc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	persistence = addPersistence(persistence, values.persistence)
	if len(persistence) != 0 {
		valueFile["persistence"] = persistence
	}
	valueChecker(t, "../testdata/pvc/output/pvc_value.yaml", valueFile)
}

func TestChartForDeployment(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/deployment/input/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	cleanUpObjectMeta(&deployment.ObjectMeta)
	cleanUpPodSpec(&deployment.Spec.Template.Spec)
	cleanUpDecorators(deployment.ObjectMeta.Annotations)
	name := deployment.Name
	template, values := deploymentTemplate(deployment)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/deployment/output/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/deployment/output/deployment_value.yaml", valueFile)
}

func TestChartForStorageClass(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/storageclass/input/storageclass.yaml")
	assert.Nil(t, err)
	storageclass := storage.StorageClass{}
	err = yaml.Unmarshal(yamlFile, &storageclass)
	assert.Nil(t, err)
	cleanUpObjectMeta(&storageclass.ObjectMeta)
	name := storageclass.Name
	template, values := storageClassTemplate(storageclass)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/storageclass/output/storageclass_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/storageclass/output/storageclass_value.yaml", valueFile)
}

func TestChartForStatefulsets(t *testing.T) {
	var valueFile = make(map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile("../testdata/statefulset/input/statefulset.yaml")
	assert.Nil(t, err)
	statefulset := apps.StatefulSet{}
	err = yaml.Unmarshal(yamlFile, &statefulset)
	assert.Nil(t, err)
	cleanUpObjectMeta(&statefulset.ObjectMeta)
	cleanUpPodSpec(&statefulset.Spec.Template.Spec)
	name := statefulset.Name
	template, values := statefulsetTemplate(statefulset)
	values.MergeInto(valueFile, generateSafeKey(name))
	expectedTemplate, err := ioutil.ReadFile("../testdata/statefulset/output/statefulset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	valueChecker(t, "../testdata/statefulset/output/statefulset_value.yaml", valueFile)
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

func valueChecker(t *testing.T, expectedPath string, valueFile map[string]interface{}) {
	valuesInfo, err := ylib.Marshal(valueFile)
	assert.Nil(t, err)
	expectedValues, err := ioutil.ReadFile(expectedPath)
	assert.Nil(t, err)
	assert.Equal(t, string(expectedValues), string(valuesInfo))
}
