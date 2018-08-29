

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sc.dforcepro.com/meter"
)

// import (
// 	"gopkg.in/mgo.v2/bson"
// )

const (
	db           = "sc"
	c_lastreport = "lastreport"
	c_devices    = "devices"

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
				"MACAddress": devID,
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
			"Timestamp":        bson.M{"$last": "$Timestamp"},
			"":                 bson.M{"$last": "$Building_Name"},
			"Device_Name":      bson.M{"$last": "$Device_Name"},
			"Facility":         bson.M{"$last": "$Facility"},
			"Device_Type":      bson.M{"$last": "$Device_Type"},
			"Building_Details": bson.M{"$last": "$Building_Details"},

			"ae.tot": bson.M{"$avg": "$ae.tot"},
			"min_Demand": bson.M{"$min": "$Pwr_Demand"},
			"max_Usage":  bson.M{"$max": "$Usage"},
			"min_Usage":  bson.M{"$min": "$Usage"},
			"avg_PF":     bson.M{"$avg": bson.M{"$abs": "$PF"}},
			"max_PF":     bson.M{"$max": bson.M{"$abs": "$PF"}},
			"min_PF":     bson.M{"$min": bson.M{"$abs": "$PF"}},
			// "avg_Usage":    bson.M{"$avg": "$Usage"},
			"avg_Demand":   bson.M{"$avg": "$Pwr_Demand"},
			"total_Usage":  bson.M{"$sum": "$Usage"},
			"CC":           bson.M{"$avg": "$CC"},
			"weather_Temp": bson.M{"$avg": "$weather_Temp"},
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

			// "avg_Usage":    1,
			"avg_Demand":   1,
			"max_Demand":   1,
			"min_Demand":   1,
			"max_Usage":    1,
			"min_Usage":    1,
			"avg_PF":       1,
			"max_PF":       1,
			"min_PF":       1,
			"weather_Temp": 1,
			"total_Usage":  1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}

func pipeDeviceHour(start time.Time, devID string) []bson.M {

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"lastReportTime": bson.M{
					"$gt": start,
				}, "Device_ID": devID,
			},
		}}
	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",
				"Hour":      bson.M{"$hour": "$lastReportTime"},
				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month":      bson.M{"$month": "$lastReportTime"},
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

			"max_Demand": bson.M{"$max": "$Pwr_Demand"},
			"min_Demand": bson.M{"$min": "$Pwr_Demand"},
			"max_Usage":  bson.M{"$max": "$Usage"},
			"min_Usage":  bson.M{"$min": "$Usage"},
			"avg_PF":     bson.M{"$avg": bson.M{"$abs": "$PF"}},
			"max_PF":     bson.M{"$max": bson.M{"$abs": "$PF"}},
			"min_PF":     bson.M{"$min": bson.M{"$abs": "$PF"}},
			// "avg_Usage":    bson.M{"$avg": "$Usage"},
			"avg_Demand":   bson.M{"$avg": "$Pwr_Demand"},
			"total_Usage":  bson.M{"$sum": "$Usage"},
			"CC":           bson.M{"$avg": "$CC"},
			"weather_Temp": bson.M{"$avg": "$weather_Temp"},
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

			// "avg_Usage":  1,
			"avg_Demand":  1,
			"max_Demand":  1,
			"min_Demand":  1,
			"max_Usage":   1,
			"min_Usage":   1,
			"total_Usage": 1,

			"avg_PF":       1,
			"max_PF":       1,
			"min_PF":       1,
			"weather_Temp": 1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
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
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	CC              float64   `json:"CC" bson:"CC"`
	Usage           float64   `json:"total_Usage" bson:"total_Usage"` // KWh
	Floor           string    `json:"Floor" bson:"Floor"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
	AvgPF           float64   `json:"avg_PF" bson:"avg_PF"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"` //KWh
	PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
	MinPF           float64   `json:"min_PF" bson:"min_PF"`
	MaxPF           float64   `json:"max_PF" bson:"max_PF"`
	MinDemand       float64   `json:"min_Demand" bson:"min_Demand"` //W
	MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"` //W
	MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`   //KWh
}

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

func AggHour() {

	if checkDBStatus(); true {

		// containerLastRecord := aggHourStruct{}
		var containerdevMan []interface{}
		// var containerdevID []interface{}
		// explainacontainer := []aggHourStruct{}

		qu := session.DB(db)
		qu.C(c_devices).Find(nil).Distinct("MACAddress", &containerdevMan)

		fmt.Println(containerdevMan)

		for _, one := range containerdevMan {

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

		// 	fmt.Println(_EOF)
		// 	session.Close()

	}
}

func AggDay() {
	status := checkDBStatus()
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := meter.AggDayStruct{}
		var containerdevMan []interface{}
		explainacontainer := []meter.AggDayStruct{}

		qu := mongo.DB(meter.DBName)
		qu.C(CdevMan).Find(nil).Distinct("devID", &containerdevMan)

		// fmt.Println(containerdevMan)
		for _, two := range containerdevMan {

			fmt.Println(two)
			qu.C(dayCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
			// fmt.Println("devID_query = ", containerLastRecord.DeviceID, two)
			// fmt.Println(containerLastRecord)

			if (meter.AggDayStruct{}) == containerLastRecord {

				thepipe := pipeDeviceDayWhole(two.(string))
				qu.C(hourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForDay(each.LastReportTime)
					each.LastReportTime = tempTime
					qu.C(dayCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
					// each = meter.AggDayStruct{}
				}
				explainacontainer = []meter.AggDayStruct{}

			} else {

				thepipe := pipeDeviceDay(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
				qu.C(hourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForDay(each.LastReportTime)
					each.LastReportTime = tempTime
					qu.C(dayCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
					// each = meter.AggDayStruct{}
				}
				explainacontainer = []meter.AggDayStruct{}

			}
			containerLastRecord = meter.AggDayStruct{}
		}

		fmt.Println(_EOF)
		qu.Close()
	}
}
func AggMonth() {
	status := checkDBStatus()
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := meter.AggDayStruct{}
		var containerdevMan []interface{}
		explainacontainer := []meter.AggDayStruct{}

		qu := mongo.DB(meter.DBName)
		qu.C(CdevMan).Find(nil).Distinct("devID", &containerdevMan)

		// fmt.Println(containerdevMan)

		for _, two := range containerdevMan {

			fmt.Println(two)
			qu.C(monthCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
			fmt.Println(containerLastRecord.DeviceID)
			if (meter.AggDayStruct{}) == containerLastRecord {

				thepipe := pipeDeviceMonthWhole(two.(string))
				qu.C(dayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

				for _, each := range explainacontainer {
					fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForMonth(each.LastReportTime)
					fmt.Println(tempTime)
					each.LastReportTime = tempTime
					qu.C(monthCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
				}
			} else {
				thepipe := pipeDeviceMonth(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
				qu.C(dayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForMonth(each.LastReportTime)

					each.LastReportTime = tempTime
					qu.C(monthCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
				}
				explainacontainer = []meter.AggDayStruct{}
			}

			containerLastRecord = meter.AggDayStruct{}
		}

		fmt.Println(_EOF)
		qu.Close()
	}
}

//SetTimeStampForHour set minute second to 0
func SetTimeStampForHour(theTime time.Time) time.Time {
	hour := theTime.Hour()
	date := theTime.Day()
	month := theTime.Month()
	year := theTime.Year()
	return time.Date(year, month, date, hour, 0, 0, 0, time.Local)

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
	AggHour()
}
