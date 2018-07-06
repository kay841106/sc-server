package common

import (
	"encoding/json"

	"dforcepro.com/resource"
	"dforcepro.com/resource/db"
)

var (
	_di *resource.Di
)

const (
	Database = "Common"
)

func SetDI(c *resource.Di) {
	_di = c
}

func getRedis(redisdb int) *db.Redis {
	return (&_di.Redis).DB(redisdb)
}

func getMongo() db.Mongo {
	return _di.Mongodb
}

func toJSONByte(obj interface{}) []byte {
	jsonByte, _ := json.Marshal(obj)
	return jsonByte
}

func toJSONStr(obj interface{}) string {
	return string(toJSONByte(obj))
}
