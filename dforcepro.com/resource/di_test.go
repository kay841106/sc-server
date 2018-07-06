package resource

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Connect(t *testing.T) {
	_, err := GetDI()
	assert.NotNil(t, err)

	filename, _ := filepath.Abs("./test_config.yml")
	IniConf(filename)
	conf, err := GetDI()
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", conf.Mongodb.Host, "they should be equal")
	assert.Equal(t, "9080", conf.APIConf.Port)
	fmt.Println(conf.APIConf.AllowedOrigins)
}
