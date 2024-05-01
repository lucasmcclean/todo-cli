package main

import (
	"fmt"
	"os"

	"github.com/ljmcclean/todo-cli/pkg/commands"
)

func main() {
	var cmd *commands.Command
	
	if len(os.Args) == 1 {
		commands.Default()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "version":
		cmd = commands.NewVersionCommand()
	default:
		fmt.Printf("%s is not a 'todo' command\n", os.Args[1])
	}

	cmd.Init(os.Args[2:])
	cmd.Run()
}