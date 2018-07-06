package util

import (
	"os"
	"path/filepath"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CreateIfNotExist(absPath string) error {
	dir := filepath.Dir(absPath)
	exist, _ := FileExists(dir)

	if !exist {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
