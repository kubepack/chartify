package pkg

const (
	// ChartfileName is the default Chart file name.
	ChartfileName = "Chart.yaml"
	// ValuesfileName is the default values file name.
	ValuesfileName = "values.yaml"
	// TemplatesDir is the relative directory name for templates.
	TemplatesDir = "templates"
	// ChartsDir is the relative directory name for charts dependencies.
	ChartsDir = "charts"
	// IgnorefileName is the name of the Helm ignore file.
	IgnorefileName = ".helmignore"
	// DeploymentName is the name of the example deployment file.
	DeploymentName = "deployment.yaml"
	// ServiceName is the name of the example service file.
	ServiceName = "service.yaml"
	// NotesName is the name of the example NOTES.txt file.
	NotesName = "NOTES.txt"
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

var (
	dir       string
	location  string
	chartName string
)

type objects struct {
	kubeContext            string
	namespace              string
	pods                   []string
	replicationControllers []string
	configMaps             []string
	services               []string
	secrets                []string
	persistentVolume       []string
	persistentVolumeClaim  []string
	petsets                []string
	jobs                   []string
	daemons                []string
	replicaSet             []string
}

type chartInfo struct {
	dir       string
	location  string
	chartName string
	yamlFiles []string
}
