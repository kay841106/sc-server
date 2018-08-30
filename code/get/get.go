package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

const (
	db           = "sc"
	c_lastreport = "lastreport"
	c_devices    = "devices"
	c_gwtstat    = "gw_status"
)

var session *mgo.Session

func init() {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN("140.118.70.136:10003", ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 2,
	}
	session, _ = mgo.DialWithInfo(dbInfo)
}

type CPMSnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string        `json:"GW_ID" bson:"GW_ID"`
}

func gogetlastreport(w http.ResponseWriter, r *http.Request) {

	container := []CPMSnd{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_lastreport)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	fmt.Println(container)
}

type devices struct {
	MACAddress string `json:"MACAddress" bson:"MACAddress"`
	DevID      int    `json:"DevID" bson:"DevID"`
	Floor      string `json:"Floor" bson:"Floor"`
	GWID       string `json:"GWID" bson:"GWID"`
	MGWID      string `json:"M_GWID" bson:"M_GWID"`
	MMAC       string `json:"M_MAC" bson:"M_MAC"`
	NUM        string `json:"NUM" bson:"NUM"`
	Place      string `json:"Place" bson:"Place"`
	Territory  string `json:"Territory" bson:"Territory"`
	Type       string `json:"Type" bson:"Type"`
}

func gogetDevices(w http.ResponseWriter, r *http.Request) {

	container := []devices{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_devices)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	fmt.Println(container)
}

type gwstat struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
}

func gogetgwstat(w http.ResponseWriter, r *http.Request) {

	container := []gwstat{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	fmt.Println(container)
}
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/meter/lastreport", gogetlastreport).Methods("GET")
	router.HandleFunc("/meter/devices", gogetDevices).Methods("GET")
	router.HandleFunc("/meter/gwstat", gogetgwstat).Methods("GET")

	log.Println(http.ListenAndServe(":8081", router))

}
