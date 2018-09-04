package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
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

	displayDataCalcCollection = "SC01_displayData_Calc_"
	displayDataCollection     = "SC01_displayData_"
	hourCollection            = "SC01_hour_All"
	dayCollection             = "SC01_day_All"
	monthCollection           = "SC01_month_All"
	_EMPTYDEST                = "DESTINATION IS EMPTY"
	_AGG                      = "AGGREGATION ->"
	_NEC                      = "NON_EMPTY_COLL"
	_EOF                      = "END_OF_FILE"

	streamCollection    = "SC01_Stream"
	statusCollection    = "SC01_Status"
	streamAllCollection = "SC01_Stream_All"

	coll = "SC01_displayData_Calc_"
	//collAll = "SC01_displayData_All"

	collHour         = "SC01_hour_All"
	lookupHourOnTime = "SC01_hour_All"

	collDay         = "SC01_day_All"
	lookupDayOnTime = "SC01_day_All"

	collMonth         = "SC01_month_All"
	lookupMonthOnTime = "SC01_month_All"

	CdevMan = "SC01_DeviceManager"
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
			"pf_avg":    bson.M{"$abs": bson.M{"$avg": "$pf_avg"}},
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
			"ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceDayWhole(devID string) []bson.M {
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
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year": bson.M{"$year": "$Timestamp"},
				"day":  bson.M{"$dayOfYear": "$Timestamp"},
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

func pipeDeviceDay(start time.Time, devID string) []bson.M {

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
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year": bson.M{"$year": "$Timestamp"},
				"day":  bson.M{"$dayOfYear": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			"max_val":   bson.M{"$avg": "$ae_tot"},
			"min_val":   bson.M{"$min": "$ae_tot"},
			"pf_avg":    bson.M{"$abs": bson.M{"$avg": "$pf_avg"}},
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
			"ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Timestamp": 1},
	})

	return pipeline
}

func pipeDeviceMonthWhole(devID string) []bson.M {
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
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":  bson.M{"$year": "$Timestamp"},
				"month": bson.M{"$month": "$Timestamp"},
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

func pipeDeviceMonth(start time.Time, devID string) []bson.M {

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
				// "Hour":        bson.M{"$hour": "$Timestamp"},
				"Year":  bson.M{"$year": "$Timestamp"},
				"month": bson.M{"$month": "$Timestamp"},
			},
			"Timestamp": bson.M{"$last": "$Timestamp"},
			"max_val":   bson.M{"$avg": "$ae_tot"},
			"min_val":   bson.M{"$min": "$ae_tot"},
			"pf_avg":    bson.M{"$abs": bson.M{"$avg": "$pf_avg"}},
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
			"ae_tot": bson.M{"$subtract": []interface{}{"$max_val", "$min_val"}},
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

func checkDBStatus() bool {
	err := session.Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		session.Close()
		time.Sleep(5 * time.Second)
		session.Refresh()
	}
	return true
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

	if checkDBStatus(); true {

		var cont []aggHourStruct
		contdata := cont
		thetempstructs := []tempstruct{}
		tempstructs := thetempstructs

		var containerdevMan []interface{}

		qu := session.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MACAddress", &containerdevMan)

		for _, one := range containerdevMan {

			count, _ := qu.C(c_hour).Find(bson.M{"MAC_Address": one.(string)}).Count()

			if count != 0 {
				err := qu.C(c_hour).Find(bson.M{"MAC_Address": one.(string)}).Limit(1).Sort("-Timestamp").All(&tempstructs)
				fmt.Print(tempstructs)
				if err != nil {
					fmt.Print(err)
				}
				for _, two := range tempstructs {
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
						fmt.Print(each)
						contdata = cont

					}

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
						fmt.Print(each)
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

		qu := session.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MACAddress", &containerdevMan)

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
					err := qu.C(c_hour).Pipe(pipeDeviceDay(two.Timestamp, one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForDay(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_day).Insert(each)
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
				err := qu.C(c_hour).Pipe(pipeDeviceDayWhole(one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForDay(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_day).Insert(each)
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

		qu := session.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MACAddress", &containerdevMan)

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
					err := qu.C(c_day).Pipe(pipeDeviceMonth(two.Timestamp, one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = SetTimeStampForMonth(each.Timestamp)

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_month).Insert(each)
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
				err := qu.C(c_day).Pipe(pipeDeviceMonthWhole(one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = SetTimeStampForMonth(each.Timestamp)

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())
						fmt.Print(each)

						qu.C(c_month).Insert(each)
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

func main() {

	c := cron.New()

	c.AddFunc("@hourly", aggHour)
	c.AddFunc("@daily", aggDay)
	c.AddFunc("@monthly", aggMonth)

	c.Start()
	select {}

}
