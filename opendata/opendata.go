package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"
	// "/Users/avbee/go/src/sc-server/meter/get/auth"

	"github.com/gorilla/mux"

	// change due to high cpu using globalsign
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:10003"
	dbbackup = "140.118.122.103:27017"

	db           = "sc"
	c_lastreport = "lastreport"
	c_devices    = "devices"
	c_gwtstat    = "gw_status"

	c_hourly = "hour"
	c_daily  = "day"
	c_month  = "month"

	cpm    = "cpm"
	aemdra = "aemdra"
)

type session struct {
	theSess *mgo.Session
}

func (s *session) startSession() *session {
	return &session{s.theSess.Clone()}
}
func db_connect() *mgo.Session {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbpublic, ",", -1),
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

///////////////////////////////////

func (s *session) gopostqueryMonthly(w http.ResponseWriter, r *http.Request) {

	if s.checkDBStatus(); true {
		vars := mux.Vars(r)

		sutarto := vars["sutarto"]
		endo := vars["endo"]

		co, err := strconv.Atoi(endo)
		// headercontainer := postAgg{}

		switch meterType := vars["meterType"]; meterType {
		case "01":
			if sutarto != "" || err != nil {

				start, e := time.ParseInLocation("2006-01-02", sutarto, time.Local)
				if e != nil {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					break
				}
				a, b, c := start.Date()
				if a < 2018 && b < 11 && c < 01 {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					break
				}
				stop := start.Add(time.Hour * 24)
				if e != nil {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					fmt.Println("1")
					break
				}
				if co < 4 {

					// container := []getAgg{}
					var container interface{}
					sess := s.startSession().theSess
					defer sess.Close()

					Mongo := sess.DB(db).C(aemdra)
					// Mongo.Find(bson.M{"Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
					iter := Mongo.Find(bson.M{"Timestamp_Unix": bson.M{"$gte": start.Unix(), "$lte": stop.Unix()}}).Batch(100).Prefetch(0.2).Iter()

					for iter.Next(&container) {
						// fmt.Println(container)
						// time.Sleep(5 * time.Second)
						// fmt.Printf("Result: %v\n", result.Id)
						json.NewEncoder(w).Encode(container)
					}
					if err := iter.Close(); err != nil {
						fmt.Print(err)
					} else {
						// WARNING := "Time more than 36 months"
						json.NewEncoder(w).Encode(http.StatusBadRequest)
						fmt.Println("2")
					}
				}
			}
		case "02":
			if sutarto != "" || err != nil {

				start, e := time.ParseInLocation("2006-01-02", sutarto, time.Local)
				if e != nil {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					break
				}
				a, b, c := start.Date()
				if a < 2018 && b < 11 && c < 01 {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					break
				}
				stop := start.Add(time.Hour * 24)

				if e != nil {
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					fmt.Println(e)
					break
				}
				if co < 4 {

					// container := []getAgg{}
					var container interface{}
					sess := s.startSession().theSess
					defer sess.Close()

					Mongo := sess.DB(db).C(cpm)
					// Mongo.Find(bson.M{"Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
					iter := Mongo.Find(bson.M{"Timestamp_Unix": bson.M{"$gte": start.Unix(), "$lte": stop.Unix()}}).Batch(100).Prefetch(0.2).Iter()

					for iter.Next(&container) {
						// fmt.Println(container)
						// time.Sleep(5 * time.Second)
						// fmt.Printf("Result: %v\n", result.Id)
						json.NewEncoder(w).Encode(container)
					}
					if err := iter.Close(); err != nil {
						fmt.Print(err)
					}
					// fmt.Println(allan)
					// json.NewEncoder(w).Encode(container)
					// fmt.Println(diff)
				} else {
					// WARNING := "Time more than 36 months"
					json.NewEncoder(w).Encode(http.StatusBadRequest)
					fmt.Println("4")
				}
			}
		default:
			json.NewEncoder(w).Encode(http.StatusBadRequest)
			fmt.Println(meterType, "5")
			break
		}
		// json.NewDecoder(r.Body).Decode(&headercontainer)

	}
}

// `````
// SMART SPACE
// `````

func (s *session) checkDBStatus() bool {
	sess := s.startSession().theSess
	err := sess.Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		sess.Close()
		time.Sleep(5 * time.Second)
		sess.Refresh()
		err = sess.Ping()

	}
	fmt.Println("DB GOOD")
	return true
}

// ```
// MAIN
// ```

func main() {

	// auth.GenAuth()
	router := mux.NewRouter()
	sess := db_connect()

	v := session{sess}

	// router.HandleFunc("/v1/opendata/{type:meterType}/{start:sutarto}/{end:endo}", gopostqueryMonthly).Methods("GET")
	router.HandleFunc("/v1/opendata/{meterType}/{sutarto}", v.gopostqueryMonthly).Methods("GET")

	log.Println(http.ListenAndServe(":20000", router))

}
