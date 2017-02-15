package pkg

import (
	"io/ioutil"
	//"os"
	//"path/filepath"
	"testing"

	"fmt"
	"github.com/ghodss/yaml"
	ylib "github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	kubeapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	"strings"
)

func TestChartForPod(t *testing.T) {
	fmt.Println("Testing for Pod...\n")
	valueFile := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
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
	persistence = addPersistence(persistence, values.persistence)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pod/output/pod_chart.yaml")
	assert.Nil(t, err)
	expectedValues, err := ioutil.ReadFile("../testdata/pod/output/pod_value.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
	if len(persistence) != 0 {
		valueFile["persistence"] = persistence
	}
	valueFileData, err := ylib.Marshal(valueFile)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(expectedValues)), strings.TrimSpace(string(valueFileData)))
}

func TestChartForRc(t *testing.T) {
	valueFile := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	fmt.Println("Testing for replication controllers...\n")
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
	persistence = addPersistence(persistence, values.persistence)
	expectedTemplate, err := ioutil.ReadFile("../testdata/rc/output/rc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(expectedTemplate)), strings.TrimSpace(string(template)))
	assert.Equal(t, string(expectedTemplate), string(template))
	expectedValues, err := ioutil.ReadFile("../testdata/rc/output/rc_value.yaml")
	if len(persistence) != 0 {
		valueFile["persistence"] = persistence
	}
	valueFileData, err := ylib.Marshal(valueFile)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(expectedValues)), strings.TrimSpace(string(valueFileData)))
}

func TestChartForReplicaSet(t *testing.T) {
	fmt.Println("Testing for Replica sets...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/replicaset/input/replicaset.yaml")
	assert.Nil(t, err)
	rcSet := ext.ReplicaSet{}
	err = yaml.Unmarshal(yamlFile, &rcSet)
	assert.Nil(t, err)
	cleanUpObjectMeta(&rcSet.ObjectMeta)
	cleanUpPodSpec(&rcSet.Spec.Template.Spec)
	cleanUpDecorators(rcSet.ObjectMeta.Annotations)
	cleanUpDecorators(rcSet.ObjectMeta.Labels)
	cleanUpDecorators(rcSet.Spec.Selector.MatchLabels)
	cleanUpDecorators(rcSet.Spec.Template.ObjectMeta.Labels)
	template, _ := replicaSetTemplate(rcSet)
	expectedTemplate, err := ioutil.ReadFile("../testdata/replicaset/output/replicaset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForJob(t *testing.T) {
	fmt.Println("Testing for jobs...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/job/input/job.yaml")
	assert.Nil(t, err)
	job := batch.Job{}
	err = yaml.Unmarshal(yamlFile, &job)
	assert.Nil(t, err)
	cleanUpObjectMeta(&job.ObjectMeta)
	cleanUpPodSpec(&job.Spec.Template.Spec)
	template, _ := jobTemplate(job)
	expectedTemplate, err := ioutil.ReadFile("../testdata/job/output/job_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartCreatForDeployment(t *testing.T) {
	fmt.Println("Testing for Deployment...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/deployment/input/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	cleanUpObjectMeta(&deployment.ObjectMeta)
	cleanUpPodSpec(&deployment.Spec.Template.Spec)
	cleanUpDecorators(deployment.ObjectMeta.Annotations)
	template, _ := deploymentTemplate(deployment)
	expectedTemplate, err := ioutil.ReadFile("../testdata/deployment/output/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForConfigMap(t *testing.T) {
	fmt.Println("Testing for Configmaps...\n")
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

func TestChartForDaemon(t *testing.T) {
	fmt.Println("Testing for DaemonSets...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/daemon/input/daemon.yaml")
	assert.Nil(t, err)
	daemonset := ext.DaemonSet{}
	err = yaml.Unmarshal(yamlFile, &daemonset)
	assert.Nil(t, err)
	cleanUpObjectMeta(&daemonset.ObjectMeta)
	cleanUpPodSpec(&daemonset.Spec.Template.Spec)
	template, _ := daemonsetTemplate(daemonset)
	expectedTemplate, err := ioutil.ReadFile("../testdata/daemon/output/daemon_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForSecret(t *testing.T) {
	fmt.Println("Testing for Secrets...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/secret/input/secret.yaml")
	assert.Nil(t, err)
	secret := kubeapi.Secret{}
	err = yaml.Unmarshal(yamlFile, &secret)
	assert.Nil(t, err)
	cleanUpObjectMeta(&secret.ObjectMeta)
	template, _ := secretTemplate(secret)
	fmt.Println(template, "\n\n")
	expectedTemplate, err := ioutil.ReadFile("../testdata/secret/output/secret_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForPv(t *testing.T) {
	fmt.Println("Testing for persistent volumes...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/pv/input/pv.yaml")
	assert.Nil(t, err)
	pv := kubeapi.PersistentVolume{}
	err = yaml.Unmarshal(yamlFile, &pv)
	assert.Nil(t, err)
	cleanUpObjectMeta(&pv.ObjectMeta)
	template, _ := pvTemplate(pv)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pv/output/pv_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForService(t *testing.T) {
	fmt.Println("Testing for Service...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/service/input/service.yaml")
	assert.Nil(t, err)
	svc := kubeapi.Service{}
	err = yaml.Unmarshal(yamlFile, &svc)
	assert.Nil(t, err)
	cleanUpObjectMeta(&svc.ObjectMeta)
	template, _ := serviceTemplate(svc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/service/output/service_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForPvc(t *testing.T) {
	fmt.Println("Testing PersistentVolumeClaim...\n")
	yamlFile, err := ioutil.ReadFile("../testdata/pvc/input/pvc.yaml")
	assert.Nil(t, err)
	pvc := kubeapi.PersistentVolumeClaim{}
	err = yaml.Unmarshal(yamlFile, &pvc)
	assert.Nil(t, err)
	cleanUpObjectMeta(&pvc.ObjectMeta)
	template, _ := pvcTemplate(pvc)
	expectedTemplate, err := ioutil.ReadFile("../testdata/pvc/output/pvc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

func TestChartForDeployment(t *testing.T) {
	fmt.Println("Testing Deployment...\n\n")
	yamlFile, err := ioutil.ReadFile("../testdata/deployment/input/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	cleanUpObjectMeta(&deployment.ObjectMeta)
	cleanUpPodSpec(&deployment.Spec.Template.Spec)
	cleanUpDecorators(deployment.ObjectMeta.Annotations)
	template, _ := deploymentTemplate(deployment)
	expectedTemplate, err := ioutil.ReadFile("../testdata/deployment/output/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(expectedTemplate), string(template))
}

/*func TestChartForMultipleObject(t *testing.T) {
	yamlFiles := readLocalFiles("../testdata/mix_objects/input")
	tmp, err := ioutil.TempDir(os.TempDir(), "test")
	assert.Nil(t, err)
	g := Generator{
		ChartName: "test",
		YamlFiles: yamlFiles,
		Location:  tmp,
	}
	chdir, _ := g.Create()
	files, err := ioutil.ReadDir("../testdata/mix_objects/output")
	assert.Nil(t, err)
	for _, v := range files {
		acturalData, err := ioutil.ReadFile(filepath.Join(chdir, "templates", v.Name()))
		assert.Nil(t, err)
		expectedData, err := ioutil.ReadFile(filepath.Join("../testdata/mix_objects/output", v.Name()))
		assert.Equal(t, string(expectedData), string(acturalData))
	}
	os.Remove(chdir)
}*/
