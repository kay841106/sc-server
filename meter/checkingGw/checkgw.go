package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	db           = "sc"
	c_lastreport = "lastreport"
	c_devices    = "devices"
	c_gwtstat    = "gw_status"

	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:10003"
)

var session *mgo.Session

func init() {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbpublic, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 2,
	}
	session, _ = mgo.DialWithInfo(dbInfo)
}

type gwstat struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
	Place         string    `json:"Place" bson:"Place"`
	MGWID         string    `json:"M_GWID" bson:"M_GWID"`
}

type gwstat2 struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GWID" bson:"GWID"`
	Place         string    `json:"Place" bson:"Place"`
	MGWID         string    `json:"M_GWID" bson:"M_GWID"`
}

type gwdata struct {
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

type lastreport struct {
	MACAddress    string    `json:"MACAddress" bson:"MACAddress"`
	DevID         int       `json:"DevID" bson:"DevID"`
	Floor         string    `json:"Floor" bson:"Floor"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
	MGWID         string    `json:"M_GWID" bson:"M_GWID"`
	MMAC          string    `json:"M_MAC" bson:"M_MAC"`
	NUM           string    `json:"NUM" bson:"NUM"`
	Place         string    `json:"Place" bson:"Place"`
	Territory     string    `json:"Territory" bson:"Territory"`
	Type          string    `json:"Type" bson:"Type"`
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
}

type tmp struct {
	GWID string `json:"GWID" bson:"GWID"`
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

func gwupload() {

	container := gwstat{}
	container2 := gwstat2{}
	container3 := []string{}
	// container4 := []gwdata{}
	// container5 := []gwdata{}
	sess := session.Clone()
	defer sess.Close()

	// Mongo := sess.DB(db).C(c_gwtstat)
	// Mongo.Find(bson.M{}).All(&container)
	// sess.DB(db).C(c_devices).Find(bson.M{}).All(&container2)
	sess.DB(db).C(c_devices).Find(bson.M{}).Distinct("M_GWID", &container3)
	for _, each := range container3 {
		sess.DB(db).C(c_devices).Find(bson.M{"M_GWID": each}).Limit(1).One(&container2)
		container.GWID = container2.GWID[0:8]
		container.MGWID = container2.MGWID
		container.Place = container2.Place
		sess.DB(db).C(c_gwtstat).Insert(container)
		fmt.Print(container2)
	}
}

func lastreportupload() {

	container := tmp{}
	container2 := lastreport{}
	container3 := []string{}
	// container4 := []gwdata{}
	// container5 := []gwdata{}
	sess := session.Clone()
	defer sess.Close()

	// Mongo := sess.DB(db).C(c_gwtstat)
	// Mongo.Find(bson.M{}).All(&container)
	// sess.DB(db).C(c_devices).Find(bson.M{}).All(&container2)
	sess.DB(db).C(c_devices).Find(bson.M{}).Distinct("MACAddress", &container3)
	for _, each := range container3 {
		sess.DB(db).C(c_devices).Find(bson.M{"MACAddress": each}).One(&container2)
		sess.DB(db).C(c_devices).Find(bson.M{"MACAddress": each}).One(&container)
		container2.GWID = container.GWID
		// container.MGWID = container2.MGWID
		// container.Place = container2.Place
		sess.DB(db).C(c_lastreport).Upsert(bson.M{"MAC_Address": container2.MACAddress}, container2)
		fmt.Print(container2)
	}
}

func gogetgwstat() {

	container := []gwstat{}
	container2 := []gwdata{}
	container3 := []string{}
	container4 := []gwdata{}
	container5 := []gwdata{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{}).All(&container)
	sess.DB(db).C(c_devices).Find(bson.M{}).All(&container2)

	for _, each := range container2 {
		// fmt.Print(each)
		for _, each2 := range container {

			if each.GWID[7] == each2.GWID[7] {
				container3 = append(container3, each.MGWID)

			}
		}

	}

	// json.NewEncoder(w).Encode(container3)

	for _, each3 := range unique(container3) {
		sess.DB(db).C(c_devices).Find(bson.M{"M_GWID": each3}).All(&container4)
		for _, each := range container4 {
			container5 = append(container5, each)
		}

	}
	fmt.Println(container5)
}

func main() {
	lastreportupload()
}
