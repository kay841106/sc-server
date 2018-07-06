package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_CreateImage(t *testing.T) {
	fileName := "/data/image/lzw/estimate/111_222_start.jpg"
	fmt.Println(filepath.Dir(fileName))
}
