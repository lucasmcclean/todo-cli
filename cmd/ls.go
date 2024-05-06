package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all tasks on todo list.",
	Long: `Lists all items on the currently active list unless a different
list is specified with the -l flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		listName, _ := cmd.Flags().GetString("list")
		fmt.Printf("ls %s\n", listName)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().StringP("list", "l", "_current", "The list which you would like to view.")
	viper.BindPFlag("list", lsCmd.Flags().Lookup("list"))
}
