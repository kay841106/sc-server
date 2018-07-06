package storage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getHdStorage() *FileStorage {
	return &FileStorage{Location: "hd", Path: "./file/"}
}

func Test_HD_SaveReader(t *testing.T) {
	r := strings.NewReader("some io.Reader stream to be read\n")
	fs := getHdStorage()
	err := fs.SaveByReader("test2.txt", r)
	assert.Nil(t, err)
	exist, _ := fs.FileExist("test2.txt")
	assert.True(t, exist)
	err = fs.Delete("test2.txt")
	assert.Nil(t, err)
}

func Test_HD_Save(t *testing.T) {
	fs := getHdStorage()
	err := fs.Save("test.txt", []byte("test bb"))
	assert.Nil(t, err)
	if err != nil {
		fmt.Println(err.Error())
	}
	exist, err := fs.FileExist("test.txt")
	assert.True(t, exist)

	exist, err = fs.FileExist("test1.txt")
	assert.False(t, exist)
}

func Test_HD_Get(t *testing.T) {
	fs := getHdStorage()
	_, err := fs.Get("test.txt")
	assert.Nil(t, err)

	_, err = fs.Get("test1.txt")
	assert.NotNil(t, err)
}

func Test_HD_Del(t *testing.T) {
	fs := getHdStorage()
	err := fs.Delete("test.txt")
	assert.Nil(t, err)
	if err != nil {
		fmt.Println(err.Error())
	}
	exist, err := fs.FileExist("test.txt")
	assert.False(t, exist)
}
