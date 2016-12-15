package pkg

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	kubeapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	"testing"
)

func TestChartForPod(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pod/input/pod.yaml")
	assert.Nil(t, err)
	pod := kubeapi.Pod{}
	err = yaml.Unmarshal(yamlFile, &pod)
	assert.Nil(t, err)
	template, _ := podTemplate(pod)
	actualTemplate, err := ioutil.ReadFile("../testdata/pod/output/pod_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForRc(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/rc/input/rc.yaml")
	assert.Nil(t, err)
	rc := kubeapi.ReplicationController{}
	err = yaml.Unmarshal(yamlFile, &rc)
	assert.Nil(t, err)
	template, _ := replicationControllerTemplate(rc)
	actualTemplate, err := ioutil.ReadFile("../testdata/rc/output/rc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForReplicaSet(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/replicaset/input/replicaset.yaml")
	assert.Nil(t, err)
	replicaset := ext.ReplicaSet{}
	err = yaml.Unmarshal(yamlFile, &replicaset)
	assert.Nil(t, err)
	template, _ := replicaSetTemplate(replicaset)
	actualTemplate, err := ioutil.ReadFile("../testdata/replicaset/output/replicaset_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForJob(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/job/input/job.yaml")
	assert.Nil(t, err)
	job := batch.Job{}
	err = yaml.Unmarshal(yamlFile, &job)
	assert.Nil(t, err)
	template, _ := jobTemplate(job)
	actualTemplate, err := ioutil.ReadFile("../testdata/job/output/job_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartCreatForDeployment(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/yaml/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	template, _ := deploymentTemplate(deployment)
	actualTemplate, err := ioutil.ReadFile("../test/chart/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForConfigMap(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/configmap/input/configmap.yaml")
	assert.Nil(t, err)
	configmap := kubeapi.ConfigMap{}
	err = yaml.Unmarshal(yamlFile, &configmap)
	assert.Nil(t, err)
	template, _ := configMapTemplate(configmap)
	actualTemplate, err := ioutil.ReadFile("../testdata/configmap/output/configmap_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForDaemon(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/daemon/input/daemon.yaml")
	assert.Nil(t, err)
	daemonset := ext.DaemonSet{}
	err = yaml.Unmarshal(yamlFile, &daemonset)
	assert.Nil(t, err)
	template, _ := daemonsetTemplate(daemonset)
	actualTemplate, err := ioutil.ReadFile("../testdata/daemon/output/daemon_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForSecret(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/secret/input/secret.yaml")
	assert.Nil(t, err)
	secret := kubeapi.Secret{}
	err = yaml.Unmarshal(yamlFile, &secret)
	assert.Nil(t, err)
	template, _ := secretTemplate(secret)
	actualTemplate, err := ioutil.ReadFile("../testdata/secret/output/secret_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForPv(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pv/input/pv.yaml")
	assert.Nil(t, err)
	pv := kubeapi.PersistentVolume{}
	err = yaml.Unmarshal(yamlFile, &pv)
	assert.Nil(t, err)
	template, _ := pvTemplate(pv)
	actualTemplate, err := ioutil.ReadFile("../testdata/pv/output/pv_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForService(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/service/input/service.yaml")
	assert.Nil(t, err)
	svc := kubeapi.Service{}
	err = yaml.Unmarshal(yamlFile, &svc)
	assert.Nil(t, err)
	template, _ := serviceTemplate(svc)
	actualTemplate, err := ioutil.ReadFile("../testdata/service/output/service_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForPvc(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/pvc/input/pvc.yaml")
	assert.Nil(t, err)
	pvc := kubeapi.PersistentVolumeClaim{}
	err = yaml.Unmarshal(yamlFile, &pvc)
	assert.Nil(t, err)
	template, _ := pvcTemplate(pvc)
	actualTemplate, err := ioutil.ReadFile("../testdata/pvc/output/pvc_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}

func TestChartForDeployment(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("../testdata/deployment/input/deployment.yaml")
	assert.Nil(t, err)
	deployment := ext.Deployment{}
	err = yaml.Unmarshal(yamlFile, &deployment)
	assert.Nil(t, err)
	template, _ := deploymentTemplate(deployment)
	actualTemplate, err := ioutil.ReadFile("../testdata/deployment/output/deployment_chart.yaml")
	assert.Nil(t, err)
	assert.Equal(t, string(actualTemplate), string(template))
}
