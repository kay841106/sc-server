package util

import (
	"image"
	"io"
	"os"
)

const (
	ImageRootPath = "/data/image/"
)

func CreateImage(fileName string, reader io.Reader) error {
	absPath := StrAppend(ImageRootPath, fileName)
	err := CreateIfNotExist(absPath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, reader)

	return nil
}

func IsImageExist(imagePath string) bool {
	exist, _ := FileExists(ImageRootPath + imagePath)
	return exist
}

func GetImage(imagePath string) *image.Image {
	if !IsImageExist(imagePath) {
		return nil
	}
	imgfile, err := os.Open(ImageRootPath + imagePath)

	if err != nil {
		return nil
	}

	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	return &img
}
