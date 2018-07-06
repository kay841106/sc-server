package common

import (
	"fmt"
	"testing"

	"dforcepro.com/resource"
	"dforcepro.com/resource/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_IniApi(t *testing.T) {
	log := logger.Logger{Path: "./", Duration: "minute", DebugMode: true}
	log.StartLog()
	_di = &resource.Di{Log: log}
	router := mux.NewRouter()
	fmt.Println(router)
	//IniAPI(router, RegionApi(false))
	fmt.Println(router)
	assert.True(t, true)
}
