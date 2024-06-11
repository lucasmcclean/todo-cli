package menu

import (
	"bufio"
	"io"
	"os"
	"os/user"
)

func GetDataFileNames() (fileNames []string, err error) {
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

func RemoveDataFile(fileName string) error {
	dataDir, err := getDataDir()
	if err != nil {
		return err
	}
	err = os.Remove(dataDir + fileName)
	return err
}

func OpenDataFile(fileName string, create bool) (file *os.File, err error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}
	if create {
		file, err = os.OpenFile(dataDir+fileName, os.O_CREATE|os.O_RDWR, 0600)
	} else {
		file, err = os.OpenFile(dataDir+fileName, os.O_RDWR, 0600)
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFileLength(file os.File) (numTasks int) {
	file.Seek(0, io.SeekStart)
	reader := bufio.NewReader(&file)
	numTasks = 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		numTasks++
	}
	return numTasks
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
