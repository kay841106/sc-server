package util

import "dforcepro.com/resource/logger"

var (
	log logger.Logger
)

func SetLog(logger logger.Logger) {
	log = logger
}
