package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a tool for managing tasks from the terminal",
	Long: `Todo is a command line utility for generating and managing
todo lists. It utilizes local text files which can be
accessed and modified outside of the program.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output. (default false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}
