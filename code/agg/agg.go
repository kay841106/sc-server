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
				"Device_ID": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID":  "$Device_ID",
				"Year":       bson.M{"$year": "$lastReportTime"},
				"day":        bson.M{"$dayOfYear": "$lastReportTime"},
				"Gateway_ID": "$Gateway_ID",
			},
			"lastReportTime":   bson.M{"$last": "$lastReportTime"},
			"Floor":            bson.M{"$last": "$Floor"},
			"Building_Name":    bson.M{"$last": "$Building_Name"},
			"Device_Name":      bson.M{"$last": "$Device_Name"},
			"Facility":         bson.M{"$last": "$Facility"},
			"Device_Type":      bson.M{"$last": "$Device_Type"},
			"Building_Details": bson.M{"$last": "$Building_Details"},

			"max_Demand":  bson.M{"$max": "$avg_Demand"},
			"min_Demand":  bson.M{"$min": "$avg_Demand"},
			"max_Usage":   bson.M{"$max": "$avg_Usage"},
			"min_Usage":   bson.M{"$min": "$avg_Usage"},
			"avg_PF":      bson.M{"$avg": bson.M{"$abs": "$avg_PF"}},
			"max_PF":      bson.M{"$max": bson.M{"$abs": "$avg_PF"}},
			"min_PF":      bson.M{"$min": bson.M{"$abs": "$avg_PF"}},
			"PF_Limit":    bson.M{"$avg": "$PF_Limit"},
			"avg_Usage":   bson.M{"$avg": "$avg_Usage"},
			"avg_Demand":  bson.M{"$avg": "$avg_Demand"},
			"total_Usage": bson.M{"$sum": "$avg_Usage"},
			"CC":          bson.M{"$avg": "$CC"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID":  "$_id.Device_ID",
			"Gateway_ID": "$_id.Gateway_ID",

			"lastReportTime":   1,
			"Floor":            1,
			"Building_Name":    1,
			"Device_Name":      1,
			"Facility":         1,
			"Device_Type":      1,
			"Building_Details": 1,

			"CC": "$CC",

			"avg_Usage":   1,
			"avg_Demand":  1,
			"max_Demand":  1,
			"min_Demand":  1,
			"max_Usage":   1,
			"min_Usage":   1,
			"total_Usage": 1,

			"max_PF":   1,
			"min_PF":   1,
			"avg_PF":   bson.M{"$abs": "$avg_PF"},
			"PF_Limit": 1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}
func pipeDeviceDay(start time.Time, devID string) []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID":  "$Device_ID",
				"Year":       bson.M{"$year": "$lastReportTime"},
				"day":        bson.M{"$dayOfYear": "$lastReportTime"},
				"Gateway_ID": "$Gateway_ID",
			},
			"lastReportTime":   bson.M{"$last": "$lastReportTime"},
			"Floor":            bson.M{"$last": "$Floor"},
			"Building_Name":    bson.M{"$last": "$Building_Name"},
			"Device_Name":      bson.M{"$last": "$Device_Name"},
			"Facility":         bson.M{"$last": "$Facility"},
			"Device_Type":      bson.M{"$last": "$Device_Type"},
			"Building_Details": bson.M{"$last": "$Building_Details"},

			"max_Demand":  bson.M{"$max": "$avg_Demand"},
			"min_Demand":  bson.M{"$min": "$avg_Demand"},
			"max_Usage":   bson.M{"$max": "$total_Usage"},
			"min_Usage":   bson.M{"$min": "$total_Usage"},
			"avg_PF":      bson.M{"$avg": bson.M{"$abs": "$avg_PF"}},
			"max_PF":      bson.M{"$max": bson.M{"$abs": "$avg_PF"}},
			"min_PF":      bson.M{"$min": bson.M{"$abs": "$avg_PF"}},
			"PF_Limit":    bson.M{"$avg": "$PF_Limit"},
			"avg_Usage":   bson.M{"$avg": "$total_Usage"},
			"avg_Demand":  bson.M{"$avg": "$avg_Demand"},
			"total_Usage": bson.M{"$sum": "$total_Usage"},
			"CC":          bson.M{"$avg": "$CC"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name": "$_id.Building_Name",
			"Gateway_ID":    "$_id.Gateway_ID",

			"CC":             "$CC",
			"lastReportTime": 1,
			"total_Usage":    1,

			"Floor":            1,
			"Device_Name":      1,
			"Facility":         1,
			"Device_Type":      1,
			"Building_Details": 1,

			"avg_Usage":  1,
			"avg_Demand": 1,
			"max_Demand": 1,
			"min_Demand": 1,
			"max_Usage":  1,
			"min_Usage":  1,

			"max_PF":   1,
			"min_PF":   1,
			"PF_Limit": 1,
			"avg_PF":   bson.M{"$abs": "$avg_PF"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}

func pipeDeviceMonthWhole(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",

				"Month":      bson.M{"$month": "$lastReportTime"},
				"Year":       bson.M{"$year": "$lastReportTime"},
				"Gateway_ID": "$Gateway_ID",
			},
			"lastReportTime":   bson.M{"$last": "$lastReportTime"},
			"Floor":            bson.M{"$last": "$Floor"},
			"Building_Name":    bson.M{"$last": "$Building_Name"},
			"Device_Name":      bson.M{"$last": "$Device_Name"},
			"Facility":         bson.M{"$last": "$Facility"},
			"Device_Type":      bson.M{"$last": "$Device_Type"},
			"Building_Details": bson.M{"$last": "$Building_Details"},

			"max_Demand":  bson.M{"$max": "$avg_Demand"},
			"min_Demand":  bson.M{"$min": "$avg_Demand"},
			"max_Usage":   bson.M{"$max": "$total_Usage"},
			"min_Usage":   bson.M{"$min": "$total_Usage"},
			"avg_PF":      bson.M{"$avg": bson.M{"$abs": "$avg_PF"}},
			"max_PF":      bson.M{"$max": bson.M{"$abs": "$avg_PF"}},
			"min_PF":      bson.M{"$min": bson.M{"$abs": "$avg_PF"}},
			"PF_Limit":    bson.M{"$avg": "$PF_Limit"},
			"avg_Usage":   bson.M{"$avg": "$total_Usage"},
			"avg_Demand":  bson.M{"$avg": "$avg_Demand"},
			"total_Usage": bson.M{"$sum": "$total_Usage"},
			"CC":          bson.M{"$avg": "$CC"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID":  "$_id.Device_ID",
			"Gateway_ID": "$_id.Gateway_ID",

			"lastReportTime":   1,
			"Floor":            1,
			"Building_Name":    1,
			"Device_Name":      1,
			"Facility":         1,
			"Device_Type":      1,
			"Building_Details": 1,

			"CC": "$CC",

			"avg_Usage":   1,
			"avg_Demand":  1,
			"max_Demand":  1,
			"min_Demand":  1,
			"max_Usage":   1,
			"min_Usage":   1,
			"total_Usage": 1,

			"max_PF":   1,
			"min_PF":   1,
			"PF_Limit": 1,
			"avg_PF":   bson.M{"$abs": "$avg_PF"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}
func pipeDeviceMonth(start time.Time, devID string) []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",

				"Month":      bson.M{"$month": "$lastReportTime"},
				"Year":       bson.M{"$year": "$lastReportTime"},
				"Gateway_ID": "$Gateway_ID",
			},
			"lastReportTime":   bson.M{"$last": "$lastReportTime"},
			"Floor":            bson.M{"$last": "$Floor"},
			"Building_Name":    bson.M{"$last": "$Building_Name"},
			"Device_Name":      bson.M{"$last": "$Device_Name"},
			"Facility":         bson.M{"$last": "$Facility"},
			"Device_Type":      bson.M{"$last": "$Device_Type"},
			"Building_Details": bson.M{"$last": "$Building_Details"},

			"max_Demand":  bson.M{"$max": "$avg_Demand"},
			"min_Demand":  bson.M{"$min": "$avg_Demand"},
			"max_Usage":   bson.M{"$max": "$avg_Usage"},
			"min_Usage":   bson.M{"$min": "$avg_Usage"},
			"avg_PF":      bson.M{"$avg": bson.M{"$abs": "$avg_PF"}},
			"max_PF":      bson.M{"$max": bson.M{"$abs": "$avg_PF"}},
			"min_PF":      bson.M{"$min": bson.M{"$abs": "$avg_PF"}},
			"avg_Usage":   bson.M{"$avg": "$avg_Usage"},
			"avg_Demand":  bson.M{"$avg": "$avg_Demand"},
			"total_Usage": bson.M{"$sum": "$total_Usage"},
			"CC":          bson.M{"$avg": "$CC"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name": "$_id.Building_Name",
			"Gateway_ID":    "$_id.Gateway_ID",

			"CC":             "$CC",
			"lastReportTime": 1,
			"total_Usage":    1,

			"Floor":            1,
			"Device_Name":      1,
			"Facility":         1,
			"Device_Type":      1,
			"Building_Details": 1,

			"avg_Usage":  1,
			"avg_Demand": 1,
			"max_Demand": 1,
			"min_Demand": 1,
			"max_Usage":  1,
			"min_Usage":  1,

			"max_PF": 1,
			"min_PF": 1,

			"PF_Limit": 1,
			"avg_PF":   bson.M{"$abs": "$avg_PF"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
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
	Timestamp time.Time `json:"Timestamp" bson:"Timestamp"`
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
		tempstructs := []tempstruct{}
		// contz := aggHourStruct{}
		// containerhour := []aggHourStruct{}
		var containerdevMan []interface{}

		qu := session.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MACAddress", &containerdevMan)

		for _, one := range containerdevMan {

			count, _ := qu.C(c_hour).Find(bson.M{}).Count()
			if count != 0 {
				qu.C(c_hour).Find(bson.M{"Mac_Address": one}).Limit(1).Sort("-Timestamp").All(&tempstructs)
				for _, two := range tempstructs {
					err := qu.C(c_cpm).Pipe(pipeDeviceHour(two.Timestamp, one.(string))).All(&contdata)
					for _, each := range contdata {
						if (each.Timestamp != time.Time{}) {

							each.Timestamp = convertHour(each.Timestamp.Unix())

							each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

							qu.C(c_hour).Insert(each)
							if err != nil {
								fmt.Print(err)
							}
						}

						// containerhour.WeatherTemp=
						// fmt.Print(containerdevID)

						// 	// 	qu.C(c_devices).Find(bson.M{"Building_Name": one}).Distinct("devID", &containerdevID)
						// 	// fmt.Println(containerdevID)
						// 	for _, two := range containerdevID {

						// 		fmt.Println("two" + two.(string))
						// 		qu.C(hourCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
						// 		fmt.Println(containerLastRecord.DeviceID)
						// 		// fmt.Println(containerLastRecord)
						// 		if (aggHourStruct{}) == containerLastRecord {

						// 			thepipe := pipeDeviceHourWhole(two.(string))
						// 			// qu.C(displayDataCalcCollection + one.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

						// 			for _, each := range explainacontainer {
						// 				fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
						// 				tempTime := SetTimeStampForHour(each.LastReportTime)
						// 				each.LastReportTime = tempTime
						// 				each.PFLimit = 0.8

						// 				qu.C(hourCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
						// 			}
						// 		} else {
						// 			thepipe := pipeDeviceHour(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
						// 			// qu.C(displayDataCalcCollection + one.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
						// 			fmt.Println(explainacontainer)
						// 			for _, each := range explainacontainer {
						// 				fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
						// 				tempTime := SetTimeStampForHour(each.LastReportTime)
						// 				each.LastReportTime = tempTime
						// 				each.PFLimit = 0.8
						// 				qu.C(hourCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
						// 			}
						// 		}
						// 		explainacontainer = []aggHourStruct{}
						// 		containerLastRecord = aggHourStruct{}
						// 	}
					}
				}
			} else {
				err := qu.C(c_cpm).Pipe(pipeDeviceHourWhole(one.(string))).All(&contdata)
				for _, each := range contdata {
					if (each.Timestamp != time.Time{}) {

						each.Timestamp = convertHour(each.Timestamp.Unix())

						each.ID = getObjectIDTwoArg(each.GWID, each.MACAddress, each.Timestamp.Unix())

						qu.C(c_hour).Insert(each)
						if err != nil {
							fmt.Print(err)
						}
					}
				}
				// 	fmt.Println(_EOF)
				session.Close()

			}
		}
	}
}

// func AggDay() {
// 	status := checkDBStatus()
// 	if status == true {
// 		mongo := meter.GetMongo()
// 		containerLastRecord := meter.AggDayStruct{}
// 		var containerdevMan []interface{}
// 		explainacontainer := []meter.AggDayStruct{}

// 		qu := mongo.DB(meter.DBName)
// 		qu.C(CdevMan).Find(nil).Distinct("devID", &containerdevMan)

// 		// fmt.Println(containerdevMan)
// 		for _, two := range containerdevMan {

// 			fmt.Println(two)
// 			qu.C(dayCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
// 			// fmt.Println("devID_query = ", containerLastRecord.DeviceID, two)
// 			// fmt.Println(containerLastRecord)

// 			if (meter.AggDayStruct{}) == containerLastRecord {

// 				thepipe := pipeDeviceDayWhole(two.(string))
// 				qu.C(hourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
// 				fmt.Println(explainacontainer)
// 				for _, each := range explainacontainer {
// 					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
// 					tempTime := SetTimeStampForDay(each.LastReportTime)
// 					each.LastReportTime = tempTime
// 					qu.C(dayCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
// 					// each = meter.AggDayStruct{}
// 				}
// 				explainacontainer = []meter.AggDayStruct{}

// 			} else {

// 				thepipe := pipeDeviceDay(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
// 				qu.C(hourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
// 				fmt.Println(explainacontainer)
// 				for _, each := range explainacontainer {
// 					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
// 					tempTime := SetTimeStampForDay(each.LastReportTime)
// 					each.LastReportTime = tempTime
// 					qu.C(dayCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
// 					// each = meter.AggDayStruct{}
// 				}
// 				explainacontainer = []meter.AggDayStruct{}

// 			}
// 			containerLastRecord = meter.AggDayStruct{}
// 		}

// 		fmt.Println(_EOF)
// 		qu.Close()
// 	}
// }
// func AggMonth() {
// 	status := checkDBStatus()
// 	if status == true {
// 		mongo := meter.GetMongo()
// 		containerLastRecord := meter.AggDayStruct{}
// 		var containerdevMan []interface{}
// 		explainacontainer := []meter.AggDayStruct{}

// 		qu := mongo.DB(meter.DBName)
// 		qu.C(CdevMan).Find(nil).Distinct("devID", &containerdevMan)

// 		// fmt.Println(containerdevMan)

// 		for _, two := range containerdevMan {

// 			fmt.Println(two)
// 			qu.C(monthCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
// 			fmt.Println(containerLastRecord.DeviceID)
// 			if (meter.AggDayStruct{}) == containerLastRecord {

// 				thepipe := pipeDeviceMonthWhole(two.(string))
// 				qu.C(dayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

// 				for _, each := range explainacontainer {
// 					fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
// 					tempTime := SetTimeStampForMonth(each.LastReportTime)
// 					fmt.Println(tempTime)
// 					each.LastReportTime = tempTime
// 					qu.C(monthCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
// 				}
// 			} else {
// 				thepipe := pipeDeviceMonth(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
// 				qu.C(dayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
// 				fmt.Println(explainacontainer)
// 				for _, each := range explainacontainer {
// 					fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
// 					tempTime := SetTimeStampForMonth(each.LastReportTime)

// 					each.LastReportTime = tempTime
// 					qu.C(monthCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
// 				}
// 				explainacontainer = []meter.AggDayStruct{}
// 			}

// 			containerLastRecord = meter.AggDayStruct{}
// 		}

// 		fmt.Println(_EOF)
// 		qu.Close()
// 	}
// }

//SetTimeStampForHour set minute second to 0
func SetTimeStampForHour(theTime time.Time) time.Time {
	year, month, day := theTime.Date()
	hour, _, _ := theTime.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

func convertHour(t int64) time.Time {
	times := time.Unix(t, 0).UTC()
	year, month, day := times.Date()
	hour, _, _ := times.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

//SetTimeStampForDay set hour minute second to 0
func SetTimeStampForDay(theTime time.Time) time.Time {
	date := theTime.Day()
	month := theTime.Month()
	year := theTime.Year()
	return time.Date(year, month, date, 0, 0, 0, 0, time.Local)
}

//SetTimeStampForMonth set day hour minute sec to 0
func SetTimeStampForMonth(theTime time.Time) time.Time {
	month := theTime.Month()
	year := theTime.Year()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func main() {
	aggHour()
}
