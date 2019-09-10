package pkg

import "fmt"

const (
	// ChartfileName is the default Chart file name.
	ChartfileName = "Chart.yaml"
	// ValuesfileName is the default values file name.
	ValuesfileName = "values.yaml"
	// TemplatesDir is the relative directory name for templates.
	TemplatesDir = "templates"
	// HelpersName is the name of the example NOTES.txt file.
	HelpersName = "_helpers.tpl"
)

const defaultHelpers = `{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 24 -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 24 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 24 -}}
{{- end -}}
`

type valueFileGenerator struct {
	value       map[string]interface{}
	persistence map[string]interface{}
}

const (
	ClusterName                    = "clusterName"
	GenerateName                   = "generateName"
	Namespace                      = "namespace"
	HostName                       = "hostName"
	Subdomain                      = "subdomain"
	Nodename                       = "nodeName"
	ServiceAccountName             = "serviceAccountName"
	Enabled                        = "enabled"
	Path                           = "path"
	EndpointsName                  = "endpointsName"
	VolumeID                       = "volumeID"
	Repository                     = "repository"
	Revision                       = "revision"
	Server                         = "Server"
	TargetPortal                   = "targetPortal"
	FSType                         = "fsType"
	Directory                      = "directory"
	User                           = "user"
	DiskName                       = "diskName"
	DataDiskURI                    = "dataDiskURI"
	VolumePath                     = "volumePath"
	ImagePullPolicy                = "imagePullPolicy"
	PDName                         = "pdName"
	IQN                            = "IQN"
	ISCSIInterface                 = "ISCSIInterface"
	Image                          = "image"
	ImageTag                       = "imageTag"
	ClusterIP                      = "clusterIP"
	ExternalName                   = "externalName"
	LoadBalancer                   = "loadBalancer"
	VolumeName                     = "volumeName"
	AccessMode                     = "accessMode"
	ServiceType                    = "serviceType"
	SessionAffinity                = "sessionAffinity"
	SecretName                     = "secretName"
	ShareName                      = "shareName"
	DatasetName                    = "datasetName"
	SecretFile                     = "secretFile"
	Group                          = "group"
	Volume                         = "volume"
	Registry                       = "registry"
	Keyring                        = "keyring"
	RadosUser                      = "radosUser"
	RBDImage                       = "RBDImage"
	RBDPool                        = "RBDPool"
	Persistence                    = "persistence"
	DeploymentStrategy             = "deploymentStrategy"
	ServiceName                    = "serviceName"
	Type                           = "type"
	Provisioner                    = "provisioner"
	RestartPolicy                  = "restartPolicy"
	ReclaimPolicy                  = "reclaimPolicy"
	MinReplicas                    = "minReplicas"
	MaxReplicas                    = "maxReplicas"
	TargetCPUUtilizationPercentage = "targetCPUUtilizationPercentage"
)

func (v *valueFileGenerator) MergeInto(dst map[string]interface{}, key string) {
	existing, found := dst[key]
	if !found {
		dst[key] = v.value
		return
	}
	if m, ok := existing.(map[string]interface{}); ok {
		for k1, v1 := range v.value {
			m[k1] = v1
		}
		dst[key] = m
	} else {
		fmt.Println("Overwriting string value with map.")
		dst[key] = v.value
		return
	}
}
