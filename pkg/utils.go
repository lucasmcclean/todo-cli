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

func IncludesInt(arr []int, x int) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == x {
			return true
		}
	}
	return false
}

func RemoveInt(arr *[]int, x int) bool {
	for i := 0; i < len(*arr); i++ {
		if (*arr)[i] == x {
			*arr = append((*arr)[:i], (*arr)[i+1:]...)
			return true
		}
	}
	return false
}
