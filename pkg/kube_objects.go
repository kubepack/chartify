package pkg

import (
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"

	"github.com/ghodss/yaml"
)

func (kubeObjects objects) readKubernetesObjects(kubeClient *client.Client) []string {
	var yamlFiles []string
	if len(kubeObjects.pods) != 0 {
		serviceFiles := kubeObjects.getPodsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, serviceFiles)

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
/*	if len(kubeObjects.jobs) != 0 {
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
	var data []string
	for k, v := range kubeObjects.pods {
		pod, err := kubeClient.Pods(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(pod)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getReplicationControllerYamlList(kubeClient *client.Client) []string {
	var data []string
	for k, v := range kubeObjects.replicationControllers {
		service, err := kubeClient.ReplicationControllers(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(service)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getServicesYamlList(kubeClient *client.Client) []string {
	var data []string
	for k, v := range kubeObjects.services {
		service, err := kubeClient.Services(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(service)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getSecretsYamlList(kubeClient *client.Client) []string {
	var data []string
	for k, v := range kubeObjects.secrets {
		secret, err := kubeClient.Secrets(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(secret)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getConfigMapsYamlList(kubeClient *client.Client) []string {
	var data []string
	for k, v := range kubeObjects.configMaps {
		configmap, err := kubeClient.ConfigMaps(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(configmap)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getPetSetsYamlList(kubeClient *client.Client) []string {
	var data []string
	for k, v := range kubeObjects.petsets {
		petset, err := kubeClient.PetSets(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(petset)
		data[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func (kubeObjects objects) getPersistentVolumeYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for k, v := range kubeObjects.persistentVolume {
		pv, err := kubeClient.PersistentVolumes().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(pv)
		yamlFiles[k] = string(dataByte)
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (kubeObjects objects) getPersistentVolumeClaimYamlList(kubeClient *client.Client) []string {
	var yamlFiles []string
	for k, v := range kubeObjects.persistentVolumeClaim {
		pvc, err := kubeClient.PersistentVolumeClaims(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(pvc)
		yamlFiles[k] = string(dataByte)
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
	for k, v := range kubeObjects.daemons {
		daemon, err := kubeClient.ExtensionsClient.DaemonSets(kubeObjects.namespace).Get(v)
		if err != nil {
			log.Fatal(err)
		}
		dataByte, err := yaml.Marshal(daemon)
		yamlFiles[k] = string(dataByte)
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
