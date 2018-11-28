package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	db                     = "SC"
	collectionremoveDevice = ""
	// collection
	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:27017"
)

type session struct {
	theSess *mgo.Session
}

func (s *session) startSession() *session {
	return &session{s.theSess.Clone()}
}

///////////////////////////

type RegisterDevID struct {
	MACAddress string `json:"MACAddress" bson:"MACAddress"`
	GWID       string `json:"GWID" bson:"GWID:" `
	User       string `json:"User" bson:"User"`
}

func (s *session) rmDevice(w http.ResponseWriter, req *http.Request) {
	container := RegisterDevID{}
	mongo := s.startSession().theSess
	defer mongo.Close()

	err := mongo.DB(db).C(collectionremoveDevice).Remove(bson.M{"MACAddress": container.MACAddress})

	if err != nil {
		log.Println(err)
	}

	err = mongo.DB(db).C(collectionremoveDevice).Insert(bson.M{"MACAddress": container.MACAddress})

	if err != nil {
		log.Println(err)
	}

}

func (s *session) addDevice(w http.ResponseWriter, req *http.Request) {
	container := RegisterDevID{}

	json.NewDecoder(req.Body).Decode(&container)

	mongo := s.startSession().theSess

	err := mongo.DB(db).C(collection).Insert(bson.M{"MACAddress": container.MACAddress, "GWID": container.GWID})

	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(container)
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
	sess := db_connect()

	v := session{sess}

	router := mux.NewRouter()
	router.HandleFunc("/meter/aemdra", v.rmDevice).Methods("POST")
	router.HandleFunc("/meter/cpm", v.addDevice).Methods("POST")

	log.Println(http.ListenAndServe(":8082", router))
}
