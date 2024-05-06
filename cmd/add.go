package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to the list.",
	Long: `Adds a task to the currently active list unless a different
list is specified with the -l tag.`,
	Run: func(cmd *cobra.Command, args []string) {
		listName, _ := cmd.Flags().GetString("list")
		fmt.Printf("add %s\n", listName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("list", "l", "_current", "The list which you would like to view.")
	viper.BindPFlag("list", addCmd.Flags().Lookup("list"))
}
