package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/appscode/go/encoding/yaml"
	ylib "github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	storage "k8s.io/client-go/pkg/apis/storage/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Generator struct {
	Location  string
	ChartName string
	YamlFiles []string
}

var ChartObject map[string][]string
var chnageObjectType = []string{"Secret", "Configmap", "PersistentVolume", "PersistentVolumeClaim"}

func (g Generator) Create() (string, error) {
	chartfile := chartMetaData(g.ChartName)
	imageTag := "" //TODO
	fmt.Println("Creating chart...")
	cdir := filepath.Join(g.Location, chartfile.Name)
	fi, err := os.Stat(cdir)
	if err == nil && !fi.IsDir() {
		return cdir, fmt.Errorf("%s already exists and is not a directory", cdir)
	}
	ChartObject = getInsideObjects(g.YamlFiles)
	if err := os.MkdirAll(cdir, 0755); err != nil {
		return cdir, err
	}
	cf := filepath.Join(cdir, ChartfileName)
	if _, err := os.Stat(cf); err != nil {
		if len(chartfile.Version) == 0 {
			chartfile.Version = imageTag
		}
		if err := SaveChartfile(cf, &chartfile); err != nil {
			return cdir, err
		}
	}
	valueFile := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	templateLocation := filepath.Join(cdir, TemplatesDir)
	err = os.MkdirAll(templateLocation, 0755)
	for _, kubeObj := range g.YamlFiles {
		kubeJson, err := yaml.ToJSON([]byte(kubeObj))
		if err != nil {
			log.Fatal(err)
		}

		var objMeta metav1.TypeMeta
		if err := json.Unmarshal(kubeJson, &objMeta); err != nil {
			log.Fatal(err)
		}

		values := valueFileGenerator{}
		var template, templateName string
		if objMeta.Kind == "Pod" {
			pod := apiv1.Pod{}
			if err := json.Unmarshal(kubeJson, &pod); err != nil {
				log.Fatal(err)
			}
			name := pod.Name
			templateName = filepath.Join(templateLocation, name+".pod.yaml")
			template, values = podTemplate(pod)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "ReplicationController" {
			rc := apiv1.ReplicationController{}
			if err := json.Unmarshal(kubeJson, &rc); err != nil {
				log.Fatal(err)
			}
			name := rc.Name
			templateName = filepath.Join(templateLocation, name+".rc.yaml")
			template, values = replicationControllerTemplate(rc)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "Deployment" {
			deployment := extensions.Deployment{}
			if err := json.Unmarshal(kubeJson, &deployment); err != nil {
				log.Fatal(err)
			}
			name := deployment.Name
			templateName = filepath.Join(templateLocation, name+".deployment.yaml")
			template, values = deploymentTemplate(deployment)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "Job" {
			job := batch.Job{}
			if err := json.Unmarshal(kubeJson, &job); err != nil {
				log.Fatal(err)
			}
			name := job.Name
			templateName = filepath.Join(templateLocation, name+".job.yaml")
			template, values = jobTemplate(job)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "DaemonSet" {
			daemonset := extensions.DaemonSet{}
			if err := json.Unmarshal(kubeJson, &daemonset); err != nil {
				log.Fatal(err)
			}
			name := daemonset.Name
			templateName = filepath.Join(templateLocation, name+".daemonset.yaml")
			template, values = daemonsetTemplate(daemonset)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "ReplicaSet" {
			rcSet := extensions.ReplicaSet{}
			if err := json.Unmarshal(kubeJson, &rcSet); err != nil {
				log.Fatal(err)
			}
			name := rcSet.Name
			templateName = filepath.Join(templateLocation, name+".rs.yaml")
			template, values = replicaSetTemplate(rcSet)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "StatefulSet" {
			statefulset := apps.StatefulSet{}
			if err := json.Unmarshal(kubeJson, &statefulset); err != nil {
				log.Fatal(err)
			}
			name := statefulset.Name
			templateName = filepath.Join(templateLocation, name+".statefulset.yaml")
			template, values = statefulsetTemplate(statefulset)
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "Service" {
			service := apiv1.Service{}
			if err := json.Unmarshal(kubeJson, &service); err != nil {
				log.Fatal(err)
			}
			template, values = serviceTemplate(service)
			name := service.Name
			templateName = filepath.Join(templateLocation, name+".svc.yaml")
			values.MergeInto(valueFile, generateSafeKey(name))
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "ConfigMap" {
			configMap := apiv1.ConfigMap{}
			if err := json.Unmarshal(kubeJson, &configMap); err != nil {
				log.Fatal(err)
			}
			name := configMap.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = configMapTemplate(configMap)
			values.MergeInto(valueFile, generateSafeKey(name))
		} else if objMeta.Kind == "Secret" {
			secret := apiv1.Secret{}
			if err := json.Unmarshal(kubeJson, &secret); err != nil {
				log.Fatal(err)
			}
			name := secret.Name
			templateName = filepath.Join(templateLocation, name+".secret.yaml")
			template, values = secretTemplate(secret)
			values.MergeInto(valueFile, generateSafeKey(name))
		} else if objMeta.Kind == "PersistentVolumeClaim" {
			pvc := apiv1.PersistentVolumeClaim{}
			if err := json.Unmarshal(kubeJson, &pvc); err != nil {
				log.Fatal(err)
			}
			name := pvc.Name
			templateName = filepath.Join(templateLocation, name+".pvc.yaml")
			template, values = pvcTemplate(pvc)
			persistence = addPersistence(persistence, values.persistence)
		} else if objMeta.Kind == "PersistentVolume" {
			pv := apiv1.PersistentVolume{}
			if err := json.Unmarshal(kubeJson, &pv); err != nil {
				log.Fatal(err)
			}
			name := pv.Name
			templateName = filepath.Join(templateLocation, name+".pv.yaml")
			template, values = pvTemplate(pv)
			values.MergeInto(valueFile, generateSafeKey(name))
		} else if objMeta.Kind == "StorageClass" {
			storageClass := storage.StorageClass{}
			if err := json.Unmarshal(kubeJson, &storageClass); err != nil {
				log.Fatal(err)
			}
			name := storageClass.Name
			templateName = filepath.Join(templateLocation, name+".storage.yaml")
			template, values = storageClassTemplate(storageClass)
			values.MergeInto(valueFile, generateSafeKey(name))
		} else {
			fmt.Printf("%v is not supported. Please add manually. Consider filing bug here: https://github.com/kubepack/chartify/issues", objMeta.Kind)
		}
		if err := ioutil.WriteFile(templateName, []byte(template), 0644); err != nil {
			log.Fatal(err)
		}
	}
	if len(persistence) != 0 {
		valueFile["persistence"] = persistence
	}
	valueFileData, err := ylib.Marshal(valueFile)
	if err != nil {
		log.Fatal(err)
	}
	helperDir := filepath.Join(templateLocation, HelpersName)
	if err := ioutil.WriteFile(helperDir, []byte(defaultHelpers), 0644); err != nil {
		log.Fatal(err)
	}
	valueDir := filepath.Join(cdir, ValuesfileName)
	if err := ioutil.WriteFile(valueDir, []byte(valueFileData), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("CREATE : SUCCESSFUL")
	return cdir, nil
}

func cleanUpObjectMeta(m *metav1.ObjectMeta) {
	var t metav1.Time
	m.GenerateName = ""
	m.SelfLink = ""
	m.UID = types.UID("")
	m.ResourceVersion = ""
	m.Generation = 0
	m.CreationTimestamp = t
	m.DeletionTimestamp = nil
}

func cleanUpDecorators(m map[string]string) {
	delete(m, "controller-uid")
	delete(m, "deployment.kubernetes.io/desired-replicas")
	delete(m, "deployment.kubernetes.io/max-replicas")
	delete(m, "deployment.kubernetes.io/revision")
	delete(m, "pod-template-hash")
	delete(m, "pv.kubernetes.io/bind-completed")
	delete(m, "pv.kubernetes.io/bound-by-controller")
}

func cleanUpPodSpec(p *apiv1.PodSpec) {
	p.DNSPolicy = apiv1.DNSPolicy("")
	p.NodeName = ""
	if p.ServiceAccountName == "default" {
		p.ServiceAccountName = ""
	}
	p.TerminationGracePeriodSeconds = nil
	for i, c := range p.Containers {
		c.TerminationMessagePath = ""
		p.Containers[i] = c
	}
	for i, c := range p.InitContainers {
		c.TerminationMessagePath = ""
		p.InitContainers[i] = c
	}
}

func podTemplate(pod apiv1.Pod) (string, valueFileGenerator) {
	cleanUpObjectMeta(&pod.ObjectMeta)
	cleanUpPodSpec(&pod.Spec)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(pod.ObjectMeta.Name)
	pod.ObjectMeta = generateObjectMetaTemplate(pod.ObjectMeta, key, value, pod.ObjectMeta.Name)
	pod.Spec = generateTemplateForPodSpec(pod.Spec, key, value)
	if len(pod.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(pod.Spec.Volumes, key, value)
		pod.Spec.Volumes = nil
	}
	tempPodByte, err := ylib.Marshal(pod)
	if err != nil {
		log.Fatal(err)
	}
	tempPod := removeEmptyFields(string(tempPodByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplateForPod(string(tempPod), volumes)
	} else {
		template = string(tempPod)
	}
	data := valueFileGenerator{
		value:       value,
		persistence: persistence,
	}
	return template, data
}

func replicationControllerTemplate(rc apiv1.ReplicationController) (string, valueFileGenerator) {
	cleanUpObjectMeta(&rc.ObjectMeta)
	cleanUpPodSpec(&rc.Spec.Template.Spec)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(rc.ObjectMeta.Name)
	rc.ObjectMeta = generateObjectMetaTemplate(rc.ObjectMeta, key, value, rc.ObjectMeta.Name)
	rc.Spec.Template.Spec = generateTemplateForPodSpec(rc.Spec.Template.Spec, key, value)
	if len(rc.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(rc.Spec.Template.Spec.Volumes, key, value)
		value[Persistence] = true
		rc.Spec.Template.Spec.Volumes = nil
	}
	tempRcByte, err := ylib.Marshal(rc)
	if err != nil {
		log.Fatal(err)
	}
	tempRc := removeEmptyFields(string(tempRcByte))

	tempRc, value = generateTemplateReplicationCtrSpec(rc.Spec, tempRc, key, value)

	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempRc, volumes)
	} else {
		template = tempRc
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func replicaSetTemplate(replicaSet extensions.ReplicaSet) (string, valueFileGenerator) {
	cleanupForReplicaSets(&replicaSet)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(replicaSet.ObjectMeta.Name)
	replicaSet.ObjectMeta = generateObjectMetaTemplate(replicaSet.ObjectMeta, key, value, replicaSet.ObjectMeta.Name)
	replicaSet.Spec.Template.Spec = generateTemplateForPodSpec(replicaSet.Spec.Template.Spec, key, value)
	if len(replicaSet.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(replicaSet.Spec.Template.Spec.Volumes, key, value)
		value[Persistence] = true
		replicaSet.Spec.Template.Spec.Volumes = nil
	}
	if replicaSet.Spec.Selector != nil {
		modifyLabelSelector(replicaSet.Spec.Selector, replicaSet.Spec.Template.Labels, replicaSet.ObjectMeta.Labels)
	}
	template := ""
	tempRcSetByte, err := ylib.Marshal(replicaSet)
	if err != nil {
		log.Fatal(err)
	}
	tempReplicaSet := removeEmptyFields(string(tempRcSetByte))

	tempReplicaSet, value = generateTemplateReplicaSetSpec(replicaSet.Spec, tempReplicaSet, key, value)

	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempReplicaSet, volumes) // RC and replica_set has volume in same layer
	} else {
		template = tempReplicaSet
	}
	return template, valueFileGenerator{
		value:       value,
		persistence: persistence,
	}
}

func deploymentTemplate(deployment extensions.Deployment) (string, valueFileGenerator) {
	cleanUpObjectMeta(&deployment.ObjectMeta)
	cleanUpPodSpec(&deployment.Spec.Template.Spec)
	cleanUpDecorators(deployment.ObjectMeta.Annotations)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(deployment.ObjectMeta.Name)
	deployment.ObjectMeta = generateObjectMetaTemplate(deployment.ObjectMeta, key, value, deployment.ObjectMeta.Name)
	deployment.Spec.Template.Spec = generateTemplateForPodSpec(deployment.Spec.Template.Spec, key, value)
	if len(deployment.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(deployment.Spec.Template.Spec.Volumes, key, value)
		deployment.Spec.Template.Spec.Volumes = nil
	}

	if len(string(deployment.Spec.Strategy.Type)) != 0 {
		if len(string(deployment.Spec.Strategy.Type)) != 0 {
			value[DeploymentStrategy] = deployment.Spec.Strategy.Type
			deployment.Spec.Strategy.Type = extensions.DeploymentStrategyType(fmt.Sprintf("{{.Values.%s.%s}}", key, DeploymentStrategy))
		}
	}

	if deployment.Spec.Selector != nil {
		modifyLabelSelector(deployment.Spec.Selector, deployment.Spec.Template.Labels, deployment.ObjectMeta.Labels)
	}

	template := ""
	tempDeploymentByte, err := ylib.Marshal(deployment)
	if err != nil {
		log.Fatal(err)
	}
	tempDeployment := removeEmptyFields(string(tempDeploymentByte))

	tempDeployment, value = generateTemplateDeplymentSpec(deployment.Spec, tempDeployment, key, value)

	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempDeployment, volumes)
	} else {
		template = tempDeployment
	}

	return template, valueFileGenerator{value: value, persistence: persistence}
}

func daemonsetTemplate(daemonset extensions.DaemonSet) (string, valueFileGenerator) {
	cleanUpObjectMeta(&daemonset.ObjectMeta)
	cleanUpPodSpec(&daemonset.Spec.Template.Spec)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(daemonset.ObjectMeta.Name)
	daemonset.ObjectMeta = generateObjectMetaTemplate(daemonset.ObjectMeta, key, value, daemonset.ObjectMeta.Name)
	daemonset.Spec.Template.Spec = generateTemplateForPodSpec(daemonset.Spec.Template.Spec, key, value)
	if len(daemonset.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(daemonset.Spec.Template.Spec.Volumes, key, value)
		value[Persistence] = true
		daemonset.Spec.Template.Spec.Volumes = nil
	}

	if daemonset.Spec.Selector != nil {
		modifyLabelSelector(daemonset.Spec.Selector, daemonset.Spec.Template.Labels, daemonset.ObjectMeta.Labels)
	}

	template := ""
	tempDaemonSetByte, err := ylib.Marshal(daemonset)
	if err != nil {
		log.Fatal(err)
	}
	tempDaemonSet := removeEmptyFields(string(tempDaemonSetByte))
	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempDaemonSet, volumes)
	} else {
		template = tempDaemonSet
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func statefulsetTemplate(statefulset apps.StatefulSet) (string, valueFileGenerator) {
	cleanUpObjectMeta(&statefulset.ObjectMeta)
	cleanUpPodSpec(&statefulset.Spec.Template.Spec)
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := generateSafeKey(statefulset.ObjectMeta.Name)
	statefulset.ObjectMeta = generateObjectMetaTemplate(statefulset.ObjectMeta, key, value, statefulset.ObjectMeta.Name)
	if len(statefulset.Spec.ServiceName) != 0 {
		value[ServiceName] = statefulset.Spec.ServiceName //generateTemplateForSingleValue(statefulset.Spec.ServiceName, "ServiceName", value)
		statefulset.Spec.ServiceName = fmt.Sprintf("{{.Values.%s.%s}}", key, ServiceName)
	}
	statefulset.Spec.Template.Spec = generateTemplateForPodSpec(statefulset.Spec.Template.Spec, key, value)
	if statefulset.Spec.Selector != nil {
		modifyLabelSelector(statefulset.Spec.Selector, statefulset.Spec.Template.Labels, statefulset.ObjectMeta.Labels)
	}
	if len(statefulset.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(statefulset.Spec.Template.Spec.Volumes, key, value)
		statefulset.Spec.Template.Spec.Volumes = nil
	}
	tempStatefulSetByte, err := ylib.Marshal(statefulset)
	if err != nil {
		log.Fatal(err)
	}
	tempStatefulSet := removeEmptyFields(string(tempStatefulSetByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempStatefulSet, volumes)
	} else {
		template = tempStatefulSet
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func jobTemplate(job batch.Job) (string, valueFileGenerator) {
	cleanUpObjectMeta(&job.ObjectMeta)
	cleanUpPodSpec(&job.Spec.Template.Spec)
	cleanUpDecorators(job.ObjectMeta.Labels)
	cleanUpDecorators(job.Spec.Template.Labels)
	cleanUpDecorators(job.Spec.Selector.MatchLabels)
	volumes := ""
	persistence := make(map[string]interface{}, 0)
	value := make(map[string]interface{}, 0)
	key := generateSafeKey(job.ObjectMeta.Name)
	job.ObjectMeta = generateObjectMetaTemplate(job.ObjectMeta, key, value, job.ObjectMeta.Name)
	job.Spec.Template.Spec = generateTemplateForPodSpec(job.Spec.Template.Spec, key, value)
	if len(job.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(job.Spec.Template.Spec.Volumes, key, value)
		value[Persistence] = true
		job.Spec.Template.Spec.Volumes = nil
	}
	if job.Spec.Selector != nil {
		modifyLabelSelector(job.Spec.Selector, job.Spec.Template.Labels, job.ObjectMeta.Labels)
	}
	tempJobByte, err := ylib.Marshal(job)
	if err != nil {
		log.Fatal(err)
	}
	tempJob := removeEmptyFields(string(tempJobByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplate(tempJob, volumes)
	} else {
		template = tempJob
	}
	return template, valueFileGenerator{value: value, persistence: persistence}

}

func serviceTemplate(svc apiv1.Service) (string, valueFileGenerator) {
	cleanUpObjectMeta(&svc.ObjectMeta)
	value := make(map[string]interface{}, 0)
	key := generateSafeKey(svc.ObjectMeta.Name)
	svc.ObjectMeta = generateObjectMetaTemplate(svc.ObjectMeta, key, value, svc.ObjectMeta.Name)
	ip := net.ParseIP(svc.Spec.ClusterIP)
	if ip != nil {
		svc.Spec.ClusterIP = ""
	}
	svc.Spec = generateServiceSpecTemplate(svc.Spec, key, value)
	if svc.Spec.Selector != nil {
		svc.Spec.Selector = modifySvcLabelSelector(svc.Spec.Selector)
	}
	svcData, err := ylib.Marshal(svc)
	if err != nil {
		log.Fatal(err)
	}
	service := removeEmptyFields(string(svcData))
	return string(service), valueFileGenerator{value: value}
}

func configMapTemplate(configMap apiv1.ConfigMap) (string, valueFileGenerator) {
	cleanUpObjectMeta(&configMap.ObjectMeta)
	value := make(map[string]interface{}, 0)
	key := generateSafeKey(configMap.ObjectMeta.Name)
	configMap.ObjectMeta = generateObjectMetaTemplate(configMap.ObjectMeta, key, value, configMap.ObjectMeta.Name)
	configMapData, err := ylib.Marshal(configMap)
	if err != nil {
		log.Fatal(err)
	}
	if len(configMap.Data) != 0 {
		for k, v := range configMap.Data {
			value[k] = v
			configMap.Data[k] = (fmt.Sprintf("{{.Values.%s.%s}}", key, k))
		}
	}
	data := removeEmptyFields(string(configMapData))
	return string(data), valueFileGenerator{value: value}
}

func secretTemplate(secret apiv1.Secret) (string, valueFileGenerator) {
	cleanUpObjectMeta(&secret.ObjectMeta)
	value := make(map[string]interface{}, 0)
	secretDataMap := make(map[string]interface{}, 0)
	key := generateSafeKey(secret.ObjectMeta.Name)
	secret.ObjectMeta = generateObjectMetaTemplate(secret.ObjectMeta, key, value, secret.ObjectMeta.Name)
	if len(secret.Data) != 0 {
		for k, v := range secret.Data {
			if strings.HasPrefix(k, ".") {
				// For values that starts with ".", the Values string get populated with ".." - error for helm
				kmod := strings.Replace(k, ".", "", 1)
				value[kmod] = v
				secretDataMap[k] = (fmt.Sprintf("{{.Values.%s.%s}}", key, kmod))
			} else {
				value[k] = v
				secretDataMap[k] = (fmt.Sprintf("{{.Values.%s.%s}}", key, k))
			}
		}
	}
	secret.Data = nil
	value[Type] = secret.Type
	secret.Type = apiv1.SecretType(fmt.Sprintf("{{.Values.%s.%s}}", key, Type))
	secretDataByte, err := ylib.Marshal(secret)
	if err != nil {
		log.Fatal(err)
	}
	secretData := removeEmptyFields(string(secretDataByte))
	secretData = addSecretData(secretData, secretDataMap, key)
	return secretData, valueFileGenerator{value: value}
}

func pvcTemplate(pvc apiv1.PersistentVolumeClaim) (string, valueFileGenerator) {
	cleanUpObjectMeta(&pvc.ObjectMeta)
	cleanUpDecorators(pvc.ObjectMeta.Annotations)
	tempValue := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	rawKey := generateSafeKey(pvc.ObjectMeta.Name)
	key := Persistence + "." + rawKey
	pvc.ObjectMeta = generateObjectMetaTemplate(pvc.ObjectMeta, key, tempValue, pvc.ObjectMeta.Name)
	pvc.Spec = generatePersistentVolumeClaimSpec(pvc.Spec, key, tempValue)
	pvcData, err := ylib.Marshal(pvc)
	if err != nil {
		log.Fatal(err)
	}
	temp := removeEmptyFields(string(pvcData))
	pvcTemplateData := fmt.Sprintf("{{- if .Values.%s.%s -}}\n%s{{- end -}}", key, Enabled, temp)
	tempValue[Enabled] = true // By Default use persistence volume true
	persistence[rawKey] = tempValue
	return pvcTemplateData, valueFileGenerator{persistence: persistence}
}

func pvTemplate(pv apiv1.PersistentVolume) (string, valueFileGenerator) {
	cleanUpObjectMeta(&pv.ObjectMeta)
	value := make(map[string]interface{}, 0)
	key := generateSafeKey(pv.ObjectMeta.Name)
	pv.ObjectMeta = generateObjectMetaTemplate(pv.ObjectMeta, key, value, pv.Name)
	pv.Spec = generatePersistentVolumeSpec(pv.Spec, key, value)
	pvData, err := ylib.Marshal(pv)
	if err != nil {
		log.Fatal(err)
	}
	temp := removeEmptyFields(string(pvData))
	return string(temp), valueFileGenerator{value: value}
}

func storageClassTemplate(storageClass storage.StorageClass) (string, valueFileGenerator) {
	cleanUpObjectMeta(&storageClass.ObjectMeta)
	value := make(map[string]interface{}, 0)
	key := generateSafeKey(storageClass.ObjectMeta.Name)
	storageClass.ObjectMeta = generateObjectMetaTemplate(storageClass.ObjectMeta, key, value, storageClass.ObjectMeta.Name)
	value[Provisioner] = storageClass.Provisioner
	storageClass.Provisioner = fmt.Sprintf("{{.Values.%s.%s}}", key, Provisioner)
	storageClass.Parameters = mapToValueMaker(storageClass.Parameters, value, key)
	storageData, err := ylib.Marshal(storageClass)
	if err != nil {
		log.Fatal(err)
	}
	return string(storageData), valueFileGenerator{value: value}
}

func addSecretData(secretData string, secretDataMap map[string]interface{}, key string) string {
	elseCondition := "{{ else }}"
	elseAction := "{{ randAlphaNum 10 | b64enc | quote }}"
	end := "{{ end }}"
	data := ""
	for k, v := range secretDataMap {
		if strings.HasPrefix(k, ".") {
			// For values that starts with ".", the Values string get populated with ".." - error for helm
			kmod := strings.Replace(k, ".", "", 1)
			ifCondition := fmt.Sprintf("{{ if .Values.%s.%s }}", key, kmod)
			data += fmt.Sprintf("  %s\n  %s: %s\n  %s\n  %s: %s\n  %s\n", ifCondition, k, v, elseCondition, k, elseAction, end)
		} else {
			ifCondition := fmt.Sprintf("{{ if .Values.%s.%s }}", key, k)
			data += fmt.Sprintf("  %s\n  %s: %s\n  %s\n  %s: %s\n  %s\n", ifCondition, k, v, elseCondition, k, elseAction, end)
		}

	}
	dataOfSecret := "data:" + "\n" + data
	return (secretData + dataOfSecret)
}

func addPersistence(persistence map[string]interface{}, elements map[string]interface{}) map[string]interface{} {
	for k, v := range elements {
		persistence[k] = v
	}
	return persistence
}

func chartMetaData(name string) chart.Metadata {
	return chart.Metadata{
		Name:        name,
		Description: "Helm chart generated by https://github.com/kubepack/chartify",
		Version:     "0.1.0",
		ApiVersion:  "v1",
	}
}

func mapToValueMaker(mp map[string]string, value map[string]interface{}, key string) map[string]string {
	for k, v := range mp {
		value[k] = v
		mp[k] = fmt.Sprintf("{{.Values.%s.%s}}", key, k)
	}
	return mp
}

func getInsideObjects(objects []string) map[string][]string {
	obj := make(map[string][]string)
	for _, v := range objects {
		kind, name := getObjectKindAndName(v)
		for _, t := range chnageObjectType {
			if kind == t {
				obj[kind] = append(obj[kind], name)
			}
		}
	}
	return obj
}

func getObjectKindAndName(yamlData string) (string, string) {
	kubeJson, err := yaml.ToJSON([]byte(yamlData))
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(kubeJson, &m)
	if err != nil {
		log.Println(err)
	}
	objMeta, ok := m["metadata"].(map[string]interface{})
	if !ok {
		log.Fatal("Name not found")
	}
	name := objMeta["name"]
	var typeMeta metav1.TypeMeta

	err = json.Unmarshal(kubeJson, &typeMeta)
	if err != nil {
		log.Fatal(err)
	}
	objName, ok := name.(string)
	if !ok {
		return typeMeta.Kind, ""
	}
	return typeMeta.Kind, objName
}

func modifyLabelSelector(selector *metav1.LabelSelector, templateLabels map[string]string, metaLabels map[string]string) {
	if len(selector.MatchLabels) == 0 {
		return
	}
	for k, v := range selector.MatchLabels {
		_, ok := templateLabels[k]
		if !ok {
			continue
		}
		_, ok = metaLabels[k]
		if !ok {
			return
		}
		selector.MatchLabels[k] = "{{.Release.Name}}-" + v
		templateLabels[k] = selector.MatchLabels[k]
		metaLabels[k] = selector.MatchLabels[k]
	}
}

func modifySvcLabelSelector(selector map[string]string) map[string]string {

	for k, v := range selector {

		selector[k] = "{{.Release.Name}}-" + v
	}

	return selector
}

func cleanupForReplicaSets(rcSet *extensions.ReplicaSet) {
	cleanUpObjectMeta(&rcSet.ObjectMeta)
	cleanUpPodSpec(&rcSet.Spec.Template.Spec)
	cleanUpDecorators(rcSet.ObjectMeta.Annotations)
	cleanUpDecorators(rcSet.ObjectMeta.Labels)
	cleanUpDecorators(rcSet.Spec.Selector.MatchLabels)
	cleanUpDecorators(rcSet.Spec.Template.ObjectMeta.Labels)
}

func ReadLocalFiles(dirName string) []string {
	var yamlFiles []string
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fileDir := filepath.Join(dirName, f.Name())
		dataByte, err := ioutil.ReadFile(fileDir)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}
