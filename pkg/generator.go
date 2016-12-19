package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/proto/hapi/chart"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	kext "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/storage"
)

type Generator struct {
	Location  string
	ChartName string
	YamlFiles []string
}

func (g Generator) Create() (string, error) {
	chartfile := chartMetaData(g.ChartName)
	imageTag := "" //TODO
	fmt.Println("Creating Custom Chart...")
	cdir := filepath.Join(g.Location, chartfile.Name)
	fi, err := os.Stat(cdir)
	if err == nil && !fi.IsDir() {
		return cdir, fmt.Errorf("%s already exists and is not a directory", cdir)
	}
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
	var resourceType unversioned.TypeMeta
	templateLocation := filepath.Join(cdir, TemplatesDir)
	err = os.MkdirAll(templateLocation, 0755)
	template := ""
	for _, yamlData := range g.YamlFiles {
		err = yaml.Unmarshal([]byte(yamlData), &resourceType)
		if err != nil {
			log.Fatal(err)
		}
		template = ""
		values := valueFileGenerator{}
		templateName := ""
		if resourceType.Kind == "Pod" {
			pod := kapi.Pod{}
			err = yaml.Unmarshal([]byte(yamlData), &pod)
			if err != nil {
				log.Fatal(err)
			}
			name := pod.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = podTemplate(pod)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "ReplicationController" {
			rc := kapi.ReplicationController{}
			err = yaml.Unmarshal([]byte(yamlData), &rc)
			if err != nil {
				log.Fatal(err)
			}
			name := rc.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = replicationControllerTemplate(rc)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "Deployment" {
			deployment := kext.Deployment{}
			err = yaml.Unmarshal([]byte(yamlData), &deployment)
			name := deployment.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			if err != nil {
				log.Fatal(err)
			}
			template, values = deploymentTemplate(deployment)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "Job" {
			job := batch.Job{}
			err = yaml.Unmarshal([]byte(yamlData), &job)
			if err != nil {
				log.Fatal(err)
			}
			name := job.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = jobTemplate(job)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "DaemonSet" {
			daemonset := kext.DaemonSet{}
			err = yaml.Unmarshal([]byte(yamlData), &daemonset)
			if err != nil {
				log.Fatal(err)
			}
			name := daemonset.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = daemonsetTemplate(daemonset)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "ReplicaSet" {
			rcSet := kext.ReplicaSet{}
			err = yaml.Unmarshal([]byte(yamlData), &rcSet)
			if err != nil {
				log.Fatal(err)
			}
			name := rcSet.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = replicaSetTemplate(rcSet)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "StatefulSet" {
			statefulset := apps.StatefulSet{}
			err := yaml.Unmarshal([]byte(yamlData), &statefulset)
			if err != nil {
				log.Fatal(err)
			}
			name := statefulset.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = statefulsetTemplate(statefulset)
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "Service" {
			service := kapi.Service{}
			err := yaml.Unmarshal([]byte(yamlData), &service)
			if err != nil {
				log.Fatal(err)
			}
			template, values = serviceTemplate(service)
			name := service.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			valueFile[removeCharactersFromName(name)] = values.value
			persistence = addPersistence(persistence, values.persistence)
		} else if resourceType.Kind == "ConfigMap" {
			configMap := kapi.ConfigMap{}
			err := yaml.Unmarshal([]byte(yamlData), &configMap)
			if err != nil {
				log.Fatal(err)
			}
			name := configMap.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = configMapTemplate(configMap)
			valueFile[removeCharactersFromName(name)] = values.value
		} else if resourceType.Kind == "Secret" {
			secret := kapi.Secret{}
			err := yaml.Unmarshal([]byte(yamlData), &secret)
			if err != nil {
				log.Fatal(err)
			}
			name := secret.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = secretTemplate(secret)
			valueFile[removeCharactersFromName(name)] = values.value
		} else if resourceType.Kind == "PersistentVolumeClaim" {
			pvc := kapi.PersistentVolumeClaim{}
			err := yaml.Unmarshal([]byte(yamlData), &pvc)
			if err != nil {
				log.Fatal(err)
			}
			name := pvc.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = pvcTemplate(pvc)
			persistence = addPersistence(persistence, values.persistence)
			//valueFile[removeCharactersFromName(name)] = values.value
		} else if resourceType.Kind == "PersistentVolume" {
			pv := kapi.PersistentVolume{}
			err := yaml.Unmarshal([]byte(yamlData), &pv)
			if err != nil {
				log.Fatal(err)
			}
			name := pv.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = pvTemplate(pv)
			valueFile[removeCharactersFromName(name)] = values.value
		} else if resourceType.Kind == "StorageClass" {
			storageClass := storage.StorageClass{}
			err := yaml.Unmarshal([]byte(yamlData), &storageClass)
			if err != nil {
				log.Fatal(err)
			}
			name := storageClass.Name
			templateName = filepath.Join(templateLocation, name+".yaml")
			template, values = storageClassTemplate(storageClass)
			valueFile[removeCharactersFromName(name)] = values.value

		} else {
			fmt.Printf("NOT IMPLEMENTED. ADD MAUALLY ")
		}
		err = ioutil.WriteFile(templateName, []byte(template), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(persistence) != 0 {
		valueFile["persistence"] = persistence
	}
	valueFileData, err := yaml.Marshal(valueFile)
	if err != nil {
		log.Fatal(err)
	}
	helperDir := filepath.Join(templateLocation, HelpersName)
	err = ioutil.WriteFile(helperDir, []byte(defaultHelpers), 0644) //TODO  change default values
	valueDir := filepath.Join(cdir, ValuesfileName)
	err = ioutil.WriteFile(valueDir, []byte(valueFileData), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CREATE : SUCCESSFULL")

	return cdir, nil
}

func podTemplate(pod kapi.Pod) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(pod.ObjectMeta.Name)
	pod.ObjectMeta = generateObjectMetaTemplate(pod.ObjectMeta, key, value, pod.ObjectMeta.Name)
	//pod.Spec.Containers = generateTemplateForContainer(pod.Spec.Containers, value)
	pod.Spec = generateTemplateForPodSpec(pod.Spec, key, value)
	if len(pod.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(pod.Spec.Volumes, key, value)
		pod.Spec.Volumes = nil
	}
	tempPodByte, err := yaml.Marshal(pod)
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

func replicationControllerTemplate(rc kapi.ReplicationController) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(rc.ObjectMeta.Name)
	rc.ObjectMeta = generateObjectMetaTemplate(rc.ObjectMeta, key, value, rc.ObjectMeta.Name)
	rc.Spec.Template.Spec = generateTemplateForPodSpec(rc.Spec.Template.Spec, key, value)
	if len(rc.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(rc.Spec.Template.Spec.Volumes, key, value)
		value["persistence"] = true
		rc.Spec.Template.Spec.Volumes = nil
	}
	tempRcByte, err := yaml.Marshal(rc)
	if err != nil {
		log.Fatal(err)
	}
	tempRc := removeEmptyFields(string(tempRcByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempRc, volumes)
	} else {
		template = tempRc
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func replicaSetTemplate(replicaSet kext.ReplicaSet) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(replicaSet.ObjectMeta.Name)
	replicaSet.ObjectMeta = generateObjectMetaTemplate(replicaSet.ObjectMeta, key, value, replicaSet.ObjectMeta.Name)
	replicaSet.Spec.Template.Spec = generateTemplateForPodSpec(replicaSet.Spec.Template.Spec, key, value)
	if len(replicaSet.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(replicaSet.Spec.Template.Spec.Volumes, key, value)
		value["persistence"] = true
		replicaSet.Spec.Template.Spec.Volumes = nil
	}
	template := ""
	tempRcSetByte, err := yaml.Marshal(replicaSet)
	if err != nil {
		log.Fatal(err)
	}
	tempReplicaSet := removeEmptyFields(string(tempRcSetByte))
	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempReplicaSet, volumes) // RC and replica_set has volume in same layer
	} else {
		template = tempReplicaSet
	}
	return template, valueFileGenerator{
		value:       value,
		persistence: persistence,
	}
}

func deploymentTemplate(deployment kext.Deployment) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(deployment.ObjectMeta.Name)
	deployment.ObjectMeta = generateObjectMetaTemplate(deployment.ObjectMeta, key, value, deployment.ObjectMeta.Name)
	deployment.Spec.Template.Spec = generateTemplateForPodSpec(deployment.Spec.Template.Spec, key, value)
	if len(deployment.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(deployment.Spec.Template.Spec.Volumes, key, value)
		deployment.Spec.Template.Spec.Volumes = nil
	}
	if len(string(deployment.Spec.Strategy.Type)) != 0 {
		deployment.Spec.Strategy.Type = kext.DeploymentStrategyType(fmt.Sprintf("{{.Values.%sDeploymentStrategy}}", key))
		//generateTemplateForSingleValue(string(deployment.Spec.Strategy.Type), "DeploymentStrategy", value)

		value["DeploymentStrategy"] = deployment.Spec.Strategy.Type //TODO test
	}
	template := ""
	tempDeploymentByte, err := yaml.Marshal(deployment)
	if err != nil {
		log.Fatal(err)
	}
	tempDeployment := removeEmptyFields(string(tempDeploymentByte))

	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempDeployment, volumes)
	} else {
		template = tempDeployment
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func daemonsetTemplate(daemonset kext.DaemonSet) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(daemonset.ObjectMeta.Name)
	daemonset.ObjectMeta = generateObjectMetaTemplate(daemonset.ObjectMeta, key, value, daemonset.ObjectMeta.Name)
	daemonset.Spec.Template.Spec = generateTemplateForPodSpec(daemonset.Spec.Template.Spec, key, value)
	if len(daemonset.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(daemonset.Spec.Template.Spec.Volumes, key, value)
		value["persistence"] = true
		daemonset.Spec.Template.Spec.Volumes = nil
	}
	template := ""
	//valueData, err := yaml.Marshal(value)

	tempDaemonSetByte, err := yaml.Marshal(daemonset)
	if err != nil {
		log.Fatal(err)
	}
	tempDaemonSet := removeEmptyFields(string(tempDaemonSetByte))
	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempDaemonSet, volumes)
	} else {
		template = tempDaemonSet
	}
	if err != nil {
		log.Fatal(err)
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func statefulsetTemplate(statefulset apps.StatefulSet) (string, valueFileGenerator) {
	volumes := ""
	value := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(statefulset.ObjectMeta.Name)
	statefulset.ObjectMeta = generateObjectMetaTemplate(statefulset.ObjectMeta, key, value, statefulset.ObjectMeta.Name)
	if len(statefulset.Spec.ServiceName) != 0 {
		statefulset.Spec.ServiceName = fmt.Sprintf("{{.Values.%s.ServiceName}}", key)
		value["ServiceName"] = statefulset.Spec.ServiceName //generateTemplateForSingleValue(statefulset.Spec.ServiceName, "ServiceName", value)
	}
	statefulset.Spec.Template.Spec = generateTemplateForPodSpec(statefulset.Spec.Template.Spec, key, value)
	if len(statefulset.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(statefulset.Spec.Template.Spec.Volumes, key, value)
		statefulset.Spec.Template.Spec.Volumes = nil
	}
	tempStatefulSetByte, err := yaml.Marshal(statefulset)
	if err != nil {
		log.Fatal(err)
	}
	tempStatefulSet := removeEmptyFields(string(tempStatefulSetByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempStatefulSet, volumes)
	} else {
		template = tempStatefulSet
	}
	return template, valueFileGenerator{value: value, persistence: persistence}
}

func jobTemplate(job batch.Job) (string, valueFileGenerator) {
	volumes := ""
	persistence := make(map[string]interface{}, 0)
	value := make(map[string]interface{}, 0)
	key := removeCharactersFromName(job.ObjectMeta.Name)
	job.ObjectMeta = generateObjectMetaTemplate(job.ObjectMeta, key, value, job.ObjectMeta.Name)
	job.Spec.Template.Spec = generateTemplateForPodSpec(job.Spec.Template.Spec, key, value)
	if len(job.Spec.Template.Spec.Volumes) != 0 {
		volumes, persistence = generateTemplateForVolume(job.Spec.Template.Spec.Volumes, key, value)
		value["persistence"] = true
		job.Spec.Template.Spec.Volumes = nil
	}
	tempJobByte, err := yaml.Marshal(job)
	if err != nil {
		log.Fatal(err)
	}
	tempJob := removeEmptyFields(string(tempJobByte))
	template := ""
	if len(volumes) != 0 {
		template = addVolumeToTemplateForRc(tempJob, volumes)
	} else {
		template = tempJob
	}
	return template, valueFileGenerator{value: value, persistence: persistence}

}

func serviceTemplate(svc kapi.Service) (string, valueFileGenerator) {
	value := make(map[string]interface{}, 0)
	key := removeCharactersFromName(svc.ObjectMeta.Name)
	svc.ObjectMeta = generateObjectMetaTemplate(svc.ObjectMeta, key, value, svc.ObjectMeta.Name)
	svc.Spec = generateServiceSpecTemplate(svc.Spec, key, value)
	svcData, err := yaml.Marshal(svc)
	if err != nil {
		log.Fatal(err)
	}
	service := removeEmptyFields(string(svcData))
	return string(service), valueFileGenerator{value: value}
}

func configMapTemplate(configMap kapi.ConfigMap) (string, valueFileGenerator) {
	value := make(map[string]interface{}, 0)
	key := removeCharactersFromName(configMap.ObjectMeta.Name)
	configMap.ObjectMeta = generateObjectMetaTemplate(configMap.ObjectMeta, key, value, configMap.ObjectMeta.Name)
	configMap.ObjectMeta.Name = key // not using release name befor configmap
	configMapData, err := yaml.Marshal(configMap)
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

func secretTemplate(secret kapi.Secret) (string, valueFileGenerator) {
	value := make(map[string]interface{}, 0)
	secretDataMap := make(map[string]interface{}, 0)
	key := removeCharactersFromName(secret.ObjectMeta.Name)
	secret.ObjectMeta = generateObjectMetaTemplate(secret.ObjectMeta, key, value, secret.ObjectMeta.Name)
	secret.ObjectMeta.Name = key
	if len(secret.Data) != 0 {
		for k, v := range secret.Data {
			value[k] = v
			secretDataMap[k] = (fmt.Sprintf("{{.Values.%s.%s}}", key, k))
		}
	}
	secret.Data = nil
	value["Type"] = secret.Type
	secret.Type = kapi.SecretType(fmt.Sprintf("{{.Values.%s.Type}}", key))
	secretDataByte, err := yaml.Marshal(secret)
	if err != nil {
		log.Fatal(err)
	}
	secretData := removeEmptyFields(string(secretDataByte))
	//dataSecret := make(map[string]interface{}, 0)
	//dataSecret["data"] = secretDataMap
	secretData = addSecretData(secretData, secretDataMap, key)
	return secretData, valueFileGenerator{value: value}
}

func pvcTemplate(pvc kapi.PersistentVolumeClaim) (string, valueFileGenerator) {
	tempValue := make(map[string]interface{}, 0)
	persistence := make(map[string]interface{}, 0)
	key := removeCharactersFromName(pvc.ObjectMeta.Name)
	pvc.ObjectMeta = generateObjectMetaTemplate(pvc.ObjectMeta, key, tempValue, pvc.ObjectMeta.Name)
	pvc.Spec = generatePersistentVolumeClaimSpec(pvc.Spec, key, tempValue)
	pvcData, err := yaml.Marshal(pvc)
	if err != nil {
		log.Fatal(err)
	}
	temp := removeEmptyFields(string(pvcData))
	pvcTemplateData := fmt.Sprintf("{{- if .Values.persistence.%s.enabled -}}\n%s{{- end -}}", key, temp)
	tempValue["enabled"] = true // By Default use persistence volume true
	persistence[key] = tempValue
	return pvcTemplateData, valueFileGenerator{persistence: persistence}
}

func pvTemplate(pv kapi.PersistentVolume) (string, valueFileGenerator) {
	value := make(map[string]interface{}, 0)
	key := removeCharactersFromName(pv.ObjectMeta.Name)
	pv.ObjectMeta = generateObjectMetaTemplate(pv.ObjectMeta, key, value, pv.Name)
	pv.Spec = generatePersistentVolumeSpec(pv.Spec, key, value)
	pvData, err := yaml.Marshal(pv)
	if err != nil {
		log.Fatal(err)
	}
	temp := removeEmptyFields(string(pvData))
	return string(temp), valueFileGenerator{value: value}
}

func storageClassTemplate(storageClass storage.StorageClass) (string, valueFileGenerator) {
	value := make(map[string]interface{}, 0)
	key := removeCharactersFromName(storageClass.ObjectMeta.Name)
	storageClass.ObjectMeta = generateObjectMetaTemplate(storageClass.ObjectMeta, key, value, storageClass.ObjectMeta.Name)
	value["Provisioner"] = storageClass.Provisioner
	storageClass.Provisioner = fmt.Sprintf("{{.Values.%s.Provisioner}}", key)
	storageClass.Parameters = mapToValueMaker(storageClass.Parameters, value, key)
	storageData, err := yaml.Marshal(storageClass)
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
		ifCondition := fmt.Sprintf("{{ if .Values.%s.%s }}", key, k)
		data += fmt.Sprintf("  %s\n  %s: %s\n  %s\n  %s: %s\n  %s\n", ifCondition, k, v, elseCondition, k, elseAction, end)
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
		Description: "Helm chart generated by https://github.com/appscode/chartify",
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
