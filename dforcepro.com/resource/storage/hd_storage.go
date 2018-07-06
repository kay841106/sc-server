package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Hd struct {
	Path string
}

func (hd *Hd) getAbsFilePath(filePath string) string {
	absFilePath, _ := filepath.Abs(hd.Path + filePath)
	return absFilePath
}

func (hd *Hd) save(fp string, file []byte) error {
	absFilePath := hd.getAbsFilePath(fp)
	err := hd.mkdir(absFilePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(absFilePath, file, 0644)
}

func (hd *Hd) saveByReader(fp string, reader io.Reader) error {
	absFilePath := hd.getAbsFilePath(fp)
	f, err := os.OpenFile(absFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	return nil
}

func (hd *Hd) delete(filePath string) error {
	absFilePath := hd.getAbsFilePath(filePath)
	exist, err := fileExist(absFilePath)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("file not exist: " + absFilePath)
	}
	return os.Remove(absFilePath)
}

func (hd *Hd) get(fp string) ([]byte, error) {
	absFilePath := hd.getAbsFilePath(fp)
	return ioutil.ReadFile(absFilePath)
}

func (hd *Hd) mkdir(absPath string) error {
	dir := filepath.Dir(absPath)
	exist, _ := fileExist(dir)

	if !exist {
		err := os.MkdirAll(dir, 0766)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hd *Hd) fileExist(fp string) bool {
	absFilePath := hd.getAbsFilePath(fp)
	exist, _ := fileExist(absFilePath)
	return exist
}

func fileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
