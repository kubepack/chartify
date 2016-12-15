package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/proto/hapi/chart"
	kapi "k8s.io/kubernetes/pkg/api"
)

func generateObjectMetaTemplate(objectMeta kapi.ObjectMeta, key string, value map[string]interface{}, extraTagForName string) kapi.ObjectMeta {
	key = checkKeyValue(key)
	if len(objectMeta.Name) != 0 {
		objectMeta.Name = fmt.Sprintf(`{{ template "fullname" . }}`)
	}
	if len(extraTagForName) != 0 {
		objectMeta.Name = fmt.Sprintf("%s-%s", objectMeta.Name, extraTagForName)
	}
	if len(objectMeta.ClusterName) != 0 {
		value["ClusterName"] = objectMeta.ClusterName
		objectMeta.ClusterName = fmt.Sprintf("{{.Values%s.ClusterName}}", key)
	}
	if len(objectMeta.GenerateName) != 0 {
		value["GenerateName"] = objectMeta.GenerateName
		objectMeta.GenerateName = fmt.Sprintf("{{.Values%s.GenerateName}}", key)
	}
	if len(objectMeta.Namespace) != 0 {
		value["Namespace"] = objectMeta.Namespace
		objectMeta.Namespace = fmt.Sprintf("{{.Values%s.Namespace}}", key)
	}
	if len(objectMeta.SelfLink) != 0 {
		value["SelfLink"] = objectMeta.SelfLink
		objectMeta.SelfLink = fmt.Sprintf("{{.Values%s.SelfLink}}", key)
	}
	objectMeta.Labels = generateTemplateForLables(objectMeta.Labels)
	return objectMeta
}

func generateTemplateForPodSpec(podSpec kapi.PodSpec, key string, value map[string]interface{}) kapi.PodSpec {
	podSpec.Containers = generateTemplateForContainer(podSpec.Containers, key, value)
	key = checkKeyValue(key)
	if len(podSpec.Hostname) != 0 {
		value["HostName"] = podSpec.Hostname
		podSpec.Hostname = fmt.Sprintf("{{.Values%s.HostName}}", key)
	}
	if len(podSpec.Subdomain) != 0 {
		value["Subdomain"] = podSpec.Subdomain
		podSpec.Subdomain = fmt.Sprintf("{{.Values%s.Subdomain}}", key)
	}
	if len(podSpec.NodeName) != 0 {
		value["Nodename"] = podSpec.NodeName
		podSpec.NodeName = fmt.Sprintf("{{.Values%s.Nodename}}", key)
	}
	if len(podSpec.ServiceAccountName) != 0 {
		value["ServiceAccountName"] = podSpec.ServiceAccountName
		podSpec.ServiceAccountName = fmt.Sprintf("{{.Values%s.ServiceAccountName}}", key)
	}
	return podSpec
}

func generateTemplateForVolume(volumes []kapi.Volume, key string, value map[string]interface{}) (string, map[string]interface{}) {
	key = checkKeyValue(key)
	volumeTemplate := ""
	ifCondition := ""
	partialvolumeTemplate := ""
	persistence := make(map[string]interface{}, 0)
	for _, volume := range volumes {
		ifCondition = ""
		volumeMap := make(map[string]interface{}, 0)
		volumeMap["enabled"] = true
		vol := []kapi.Volume{}
		vol = append(vol, volume)
		if volume.PersistentVolumeClaim != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volume.PersistentVolumeClaim.ClaimName = fmt.Sprintf(`{{template "fullname"}}-%s`, volume.PersistentVolumeClaim.ClaimName)
		} else if volume.ConfigMap != nil {
			//volume.ConfigMap.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, volume.ConfigMap.Name) // TODO if config map is deployed by helm map name will be like that
		} else if volume.Secret != nil {
			//volume.Secret.SecretName = fmt.Sprintf(`{{ template "fullname" . }}-%s`, volume.Secret.SecretName) // TODO if secret is deployed by helm map name will be like that
			//TODO add items
		} else if volume.Glusterfs != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Path"] = volume.Glusterfs.Path

			volumeMap["EndpointsName"] = volume.Glusterfs.EndpointsName
			volume.Glusterfs.EndpointsName = VolumeTemplateForElement(volume.Name, "EndpointsName")
			volume.Glusterfs.Path = VolumeTemplateForElement(volume.Name, "Path")
			persistence[volume.Name] = volumeMap
		} else if volume.HostPath != nil {
			volumeMap["Path"] = volume.HostPath.Path
			volume.HostPath.Path = VolumeTemplateForElement(volume.Name, "Path")
			persistence[volume.Name] = volumeMap
		} else if volume.GCEPersistentDisk != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["PDName"] = volume.GCEPersistentDisk.PDName
			volumeMap["FSType"] = volume.GCEPersistentDisk.FSType
			volume.GCEPersistentDisk.PDName = VolumeTemplateForElement(volume.Name, "PDName")
			volume.GCEPersistentDisk.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			persistence[volume.Name] = volumeMap
		} else if volume.AWSElasticBlockStore != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["FSType"] = volume.GCEPersistentDisk.FSType
			volumeMap["VolumeID"] = volume.AWSElasticBlockStore.VolumeID
			volume.AWSElasticBlockStore.VolumeID = VolumeTemplateForElement(volume.Name, "VolumeID")
			volume.AWSElasticBlockStore.FSType = VolumeTemplateForElement(volume.Name, "FSType")
		} else if volume.GitRepo != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Repository"] = volume.GitRepo.Repository
			volumeMap["Revision"] = volume.GitRepo.Revision
			volumeMap["Directory"] = volume.GitRepo.Directory
			volume.GitRepo.Revision = VolumeTemplateForElement(volume.Name, "Revision")
			volume.GitRepo.Repository = VolumeTemplateForElement(volume.Name, "Repository")
			volume.GitRepo.Directory = VolumeTemplateForElement(volume.Name, "Directory")
			persistence[volume.Name] = volumeMap
		} else if volume.NFS != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Server"] = volume.NFS.Server
			volumeMap["Path"] = volume.NFS.Path
			volume.NFS.Path = fmt.Sprintf(`{{.Values.%s.Path}}`, volume.Name)
			volume.NFS.Server = fmt.Sprintf(`{{.Values.%s.Server}}`, volume.Name)
			persistence[volume.Name] = volumeMap
		} else if volume.ISCSI != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["TargetPortal"] = volume.ISCSI.TargetPortal
			volumeMap["IQN"] = volume.ISCSI.IQN
			volumeMap["ISCSIInterface"] = volume.ISCSI.ISCSIInterface
			volumeMap["FSType"] = volume.ISCSI.FSType
			volume.ISCSI.TargetPortal = VolumeTemplateForElement(volume.Name, "TargetPortal")
			volume.ISCSI.IQN = VolumeTemplateForElement(volume.Name, "IQN")
			volume.ISCSI.FSType = fmt.Sprintf(`{{.Values.%s.FSType}}`, volume.Name)
			volume.ISCSI.ISCSIInterface = VolumeTemplateForElement(volume.Name, "ISCSIInterface")
			persistence[volume.Name] = volumeMap
		} else if volume.RBD != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["FSType"] = volume.RBD.FSType
			volumeMap["RBDImage"] = volume.RBD.RBDImage
			volumeMap["RBDPool"] = volume.RBD.RBDPool
			volumeMap["RadosUser"] = volume.RBD.RadosUser
			volumeMap["Keyring"] = volume.RBD.Keyring
			volume.RBD.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			volume.RBD.RBDImage = VolumeTemplateForElement(volume.Name, "RBDImage")
			volume.RBD.RBDPool = VolumeTemplateForElement(volume.Name, "RBDPool")
			volume.RBD.RadosUser = VolumeTemplateForElement(volume.Name, "RadosUser")
			volume.RBD.Keyring = VolumeTemplateForElement(volume.Name, "Keyring")
			persistence[volume.Name] = volumeMap
		} else if volume.Quobyte != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Registry"] = volume.Quobyte.Registry
			volumeMap["Volume"] = volume.Quobyte.Volume
			volumeMap["Group"] = volume.Quobyte.Group
			volumeMap["User"] = volume.Quobyte.User
			volume.Quobyte.Registry = VolumeTemplateForElement(volume.Name, "Registry")
			volume.Quobyte.Volume = VolumeTemplateForElement(volume.Name, "Volume")
			volume.Quobyte.Group = VolumeTemplateForElement(volume.Name, "Group")
			volume.Quobyte.User = VolumeTemplateForElement(volume.Name, "User")
			persistence[volume.Name] = volumeMap
		} else if volume.FlexVolume != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Driver"] = volume.FlexVolume.Driver
			volumeMap["FSType"] = volume.FlexVolume.FSType
			// TODO secret reference
			volume.FlexVolume.Driver = VolumeTemplateForElement(volume.Name, "Driver")
			volume.FlexVolume.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			persistence[volume.Name] = volumeMap
		} else if volume.Cinder != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["FSType"] = volume.Cinder.FSType
			volumeMap["VolumeID"] = volume.Cinder.VolumeID
			volume.Cinder.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			volume.Cinder.VolumeID = VolumeTemplateForElement(volume.Name, "VolumeID")
			persistence[volume.Name] = volumeMap
		} else if volume.CephFS != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["Path"] = volume.CephFS.Path
			volumeMap["SecretFile"] = volume.CephFS.SecretFile
			volumeMap["User"] = volume.CephFS.User
			volume.CephFS.Path = VolumeTemplateForElement(volume.Name, "Path")
			volume.CephFS.SecretFile = VolumeTemplateForElement(volume.Name, "SecretFile")
			volume.CephFS.User = VolumeTemplateForElement(volume.Name, "User")
			persistence[volume.Name] = volumeMap
		} else if volume.Flocker != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["DatasetName"] = volume.Flocker.DatasetName
			volume.Flocker.DatasetName = VolumeTemplateForElement(volume.Name, "DatasetName")
			persistence[volume.Name] = volumeMap
		} else if volume.DownwardAPI != nil {
			//TODO
		} else if volume.FC != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["FSType"] = volume.FC.FSType
			volume.FC.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			persistence[volume.Name] = volumeMap
		} else if volume.AzureFile != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["SecretName"] = volume.AzureFile.SecretName
			volumeMap["ShareName"] = volume.AzureFile.ShareName
			volume.AzureFile.ShareName = VolumeTemplateForElement(volume.Name, "ShareName")
			volume.AzureFile.SecretName = VolumeTemplateForElement(volume.Name, "SecretName")
			persistence[volume.Name] = volumeMap
		} else if volume.AzureDisk != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["DiskName"] = volume.AzureDisk.DiskName
			volumeMap["DataDiskURI"] = volume.AzureDisk.DataDiskURI
			volumeMap["FSType"] = volume.AzureDisk.FSType
			volume.AzureDisk.DiskName = VolumeTemplateForElement(volume.Name, "DiskName")
			volume.AzureDisk.DataDiskURI = VolumeTemplateForElement(volume.Name, "DataDiskURI")
			//volume.AzureDisk.FSType = *string(VolumeTemplateForElement(volume.Name, "FSType"))
			persistence[volume.Name] = volumeMap
		} else if volume.VsphereVolume != nil {
			ifCondition = buildIfConditionForVolume(volume.Name)
			volumeMap["FSType"] = volume.VsphereVolume.FSType
			volumeMap["VolumePath"] = volume.VsphereVolume.VolumePath
			volume.VsphereVolume.FSType = VolumeTemplateForElement(volume.Name, "FSType")
			volume.VsphereVolume.VolumePath = VolumeTemplateForElement(volume.Name, "VolumePath")
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

func generateTemplateForContainer(containers []kapi.Container, key string, value map[string]interface{}) []kapi.Container {
	containterValue := make(map[string]interface{}, 0)
	result := make([]kapi.Container, len(containers))
	for k := range containterValue {
		delete(containterValue, k)
	}
	for i, container := range containers {
		containerName := removeCharactersFromName(container.Name)
		container.Image = addTemplateImageValue(containerName, container.Image, key, containterValue)
		if len(container.ImagePullPolicy) != 0 {
			containterValue["ImagePullPolicy"] = string(container.ImagePullPolicy)
			container.ImagePullPolicy = kapi.PullPolicy(addContainerValue(key, containerName, "imagePullPolicy"))
		}
		if len(container.Env) != 0 {
			for k, v := range container.Env {
				if len(v.Value) != 0 {
					tmp := removeCharactersFromName(v.Name)
					containterValue[tmp] = v.Value
					key := checkKeyValue(key)
					container.Env[k].Value = fmt.Sprintf("{{.Values%s.%s.%s}}", key, removeCharactersFromName(container.Name), tmp) //TODO test
				}
				// Secret of Configmap has to be deployed by chart. else value from wont work.
				/*				if v.ValueFrom != nil {
								if v.ValueFrom.ConfigMapKeyRef != nil {
									container.Env[k].ValueFrom.ConfigMapKeyRef.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, v.ValueFrom.ConfigMapKeyRef.Name)
								} else if v.ValueFrom.SecretKeyRef != nil {
									container.Env[k].ValueFrom.SecretKeyRef.Name = fmt.Sprintf(`{{ template "fullname" . }}-%s`, v.ValueFrom.SecretKeyRef.Name)
								}
							}*/
			}
		}
		result[i] = container
		value[removeCharactersFromName(container.Name)] = containterValue
	}
	return result
}

func generateTemplateForLables(labels map[string]string) map[string]string { // Add levels needed for chart
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
	key = checkKeyValue(key)
	value := fmt.Sprintf("{{.Values%s.%s.%s}}", key, s1, s2)
	return value
}

func addTemplateImageValue(containerName string, image string, key string, containerValue map[string]interface{}) string {
	img := strings.Split(image, ":")
	imageName := ""
	imageTag := "latest"
	key = checkKeyValue(key)
	imageNameTemplate := fmt.Sprintf("{{.Values%s.%s.image}}", key, containerName)
	imageTagTemplate := fmt.Sprintf("{{.Values%s.%s.imageTag}}", key, containerName)
	if len(img) == 2 {
		imageName = img[0]
		imageTag = img[1]
	} else {
		imageName = img[0]
	}
	containerValue["image"] = imageName
	containerValue["imageTag"] = imageTag
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
			err = json.Unmarshal(data, &newMap)
			if err != nil {
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
			if err == nil {
				var newMap map[string]interface{}
				err = json.Unmarshal(data, &newMap)
				if err != nil {
					log.Fatal(err)
				}
				for k1, v1 := range newMap {
					omitEmptyMap(newMap, k1, v1)
				}
				z = append(z, newMap)
			} else {
				log.Fatal(err)
			}
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

func addVolumeToTemplateForRc(rc string, volumes string) string {
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

func generateServiceSpecTemplate(svc kapi.ServiceSpec, key string, value map[string]interface{}) kapi.ServiceSpec {
	if len(key) != 0 {
		key = "." + key
	}
	if len(svc.ClusterIP) != 0 {
		value["ClusterIp"] = svc.ClusterIP
		svc.ClusterIP = fmt.Sprintf("{{.Values%s.ClusterIp}}", key)
	}
	if len(svc.ExternalName) != 0 {
		value["ExternalName"] = svc.ExternalName
		svc.ExternalName = fmt.Sprintf("{{.Values%s.ExternalName}}", key)
	}
	if len(svc.LoadBalancerIP) != 0 {
		value["LoadBalancer"] = svc.LoadBalancerIP
		svc.LoadBalancerIP = fmt.Sprintf("{{.Values%s.LoadBalancer}}", key)
	}
	if len(string(svc.Type)) != 0 {
		value["ServiceType"] = string(svc.Type)
		svc.Type = kapi.ServiceType(fmt.Sprintf("{{.Values%s.ServiceType}}", key))
	}
	if len(string(svc.SessionAffinity)) != 0 {
		value["SessionAffinity"] = string(svc.SessionAffinity)
		svc.SessionAffinity = kapi.ServiceAffinity(fmt.Sprintf("{{.Values%s.SessionAffinity}}", key))
	}
	return svc
}

func generatePersistentVolumeClaimSpec(pvcspec kapi.PersistentVolumeClaimSpec, key string, value map[string]interface{}) kapi.PersistentVolumeClaimSpec {
	key = checkKeyValue(key)
	if len(pvcspec.VolumeName) != 0 {
		value["VolumeName"] = pvcspec.VolumeName
		pvcspec.VolumeName = fmt.Sprintf("{{.Values%s.VolumeName}}", key)
	}
	if len(pvcspec.AccessModes) != 0 {
		value["AccessMode"] = pvcspec.AccessModes[0] //TODO sauman (multiple access mode)
		pvcspec.AccessModes = nil
		pvcspec.AccessModes = append(pvcspec.AccessModes, kapi.PersistentVolumeAccessMode(fmt.Sprintf("{{.Values.persistence%s.AccessMode}}", key)))
	}
	if pvcspec.Resources.Requests != nil {
		//TODO sauman
	}
	return pvcspec
}

func generatePersistentVolumeSpec(spec kapi.PersistentVolumeSpec, key string, value map[string]interface{}) kapi.PersistentVolumeSpec {
	value["ReclaimPolicy"] = spec.PersistentVolumeReclaimPolicy
	spec.PersistentVolumeReclaimPolicy = kapi.PersistentVolumeReclaimPolicy(fmt.Sprintf("{{.Values.%s.ReclaimPolicy}}", key))
	if len(spec.AccessModes) != 0 {
		value["AccessMode"] = spec.AccessModes[0] //TODO sauman (multiple access mode)
		spec.AccessModes = nil
		spec.AccessModes = append(spec.AccessModes, kapi.PersistentVolumeAccessMode(fmt.Sprintf("{{.Values.%s.AccessMode}}", key)))
	}
	return spec
}

func checkKeyValue(key string) string {
	// TODO From key have to remove unnecessary characters
	if len(key) != 0 {
		key = "." + key
	}
	return key
}

func removeCharactersFromName(name string) string {
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
	return fmt.Sprintf("{{- if .Values.persistence.%s.enabled}}", volumeName)
}
