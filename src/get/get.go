package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

const (
	db           = "sc"
	c_lastreport = "lastreport"
	c_devices    = "devices"
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
	CPURate       float64       `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   int           `json:"Storage_rate" bson:"Storage_rate"`
	GET11         float64       `json:"ae.tot" bson:"ae.tot"`
	GET12         float64       `json:"wire" bson:"wire"`
	GET13         float64       `json:"freq" bson:"freq"`
	GET14         float64       `json:"ua" bson:"ua"`
	GET15         float64       `json:"ub" bson:"ub"`
	GET16         float64       `json:"uc" bson:"uc"`
	GET17         float64       `json:"u.avg" bson:"u.avg"`
	GET18         float64       `json:"uab" bson:"uab"`
	GET19         float64       `json:"ubc" bson:"ubc"`
	GET110        float64       `json:"uca" bson:"uca"`
	GET111        float64       `json:"uln.avg" bson:"uln.avg"`
	GET112        float64       `json:"ia" bson:"ia"`
	GET113        float64       `json:"ib" bson:"ib"`
	GET114        float64       `json:"ic" bson:"ic"`
	GET115        float64       `json:"i.avg" bson:"i.avg"`
	GET116        float64       `json:"pa" bson:"pa"`
	GET117        float64       `json:"pb" bson:"pb"`
	GET118        float64       `json:"pc" bson:"pc"`
	GET119        float64       `json:"p.sum" bson:"p.sum"`
	GET120        float64       `json:"qa" bson:"qa"`
	GET121        float64       `json:"qb" bson:"qb"`
	GET122        float64       `json:"qc" bson:"qc"`
	GET123        float64       `json:"q.sum" bson:"q.sum"`
	GET124        float64       `json:"sa" bson:"sa"`
	GET125        float64       `json:"sb" bson:"sb"`
	GET126        float64       `json:"sc" bson:"sc"`
	GET127        float64       `json:"s.sum" bson:"s.sum"`
	GET128        float64       `json:"pfa" bson:"pfa"`
	GET129        float64       `json:"pfb" bson:"pfb"`
}

func goget(w http.ResponseWriter, r *http.Request) {

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

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/meter/test", goget).Methods("GET")
	router.HandleFunc("/meter/devices", gogetDevices).Methods("GET")

	log.Println(http.ListenAndServe(":8081", router))

}
