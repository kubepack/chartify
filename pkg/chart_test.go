package pkg

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateCustomChartFromYaml(t *testing.T) {
	c := chartInfo{
		dir:       "/home/sauman/yamlDir",
		location:  "/home/sauman/helm-test",
		chartName: "test",
	}
	c.yamlFiles = readLocalFiles(c.dir)
	_, err := c.Create()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func TestResourceEmpty(t *testing.T) {
	yamlData := `apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    chart: '{{.Chart.Name}}-{{.Chart.Version}}'
    heritage: '{{.Release.Service}}'
    release: '{{.Release.Name}}'
  name: '{{.Release.Name}}-{{.Values.Name}}'
spec:
  containers:
  - image: '{{.Values.myfrontend.image}}:{{.Values.myfrontend.imageTag}}'
    imagePullPolicy: '{{.Values.myfrontend.imagePullPolicy}}'
    name: myfrontend
    resources: {}
    volumeMounts:
    - mountPath: /var/www/html
      name: mypd
  serviceAccountName: ""
  volumes: []
status: {}
`
	s := removeEmptyFields(yamlData)
	fmt.Println(s)
}
