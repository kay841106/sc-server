package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"

	// change due to high cpu using globalsign
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

const (
	db_airbox = "airbox"
	dblocal   = "172.16.0.132:27017"
	dbpublic  = "140.118.70.136:10003"
	dbleoass  = "140.118.123.95:27017"
	// c            = "testing"

	c_airboxraw = "airbox_raw"
)

var session *mgo.Session

var bannr = `
Program name : airboxpost

maintainer   : avbee.lab@gmail.com

Date         : October, 01 2018

`

// type airbox struct {
// 	Key      string  `json:"Key" bson:"Key"`
// 	DevID    string  `json:"Device_ID" bson:"Device_ID"`
// 	Temp     float32 `json:"Temp" bson:"Temp"`
// 	Humidity int     `json:"Humidity" bson:"Humidity"`
// 	PM2_5    float32 `json:"PM2.5" bson:"PM2.5"`
// 	CO       float32 `json:"CO" bson:"CO"`
// 	CO2      float32 `json:"CO2" bson:"CO2"`
// 	Noise    int     `json:"Noise" bson:"Noise"`
// }

type airboxSnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          *string       `json:"GW_ID" bson:"GW_ID"`
	CPURate       *float64      `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   *int          `json:"Storage_rate" bson:"Storage_rate"`
	GET11         *float64      `json:"Temp" bson:"Temp"`
	GET12         *float64      `json:"Humidity" bson:"Humidity"`
	GET13         *float64      `json:"PM2_5" bson:"PM2_5"`
	GET14         *float64      `json:"CO" bson:"CO"`
	GET15         *float64      `json:"CO2" bson:"CO2"`
	GET16         *float64      `json:"Noise" bson:"Noise"`
}

type airboxRcv struct {
	Timestamp     string   `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64    `json:"Timestamp_Unix"`
	MACAddress    string   `json:"MAC_Address" bson:"MAC_Address"`
	GWID          *string  `json:"GW_ID" bson:"GW_ID"`
	CPURate       *float64 `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   *int     `json:"Storage_rate" bson:"Storage_rate"`
	GET11         *float64 `json:"GET_1_1" bson:"GET_1_1"`
	GET12         *float64 `json:"GET_1_2" bson:"GET_1_2"`
	GET13         *float64 `json:"GET_1_3" bson:"GET_1_3"`
	GET14         *float64 `json:"GET_1_4" bson:"GET_1_4"`
	GET15         *float64 `json:"GET_1_5" bson:"GET_1_5"`
	GET16         *float64 `json:"GET_1_6" bson:"GET_1_6"`
}

// type airboxSnd struct {
// 	ID             bson.ObjectId `json:"_id" bson:"_id"`
// 	DevID          string        `json:"Device_ID" bson:"Device_ID"`
// 	Timestamp      time.Time     `json:"Timestamp" bson:"Timestamp"`
// 	Timestamp_Unix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
// 	Temp           float32       `json:"Temp" bson:"Temp"`
// 	Humidity       int           `json:"Humidity" bson:"Humidity"`
// 	PM2_5          float32       `json:"PM2_5" bson:"PM2_5"`
// 	CO             float32       `json:"CO" bson:"CO"`
// 	CO2            float32       `json:"CO2" bson:"CO2"`
// 	Noise          int           `json:"Noise" bson:"Noise"`
// }

// type thekey struct {
// 	Key string `json:"Key" bson:"Key"`
// }

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

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbpublic, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Minute * 15,
	}
	sess, _ := mgo.DialWithInfo(dbInfo)

	// sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db_airbox)

	// container := airboxRcv{}
	containertemp := airboxRcv{}
	json.NewDecoder(r.Body).Decode(&containertemp)

	fmt.Print(containertemp)
	// fmt.Println(reflect.DeepEqual(container, containertemp))
	fmt.Println(r.Header)

	container2 := airboxSnd{
		ID:            getObjectIDTwoArg(containertemp.MACAddress, containertemp.MACAddress, containertemp.TimestampUnix),
		Timestamp:     time.Unix(containertemp.TimestampUnix, 0).UTC(),
		TimestampUnix: containertemp.TimestampUnix,
		MACAddress:    containertemp.MACAddress,
		GWID:          containertemp.GWID,
		CPURate:       containertemp.CPURate,
		StorageRate:   containertemp.StorageRate,
		GET11:         containertemp.GET11,
		GET12:         containertemp.GET12,
		GET13:         containertemp.GET13,
		GET14:         containertemp.GET14,
		GET15:         containertemp.GET15,
		GET16:         containertemp.GET16,
	}

	err := Mongo.C(c_airboxraw).Insert(container2)
	if err != nil {
		log.Println(err)
	}
	r.Body.Close()
	// fmt.Print(container)
	json.NewEncoder(w).Encode(&container2)
}

// }

// func db_connect() *Session {

// 	dbInfo := &mgo.DialInfo{
// 		Addrs:    strings.SplitN(dblocal, ",", -1),
// 		Database: "admin",
// 		Username: "dontask",
// 		Password: "idontknow",
// 		Timeout:  time.Second * 120,
// 	}
// 	session, _ = mgo.DialWithInfo(dbInfo)
// 	return session
// }

func main() {
	// db_connect()
	fmt.Println(bannr)
	router := mux.NewRouter()
	router.HandleFunc("/airbox/rawdata", airboxPost).Methods("POST")

	log.Println(http.ListenAndServe(":8090", router))
}
