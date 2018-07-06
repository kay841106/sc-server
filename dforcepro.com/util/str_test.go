package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StrAppend(t *testing.T) {
	var a = "hello"
	var b = "go"
	result := StrAppend(a, " ", b, "!!")
	assert.Equal(t, "hello go!!", result, "they should be equal")
}

func Test_MD5(t *testing.T) {
	// var a = "aaa"
	// var b = "bbbb"
	// result := MD5(a + b)
	// assert.Equal(t, 12, len(result), "they should be equal")

	//assert.Equal(t, "49fb00e00a86e7680adacfec3700aaaf", time.RFC1123, "they should be equal")
}
