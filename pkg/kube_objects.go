package pkg

import (
	"github.com/ghodss/yaml"
	kubeapi "k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"strings"
)

func (kubeObjects objects) readKubernetesObjects(kubeClient *client.Client) []string {
	var yamlFiles []string
	if len(kubeObjects.pods) != 0 {
		podFiles := kubeObjects.getPodsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, podFiles)
	}
	if len(kubeObjects.services) != 0 {
		serviceFiles := kubeObjects.getServicesYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, serviceFiles)
	}
	if len(kubeObjects.replicationControllers) != 0 {
		rcFiles := kubeObjects.getReplicationControllerYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, rcFiles)
	}
	if len(kubeObjects.secrets) != 0 {
		secretFiles := kubeObjects.getSecretsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, secretFiles)
	}
	if len(kubeObjects.configMaps) != 0 {
		configMapsFiles := kubeObjects.getConfigMapsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, configMapsFiles)
	}
	if len(kubeObjects.petsets) != 0 {
		petSetsFiles := kubeObjects.getPetSetsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, petSetsFiles)
	}
	if len(kubeObjects.persistentVolume) != 0 {
		pvFiles := kubeObjects.getPersistentVolumeYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, pvFiles)
	}
	if len(kubeObjects.persistentVolumeClaim) != 0 {
		pvcFiles := kubeObjects.getPersistentVolumeClaimYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, pvcFiles)
	}
	/*	if len(kubeObjects.jobs) != 0 {  //TODO sauman
		jobFiles := kubeObjects.getJobsyamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, jobFiles)
	}*/
	if len(kubeObjects.daemons) != 0 {
		daemonFiles := kubeObjects.getDaemonsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, daemonFiles)
	}

	return yamlFiles
}

func (kubeObjects objects) getPodsYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.pods {
		objectName, namespace := splitnamespace(v)
		pod, err := kubeClient.Pods(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(pod)
		if err != nil {
			log.Fatal(err)
		}
		if pod.Kind == "" {
			pod.Kind = ref.Kind
		}
		if pod.APIVersion == "" {
			pod.APIVersion = ref.APIVersion
		}
		pod.Status = kubeapi.PodStatus{} //TODO check
		dataByte, err := yaml.Marshal(pod)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getReplicationControllerYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.replicationControllers {
		objectName, namespace := splitnamespace(v)
		rc, err := kubeClient.ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(rc)
		if err != nil {
			log.Fatal(err)
		}
		if rc.Kind == "" {
			rc.Kind = ref.Kind
		}
		if rc.APIVersion == "" {
			rc.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(rc)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getServicesYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.services {
		objectName, namespace := splitnamespace(v)
		service, err := kubeClient.Services(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(service)
		if service.Kind == "" {
			service.Kind = ref.Kind
		}
		if service.APIVersion == "" {
			service.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(service)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getSecretsYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.secrets {
		objectName, namespace := splitnamespace(v)
		secret, err := kubeClient.Secrets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(secret)
		if secret.Kind == "" {
			secret.Kind = ref.Kind
		}
		if secret.APIVersion == "" {
			secret.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(secret)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getConfigMapsYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.configMaps {
		objectName, namespace := splitnamespace(v)
		configmap, err := kubeClient.ConfigMaps(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(configmap)
		if configmap.Kind == "" {
			configmap.Kind = ref.Kind
		}
		if configmap.APIVersion == "" {
			configmap.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(configmap)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getPetSetsYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.petsets {
		objectName, namespace := splitnamespace(v)
		petset, err := kubeClient.PetSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(petset)
		if petset.Kind == "" {
			petset.Kind = ref.Kind
		}
		if petset.APIVersion == "" {
			petset.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(petset)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getPersistentVolumeYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.persistentVolume {
		pv, err := kubeClient.PersistentVolumes().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(pv)
		if pv.Kind == "" {
			pv.Kind = ref.Kind
		}
		if pv.APIVersion == "" {
			pv.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(pv)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getPersistentVolumeClaimYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.persistentVolumeClaim {
		objectName, namespace := splitnamespace(v)
		pvc, err := kubeClient.PersistentVolumeClaims(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(pvc)
		if pvc.Kind == "" {
			pvc.Kind = ref.Kind
		}
		if pvc.APIVersion == "" {
			pvc.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(pvc)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

/*
func (kubeObjects objects) getJobsyamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for k, v := range kubeObjects.jobs {
		job, err := kubeClient.Jobs(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(job)
		yamlFiles[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}*/

func (kubeObjects objects) getDaemonsYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for _, v := range kubeObjects.daemons {
		objectName, namespace := splitnamespace(v)
		daemon, err := kubeClient.ExtensionsClient.DaemonSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := kubeapi.GetReference(daemon)
		if daemon.Kind == "" {
			daemon.Kind = ref.Kind
		}
		if daemon.APIVersion == "" {
			daemon.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(daemon)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func appendSlice(mainSlice []string, subSlice []string) []string {
	for _, v := range subSlice {
		mainSlice = append(mainSlice, v)
	}
	return mainSlice
}

func splitnamespace(s string) (string, string) {
	str := strings.Split(s, ".")
	if len(str) == 2 {
		return str[0], str[1]
	}else if len(str) < 2 {
		log.Fatal("Namespace not given")
	}else {
		log.Fatal("Can not detect Namespace")
	}
	return "",""
}