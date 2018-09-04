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

type gwstat struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
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
	gogetgwstat()
}
