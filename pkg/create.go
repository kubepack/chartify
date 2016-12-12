package pkg

import (
	"fmt"
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
			if chartName == "" {
				fmt.Println("ERROR : Provide a CharName")
				os.Exit(1)
			}
			if location == "" {
				fmt.Println("ERROR : Provide a location for the chart file")
			}
			chartData := chartInfo{
				location:  checkLocation(location),
				dir:       dir,
				chartName: chartName,
			}
			if len(dir) != 0 {
				chartData.yamlFiles = readLocalFiles(dir)
			} else {
				chartData.yamlFiles = kubeObjects.makeYamlListFromKube()
			}
			chartData.Create()
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "", "specify the directory of the yaml files")
	cmd.Flags().StringVar(&location, "location", "", "specify the location where charts will be created. By default Home directory")
	cmd.Flags().StringVar(&chartName, "chart_name", "", "specify the chart name")
	cmd.Flags().StringVar(&kubeObjects.kubeContext, "kube_context", "", "specify the kube context name")
	//cmd.Flags().StringVar(&kubeObjects.namespace, "namespace", "default", "specify the namespace for the selected objects")
	cmd.Flags().StringSliceVar(&kubeObjects.pods, "pods", kubeObjects.pods, "specify the names of pods (podname.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&kubeObjects.replicationControllers, "replicationControllers", kubeObjects.replicationControllers, "specify the names of replication cotrollers (rcname.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&kubeObjects.services, "services", kubeObjects.services, "specify the names of services to include them in chart")
	cmd.Flags().StringSliceVar(&kubeObjects.configMaps, "config_maps", kubeObjects.configMaps, "specify the names of configmaps(configmap.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&kubeObjects.secrets, "secrets", kubeObjects.configMaps, "specify the names of secrets(secret_name.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&kubeObjects.persistentVolume, "pv", kubeObjects.persistentVolume, "specify names of persistent volumes")
	cmd.Flags().StringSliceVar(&kubeObjects.persistentVolumeClaim, "pvc", kubeObjects.persistentVolumeClaim, "specify names of persistent volume claim")
	cmd.Flags().StringSliceVar(&kubeObjects.petsets, "petsets", kubeObjects.petsets, "specify names of petsets(petset_name.namespace)")
	cmd.Flags().StringSliceVar(&kubeObjects.jobs, "jobs", kubeObjects.jobs, "specify names of jobs")
	cmd.Flags().StringSliceVar(&kubeObjects.replicaSet, "replica_sets", kubeObjects.replicaSet, "specify names of replica sets(replicaset_name.namespace)")
	cmd.Flags().StringSliceVar(&kubeObjects.daemons, "daemons", kubeObjects.daemons, "specify names of daemon sets")

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

func (kubeObjects objects) makeYamlListFromKube() []string {
	kubeClient, err := NewKubeClient(kubeObjects.kubeContext)
	if err != nil {
		log.Fatal(err)
	}
	yamlFiles := kubeObjects.readKubernetesObjects(kubeClient)
	return yamlFiles
}
