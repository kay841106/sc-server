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

	db        = "sc"
	db_airbox = "airbox"
	// c            = "testing"

	// METER
	c_lastreport   = "lastreport"
	c_aemdra       = "aemdra"
	c_cpm          = "cpm"
	c_gw_status    = "gw_status"
	c_devices      = "devices"
	c_hour         = "hour"
	c_day          = "day"
	c_month        = "month"
	c_downtime     = "downtime"
	c_offlinechart = "offline_chart"

	// AIRBOX
	c_airboxraw  = "airbox_raw"
	c_airboxhour = "airbox_hour"

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
type session struct {
	theSess *mgo.Session
}

func (s *session) startSession() *session {
	return &session{s.theSess.Clone()}
}

func (s *session) checkDBStatus() bool {
	err := s.startSession().theSess.Ping()

	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		s.startSession().theSess.Close()
		time.Sleep(5 * time.Second)
		s.startSession().theSess.Refresh()
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

func (s *session) aggHour() {
	// fmt.Println("ASU")
	log.Println("== Hour running ==")
	if s.checkDBStatus(); true {
		// fmt.Println("Aaaaaaa")
		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := s.startSession().theSess.Clone().DB(db)

		// Mongo := session2.DB(db)

		qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)
		// // DEBUG
		// qu.C(c_devices).Find(bson.M{"MAC_Address": "aa:bb:02:03:01:01"}).Distinct("MAC_Address", &containerdevMan)
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
		qu.Session.Close()
	}
	log.Println("== Hour finish ==")
}

func (s *session) aggDay() {
	log.Println("== Day running ==")
	if s.checkDBStatus(); true {

		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := s.startSession().theSess.Clone().DB(db)
		// defer qu.Session.Close()
		// Mongo := session2.DB(db)

		qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)
		// qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)

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
		qu.Session.Close()
	}
	log.Println("== Day finish ==")
}

func (s *session) aggMonth() {
	log.Println("== Month running ==")
	if s.checkDBStatus(); true {

		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		// // Backup DB
		// session2 = db_connect2()

		qu := s.startSession().theSess.Clone().DB(db)
		// defer qu.Session.Close()
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
		qu.Session.Close()
	}
	log.Println("== Month finish ==")
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
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

// func init() {

// 	dbInfo := &mgo.DialInfo{
// 		Addrs:    strings.SplitN(dblocal, ",", -1),
// 		Database: "admin",
// 		Username: "dontask",
// 		Password: "idontknow",
// 		Timeout:  time.Second * 2,
// 	}
// 	session, _ = mgo.DialWithInfo(dbInfo)
// }

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

type gwstatus struct {
	TimestampUnix int64  `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string `json:"GW_ID" bson:"GW_ID"`
	Place         string `json:"Place" bson:"Place"`
	MGWID         string `json:"M_GWID" bson:"M_GWID"`
}

type gwdowntime struct {
	Type     string `json:"Type" bson:"Type"`
	Duration int64  `json:"Duration" bson:"Duration"`
	GWID     string `json:"GW_ID" bson:"GW_ID"`
	Place    string `json:"Place" bson:"Place"`
	MGWID    string `json:"M_GWID" bson:"M_GWID"`
	Status   bool   `json:"Status" bson:"Status"`
}

type meterdowntime struct {
	Type       string `json:"Type" bson:"Type"`
	Duration   int64  `json:"Duration" bson:"Duration"`
	MACAddress string `json:"MAC_Address" bson:"MAC_Address"`
	Place      string `json:"Place" bson:"Place"`
	MMAC       string `json:"M_MAC" bson:"M_MAC"`
	Territory  string `json:"TERRITORY" bson:"TERRITORY"`
	ID         int    `json:"ID" bson:"ID"`
}

func (s *session) downtime() {
	if s.checkDBStatus(); true {
		// fmt.Println("SS")
		qu := s.startSession().theSess.DB(db)

		container := gwstatus{}
		container2 := gwdowntime{}
		container3 := meterdowntime{}
		var container4 []interface{}
		var container5 []interface{}
		qu.C(c_gw_status).Find(bson.M{}).Distinct("GW_ID", &container4)
		fmt.Println(container4)

		for _, x := range container4 {
			if x == nil {
				break
			}
			qu.C(c_gw_status).Find(bson.M{"GW_ID": x.(string)[0:8]}).One(&container)
			// qu.C(c_gw_status).Find(bson.M{"GW_ID": x.(string)[0:8]}).One(&container)
			qu.C(c_downtime).Find(bson.M{"GW_ID": x.(string)[:8]}).One(&container2)
			if time.Duration(time.Now().Unix()-container.TimestampUnix) <= time.Hour {

			}
			container6 := gwdowntime{
				Type:     "gateway",
				Duration: container2.Duration + (time.Now().Unix() - container.TimestampUnix),
				GWID:     container2.GWID,
				Place:    container2.Place,
				MGWID:    container2.MGWID,
			}
			log.Println(time.Now().Unix(), container.TimestampUnix, container6.GWID[0:8])
			qu.C(c_downtime).Update(bson.M{"GW_ID": container6.GWID[0:8]}, bson.M{"$set": container6})

		}
		qu.C(c_devices).Find(bson.M{}).Distinct("MAC_Address", container5)
		for _, x := range container5 {
			qu.C(c_lastreport).Find(bson.M{"MAC_Address": x}).One(&container2)
			qu.C(c_downtime).Find(bson.M{"MAC_Address": x}).One(&container2)

			container7 := meterdowntime{
				Type:       "meter",
				Duration:   container3.Duration + (time.Now().Unix() - container.TimestampUnix),
				MACAddress: container3.MACAddress,
				Place:      container3.Place,
				MMAC:       container3.MMAC,
				Territory:  container3.Territory,
			}
			qu.C(c_downtime).Update(bson.M{"Type": container7.Type, "MAC_Address": container7.MACAddress}, bson.M{"$set": container7})
			log.Println(container7)
		}

		// if
		// qu.C(c_downtime)
		qu.Session.Close()

	}
}

// // AIRBOX

func pipeAirboxHourWhole(devID string) []bson.M {
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
			"Temp":      bson.M{"$avg": "$Temp"},
			"Humidity":  bson.M{"$avg": "$Humidity"},
			"PM2_5":     bson.M{"$avg": "$PM2_5"},
			"CO":        bson.M{"$avg": "$CO"},
			"CO2":       bson.M{"$avg": "$CO2"},
			"Noise":     bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,
			"Temp":      1,
			"Humidity":  1,
			"PM2_5":     1,
			"CO":        1,
			"CO2":       1,
			"Noise":     1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeAirboxHour(start time.Time, devID string) []bson.M {

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
			"Temp":      bson.M{"$avg": "$Temp"},
			"Humidity":  bson.M{"$avg": "$Humidity"},
			"PM2_5":     bson.M{"$avg": "$PM2_5"},
			"CO":        bson.M{"$avg": "$CO"},
			"CO2":       bson.M{"$avg": "$CO2"},
			"Noise":     bson.M{"$avg": "$p_sum"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"MAC_Address": "$_id.MAC_Address",
			"GW_ID":       "$_id.GW_ID",

			"Timestamp": 1,
			"Temp":      1,
			"Humidity":  1,
			"PM2_5":     1,
			"CO":        1,
			"CO2":       1,
			"Noise":     1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

type airboxAgg struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          *string       `json:"GW_ID" bson:"GW_ID"`
	GET11         *float64      `json:"Temp" bson:"Temp"`
	GET12         *float64      `json:"Humidity" bson:"Humidity"`
	GET13         *float64      `json:"PM2_5" bson:"PM2_5"`
	GET14         *float64      `json:"CO" bson:"CO"`
	GET15         *float64      `json:"CO2" bson:"CO2"`
	GET16         *float64      `json:"Noise" bson:"Noise"`
}

func (s *session) AirboxHour() {
	// fmt.Println("ASU")
	log.Println("== Hour running ==")
	if s.checkDBStatus(); true {
		// fmt.Println("Aaaaaaa")
		var cont []airboxAgg
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs
		var containerdevMan [2]string
		// var containerdevMan []interface{}
		containerdevMan[0] = "58:7a:62:31:32:99"

		// // Backup DB
		// session2 = db_connect2()

		qu := s.startSession().theSess.DB(db_airbox)

		// Mongo := session2.DB(db)

		// qu.C(c_devices).Find(nil).Distinct("MAC_Address", &containerdevMan)
		// // DEBUG
		// qu.C(c_devices).Find(bson.M{"MAC_Address": "aa:bb:02:03:01:01"}).Distinct("MAC_Address", &containerdevMan)
		// fmt.Println(containerdevMan)
		for _, one := range containerdevMan {
			fmt.Println(one)
			count, _ := qu.C(c_airboxraw).Find(bson.M{"MAC_Address": one}).Count()
			fmt.Println("COUNT = ", count)
			if count != 0 {
				err := qu.C(c_airboxraw).Find(bson.M{"MAC_Address": one}).Limit(1).Sort("Timestamp").All(&tempstructs)
				fmt.Print(tempstructs)
				if err != nil {
					fmt.Print(err)
				}
				for _, two := range tempstructs {
					fmt.Println("timestamp=", two.Timestamp)
					err := qu.C(c_airboxraw).Pipe(pipeAirboxHour(two.Timestamp, one)).All(&contdata)
					fmt.Println(contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForHour(each.Timestamp)

							each.ID = getObjectIDTwoArg(*each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_airboxhour).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}
						fmt.Println(each)
						contdata = cont

					}

					// time.Sleep(time.Second * 1)
					fmt.Println("AIRBOX")
					err = qu.C(c_airboxraw).Pipe(pipeAirboxHour(two.Timestamp, one)).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForHour(each.Timestamp)

							each.ID = getObjectIDTwoArg(*each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_airboxhour).Insert(each)
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
				err := qu.C(c_airboxraw).Pipe(pipeAirboxHourWhole(one)).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForHour(each.Timestamp)

						each.ID = getObjectIDTwoArg(*each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_airboxhour).Insert(each)
						// Mongo.C(c_hour).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}
				err = qu.C(c_airboxraw).Pipe(pipeAirboxHourWhole(one)).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForHour(each.Timestamp)

						each.ID = getObjectIDTwoArg(*each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_airboxhour).Insert(each)
						// Mongo.C(c_hour).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
						contdata = cont
					}
				}

			}
		}
		qu.Session.Close()
	}
	log.Println("== Hour finish ==")
}

type structMeterOnlineChart struct {
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp      time.Time     `json:"Timestamp" bson:"Timestamp"`
	Timestamp_Unix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MeterOffline   int           `json:"Meter_Offline" bson:"Meter_Offline"`
	GWOffline      int           `json:"GW_Offline" bson:"GW_Offline"`
}

type structMeterReport struct {
	Timestamp      time.Time `json:"Timestamp" bson:"Timestamp"`
	Timestamp_Unix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress     string    `json:"MAC_Address" bson:"MAC_Address"`
}

type structGWReport struct {
	Timestamp      time.Time `json:"Timestamp" bson:"Timestamp"`
	Timestamp_Unix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID           string    `json:"GW_ID" bson:"GW_ID"`
}

type structGWStatus struct {
	Timestamp      time.Time `json:"Timestamp" bson:"Timestamp"`
	Timestamp_Unix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWOffline      int       `json:"GW_Offline" bson:"GW_Offline"`
}

func (s *session) meterOnlineChart() {
	// fmt.Println("ASU")
	log.Println("== meterOnlineChart running ==")
	if s.checkDBStatus(); true {

		var cont []structMeterReport
		var cont2 []structGWReport

		// INITIATE DB CONNECTION
		qu := s.startSession().theSess.DB(db)

		// ACCUMULATE METER OFFLINE
		qu.C(c_lastreport).Find(bson.M{}).All(&cont)
		var i int
		// fmt.Print(cont)
		for _, each := range cont {
			if time.Now().Sub(each.Timestamp) > time.Duration(time.Hour*1) {
				i++
			}
		}
		// fmt.Print(i)
		// ACCUMULATE GW OFFLINE
		qu.C(c_gw_status).Find(bson.M{}).All(&cont2)
		var j int
		for _, each2 := range cont2 {
			if time.Now().Sub(each2.Timestamp) > time.Duration(time.Hour*1) {
				j++
			}
		}

		sendContainer := structMeterOnlineChart{
			ID:             getObjectIDTwoArg("GW", "MAC", time.Now().Unix()),
			Timestamp_Unix: time.Now().Unix(),
			Timestamp:      time.Now(),
			MeterOffline:   i,
			GWOffline:      j,
		}

		fmt.Print(sendContainer)
		e := qu.C(c_offlinechart).Insert(sendContainer)
		fmt.Print(e)
		i, j = 0, 0
		log.Println("== meterOnlineChart finish ==")
	}
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

	c := cron.New()

	sess := db_connect()

	v := session{sess}

	c.AddFunc("@hourly", v.aggHour)
	c.AddFunc("@daily", v.aggDay)
	c.AddFunc("@monthly", v.aggMonth)

	c.AddFunc("@hourly", v.meterOnlineChart)
	c.AddFunc("@hourly", v.AirboxHour)

	go c.Start()
	sig := make(chan os.Signal)
	fmt.Println("end")
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	// // DEBUG
	// v.aggMonth()
	// // v.downtime()
	// sig := make(chan os.Signal)
	// signal.Notify(sig, os.Interrupt, os.Kill)
	// <-sig
}
