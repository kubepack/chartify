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
	ConfigMaps             []string
	Deployments            []string
	Daemons                []string
	Jobs                   []string
	PersistentVolumes      []string
	PersistentVolumeClaims []string
	Pods                   []string
	ReplicaSets            []string
	ReplicationControllers []string
	Secrets                []string
	Services               []string
	StatefulSets           []string
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
		yamlFiles = appendSlice(yamlFiles, ko.getPods(kubeClient))
	}
	if len(ko.Services) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getServices(kubeClient))
	}
	if len(ko.ReplicationControllers) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getReplicationControllers(kubeClient))
	}
	if len(ko.Secrets) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getSecrets(kubeClient))
	}
	if len(ko.ConfigMaps) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getConfigMaps(kubeClient))
	}
	if len(ko.StatefulSets) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getStatefulSets(kubeClient))
	}
	if len(ko.PersistentVolumes) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getPersistentVolumes(kubeClient))
	}
	if len(ko.PersistentVolumeClaims) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getPersistentVolumeClaims(kubeClient))
	}
	if len(ko.Jobs) != 0 { //TODO sauman
		yamlFiles = appendSlice(yamlFiles, ko.getJobs(kubeClient))
	}
	if len(ko.Daemons) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getDaemons(kubeClient))
	}
	if len(ko.Deployments) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getDeployments(kubeClient))
	}
	if len(ko.ReplicaSets) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getReplicaSets(kubeClient))
	}
	if len(ko.StorageClasses) != 0 {
		yamlFiles = appendSlice(yamlFiles, ko.getStorageClasses(kubeClient))
	}
	return yamlFiles
}

func (ko KubeObjects) getPods(kubeClient clientset.Interface) []string {
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
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getReplicationControllers(kubeClient clientset.Interface) []string {
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
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getServices(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Services {
		objectName, namespace := splitNamespace(v)
		service, err := kubeClient.Core().Services(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(service)
		if err != nil {
			log.Fatal(err)
		}
		if service.Kind == "" {
			service.Kind = ref.Kind
		}
		if service.APIVersion == "" {
			service.APIVersion = ref.APIVersion
		}
		service.Status = api.ServiceStatus{}
		dataByte, err := yaml.Marshal(service)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getSecrets(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.Secrets {
		objectName, namespace := splitNamespace(v)
		secret, err := kubeClient.Core().Secrets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(secret)
		if err != nil {
			log.Fatal(err)
		}
		if secret.Kind == "" {
			secret.Kind = ref.Kind
		}
		if secret.APIVersion == "" {
			secret.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(secret)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getConfigMaps(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.ConfigMaps {
		objectName, namespace := splitNamespace(v)
		configmap, err := kubeClient.Core().ConfigMaps(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(configmap)
		if err != nil {
			log.Fatal(err)
		}
		if configmap.Kind == "" {
			configmap.Kind = ref.Kind
		}
		if configmap.APIVersion == "" {
			configmap.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(configmap)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getStatefulSets(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.StatefulSets {
		objectName, namespace := splitNamespace(v)
		statefulset, err := kubeClient.Apps().StatefulSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(statefulset)
		if err != nil {
			log.Fatal(err)
		}
		if statefulset.Kind == "" {
			statefulset.Kind = ref.Kind
		}
		if len(statefulset.APIVersion) == 0 {
			statefulset.APIVersion = makeAPIVersion(statefulset.GetSelfLink())
		}
		statefulset.Status = apps.StatefulSetStatus{}
		dataByte, err := yaml.Marshal(statefulset)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getPersistentVolumes(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.PersistentVolumes {
		pv, err := kubeClient.Core().PersistentVolumes().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(pv)
		if err != nil {
			log.Fatal(err)
		}
		if pv.Kind == "" {
			pv.Kind = ref.Kind
		}
		if pv.APIVersion == "" {
			pv.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(pv)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getPersistentVolumeClaims(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.PersistentVolumeClaims {
		objectName, namespace := splitNamespace(v)
		pvc, err := kubeClient.Core().PersistentVolumeClaims(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(pvc)
		if err != nil {
			log.Fatal(err)
		}
		if pvc.Kind == "" {
			pvc.Kind = ref.Kind
		}
		if pvc.APIVersion == "" {
			pvc.APIVersion = ref.APIVersion
		}
		dataByte, err := yaml.Marshal(pvc)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getJobs(kubeClient clientset.Interface) []string {
	var jobFiles []string
	for _, v := range ko.Jobs {
		objectName, namespace := splitNamespace(v)
		job, err := kubeClient.Batch().Jobs(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(job)
		if err != nil {
			log.Fatal(err)
		}
		if job.Kind == "" {
			job.Kind = ref.Kind
		}
		if job.APIVersion == "" {
			job.APIVersion = makeAPIVersion(job.GetSelfLink())
		}
		job.Status = batch.JobStatus{}
		dataByte, err := yaml.Marshal(job)
		if err != nil {
			log.Fatal(err)
		}
		jobFiles = append(jobFiles, string(dataByte))
	}
	return jobFiles
}

func (ko KubeObjects) getDaemons(kubeClient clientset.Interface) []string {
	var daemonFiles []string
	for _, v := range ko.Daemons {
		objectName, namespace := splitNamespace(v)
		daemon, err := kubeClient.Extensions().DaemonSets(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(daemon)
		if err != nil {
			log.Fatal(err)
		}
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

func (ko KubeObjects) getDeployments(kubeClient clientset.Interface) []string {
	var files []string
	for _, v := range ko.Deployments {
		objectName, namespace := splitNamespace(v)
		deployment, err := kubeClient.Extensions().Deployments(namespace).Get(objectName)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(deployment)
		if err != nil {
			log.Fatal(err)
		}
		if deployment.Kind == "" {
			deployment.Kind = ref.Kind
		}
		if deployment.APIVersion == "" {
			deployment.APIVersion = makeAPIVersion(deployment.GetSelfLink())
		}
		deployment.Status = extensions.DeploymentStatus{}
		dataByte, err := yaml.Marshal(deployment)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, string(dataByte))

	}
	return files
}

func (ko KubeObjects) getReplicaSets(kubeClient clientset.Interface) []string {
	var yamlFiles []string
	for _, v := range ko.ReplicaSets {
		objectName, namespace := splitNamespace(v)
		rc, err := kubeClient.Extensions().ReplicaSets(namespace).Get(objectName)
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
		rc.Status = extensions.ReplicaSetStatus{}
		dataByte, err := yaml.Marshal(rc)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (ko KubeObjects) getStorageClasses(kubeClient clientset.Interface) []string {
	var storageFiles []string
	for _, v := range ko.StorageClasses {
		//objectsName, namespace := splitnamespace(v)
		storageClass, err := kubeClient.Storage().StorageClasses().Get(v)
		if err != nil {
			log.Fatal(err)
		}
		ref, err := api.GetReference(storageClass)
		if err != nil {
			log.Fatal(err)
		}
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
	str := strings.Split(s, "@")
	if len(str) == 2 {
		return str[0], str[1]
	} else if len(str) == 1 {
		return str[0], api.NamespaceDefault
	}
	log.Fatal("ERROR : Can not detect Namespace")
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
