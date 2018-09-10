package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

const (
	db_airbox = "airbox"
	dblocal   = "172.16.0.132:27017"
	dbpublic  = "140.118.70.136:10003"
	// c            = "testing"

	c_airboxraw = "airbox_raw"
)

var session *mgo.Session

type airbox struct {
	Key      string  `json:"Key" bson:"Key"`
	DevID    string  `json:"Device_ID" bson:"Device_ID"`
	Temp     float32 `json:"Temp" bson:"Temp"`
	Humidity int     `json:"Humidity" bson:"Humidity"`
	PM2_5    float32 `json:"PM2.5" bson:"PM2.5"`
	CO       float32 `json:"CO" bson:"CO"`
	CO2      float32 `json:"CO2" bson:"CO2"`
	Noise    int     `json:"Noise" bson:"Noise"`
}

type airboxSnd struct {
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	DevID          string        `json:"Device_ID" bson:"Device_ID"`
	Timestamp      time.Time     `json:"Timestamp" bson:"Timestamp"`
	Timestamp_Unix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	Temp           float32       `json:"Temp" bson:"Temp"`
	Humidity       int           `json:"Humidity" bson:"Humidity"`
	PM2_5          float32       `json:"PM2_5" bson:"PM2_5"`
	CO             float32       `json:"CO" bson:"CO"`
	CO2            float32       `json:"CO2" bson:"CO2"`
	Noise          int           `json:"Noise" bson:"Noise"`
}

func getObjectIDTwoArg(GWID string, macID string, timestamp int64) bson.ObjectId {
	var b [12]byte
	var sum [8]byte
	var c [4]byte
	// timestamp := time.Unix(LastReportTime, 0)
	// binary.BigEndian.PutUint32(b[:], uint32(timestamp))
	binary.BigEndian.PutUint32(c[:], uint32(timestamp))

	did := sum[:]
	gid := sum[:]

	hw := md5.New()
	hw.Write([]byte(GWID))
	copy(did, hw.Sum(nil))
	hw.Write([]byte(macID))
	copy(gid, hw.Sum(nil))
	// b[0] = c[:1]
	b[0] = c[0]
	b[1] = c[1]
	b[2] = c[2]
	b[3] = c[3]
	b[4] = did[4]
	b[5] = did[5]
	b[6] = did[6]
	b[7] = did[7]
	b[8] = gid[4]
	b[9] = gid[5]
	b[10] = gid[6]
	b[11] = gid[7]

	dst := hex.EncodeToString(b[:])
	theid := bson.ObjectIdHex(dst)

	fmt.Println(theid, uint32(timestamp))
	return theid

}

func airboxPost(w http.ResponseWriter, r *http.Request) {

	// sess := session.Clone()
	// defer sess.Close()

	// Mongo := sess.DB(db_airbox)

	container := airbox{}
	var containertemp interface{}
	json.NewDecoder(r.Body).Decode(&containertemp)

	// fmt.Println(r.Body.Read.(string))

	// if containertemp != container {
	// 	fmt.Println("asu")
	// }
	fmt.Print(containertemp)
	fmt.Println(reflect.DeepEqual(container, containertemp))
	fmt.Println(r.Header)

	// v := reflect.ValueOf(containertemp)
	s := reflect.ValueOf(containertemp).NumField()

	// for _, v := range s.MapKeys() {

	fmt.Println(s)
	// }
	// for _, v := range container {

	// }
	// for i := 0; i < v.NumField(); i++ {
	// 	fmt.Println(v.Field(i))
	// }
	// fmt.Println(v)

	// // var container interface{}
	// container2 := airboxSnd{
	// 	ID:             getObjectIDTwoArg(container.DevID, container.DevID, time.Now().Unix()),
	// 	Timestamp:      time.Now().UTC(),
	// 	Timestamp_Unix: time.Now().Unix(),
	// 	Temp:           container.Temp,
	// 	Humidity:       container.Humidity,
	// 	CO:             container.CO,
	// 	CO2:            container.CO2,
	// 	Noise:          container.Noise,
	// 	PM2_5:          container.PM2_5,
	// 	DevID:          container.DevID,
	// }

	// Mongo.C(c_airboxraw).Insert(container2)
	r.Body.Close()
	fmt.Print(container)

}

func init() {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbpublic, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 60,
	}
	session, _ = mgo.DialWithInfo(dbInfo)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/airbox_post", airboxPost).Methods("POST")

	log.Println(http.ListenAndServe(":8090", router))

}
