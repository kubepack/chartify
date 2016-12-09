package pkg

import (
	"github.com/spf13/cobra"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CreateChart() *cobra.Command {
	kubeObjects := objects{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "chartify create",
		Run: func(cmd *cobra.Command, args []string) {
			chartData := chartInfo{
				location:  location,
				dir:       dir,
				chartName: chartName,
			}
			if len(dir) != 0 {
				chartData.yamlFiles =  readLocalFiles(dir)
				location = checkLocation(location)

				chartData.Create()
			} else {
				chartData.yamlFiles = kubeObjects.makeYamlListFromKube()

			}
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "", "specify the directory of the yaml files")
	cmd.Flags().StringVar(&location, "location", "", "specify the location where charts will be created. By default Home directory")
	cmd.Flags().StringVar(&chartName, "chart_name", "", "specify the chart name")
	cmd.Flags().StringVar(&kubeObjects.kubeContext, "kube_context", "", "specify the kube context name")
	cmd.Flags().StringVar(&kubeObjects.namespace, "namespace", "default", "specify the namespace for the selected objects")
	cmd.Flags().StringArrayVar(&kubeObjects.pods,"pods",kubeObjects.pods,"specify the names of pods to incluse them in chart")
	cmd.Flags().StringArrayVar(&kubeObjects.replicationControllers, "replicationControllers", kubeObjects.replicationControllers, "specify the names if pods to include them in chart")
	cmd.Flags().StringArrayVar(&kubeObjects.services, "services", kubeObjects.services, "specify the names of services to include them in chart")
	cmd.Flags().StringArrayVar(&kubeObjects.configMaps, "config_maps", kubeObjects.configMaps, "specify the names of secrets to include them in chart")
	cmd.Flags().StringArrayVar(&kubeObjects.configMaps, "secrets", kubeObjects.configMaps, "specify the names of secrets to include them in chart")
	cmd.Flags().StringArrayVar(&kubeObjects.persistentVolume, "pv", kubeObjects.persistentVolume, "specify names of persistent volumes")
	cmd.Flags().StringArrayVar(&kubeObjects.persistentVolumeClaim,"pvc", kubeObjects.persistentVolumeClaim, "specify names of persistent volume claim")
	cmd.Flags().StringArrayVar(&kubeObjects.petsets, "petsets", kubeObjects.petsets, "specify specify names of petsets")
	cmd.Flags().StringArrayVar(&kubeObjects.jobs, "jobs", kubeObjects.jobs, "specify names of jobs")
	cmd.Flags().StringArrayVar(&kubeObjects.daemons)

	return cmd
}

func checkLocation(location string) string {
	var err error
	if len(location) == 0 {
		log.Fatal("Location fo the chart file not given")
	} else {
		_, err = os.Stat(location)
		if err != nil {
			log.Fatal(err)
		}
	}
	return location
}

func readLocalFiles(dirName string) []string {
	var yamlFiles []string
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fileDir := filepath.Join(dirName, f.Name())
		dataByte, err := ioutil.ReadFile(fileDir)
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, string(dataByte))
	}
	return yamlFiles
}

func (kubeObjects objects)makeYamlListFromKube() []string {
	kubeClient, err := NewKubeClient(kubeObjects.kubeContext)
	if err != nil {
		log.Fatal(err)
	}
	yamlFiles := kubeObjects.readKubernetesObjects(kubeClient)
	return yamlFiles
}
