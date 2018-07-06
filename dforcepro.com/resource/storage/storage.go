package storage

import (
	"errors"
	"io"
)

type storage interface {
	save(filePath string, file []byte) error
	saveByReader(fp string, reader io.Reader) error
	delete(filePath string) error
	get(filePath string) ([]byte, error)
	fileExist(fp string) bool
}

type FileStorage struct {
	Location string `yaml:"location"`
	Path     string `yaml:"path"`
}

func getStorage(fs *FileStorage) (storage, error) {
	switch fs.Location {
	case "hd":
		return &Hd{Path: fs.Path}, nil
	default:
		return nil, errors.New("not support location")
	}
}

func (fs *FileStorage) Save(filePath string, file []byte) error {
	sto, err := getStorage(fs)
	if err != nil {
		return err
	}
	err = sto.save(filePath, file)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) SaveByReader(filePath string, reader io.Reader) error {
	sto, err := getStorage(fs)
	if err != nil {
		return err
	}
	err = sto.saveByReader(filePath, reader)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Delete(filePath string) error {
	sto, err := getStorage(fs)
	if err != nil {
		return err
	}
	err = sto.delete(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Get(filePath string) ([]byte, error) {
	sto, err := getStorage(fs)
	if err != nil {
		return nil, err
	}
	return sto.get(filePath)
}

func (fs *FileStorage) FileExist(filePath string) (bool, error) {
	sto, err := getStorage(fs)
	if err != nil {
		return false, err
	}
	exist := sto.fileExist(filePath)

	return exist, nil
}
