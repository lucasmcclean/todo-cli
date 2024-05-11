package utils

import (
	"os"
	"os/user"
)

func MustGetDataDir() string {
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
	return dir
}
