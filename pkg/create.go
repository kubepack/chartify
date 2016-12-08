package pkg

import (
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/*func CreateChart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "kitten create",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(Chart())
	return cmd
}*/

func CreateChart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "chartify create",
		Run: func(cmd *cobra.Command, args []string) {
			if len(dir) != 0 {
				/*				files, err := ioutil.ReadDir(dir)
								if err != nil {
									log.Fatal(err)
								}*/
				yamlFiles := readLocalFiles(dir)
				location = checkLocation(location)
				chartData := chartInfo{
					location:  location,
					dir:       dir,
					chartName: chartName,
					yamlFiles: yamlFiles,
				}
				chartData.Create()
			}
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "", "specify the directory of the yaml files")
	cmd.Flags().StringVar(&location, "location", "", "specify the location where charts will be created. By default Home directory")
	cmd.Flags().StringVar(&chartName, "chart_name", "", "specify the location where charts will be created. By default Home directory")

	return cmd
}

func checkLocation(location string) string {
	var err error
	if len(location) == 0 {
		location, err = homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
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
