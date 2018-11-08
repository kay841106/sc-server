package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/robfig/cron"
	// "sc.dforcepro.com/meter"
)

// import (
// 	"gopkg.in/mgo.v2/bson"
// )

const (
	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:10003"

	dbbackup = "140.118.122.103:27017"

	db = "sc"
	// c            = "testing"
	c_lastreport = "lastreport"
	c_aemdra     = "aemdra"
	c_cpm        = "cpm"
	c_gw_status  = "gw_status"
	c_devices    = "devices"
	c_hour       = "hour"
	c_day        = "day"
	c_month      = "month"

	_EMPTYDEST = "DESTINATION IS EMPTY"
	_AGG       = "AGGREGATION ->"
	_NEC       = "NON_EMPTY_COLL"
	_EOF       = "END_OF_FILE"

	// weatherCollection = "SC01_Weather1"

)

func pipeDeviceHourWhole(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"MAC_Address": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				"Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":        bson.M{"$year": "$Timestamp"},
				"day":         bson.M{"$dayOfYear": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			"max_val":   bson.M{"$max": "$ae_tot"},
			"min_val":   bson.M{"$min": "$ae_tot"},
			"pf_avg":    bson.M{"$avg": "$pf_avg"},
			"p_sum":     bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,

			"ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceHour(start time.Time, devID string) []bson.M {

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Timestamp": bson.M{
					"$gt": start,
				}, "MAC_Address": devID,
			},
		}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				"Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":        bson.M{"$year": "$Timestamp"},
				"day":         bson.M{"$dayOfYear": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			"max_val":   bson.M{"$avg": "$ae_tot"},
			"min_val":   bson.M{"$min": "$ae_tot"},
			"pf_avg":    bson.M{"$avg": "$pf_avg"},
			"p_sum":     bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,

			"ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceDayWhole(stop time.Time, devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"MAC_Address": devID,
				"Timestamp":   bson.M{"$lte": stop},
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year": bson.M{"$year": "$Timestamp"},
				"day":  bson.M{"$dayOfYear": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},

			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "max_val":   bson.M{"$max": "$ae_tot"},
			// "min_val":   bson.M{"$min": "$ae_tot"},

			"ae_tot": bson.M{"$sum": "$ae_tot"},
			"pf_avg": bson.M{"$avg": "$pf_avg"},
			"p_sum":  bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,
			"ae_tot": 1,
			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},

		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceDay(start time.Time, stop time.Time, devID string) []bson.M {

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Timestamp": bson.M{
					"$gt":  start,
					"$lte": stop,
				}, "MAC_Address": devID,
			},
		}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year": bson.M{"$year": "$Timestamp"},
				"day":  bson.M{"$dayOfYear": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},

			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "max_val":   bson.M{"$max": "$ae_tot"},
			// "min_val":   bson.M{"$min": "$ae_tot"},

			"ae_tot": bson.M{"$sum": "$ae_tot"},
			"pf_avg": bson.M{"$avg": "$pf_avg"},
			"p_sum":  bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,
			"ae_tot": 1,

			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},

		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceMonthWhole(stop time.Time, devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"MAC_Address": devID,
				"Timestamp":   bson.M{"$lte": stop},
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":  bson.M{"$year": "$Timestamp"},
				"month": bson.M{"$month": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "max_val":   bson.M{"$max": "$ae_tot"},
			// "min_val":   bson.M{"$min": "$ae_tot"},

			"ae_tot": bson.M{"$sum": "$ae_tot"},
			"pf_avg": bson.M{"$avg": "$pf_avg"},
			"p_sum":  bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,

			"ae_tot": 1,

			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceMonth(start time.Time, stop time.Time, devID string) []bson.M {

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Timestamp": bson.M{
					"$gt":  start,
					"$lte": stop,
				}, "MAC_Address": devID,
			},
		}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"MAC_Address": "$MAC_Address",
				"GW_ID":       "$GW_ID",
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":  bson.M{"$year": "$Timestamp"},
				"month": bson.M{"$month": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "max_val":   bson.M{"$max": "$ae_tot"},
			// "min_val":   bson.M{"$min": "$ae_tot"},

			"ae_tot": bson.M{"$sum": "$ae_tot"},
			"pf_avg": bson.M{"$avg": "$pf_avg"},
			"p_sum":  bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,

			"pf_avg": 1,
			"p_sum":  1,

			"ae_tot": 1,

			// // Below is wrong aggregation. This is used if query data from rawdata collection. Get MAX and MIN
			// "ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

type aggHourStruct struct {
	ID         bson.ObjectId `json:"_id" bson:"_id"`
	GWID       string        `json:"GW_ID" bson:"GW_ID"`
	Timestamp  time.Time     `json:"Timestamp" bson:"Timestamp"`
	MACAddress string        `json:"MAC_Address" bson:"MAC_Address"`
	GET11      float64       `json:"pf_avg" bson:"pf_avg"` // KWh
	GET12      float64       `json:"ae_tot" bson:"ae_tot"`
	GET13      float64       `json:"p_sum" bson:"p_sum"`
	// GET14      float64       `json:"max_val" bson:"max_val"`
	// WeatherTemp int           `json:"weather_temp" bson:"weather_temp"`
}

type tempstruct struct {
	MACAddress string    `json:"MAC_Address" bson:"MAC_Address"`
	Timestamp  time.Time `json:"Timestamp" bson:"Timestamp"`
}

var session *mgo.Session

var session2 *mgo.Session

func checkDBStatus() bool {
	err := session.Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		session.Close()
		time.Sleep(5 * time.Second)
		session.Refresh()
	}
	fmt.Println("DB GOOD")
	return true
}

func getObjectIDTwoArg(GWID string, macID string, timestamp int64) bson.ObjectId {
	var b [12]byte
	var sum [8]byte
	var c [4]byte
	// timestamp := time.Unix(LastReportTime, 0)
	// binary.BigEndian.PutUint32(b[:], uint32(timestamp))
	binary.BigEndian.PutUint32(c[:], uint32(timestamp))

	did := sum[:]
	gid := sum[:]

	hw := md5.New()
	hw.Write([]byte(GWID))
	copy(did, hw.Sum(nil))
	hw.Write([]byte(macID))
	copy(gid, hw.Sum(nil))
	// b[0] = c[:1]
	b[0] = c[0]
	b[1] = c[1]
	b[2] = c[2]
	b[3] = c[3]
	b[4] = did[4]
	b[5] = did[5]
	b[6] = did[6]
	b[7] = did[7]
	b[8] = gid[4]
	b[9] = gid[5]
	b[10] = gid[6]
	b[11] = gid[7]

	dst := hex.EncodeToString(b[:])
	theid := bson.ObjectIdHex(dst)

	fmt.Println(theid, uint32(timestamp))
	return theid

}

func aggHour() {
	// fmt.Println("ASU")

	if checkDBStatus(); true {
		// fmt.Println("Aaaaaaa")
		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := session.DB(db)

		// Mongo := session2.DB(db)

		qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)
		// fmt.Println(containerdevMan)
		for _, one := range containerdevMan {

			count, _ := qu.C(c_hour).Find(bson.M{"MAC_Address": one.(string)}).Count()
			fmt.Println("COUNT = ", count)
			if count != 0 {
				err := qu.C(c_hour).Find(bson.M{"MAC_Address": one.(string)}).Limit(1).Sort("-Timestamp").All(&tempstructs)
				fmt.Print(tempstructs)
				if err != nil {
					fmt.Print(err)
				}
				for _, two := range tempstructs {
					fmt.Println("timestamp=", two.Timestamp)
					err := qu.C(c_cpm).Pipe(pipeDeviceHour(two.Timestamp, one.(string))).All(&contdata)

					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForHour(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_hour).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}
						fmt.Println(each)
						contdata = cont

					}

					// time.Sleep(time.Second * 1)
					fmt.Println("AEM_DRA")
					err = qu.C(c_aemdra).Pipe(pipeDeviceHour(two.Timestamp, one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForHour(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_hour).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}
						fmt.Println(each)
						contdata = cont

					}
				}
				tempstructs = thetempstructs
			} else {
				err := qu.C(c_cpm).Pipe(pipeDeviceHourWhole(one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForHour(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_hour).Insert(each)
						// Mongo.C(c_hour).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}
				err = qu.C(c_aemdra).Pipe(pipeDeviceHourWhole(one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForHour(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_hour).Insert(each)
						// Mongo.C(c_hour).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}

			}
		}
	}

}

func aggDay() {

	if checkDBStatus(); true {

		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := session.DB(db)
		// Mongo := session2.DB(db)

		qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)

		for _, one := range containerdevMan {

			count, _ := qu.C(c_day).Find(bson.M{"MAC_Address": one.(string)}).Count()
			// fmt.Print(count)
			if count != 0 {
				err := qu.C(c_day).Find(bson.M{"MAC_Address": one.(string)}).Limit(1).Sort("-Timestamp").All(&tempstructs)
				fmt.Print(tempstructs)
				if err != nil {
					fmt.Print(err)
				}
				for _, two := range tempstructs {
					err := qu.C(c_hour).Pipe(pipeDeviceDay(two.Timestamp, SetTimeStampForDay(time.Now()), one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForDay(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_day).Insert(each)
							// Mongo.C(c_day).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}
						fmt.Print(each)
						contdata = cont

					}
				}
				tempstructs = thetempstructs
			} else {
				err := qu.C(c_hour).Pipe(pipeDeviceDayWhole(SetTimeStampForDay(time.Now()), one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForDay(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_day).Insert(each)
						// Mongo.C(c_day).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}

			}
		}
	}

}

func aggMonth() {

	if checkDBStatus(); true {

		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := session.DB(db)

		// Mongo := session2.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)

		for _, one := range containerdevMan {

			count, _ := qu.C(c_month).Find(bson.M{"MAC_Address": one.(string)}).Count()
			// fmt.Print(count)
			if count != 0 {
				err := qu.C(c_month).Find(bson.M{"MAC_Address": one.(string)}).Limit(1).Sort("-Timestamp").All(&tempstructs)
				fmt.Print(tempstructs)
				if err != nil {
					fmt.Print(err)
				}
				for _, two := range tempstructs {
					err := qu.C(c_day).Pipe(pipeDeviceMonth(two.Timestamp, SetTimeStampForMonth(time.Now()), one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForMonth(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_month).Insert(each)
							// Mongo.C(c_month).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}
						fmt.Print(each)
						contdata = cont

					}
				}
				tempstructs = thetempstructs
			} else {
				err := qu.C(c_day).Pipe(pipeDeviceMonthWhole(SetTimeStampForMonth(time.Now()), one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForMonth(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_month).Insert(each)
						// Mongo.C(c_month).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}

			}
		}
	}

}

//SetTimeStampForHour set minute second to 0
func SetTimeStampForHour(theTime time.Time) time.Time {
	year, month, day := theTime.Date()
	hour, _, _ := theTime.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

//SetTimeStampForDay set hour minute second to 0
func SetTimeStampForDay(theTime time.Time) time.Time {
	year, month, day := theTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//SetTimeStampForMonth set day hour minute sec to 0
func SetTimeStampForMonth(theTime time.Time) time.Time {
	year, month, _ := theTime.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

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

// func db_connect2() *mgo.Session {

// 	dbInfo := &mgo.DialInfo{
// 		Addrs:    strings.SplitN(dblocal, ",", -1),
// 		Database: "admin",
// 		Username: "dontask",
// 		Password: "idontknow",
// 		Timeout:  time.Second * 10,
// 	}

// 	sess, err := mgo.DialWithInfo(dbInfo)
// 	if err != nil {
// 		os.Exit(1)
// 	}
// 	return sess
// }

func main() {

	c := cron.New()
	db_connect()
	// db_connect2()

	fmt.Print("start")

	c.AddFunc("@hourly", aggHour)
	c.AddFunc("@daily", aggDay)
	c.AddFunc("@monthly", aggMonth)

	go c.Start()
	sig := make(chan os.Signal)
	fmt.Println("end")
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// select {}

	// // DEBUG
	// aggMonth()
	// sig := make(chan os.Signal)
	// signal.Notify(sig, os.Interrupt, os.Kill)
	// <-sig
}
