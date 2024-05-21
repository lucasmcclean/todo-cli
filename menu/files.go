package menu

import (
	"os"
	"os/user"
)

func GetFileNames() (fileNames []string, err error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(dataDir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileNames, err = file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}

func getDataDir() (dataDir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	home := usr.HomeDir

	dataDir = os.Getenv("$XDG_DATA_HOME")
	if dataDir == "" {
		dataDir = "/.local/share"
	}
	dataDir += "/todo/"
	dataDir = home + dataDir
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0600)
	}
	return dataDir, nil
}
