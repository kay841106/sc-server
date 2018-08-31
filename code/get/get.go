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

type gwdata struct {
	GWID  string `json:"GWID" bson:"GWID"`
	MGWID string `json:"M_GWID" bson:"M_GWID"`
	// MMAC      string `json:"M_MAC" bson:"M_MAC"`
	// NUM       string `json:"NUM" bson:"NUM"`
	Place string `json:"Place" bson:"Place"`
	// Territory string `json:"Territory" bson:"Territory"`
	// Type      string `json:"Type" bson:"Type"`
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

func gogetgwdetail(w http.ResponseWriter, r *http.Request) {

	container := []gwstat{}
	container2 := []gwdata{}
	container3 := []string{}
	container4 := []gwdata{}
	container5 := []gwdata{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{}).All(&container)
	sess.DB(db).C(c_devices).Find(bson.M{}).Distinct("GWID", &container2)

	for _, each := range unique(container2) {
		// fmt.Print(each)
		for _, each2 := range container {

			if each2.GWID[0:7] != each.GWID[0:7] {
				fmt.Println(each2.GWID)
				// continue
				container3 = append(container3, each.MGWID)
			}
		}

	}

	// json.NewEncoder(w).Encode(container3)

	for _, each3 := range unique(container3) {
		sess.DB(db).C(c_devices).Find(bson.M{"M_GWID": each3}).Limit(1).All(&container4)
		for _, each := range container4 {
			container5 = append(container5, each)
		}

	}
	fmt.Println(container5)
	json.NewEncoder(w).Encode(container5)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/meter/lastreport", gogetlastreport).Methods("GET")
	router.HandleFunc("/meter/devices", gogetDevices).Methods("GET")
	router.HandleFunc("/meter/gwstat", gogetgwstat).Methods("GET")
	router.HandleFunc("/meter/gwdetail", gogetgwdetail).Methods("GET")

	log.Println(http.ListenAndServe(":8081", router))

}
