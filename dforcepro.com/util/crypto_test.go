package util

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EnAndDecodeMap(t *testing.T) {

	data := map[string]interface{}{
		"test": "aa",
		"name": "bbb",
	}
	token, err := EncodeMap(&data)
	assert.Nil(t, err, "they should be nil")
	fmt.Println(token)

	recover, err := DecodeMap(token)
	assert.Nil(t, err, "they should be nil")
	fmt.Println(recover)
	val, ok := (*recover)["test"]
	assert.True(t, ok)
	assert.Equal(t, "aa", val)
}

func Test_ReadPem(t *testing.T) {
	size := 2048
	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, size)
	assert.Nil(t, err)
	err = CreatePrivateKeyPem("./key.pem", serverPrivateKey)
	assert.Nil(t, err)

	rsaKey, err := ReadPrivateKeyPem("./key.pem")
	assert.Nil(t, err)
	_, err = EncodePublicKey(rsaKey.Public())
	assert.Nil(t, err)

	err = os.Remove("./key.pem")
	assert.Nil(t, err)
}
