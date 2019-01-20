package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var bannr = `
Program name : WEATHER GET API

maintainer   : avbee.lab@gmail.com

Date         : 17/01/2019

`

const (
	DBName         = "sc"
	CollectionName = "weather"
)

// Environment details
func getEnvVar() ENV {

	theenv := ENV{}
	theenv.Mongo = os.Getenv("MONGO_URL")
	if theenv.Mongo == "" {
		theenv.Mongo = "172.16.0.132:27017"
		// theenv.Mongo = "140.118.70.136:10003"
	}
	theenv.Database = os.Getenv("DB_AUTH")
	if theenv.Database == "" {
		theenv.Database = "admin"
	}
	theenv.Username = os.Getenv("USERNAME")
	if theenv.Username == "" {
		theenv.Username = "dontask"
	}
	theenv.Password = os.Getenv("PASS")
	if theenv.Password == "" {
		theenv.Password = "idontknow"
	}

	return theenv
}

type ENV struct {
	Port     string
	Mongo    string
	Database string
	Username string
	Password string
}

// Adapter function
type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {

	//return the Adapter
	return func(h http.Handler) http.Handler {
		// the adapter (when called)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//coppt the database session
			dbsession := db.Copy()
			defer dbsession.Close() // clean up

			context.Set(r, DBName, dbsession)

			h.ServeHTTP(w, r)
		})
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

type comment struct {
	Start string `json:"Start" bson:"Start"`
	Stop  string `json:"Stop" bson:"Stop"`
}

type queryInfo struct {
	Town         string    `json:"town" bson:"town"`
	Start        time.Time `json:"startTime" bson:"startTime"`
	Stop         time.Time `json:"endTime" bson:"endTime"`
	ElementValue string    `json:"ElementValue" bson:"ElementValue"`
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	// db := context.Get(r, "database").(*mgo.Session)

	var c comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, DBName).(*mgo.Session)

	var comments *comment
	var container []queryInfo
	if err := json.NewDecoder(r.Body).Decode(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// start := comments.Start
	// fmt.Println(start)

	if comments.Start != "" && comments.Stop != "" {
		start, e := time.ParseInLocation("2006-01-02T15", comments.Start, time.UTC)
		stop, er := time.ParseInLocation("2006-01-02T15", comments.Stop, time.UTC)

		if e != nil || er != nil {
			log.Fatal(e, er)
		}

		if err := db.DB(DBName).C(CollectionName).Find(bson.M{"startTime": bson.M{"$gt": start, "$lt": stop}}).Sort("-startTime").All(&container); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Println(container)
		if err := json.NewEncoder(w).Encode(container); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Println(r.Header, r.Host)
}

func main() {

	inpEnv := getEnvVar()
	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(inpEnv.Mongo, ",", -1),
		Database: inpEnv.Database,
		Username: inpEnv.Username,
		Password: inpEnv.Password,
		Timeout:  time.Second * 10,
	}

	db, err := mgo.DialWithInfo(dbInfo)

	if err != nil {
		log.Fatal("cannot dial mongo ", err)

	}
	defer db.Close()

	h := Adapt(http.HandlerFunc(handle), withDB(db))

	http.Handle("/weather", context.ClearHandler(h))

	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal(err)
	}
}
