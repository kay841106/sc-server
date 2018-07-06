package flux

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"sc.dforcepro.com/meter"

	"dforcepro.com/api"
	"dforcepro.com/resource"
	client "github.com/influxdata/influxdb/client/v2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	myDB     = "meter"
	username = "ntusttest"
	password = "123"
)

var _di *resource.Di

func SetDI(c *resource.Di) {
	_di = c
}

type Influx bool

//Enable the API
func (sca Influx) Enable() bool {
	return bool(sca)
}

//GetAPIs router
func (sca Influx) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/tesss", Next: sca.testes, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/testos", Next: sca.testos, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/we", Next: sca.wewewe, Method: "GET", Auth: false},
	}

}

func (sca Influx) testos(w http.ResponseWriter, req *http.Request) {
	var a []interface{}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"140.118.70.136:27017"},
		Database: "admin",
		Username: "root",
		Password: "123",
	}
	tick := time.Now()
	tock := tick.Add(time.Hour * -24)
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	thepipe := pipeTest(tick, tock)
	session.DB("Bimo_test").C("SC01_displayData_TR").Pipe(thepipe).All(&a)
	// session.DB("Bimo_test").C("SC01_displayData_TR").Find(bson.M{"lastReportTime": bson.M{"$lte": tick, "$gt": tock}}).Select(bson.M{"Pwr_Demand": 1, "lastReportTime": 1, "_id": 0}).All(&a)
	// fmt.Println(a)
	json.NewEncoder(w).Encode(&a)
}

func (sca Influx) testes(w http.ResponseWriter, req *http.Request) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://140.118.123.109:8086",
		Username: username,
		Password: password,
	})
	tick := time.Now()
	tock := tick.Add(time.Hour * -24)
	tock2 := tock.Unix()
	a := strconv.FormatInt(tock2*1000*1000*1000, 10)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	q := client.NewQuery(`SELECT mean(pwrDemand) FROM "autogen"."IIC3NTUST-0007" WHERE time > `+a+` AND ("devID"='33000509b52f1002') GROUP BY time(1h)`, "meter", "ns")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		json.NewEncoder(w).Encode(response.Results)
		fmt.Println(a)
	} else {
		fmt.Println(a)
		fmt.Println(response.Error())
	}

}
func pipeTest(tick time.Time, tock time.Time) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": "33000509b52f1002",
				"lastReportTime": bson.M{
					"$lte": tick,
					"$gt":  tock,
				},
			}}}
	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"hour": bson.M{"$hour": "$lastReportTime"},
				"Year": bson.M{"$year": "$lastReportTime"},
				"day":  bson.M{"$dayOfMonth": "$lastReportTime"},
			},
			"lastReportTime": bson.M{"$last": "$lastReportTime"},
			"Pwr_Demand":     bson.M{"$max": "$Pwr_Demand"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Pwr_Demand":     1,
			"lastReportTime": 1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}
func (sca Influx) wewewe(w http.ResponseWriter, req *http.Request) {
	devID := "781463DA0104F4F3"
	tick, _ := time.Parse(time.RFC3339, "2017-08-14T10:25:26.857+0800")

	// meter.DBName
	a := meter.GetObjectIDOneArg(devID, tick)
	fmt.Println(a)

}
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "meter",
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
