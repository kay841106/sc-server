package db

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

func _Test_EnsureIndexKey(t *testing.T) {
	mongo := Mongo{"127.0.0.1", "27017", "peter", "pass", "", ""}
	err := mongo.DB("ytz").C("Case").EnsureIndexKey("devices")
	if err != nil {
		fmt.Println(err.Error())
	}
}

// func Test_Mongo_Insert(t *testing.T) {
// 	mongo := Mongo{"127.0.0.1", "27017", "Peter", "Pass", "", ""}
// defer mongo.Close()
// 	result := mongo.DB("testdb").C("people").Insert(&Person{Name: "Ale", Phone: "+55 53 1234 4321", Timestamp: time.Now()},
// 		&Person{Name: "Cla", Phone: "+66 33 1234 5678", Timestamp: time.Now()})
// 	assert.Nil(t, result, "must nil")
// }

// func Test_Mongo_FindWithLimit(t *testing.T) {
// 	var results []Person
// 	mongo := Mongo{"127.0.0.1", "27017", "Peter", "Pass", "", ""}
// 	defer mongo.Close()
// 	err := mongo.DB("testdb").C("people").Find(bson.M{}).Sort("-timestamp").Limit(1).All(&results)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(results)
// }
