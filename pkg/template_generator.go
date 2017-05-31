package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func generateObjectMetaTemplate(objectMeta v1.ObjectMeta, key string, value map[string]interface{}, extraTagForName string) v1.ObjectMeta {
	if len(objectMeta.Name) != 0 {
		objectMeta.Name = fmt.Sprintf(`{{ template "fullname" . }}`)
	}
	if len(extraTagForName) != 0 {
		objectMeta.Name = fmt.Sprintf("%s-%s", objectMeta.Name, extraTagForName)
	}
	if len(objectMeta.ClusterName) != 0 {
		value[ClusterName] = objectMeta.ClusterName
		objectMeta.ClusterName = fmt.Sprintf("{{.Values.%s.%s}}", key, ClusterName)
	}
	if len(objectMeta.GenerateName) != 0 {
		value[GenerateName] = objectMeta.GenerateName
		objectMeta.GenerateName = fmt.Sprintf("{{.Values.%s.%s}}", key, GenerateName)
	}
	if len(objectMeta.Namespace) != 0 {
		value[Namespace] = objectMeta.Namespace
		objectMeta.Namespace = fmt.Sprintf("{{.Values.%s.%s}}", key, Namespace)
	}
	objectMeta.Labels = generateTemplateForLables(objectMeta.Labels)
	return objectMeta
}

func generateTemplateReplicationCtrSpec(rcSpec apiv1.ReplicationControllerSpec, rcSpecStr string, key string, value map[string]interface{}) (string, map[string]interface{}) {
	templateDeployment := rcSpecStr

	templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "replicas")
	value["replicas"] = rcSpec.Replicas

	if rcSpec.MinReadySeconds != 0 {
		templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "minReadySeconds")
		value["minReadySeconds"] = rcSpec.MinReadySeconds
	}

	return templateDeployment, value
}

func generateTemplateReplicaSetSpec(rsSpec extensions.ReplicaSetSpec, rsSpecStr string, key string, value map[string]interface{}) (string, map[string]interface{}) {
	templateDeployment := rsSpecStr

	templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "replicas")
	value["replicas"] = rsSpec.Replicas

	if rsSpec.MinReadySeconds != 0 {
		templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "minReadySeconds")
		value["minReadySeconds"] = rsSpec.MinReadySeconds
	}

	return templateDeployment, value
}

func generateTemplateDeplymentSpec(dcSpec extensions.DeploymentSpec, dcSpecStr string, key string, value map[string]interface{}) (string, map[string]interface{}) {
	templateDeployment := dcSpecStr
	templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "replicas")
	value["replicas"] = dcSpec.Replicas

	if dcSpec.MinReadySeconds != 0 {
		templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "minReadySeconds")
		value["minReadySeconds"] = dcSpec.MinReadySeconds
	}

	if dcSpec.RevisionHistoryLimit != nil {
		templateDeployment = updateIntParamAsStringInTemplate(templateDeployment, key, "revisionHistoryLimit")
		value["revisionHistoryLimit"] = dcSpec.RevisionHistoryLimit
	}

	return templateDeployment, value
}

func updateIntParamAsStringInTemplate(spec string, key string, replace string) string {
	str := strings.Split(spec, "\n")
	var tpl bytes.Buffer
	for _, l := range str {
		if len(l) == 0 {
			continue
		}
		if strings.Contains(l, replace+": ") {
			str1 := fmt.Sprintf("{{.Values.%s.%s}}", key, replace)
			regex := regexp.MustCompile(".*" + replace + ":")
			extrStr := regex.FindString(l)
			tpl.WriteString(extrStr)
			tpl.WriteString(" ")
			tpl.WriteString(str1)
		} else {
			tpl.WriteString(l)
		}
		tpl.WriteRune('\n')
	}
	return tpl.String()
}

func generateTemplateForPodSpec(podSpec apiv1.PodSpec, key string, value map[string]interface{}) apiv1.PodSpec {
	podSpec.Containers = generateTemplateForContainer(podSpec.Containers, key, value)
	if len(podSpec.Hostname) != 0 {
		value[HostName] = podSpec.Hostname
		podSpec.Hostname = fmt.Sprintf("{{.Values.%s.%s}}", key, HostName)
	}
	if len(podSpec.Subdomain) != 0 {
		value[Subdomain] = podSpec.Subdomain
		podSpec.Subdomain = fmt.Sprintf("{{.Values.%s.%s}}", key, Subdomain)
	}
	if len(podSpec.NodeName) != 0 {
		value[Nodename] = podSpec.NodeName
		podSpec.NodeName = fmt.Sprintf("{{.Values.%s.%s}}", key, Nodename)
	}
	if len(podSpec.ServiceAccountName) != 0 {
		value[ServiceAccountName] = podSpec.ServiceAccountName
		podSpec.ServiceAccountName = fmt.Sprintf("{{.Values.%s.%s}}", key, ServiceAccountName)
	}
	if len(string(podSpec.RestartPolicy)) != 0 {
		value[RestartPolicy] = string(podSpec.RestartPolicy)
		podSpec.RestartPolicy = apiv1.RestartPolicy(fmt.Sprintf("{{.Values.%s.%s}}", key, RestartPolicy))
	}

	if len(podSpec.ImagePullSecrets) != 0 {
		imagePullSecretsObj := []apiv1.LocalObjectReference{}

		for _, imagePullSecrets := range podSpec.ImagePullSecrets {
			if checkIfNameExist(imagePullSecrets.Name, "Secret") {
				value["imagePullSecrets"] = fmt.Sprintf(`{{ template "fullname" . }}-%s`, imagePullSecrets.Name)
			} else {
				value["imagePullSecrets"] = imagePullSecrets.Name
			}
			secName := fmt.Sprintf("{{.Values.%s.%s}}", key, "imagePullSecrets")
			imagePullSecretsObj = append(imagePullSecretsObj, apiv1.LocalObjectReference{Name: secName})
		}

		podSpec.ImagePullSecrets = imagePullSecretsObj

	}

	return podSpec
}

func generateTemplateForVolume(volumes []apiv1.Volume, key string, value map[string]interface{}) (string, map[string]interface{}) {
	volumeTemplate := ""
	ifCondition := ""
	partialvolumeTemplate := ""
	persistence := make(map[string]interface{}, 0)
	for _, volume := range volumes {
		ifCondition = ""
		volumeMap := make(map[string]interface{}, 0)
		volumeMap[Enabled] = true
		vol := []apiv1.Volume{}
		vol = append(vol, volume)
		if volume.PersistentVolumeClaim != nil {
			ifCondition = buildIfConditionForVolume(volume.PersistentVolumeClaim.ClaimName)
			if checkIfNameExist(volume.PersistentVolumeClaim.ClaimName, "PersistentVolumeClaim") {
				volume.PersistentVolumeClaim.ClaimName = fmt.Sprintf(`{{template "fullname"}}-%s`, volume.PersistentVolumeClaim.ClaimName)
			}
		} else if volume.ConfigMap != nil {
			if checkIfNameExist(volume.ConfigMap.Name, "Configmap") {
				volume.ConfigMap.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, volume.ConfigMap.Name)
			}
		} else if volume.Secret != nil {
			if checkIfNameExist(volume.Secret.SecretName, "Secret") {
				volume.Secret.SecretName = fmt.Sprintf(`{{ template "fullname" . }}-%s`, volume.Secret.SecretName)
			} //TODO add items
		} else if volume.Glusterfs != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[Path] = volume.Glusterfs.Path
			volumeMap[EndpointsName] = volume.Glusterfs.EndpointsName
			volume.Glusterfs.EndpointsName = VolumeTemplateForElement(volume.Name, EndpointsName)
			volume.Glusterfs.Path = VolumeTemplateForElement(volume.Name, Path)
			persistence[volume.Name] = volumeMap
		} else if volume.HostPath != nil {
			volumeMap[Path] = volume.HostPath.Path
			volume.HostPath.Path = VolumeTemplateForElement(volume.Name, Path)
			persistence[volume.Name] = volumeMap
		} else if volume.GCEPersistentDisk != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[PDName] = volume.GCEPersistentDisk.PDName
			volumeMap[FSType] = volume.GCEPersistentDisk.FSType
			volume.GCEPersistentDisk.PDName = VolumeTemplateForElement(volume.Name, PDName)
			volume.GCEPersistentDisk.FSType = VolumeTemplateForElement(volume.Name, FSType)
			persistence[volume.Name] = volumeMap
		} else if volume.AWSElasticBlockStore != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[FSType] = volume.GCEPersistentDisk.FSType
			volumeMap[VolumeID] = volume.AWSElasticBlockStore.VolumeID
			volume.AWSElasticBlockStore.VolumeID = VolumeTemplateForElement(volume.Name, VolumeID)
			volume.AWSElasticBlockStore.FSType = VolumeTemplateForElement(volume.Name, FSType)
		} else if volume.GitRepo != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[Repository] = volume.GitRepo.Repository
			volumeMap[Revision] = volume.GitRepo.Revision
			volumeMap[Directory] = volume.GitRepo.Directory
			volume.GitRepo.Revision = VolumeTemplateForElement(volume.Name, Revision)
			volume.GitRepo.Repository = VolumeTemplateForElement(volume.Name, Repository)
			volume.GitRepo.Directory = VolumeTemplateForElement(volume.Name, Directory)
			persistence[volume.Name] = volumeMap
		} else if volume.NFS != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[Server] = volume.NFS.Server
			volumeMap[Path] = volume.NFS.Path
			volume.NFS.Path = fmt.Sprintf(`{{.Values.%s.%s}}`, volume.Name, Path)
			volume.NFS.Server = fmt.Sprintf(`{{.Values.%s.%s}}`, volume.Name, Server)
			persistence[volume.Name] = volumeMap
		} else if volume.ISCSI != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[TargetPortal] = volume.ISCSI.TargetPortal
			volumeMap[IQN] = volume.ISCSI.IQN
			volumeMap[ISCSIInterface] = volume.ISCSI.ISCSIInterface
			volumeMap[FSType] = volume.ISCSI.FSType
			volume.ISCSI.TargetPortal = VolumeTemplateForElement(volume.Name, TargetPortal)
			volume.ISCSI.IQN = VolumeTemplateForElement(volume.Name, IQN)
			volume.ISCSI.FSType = fmt.Sprintf(`{{.Values.%s.%s}}`, volume.Name, FSType)
			volume.ISCSI.ISCSIInterface = VolumeTemplateForElement(volume.Name, ISCSIInterface)
			persistence[volume.Name] = volumeMap
		} else if volume.RBD != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[FSType] = volume.RBD.FSType
			volumeMap[RBDImage] = volume.RBD.RBDImage
			volumeMap[RBDPool] = volume.RBD.RBDPool
			volumeMap[RadosUser] = volume.RBD.RadosUser
			volumeMap[Keyring] = volume.RBD.Keyring
			volume.RBD.FSType = VolumeTemplateForElement(volume.Name, FSType)
			volume.RBD.RBDImage = VolumeTemplateForElement(volume.Name, RBDImage)
			volume.RBD.RBDPool = VolumeTemplateForElement(volume.Name, RBDPool)
			volume.RBD.RadosUser = VolumeTemplateForElement(volume.Name, RadosUser)
			volume.RBD.Keyring = VolumeTemplateForElement(volume.Name, Keyring)
			persistence[volume.Name] = volumeMap
		} else if volume.Quobyte != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[Registry] = volume.Quobyte.Registry
			volumeMap[Volume] = volume.Quobyte.Volume
			volumeMap[Group] = volume.Quobyte.Group
			volumeMap[User] = volume.Quobyte.User
			volume.Quobyte.Registry = VolumeTemplateForElement(volume.Name, Registry)
			volume.Quobyte.Volume = VolumeTemplateForElement(volume.Name, Volume)
			volume.Quobyte.Group = VolumeTemplateForElement(volume.Name, Group)
			volume.Quobyte.User = VolumeTemplateForElement(volume.Name, User)
			persistence[volume.Name] = volumeMap
		} else if volume.FlexVolume != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Driver"] = volume.FlexVolume.Driver
			volumeMap[FSType] = volume.FlexVolume.FSType
			// TODO secret reference
			volume.FlexVolume.Driver = VolumeTemplateForElement(volume.Name, "Driver")
			volume.FlexVolume.FSType = VolumeTemplateForElement(volume.Name, FSType)
			persistence[volume.Name] = volumeMap
		} else if volume.Cinder != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[FSType] = volume.Cinder.FSType
			volumeMap[VolumeID] = volume.Cinder.VolumeID
			volume.Cinder.FSType = VolumeTemplateForElement(volume.Name, FSType)
			volume.Cinder.VolumeID = VolumeTemplateForElement(volume.Name, VolumeID)
			persistence[volume.Name] = volumeMap
		} else if volume.CephFS != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[Path] = volume.CephFS.Path
			volumeMap[SecretFile] = volume.CephFS.SecretFile
			volumeMap[User] = volume.CephFS.User
			volume.CephFS.Path = VolumeTemplateForElement(volume.Name, Path)
			volume.CephFS.SecretFile = VolumeTemplateForElement(volume.Name, SecretFile)
			volume.CephFS.User = VolumeTemplateForElement(volume.Name, User)
			persistence[volume.Name] = volumeMap
		} else if volume.Flocker != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[DatasetName] = volume.Flocker.DatasetName
			volume.Flocker.DatasetName = VolumeTemplateForElement(volume.Name, DatasetName)
			persistence[volume.Name] = volumeMap
		} else if volume.DownwardAPI != nil {
			//TODO
		} else if volume.FC != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[FSType] = volume.FC.FSType
			volume.FC.FSType = VolumeTemplateForElement(volume.Name, FSType)
			persistence[volume.Name] = volumeMap
		} else if volume.AzureFile != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[SecretName] = volume.AzureFile.SecretName
			volumeMap[ShareName] = volume.AzureFile.ShareName
			volume.AzureFile.ShareName = VolumeTemplateForElement(volume.Name, ShareName)
			volume.AzureFile.SecretName = VolumeTemplateForElement(volume.Name, SecretName)
			persistence[volume.Name] = volumeMap
		} else if volume.AzureDisk != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[DiskName] = volume.AzureDisk.DiskName
			volumeMap[DataDiskURI] = volume.AzureDisk.DataDiskURI
			//volumeMap[FSType] = volume.AzureDisk.FSType
			volume.AzureDisk.DiskName = VolumeTemplateForElement(volume.Name, DiskName)
			volume.AzureDisk.DataDiskURI = VolumeTemplateForElement(volume.Name, DataDiskURI)
			//volume.AzureDisk.FSType = *string(VolumeTemplateForElement(volume.Name, "FSType"))
			persistence[volume.Name] = volumeMap
		} else if volume.VsphereVolume != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap[FSType] = volume.VsphereVolume.FSType
			volumeMap[VolumePath] = volume.VsphereVolume.VolumePath
			volume.VsphereVolume.FSType = VolumeTemplateForElement(volume.Name, FSType)
			volume.VsphereVolume.VolumePath = VolumeTemplateForElement(volume.Name, VolumePath)
			persistence[volume.Name] = volumeMap
		}
		volumeData, err := yaml.Marshal(vol)
		if err != nil {
			log.Fatal(err)
		}
		if len(ifCondition) != 0 {
			partialvolumeTemplate = partialVolumeTemplate(string(volumeData), ifCondition)
		} else {
			partialvolumeTemplate = string(volumeData)
		}
		volumeTemplate = volumeTemplate + partialvolumeTemplate
	}
	return volumeTemplate, persistence
}

func generateTemplateForContainer(containers []apiv1.Container, key string, value map[string]interface{}) []apiv1.Container {
	result := make([]apiv1.Container, len(containers))
	for i, container := range containers {
		containterValue := make(map[string]interface{}, 0)
		containerName := generateSafeKey(container.Name)
		container.Image = addTemplateImageValue(containerName, container.Image, key, containterValue)
		if len(container.ImagePullPolicy) != 0 {
			containterValue[ImagePullPolicy] = string(container.ImagePullPolicy)
			container.ImagePullPolicy = apiv1.PullPolicy(addContainerValue(key, containerName, ImagePullPolicy))
		}
		if len(container.Env) != 0 {
			for k, v := range container.Env {
				envName := generateSafeKey(v.Name)
				if len(v.Value) != 0 {
					containterValue[envName] = v.Value
				}
				if v.ValueFrom != nil {
					if v.ValueFrom.ConfigMapKeyRef != nil {
						if checkIfNameExist(v.ValueFrom.ConfigMapKeyRef.Name, "Configmap") {
							container.Env[k].ValueFrom.ConfigMapKeyRef.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, v.ValueFrom.ConfigMapKeyRef.Name)
							containterValue[envName] = v.ValueFrom.ConfigMapKeyRef.Key
						}
					} else if v.ValueFrom.SecretKeyRef != nil {
						if checkIfNameExist(v.ValueFrom.SecretKeyRef.Name, "Secret") {
							container.Env[k].ValueFrom.SecretKeyRef.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, v.ValueFrom.SecretKeyRef.Name)
							containterValue[envName] = v.ValueFrom.SecretKeyRef.Key
						}
					}
				}
				container.Env[k].Value = fmt.Sprintf("{{.Values.%s.%s.%s}}", key, generateSafeKey(container.Name), envName)
			}
		}

		result[i] = container
		value[generateSafeKey(container.Name)] = containterValue
	}
	return result
}

func generateTemplateForLables(labels map[string]string) map[string]string { // Add labels needed for chart
	if labels == nil {
		labels = make(map[string]string, 0)
	}
	labels["chart"] = `{{.Chart.Name}}-{{.Chart.Version}}`
	labels["release"] = `{{.Release.Name}}`
	labels["heritage"] = `{{.Release.Service}}`
	return labels
}

func partialVolumeTemplate(data string, ifCondition string) string {
	volumeElse := `{{- else }}
  emptyDir: {}
{{- end }}
`
	templateData := ""
	str := strings.Split(data, "\n")
	templateData = templateData + ifCondition + "\n"
	for _, l := range str {
		if len(l) == 0 {
			continue
		}
		templateData = templateData + l + "\n"
	}
	templateData = templateData + volumeElse
	return templateData
}

func SaveChartfile(filename string, cf *chart.Metadata) error {
	out, err := yaml.Marshal(cf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0755)
}

func addContainerValue(key string, s1 string, s2 string) string {
	value := fmt.Sprintf("{{.Values.%s.%s.%s}}", key, s1, s2)
	return value
}

func addTemplateImageValue(containerName string, image string, key string, containerValue map[string]interface{}) string {
	// Example: appscode/voyager:1.5.1                 , appscode/voyager
	// Example: docker.appscode.com/ark:0.1.0          , docker.appscode.com/ark
	// Example: localhost.localdomain:5000/ubuntu:16.04, localhost.localdomain:5000/ubuntu
	indexSlash := strings.LastIndex(image, "/")
	indexColon := strings.LastIndex(image, ":")
	if indexColon > indexSlash {
		// user used image tag
		containerValue[Image] = image[:indexColon]
		containerValue[ImageTag] = image[indexColon+1:]
	} else {
		containerValue[Image] = image
		containerValue[ImageTag] = "latest"
	}
	key = generateSafeKey(key)
	imageNameTemplate := fmt.Sprintf("{{.Values.%s.%s.%s}}", key, containerName, Image)
	imageTagTemplate := fmt.Sprintf("{{.Values.%s.%s.%s}}", key, containerName, ImageTag)
	imageTemplate := fmt.Sprintf("%s:%s", imageNameTemplate, imageTagTemplate)
	return imageTemplate
}

func addVolumeToTemplateForPod(templatePod string, templatevolumes string) string {
	templatevolumes = makeSpaceForVolume(templatevolumes, "  ")
	template := addVolumeInPodTemplate(templatePod, templatevolumes)
	return template
}

func removeEmptyFields(temp string) string {
	var resource map[string]interface{}
	err := yaml.Unmarshal([]byte(temp), &resource)
	if err != nil {
		log.Fatal(err)
	}
	delete(resource, "status")
	for k, v := range resource {
		omitEmptyMap(resource, k, v)
	}
	yamlData, err := yaml.Marshal(resource)
	if err != nil {
		log.Fatal(err)
	}
	return string(yamlData)
}

func omitEmptyMap(mp map[string]interface{}, k string, v interface{}) {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem()
	}
	if isEmptyValue(reflect.ValueOf(v)) {
		delete(mp, k)
	} else if !reflect.ValueOf(v).IsValid() {
		delete(mp, k)
	} else if reflect.ValueOf(v).Kind() == reflect.Map || reflect.ValueOf(v).Kind() == reflect.Struct {
		data, err := json.Marshal(reflect.ValueOf(v).Interface())
		if err == nil {
			var newMap map[string]interface{}
			if err := json.Unmarshal(data, &newMap); err != nil {
				log.Fatal(err)
			}
			for k1, v1 := range newMap {
				omitEmptyMap(newMap, k1, v1)
			}
			mp[k] = newMap
		}
	} else if reflect.ValueOf(v).Kind() == reflect.Slice {
		mp[k] = omitEmptySlice(InterfaceToSlice(v))
	}
}

func omitEmptySlice(i []interface{}) []interface{} {
	var z []interface{}
	for _, v := range i {
		if reflect.ValueOf(v).Kind() == reflect.Ptr {
			v = reflect.ValueOf(v).Elem()
		}
		if isEmptyValue(reflect.ValueOf(v)) {

		} else if !reflect.ValueOf(v).IsValid() {

		} else if reflect.ValueOf(v).Kind() == reflect.Map || reflect.ValueOf(v).Kind() == reflect.Struct {
			data, err := json.Marshal(reflect.ValueOf(v).Interface())
			if err != nil {
				log.Fatal(err)
			}
			var newMap map[string]interface{}
			if err := json.Unmarshal(data, &newMap); err != nil {
				log.Fatal(err)
			}
			for k1, v1 := range newMap {
				omitEmptyMap(newMap, k1, v1)
			}
			z = append(z, newMap)
		} else if reflect.ValueOf(v).Kind() == reflect.Slice {
			v1 := omitEmptySlice(InterfaceToSlice(v))
			z = append(z, v1)

		} else {
			z = append(z, v)
		}
	}
	return z
}

func InterfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Fatal("Not slice")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func makeSpaceForVolume(templatevolumes string, space string) string {
	s := space + "volumes:\n"
	str := strings.Split(templatevolumes, "\n")
	for _, l := range str {
		if len(l) == 0 {
			continue
		}
		s = s + space + l + "\n"
	}
	return s
}

func addVolumeInPodTemplate(pod string, volume string) string {
	str := strings.Split(pod, "\n")
	templateForPod := ""
	for _, l := range str {
		if len(l) == 0 {
			continue
		}
		if strings.HasPrefix("spec:", l) {
			templateForPod = templateForPod + l + "\n" + volume
		} else {
			templateForPod = templateForPod + l + "\n"
		}
	}
	return templateForPod
}

func addVolumeToTemplate(rc string, volumes string) string {
	volumes = makeSpaceForVolume(volumes, "      ")
	template := addVolumeInRcTemplate(rc, volumes)
	return template
}

func addVolumeInRcTemplate(rc string, volumes string) string {
	str := strings.Split(rc, "\n")
	templateForPod := ""
	for _, l := range str {
		if len(l) == 0 {
			continue
		}
		if strings.HasPrefix("    spec:", l) {
			templateForPod = templateForPod + l + "\n" + volumes
		} else {
			templateForPod = templateForPod + l + "\n"
		}

	}
	return templateForPod
}

func generateServiceSpecTemplate(svc apiv1.ServiceSpec, key string, value map[string]interface{}) apiv1.ServiceSpec {
	if len(svc.ClusterIP) != 0 {
		value[ClusterIP] = svc.ClusterIP
		svc.ClusterIP = fmt.Sprintf("{{.Values.%s.%s}}", key, ClusterIP)
	}
	if len(svc.ExternalName) != 0 {
		value[ExternalName] = svc.ExternalName
		svc.ExternalName = fmt.Sprintf("{{.Values.%s.%s}}", key, ExternalName)
	}
	if len(svc.LoadBalancerIP) != 0 {
		value[LoadBalancer] = svc.LoadBalancerIP
		svc.LoadBalancerIP = fmt.Sprintf("{{.Values.%s.%s}}", key, LoadBalancer)
	}
	if len(string(svc.Type)) != 0 {
		value[ServiceType] = string(svc.Type)
		svc.Type = apiv1.ServiceType(fmt.Sprintf("{{.Values.%s.%s}}", key, ServiceType))
	}
	if len(string(svc.SessionAffinity)) != 0 {
		value[SessionAffinity] = string(svc.SessionAffinity)
		svc.SessionAffinity = apiv1.ServiceAffinity(fmt.Sprintf("{{.Values.%s.%s}}", key, SessionAffinity))
	}
	return svc
}

func generatePersistentVolumeClaimSpec(pvcspec apiv1.PersistentVolumeClaimSpec, key string, value map[string]interface{}) apiv1.PersistentVolumeClaimSpec {
	if len(pvcspec.VolumeName) != 0 {
		value[VolumeName] = pvcspec.VolumeName
		pvcspec.VolumeName = fmt.Sprintf("{{.Values.%s.%s}}", key, VolumeName)
	}
	if len(pvcspec.AccessModes) != 0 {
		value[AccessMode] = pvcspec.AccessModes[0] //TODO sauman (multiple access mode)
		pvcspec.AccessModes = nil
		pvcspec.AccessModes = append(pvcspec.AccessModes, apiv1.PersistentVolumeAccessMode(fmt.Sprintf("{{.Values.%s.%s}}", key, AccessMode)))
	}
	if pvcspec.Resources.Requests != nil {
		//TODO sauman
	}
	return pvcspec
}

func generatePersistentVolumeSpec(spec apiv1.PersistentVolumeSpec, key string, value map[string]interface{}) apiv1.PersistentVolumeSpec {
	value[ReclaimPolicy] = spec.PersistentVolumeReclaimPolicy
	spec.PersistentVolumeReclaimPolicy = apiv1.PersistentVolumeReclaimPolicy(fmt.Sprintf("{{.Values.%s.%s}}", key, ReclaimPolicy))
	if len(spec.AccessModes) != 0 {
		value[AccessMode] = spec.AccessModes[0] //TODO sauman (multiple access mode)
		spec.AccessModes = nil
		spec.AccessModes = append(spec.AccessModes, apiv1.PersistentVolumeAccessMode(fmt.Sprintf("{{.Values.%s.%s}}", key, AccessMode)))
	}
	return spec
}

func generateSafeKey(name string) string {
	newName := ""
	for _, v := range name {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			newName = newName + string(v)
		}
	}
	if len(newName) == 0 {
		newName = name
	}
	return newName
}

func VolumeTemplateForElement(volumeName string, element string) string {
	return fmt.Sprintf(`{{.Values.%s.%s}}`, volumeName, element)
}

func buildIfConditionForVolume(volumeName string) string {
	return fmt.Sprintf("{{- if .Values.persistence.%s.%s}}", volumeName, Enabled)
}

func checkIfNameExist(name string, objType string) bool {
	flag := false
	for _, v := range ChartObject[objType] {
		if v == name {
			flag = true
			break
		}
	}
	return flag
}
