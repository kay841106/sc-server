package logger

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Aa string `json:"AA,omitempty"`
	Bb string `json:"BB,omitempty"`
}

func getLogger(debugMode bool) Logger {
	logPath, _ := filepath.Abs("./")
	logger := Logger{logPath, "minute", debugMode}
	logger.StartLog()
	return logger
}

func Test_Log(t *testing.T) {
	logger := getLogger(true)
	logger.Info("myInfo")
	logger.Debug("myDebug")
	logger.Err("myErr")
	logger.Warn("myWarn")
}

func Test_WriteFile(t *testing.T) {
	logger := getLogger(true)
	testData := TestData{"test1", "test2"}
	b, err := json.Marshal(testData)
	if err != nil {
		logger.Err(err.Error())
	}
	file, err := logger.WriteFile("document/ddaasd.json", b)
	exist, _ := pathExist(file)
	assert.True(t, exist)
	dir := filepath.Dir(file)
	removeContents(dir)
}

func Test_LogToFile(t *testing.T) {
	logger := getLogger(false)
	logger.Info("myInfo")
	logger.Debug("myDebug")
	logger.Err("myErr")
	logger.Warn("myWarn")

	// exist, _ := pathExist(logFile)
	// assert.True(t, exist)
	// err = os.RemoveAll(logFile)
	// if err != nil {
	// 	assert.True(t, false, err.Error())
	// }
}

func Test_GetPattern(t *testing.T) {
	pattern1, time1 := getPatternAndDuration(time.Hour)
	assert.Equal(t, "%Y%m%d%H", pattern1, "shoult be same")
	assert.Equal(t, time.Hour, time1, "shoult be same")
	pattern2, time2 := getPatternAndDuration(time.Minute)
	assert.Equal(t, "%Y%m%d%H%M", pattern2, "shoult be same")
	assert.Equal(t, time.Minute, time2, "shoult be same")
	pattern3, time3 := getPatternAndDuration(time.Microsecond)
	assert.Equal(t, "%Y%m%d", pattern3, "shoult be same")
	assert.Equal(t, time.Hour*24, time3, "shoult be same")
}

func Test_StartLog(t *testing.T) {
	logger := getLogger(true)
	logger.Info("StartLog")
}
