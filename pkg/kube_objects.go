package pkg

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

type KubeObjects struct {
	Pods                   []string
	ReplicationControllers []string
	ConfigMaps             []string
	Services               []string
	Secrets                []string
	PersistentVolume       []string
	PersistentVolumeClaim  []string
	Statefulsets           []string
	Jobs                   []string
	Daemons                []string
	ReplicaSet             []string
	StorageClasses         []string
}

func (ko KubeObjects) Extract() []string {
	kubeClient, err := newKubeClient()
	if err != nil {
		log.Fatal(err)
	}
	yamlFiles := ko.readKubernetesObjects(kubeClient)
	return yamlFiles
}

func (ko KubeObjects) CheckFlags() bool {
	v := reflect.ValueOf(ko)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Len() > 0 {
			return true
		}
	}
	return false
}

func (ko KubeObjects) readKubernetesObjects(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	if len(ko.Pods) != 0 {
		podFiles := ko.getPodsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, podFiles)
	}
	if len(ko.Services) != 0 {
		serviceFiles := ko.getServicesYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, serviceFiles)
	}
	if len(ko.ReplicationControllers) != 0 {
		rcFiles := ko.getReplicationControllerYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, rcFiles)
	}
	if len(ko.Secrets) != 0 {
		secretFiles := ko.getSecretsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, secretFiles)
	}
	if len(ko.ConfigMaps) != 0 {
		configMapsFiles := ko.getConfigMapsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, configMapsFiles)
	}
	if len(ko.Statefulsets) != 0 {
		statefulSetFiles := ko.getStatefulSetsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, statefulSetFiles)
	}
	if len(ko.PersistentVolume) != 0 {
		pvFiles := ko.getPersistentVolumeYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, pvFiles)
	}
	if len(ko.PersistentVolumeClaim) != 0 {
		pvcFiles := ko.getPersistentVolumeClaimYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, pvcFiles)
	}
	if len(ko.Jobs) != 0 { //TODO sauman
		jobFiles := ko.getJobsyamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, jobFiles)
	}
	if len(ko.Daemons) != 0 {
		daemonFiles := ko.getDaemonsYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, daemonFiles)
	}
	if len(ko.StorageClasses) != 0 {
		storageClassFiles := ko.getStorageClassYamlList(kubeClient)
		yamlFiles = appendSlice(yamlFiles, storageClassFiles)
	}

	return yamlFiles
}

func (ko KubeObjects) getPodsYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Pods {
		objectName, namespace := splitNamespace(v)
		pod, err := kubeClient.Core().Pods(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(pod)
		if err != nil {
			log.Fatal(err)
		}
		if pod.Kind == "" {
			pod.Kind = ref.Kind
		}
		if pod.APIVersion == "" {
			pod.APIVersion = ref.APIVersion
		}
		pod.Status = api.PodStatus{}
		dataByte, err := yaml.Marshal(pod)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (ko KubeObjects) getReplicationControllerYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.ReplicationControllers {
		objectName, namespace := splitNamespace(v)
		rc, err := kubeClient.Core().ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(rc)
		if err != nil {
			log.Fatal(err)
		}
		if rc.Kind == "" {
			rc.Kind = ref.Kind
		}
		if rc.APIVersion == "" {
			rc.APIVersion = ref.APIVersion
		}
		rc.Status = api.ReplicationControllerStatus{}
		dataByte, err := yaml.Marshal(rc)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (ko KubeObjects) getServicesYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Services {
		objectName, namespace := splitNamespace(v)
		service, err := kubeClient.Core().Services(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(service)
		if service.Kind == "" {
			service.Kind = ref.Kind
		}
		if service.APIVersion == "" {
			service.APIVersion = ref.APIVersion
		}
		service.Status = api.ServiceStatus{}
		dataByte, err := yaml.Marshal(service)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (ko KubeObjects) getSecretsYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Secrets {
		objectName, namespace := splitNamespace(v)
		secret, err := kubeClient.Core().Secrets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(secret)
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

func (ko KubeObjects) getConfigMapsYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.ConfigMaps {
		objectName, namespace := splitNamespace(v)
		configmap, err := kubeClient.Core().ConfigMaps(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(configmap)
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

func (ko KubeObjects) getStatefulSetsYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Statefulsets {
		objectName, namespace := splitNamespace(v)
		statefulset, err := kubeClient.Apps().StatefulSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(statefulset)
		if statefulset.Kind == "" {
			statefulset.Kind = ref.Kind
		}
		if len(statefulset.APIVersion) == 0 {
			statefulset.APIVersion = makeAPIVersion(statefulset.GetSelfLink())
		}
		statefulset.Status = apps.StatefulSetStatus{}
		dataByte, err := yaml.Marshal(statefulset)
		yamlFiles = append(yamlFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return yamlFiles
}

func (ko KubeObjects) getPersistentVolumeYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.PersistentVolume {
		pv, err := kubeClient.Core().PersistentVolumes().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(pv)
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

func (ko KubeObjects) getPersistentVolumeClaimYamlList(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.PersistentVolumeClaim {
		objectName, namespace := splitNamespace(v)
		pvc, err := kubeClient.Core().PersistentVolumeClaims(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(pvc)
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

func (ko KubeObjects) getJobsyamlList(kubeClient clientset.Interface) []string {
	var jobFiles []string
	for _, v := range ko.Jobs {
		objectName, namespace := splitNamespace(v)
		job, err := kubeClient.Batch().Jobs(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(job)
		if job.Kind == "" {
			job.Kind = ref.Kind
		}
		if job.APIVersion == "" {
			job.APIVersion = makeAPIVersion(job.GetSelfLink())
		}
		job.Status = batch.JobStatus{}
		dataByte, err := yaml.Marshal(job)
		jobFiles = append(jobFiles, string(dataByte))
		if err != nil {
			log.Fatal(err)
		}
	}
	return jobFiles
}

func (ko KubeObjects) getDaemonsYamlList(kubeClient clientset.Interface) []string {
	var daemonFiles []string
	for _, v := range ko.Daemons {
		objectName, namespace := splitNamespace(v)
		daemon, err := kubeClient.Extensions().DaemonSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(daemon)
		if daemon.Kind == "" {
			daemon.Kind = ref.Kind
		}
		if daemon.APIVersion == "" {
			daemon.APIVersion = makeAPIVersion(daemon.GetSelfLink())
		}
		daemon.Status = extensions.DaemonSetStatus{}
		dataByte, err := yaml.Marshal(daemon)
		if err != nil {
			log.Fatal(err)
		}
		daemonFiles = append(daemonFiles, string(dataByte))

	}
	return daemonFiles
}

func (ko KubeObjects) getStorageClassYamlList(kubeClient clientset.Interface) []string {
	var storageFiles []string
	for _, v := range ko.StorageClasses {
		//objectsName, namespace := splitnamespace(v)
		storageClass, err := kubeClient.Storage().StorageClasses().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(storageClass)
		if storageClass.Kind == "" {
			storageClass.Kind = ref.Kind
		}
		if storageClass.APIVersion == "" {
			storageClass.APIVersion = makeAPIVersion(storageClass.GetSelfLink())
		}
		dataByte, err := yaml.Marshal(storageClass)
		if err != nil {
			log.Fatal(err)
		}
		storageFiles = append(storageFiles, string(dataByte))

	}

	return storageFiles

}

func newKubeClient() (clientset.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not get kubernetes config: %s", err)
	}
	return clientset.NewForConfig(config)
}

func appendSlice(mainSlice []string, subSlice []string) []string {
	for _, v := range subSlice {
		mainSlice = append(mainSlice, v)
	}
	return mainSlice
}

func splitNamespace(s string) (string, string) {
	str := strings.Split(s, ".")
	if len(str) == 2 {
		return str[0], str[1]
	} else if len(str) == 1 {
		return str[0], "default"
	} else {
		log.Fatal("ERROR : Can not detect Namespace")
	}
	return "", ""
}

func makeAPIVersion(selfLink string) string {
	str := strings.Split(selfLink, "/")
	if len(str) > 2 {
		return (str[2] + "/" + str[3])
	}
	log.Fatal("api version not found")
	return ""
}
