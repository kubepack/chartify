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
		yamlFiles = append(yamlFiles, serviceFiles)

	}
	if len(kubeObjects.services) != 0 {
		serviceFiles := kubeObjects.getServicesYamlList(kubeClient)
		yamlFiles = append(yamlFiles, serviceFiles)
	}
	if len(kubeObjects.replicationControllers) != 0 {
		rcFiles := kubeObjects.getReplicationControllerYamlList(kubeClient)
		yamlFiles = append(yamlFiles, rcFiles)
	}
	if len(kubeObjects.secrets) != 0 {
		secretFiles := kubeObjects.getSecretsYamlList(kubeClient)
		yamlFiles = append(yamlFiles, secretFiles)
	}
	if len(kubeObjects.configMaps) != 0 {
		configMapsFiles := kubeObjects.getConfigMapsYamlList(kubeClient)
		yamlFiles = append(yamlFiles, configMapsFiles)
	}
	if len(kubeObjects.petsets) != 0 {
		petSetsFiles := kubeObjects.getPetSetsYamlList(kubeClient)
		yamlFiles = append(yamlFiles, petSetsFiles)
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
	for k, v := range kubeObjects.services {
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
	for k, v := range kubeObjects.secrets {
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
	for k, v := range kubeObjects.secrets {
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
