package db

import (
	"context"
	"errors"

	"dforcepro.com/util"
	search "gopkg.in/olivere/elastic.v5"
)

var (
	_client *search.Client
)

type Elastic struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (e Elastic) GetClient() *search.Client {
	if _client != nil {
		return _client
	}

	url := util.StrAppend("http://", e.Host, ":", e.Port)
	_client, err := search.NewClient(search.SetURL(url), search.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return _client
}

func (e Elastic) CreateIndex(index string, mapping string) (bool, error) {
	client := e.GetClient()
	exists, err := client.IndexExists(index).Do(context.Background())

	if err != nil {
		return false, err
	}

	// 不重覆建立 Index
	if exists {
		return true, nil
	}

	createIndex, err := client.CreateIndex(index).Body(mapping).Do(context.Background())
	if err != nil {
		return false, err
	}

	if !createIndex.Acknowledged {
		return false, errors.New("cant create index.")
	}

	return true, nil
}

func (e Elastic) Put(index string, eType string, id string, body string) (bool, error) {
	client := e.GetClient()
	_, err := client.Index().
		Index(index).
		Type(eType).
		Id(id).
		BodyJson(body).Do(context.Background())

	if err != nil {
		return false, err
	}

	// Not Support In v5
	// if !put.Created {
	// 	return false, fmt.Errorf("fail to create %s in index %s and type: %s.",
	// 		put.Id, put.Index, put.Type)
	// }

	_, err = client.Flush().Index(index).Do(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (e Elastic) Delete(index string, eType string, id string) (bool, error) {
	client := e.GetClient()
	res, err := client.Delete().Index(index).Type(eType).Id(id).Do(context.Background())
	return res.Found, err
}

type Index struct {
	Setting *_Setting               `json:"settings"`
	Mappins *map[string]*Properties `json:"mappings"`
}

type _Setting struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type Properties struct {
	All        map[string]interface{} `json:"_all"`
	Properties map[string]interface{} `json:"properties"`
}

type Property struct {
	Key   string
	Value string
}

func (p Properties) AddProperty(key string, propertyAry ...*Property) *Properties {
	inter := make(map[string]string)
	for _, prop := range propertyAry {
		inter[prop.Key] = prop.Value
	}
	p.Properties[key] = inter
	return &p
}

var typePropertyPool = make(map[string]*Property)

func GetTypeProperty(value string) *Property {
	property, ok := typePropertyPool[value]

	if !ok {

		property = &Property{"type", value}

		typePropertyPool[value] = property

	}
	return property
}

func GetIndex(shards int, replicas int, mappins *map[string]*Properties) *Index {
	return &Index{&_Setting{shards, replicas}, mappins}
}

func GetMapping(keys ...string) map[string]*Properties {
	mappings := make(map[string]*Properties)
	_all := make(map[string]interface{})
	_all["enabled"] = false
	for _, key := range keys {
		mappings[key] = &Properties{_all, make(map[string]interface{})}
	}
	return mappings
}
