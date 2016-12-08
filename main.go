package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/appscode/chartify/pkg"
	"os"
)

func main() {
	cmd := NewCmd()
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chartify",
		Example: `chartify create`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(pkg.CreateChart())
	return cmd

}
