// package excom

// import (
// 	"time"

// 	"github.com/olivere/elastic"
// )

// const (
// 	indexName    = "applications"
// 	docType      = "log"
// 	appName      = "myApp"
// 	indexMapping = `{
// 		"mappings" : {
// 			"log": {
// 				"properties": {
// 					"app": { "type" : "string", "index" : "not_analyzed"},
// 					"message": {"type" : "string","index" : "not_analyzed"},
// 					"time": {"type":"date"}
// 								}
// 							}
// 						}
// 					}`
// )

// type Log struct {
// 	App     string    `json:"app"`
// 	Message string    `json:"message"`
// 	Time    time.Time `json:"time"`
// }

// func main() {
// 	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = createIndexWithLogsIfDoesNotExist(client)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = findAndPrintAppLogs(client)
// 	if err != nil {
// 		panic(err)
// 	}
// }

package main

type Ampas struct {
	ID int `json:"ID,omitempty"`
}

func main( x int, y int) int {

}

type aaa interface {
	
}