package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_GetClient(t *testing.T) {
// 	elastic := Elastic{"127.0.0.1", "9200"}
// 	elastic.GetClient()
// }

func Test_Map(t *testing.T) {
	var expectJSON = `{"settings":{"number_of_shards":1,"number_of_replicas":0},"mappings":{"ColdWater":{"_all":{"enabled":false},"properties":{"ID":{"type":"keyword"},"Timestamp":{"type":"keyword"}}}}}`
	mapping := GetMapping("ColdWater")
	properties := mapping["ColdWater"]

	properties.AddProperty("ID", GetTypeProperty("keyword")).
		AddProperty("Timestamp", GetTypeProperty("keyword"))

	index := GetIndex(1, 0, &mapping)
	indexJSON, _ := json.Marshal(index)
	assert.Equal(t, expectJSON, string(indexJSON), "should be equal.")
}
