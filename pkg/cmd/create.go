package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/appscode/chartify/pkg"
	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	var (
		kubeDir  string
		chartDir string
	)
	ko := pkg.KubeObjects{}

	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create Helm Charts from Kubernetes api objects",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("ERROR : Provide a CharName")
				os.Exit(1)
			}
			gen := pkg.Generator{
				Location:  checkLocation(chartDir),
				ChartName: args[0],
			}
			if len(kubeDir) != 0 {
				gen.YamlFiles = readLocalFiles(kubeDir)
			} else {
				ok := ko.CheckFlags()
				if !ok {
					fmt.Println("No object given.")
					os.Exit(1)
				}
				gen.YamlFiles = ko.Extract()
			}
			gen.Create()
		},
	}
	cmd.Flags().StringVar(&kubeDir, "kube-dir", "", "Specify the directory of the yaml files for Kubernetes objects")
	cmd.Flags().StringVar(&chartDir, "chart-dir", "charts", "Specify the location where charts will be created")
	cmd.Flags().StringSliceVar(&ko.Pods, "pods", ko.Pods, "Specify the names of pods (podname.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&ko.ReplicationControllers, "rcs", ko.ReplicationControllers, "Specify the names of replication cotrollers (rcname.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&ko.Services, "services", ko.Services, "Specify the names of services to include them in chart")
	cmd.Flags().StringSliceVar(&ko.ConfigMaps, "configmaps", ko.ConfigMaps, "Specify the names of configmaps(configmap.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&ko.Secrets, "secrets", ko.ConfigMaps, "Specify the names of secrets(secret_name.namespace) to include them in chart")
	cmd.Flags().StringSliceVar(&ko.PersistentVolume, "pvs", ko.PersistentVolume, "Specify names of persistent volumes")
	cmd.Flags().StringSliceVar(&ko.PersistentVolumeClaim, "pvcs", ko.PersistentVolumeClaim, "Specify names of persistent volume claim")
	cmd.Flags().StringSliceVar(&ko.Statefulsets, "statefulsets", ko.Statefulsets, "Specify names of statefulsets(statefulset_name.namespace)")
	cmd.Flags().StringSliceVar(&ko.Jobs, "jobs", ko.Jobs, "Specify names of jobs")
	cmd.Flags().StringSliceVar(&ko.ReplicaSet, "replicasets", ko.ReplicaSet, "Specify names of replica sets(replicaset_name.namespace)")
	cmd.Flags().StringSliceVar(&ko.Daemons, "daemons", ko.Daemons, "Specify names of daemon sets(daemons.namespace)")
	cmd.Flags().StringSliceVar(&ko.StorageClasses, "storageclasses", ko.StorageClasses, "Specify names of storageclasses")

	return cmd
}

func checkLocation(location string) string {
	if len(location) == 0 {
		log.Fatalln("ERROR : Provide a chart directory")
	}
	fi, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(location, 0755)
		}
		if err != nil {
			log.Fatal(err)
		}
	} else if !fi.IsDir() {
		log.Fatalln(location, "is not a directory.")
	}
	location, err = filepath.Abs(location)
	if err != nil {
		log.Fatal(err)
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
