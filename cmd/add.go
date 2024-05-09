package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list",
	Long: `Adds a task to the currently active list unless a different
list is specified with the -l tag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("add %s\n", ListName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
