package main

import (
	"fmt"
	"os"

	"github.com/appscode/chartify/pkg"
	"github.com/spf13/cobra"
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
