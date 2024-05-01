package commands

import (
	"flag"
	"fmt"
	"os"
)

var (
	currentVersion = "0.0.0"
	buildInfo = "Not functional, under development"
	short = false
)

func NewVersionCommand() *Command {
	cmd := &Command {
		flags: flag.NewFlagSet("version", flag.ExitOnError),
		Execute: version,
	}
	cmd.flags.BoolVar(&short, "s", false, "")
	return cmd
}

func version(cmd *Command, args []string) {
	if short {
		fmt.Printf("Version: %s", currentVersion)
	} else {
		fmt.Printf("Version: %s\nBuild: %s", currentVersion, buildInfo)
	}
	os.Exit(0)
}