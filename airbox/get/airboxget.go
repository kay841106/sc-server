package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	db_airbox    = "airbox"
	dblocal      = "172.16.0.132:27017"
	dbpublic     = "140.118.70.136:10003"
	dbleoass     = "140.118.123.95:27017"
	c_airboxlast = "airbox_lastreport"
	c_airboxraw  = "airbox_raw"
	c_airboxhour = "airbox_hour"
)

type session struct {
	theSess *mgo.Session
}

func (s *session) startSession() *session {
	return &session{s.theSess.Clone()}
}

var bannr = `
Program name : airboxpost

maintainer   : avbee.lab@gmail.com

Date         : November, 30 2018

`

type airboxSnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          *string       `json:"GW_ID" bson:"GW_ID"`
	GET11         *float64      `json:"Temp" bson:"Temp"`
	GET12         *float64      `json:"Humidity" bson:"Humidity"`
	GET13         *float64      `json:"PM2_5" bson:"PM2_5"`
	GET14         *float64      `json:"CO" bson:"CO"`
	GET15         *float64      `json:"CO2" bson:"CO2"`
	GET16         *float64      `json:"Noise" bson:"Noise"`
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

func (s *session) airboxLastReport(w http.ResponseWriter, r *http.Request) {
	sess := s.startSession().theSess
	err := sess.Ping()

	container2 := []airboxSnd{}
	// sess := session.Clone()

	Mongo := sess.DB(db_airbox)
	defer sess.Close()
	// container := airboxRcv{}

	// json.NewDecoder(r.Body).Decode(&containertemp)

	// fmt.Print(containertemp)
	// fmt.Println(reflect.DeepEqual(container, containertemp))
	// fmt.Println(r.Header)

	err = Mongo.C(c_airboxlast).Find(bson.M{}).All(&container2)

	if err != nil {
		log.Println(err)
	}
	r.Body.Close()
	// fmt.Print(container)
	json.NewEncoder(w).Encode(&container2)
}

type headers struct {
	Start      *string `json:"Start" bson:"Start"`
	Stop       *string `json:"Stop" bson:"Stop"`
	MACAddress *string `json:"MAC_Address" bson:"MAC_Address"`
}

func (s *session) airboxHour(w http.ResponseWriter, r *http.Request) {

	sess := s.startSession().theSess

	headercontainer := headers{}

	// sess := session.Clone()

	Mongo := sess.DB(db_airbox).C(c_airboxhour)
	defer sess.Close()

	json.NewDecoder(r.Body).Decode(&headercontainer)
	if &headercontainer.MACAddress != nil {
		if &headercontainer.Start != nil && &headercontainer.Stop != nil {

			start, e := time.ParseInLocation("2006-01-02T15", *headercontainer.Start, time.Local)
			stop, er := time.ParseInLocation("2006-01-02T15", *headercontainer.Stop, time.Local)

			if e != nil || er != nil {
				log.Println(e, er)
			}

			// fmt.Println(headercontainer)
			container := []airboxSnd{}

			Mongo.Find(bson.M{"MAC_Address": headercontainer.MACAddress, "Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
			json.NewEncoder(w).Encode(container)
			// fmt.Println(container)
		}
	}
}
func SetTimeStampForHour(theTime time.Time) time.Time {
	year, month, day := theTime.Date()
	return time.Date(year, month, day, 23, 59, 59, 59, time.UTC)
}

func (s *session) airboxOneDayData(w http.ResponseWriter, r *http.Request) {
	sess := s.startSession().theSess

	// sess := session.Clone()

	Mongo := sess.DB(db_airbox).C(c_airboxraw)
	defer sess.Close()

	headercontainer := headers{}

	json.NewDecoder(r.Body).Decode(&headercontainer)
	fmt.Println(&headercontainer)
	if &headercontainer.MACAddress != nil {
		if &headercontainer.Start != nil && &headercontainer.Stop != nil {

			start, e := time.ParseInLocation("2006-01-02", *headercontainer.Start, time.Local)
			stop := SetTimeStampForHour(start)
			// stop, er := time.ParseInLocation("2006-01-02", *headercontainer.Stop, time.Local)
			diff := stop.Sub(start)

			if e != nil {
				log.Println(e)

			}
			if diff <= time.Hour*24 {

				container := []airboxSnd{}

				Mongo.Find(bson.M{"MAC_Address": *headercontainer.MACAddress, "Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
				json.NewEncoder(w).Encode(container)
				fmt.Println(diff)
			} else {
				// WARNING := "Time more than 24 hours"
				// json.NewEncoder(w).Encode(WARNING)
			}
		}
	}
}

func db_connect() *mgo.Session {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dblocal, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 10,
	}

	sess, err := mgo.DialWithInfo(dbInfo)
	if err != nil {
		os.Exit(1)
	}
	return sess
}

func main() {
	// db_connect()
	fmt.Println(bannr)
	sess := db_connect()

	v := session{sess}
	router := mux.NewRouter()

	router.HandleFunc("/airbox/last", v.airboxLastReport).Methods("GET")
	router.HandleFunc("/airbox/oneday", v.airboxOneDayData).Methods("POST")
	router.HandleFunc("/airbox/hourly", v.airboxHour).Methods("POST")

	log.Println(http.ListenAndServe(":8090", router))
}
