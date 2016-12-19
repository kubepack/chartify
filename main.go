package main

import (
	"fmt"
	"os"

	"github.com/appscode/chartify/pkg/cmd"
	"github.com/spf13/cobra"
)

var (
	Version         string
	VersionStrategy string
	Os              string
	Arch            string
	CommitHash      string
	GitBranch       string
	GitTag          string
	CommitTimestamp string
	BuildTimestamp  string
	BuildHost       string
	BuildHostOs     string
	BuildHostArch   string
)

func init() {
	cmd.Version.Version = Version
	cmd.Version.VersionStrategy = VersionStrategy
	cmd.Version.Os = Os
	cmd.Version.Arch = Arch
	cmd.Version.CommitHash = CommitHash
	cmd.Version.GitBranch = GitBranch
	cmd.Version.GitTag = GitTag
	cmd.Version.CommitTimestamp = CommitTimestamp
	cmd.Version.BuildTimestamp = BuildTimestamp
	cmd.Version.BuildHost = BuildHost
	cmd.Version.BuildHostOs = BuildHostOs
	cmd.Version.BuildHostArch = BuildHostArch
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "chartify [command]",
		Short: `Generate Helm Charts from Kubernetes api objects`,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	rootCmd.AddCommand(cmd.NewCmdCreate())
	rootCmd.AddCommand(cmd.NewCmdVersion())
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
