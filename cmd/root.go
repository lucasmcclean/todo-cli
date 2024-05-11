package cmd

import (
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Verbose   bool
	ListName  string
	Directory string
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a tool for managing tasks from the terminal",
	Long: `Todo is a command line utility for generating and managing
todo lists.
A call to the root command lists the tasks from your active
todo list unless a seperate path is specified with -l.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := usr.HomeDir

	dir := os.Getenv("$XDG_DATA_HOME")
	if dir == "" {
		dir = "/.local/share"
	}
	dir += "/todo/"
	dir = home + dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0600)
	}

	rootCmd.PersistentFlags().StringVarP(&Directory, "directory", "d", dir, "specify directory to store lists")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display more verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringVarP(&ListName, "list", "l", "_current", "specifiy target list")
	viper.BindPFlag("list", rootCmd.PersistentFlags().Lookup("list"))
}
