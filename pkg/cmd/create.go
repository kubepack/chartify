package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kubepack.dev/chartify/pkg"
)

func NewCmdCreate() *cobra.Command {
	var (
		kubeDir      string
		chartDir     string
		preserveName bool
	)
	ko := pkg.KubeObjects{}

	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create Helm Charts from Kubernetes api objects",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("ERROR : Provide a ChartName")
				os.Exit(1)
			}
			gen := pkg.Generator{
				Location:  checkLocation(chartDir),
				ChartName: args[0],
			}
			pkg.PreserveName = preserveName
			if len(kubeDir) != 0 {
				gen.YamlFiles = pkg.ReadLocalFiles(kubeDir)
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
	cmd.Flags().BoolVar(&preserveName, "preserve-name", false, "Specify if you want to preserve resources name from input yaml true/false (default: false)")
	cmd.Flags().StringSliceVar(&ko.ConfigMaps, "configmaps", ko.ConfigMaps, "Specify the names of configmaps(configmap@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Daemons, "daemons", ko.Daemons, "Specify the names of daemons(daemon@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Deployments, "deployments", ko.Deployments, "Specify the names of deployments(deployments@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Jobs, "jobs", ko.Jobs, "Specify the names of jobs(job@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.PersistentVolumes, "pvs", ko.PersistentVolumes, "Specify the names of persistent volumes(pv@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.PersistentVolumeClaims, "pvcs", ko.PersistentVolumeClaims, "Specify the names of persistent volume claims(pvc@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Pods, "pods", ko.Pods, "Specify the names of pods(pod@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.ReplicaSets, "replicasets", ko.ReplicaSets, "Specify the names of replica sets(rs@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.ReplicationControllers, "rcs", ko.ReplicationControllers, "Specify the names of replication cotrollers(rc@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Secrets, "secrets", ko.Secrets, "Specify the names of secrets(secret@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.Services, "services", ko.Services, "Specify the names of services(service@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.StatefulSets, "statefulsets", ko.StatefulSets, "Specify the names of statefulsets(statefulset@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.StorageClasses, "storageclasses", ko.StorageClasses, "Specify the names of storageclasses(storageclass@namespace) to include in chart")
	cmd.Flags().StringSliceVar(&ko.HorizontalPodAutoscalers, "horizontalpodautoscalers", ko.HorizontalPodAutoscalers, "Specify the names of horizontalpodautoscalers(horizontalpodautoscaler@namespace) to include in chart")

	return cmd
}

func checkLocation(location string) string {
	if len(location) == 0 {
		log.Fatalln("ERROR : Provide a chart directory")
	}
	fi, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(location, 0755); err != nil {
				log.Fatal(err)
			}
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
