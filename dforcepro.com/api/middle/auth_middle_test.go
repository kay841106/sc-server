package middle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddAuthPath(t *testing.T) {
	AddAuthPath("/:POST", true)
	result := IsAuth("/", "POST")
	assert.True(t, result)
}
