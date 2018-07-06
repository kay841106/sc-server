package cron

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/meter"
)

// import (
// 	"gopkg.in/mgo.v2/bson"
// )

const (
	meterMapping              = "SC01_DeviceManager"
	displayDataCalcCollection = "SC01_displayData_Calc_"
	displayDataCollection     = "SC01_displayData_"
	hourCollection            = "SC01_hour_All"
	dayCollection             = "SC01_day_All"
	monthCollection           = "SC01_month_All"
	_EMPTYDEST                = "DESTINATION IS EMPTY"
	_AGG                      = "AGGREGATION ->"
	_NEC                      = "NON_EMPTY_COLL"
	_EOF                      = "END_OF_FILE"
)

func pipeDeviceHourWhole(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			}}}

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

// func (mycron LastStream) displayDataCalc() {
// 	//connect to DB
// 	var buildingID []interface{}
// 	var devID []interface{}

// 	// var thecount int
// 	mappingDevice := devMan{}
// 	mongo := meter.GetMongo()
// 	q := mongo.DB(meter.DBName)
// 	// displayDataMongos := displayDataMongo{}
// 	displayDataCalcMongos := displayDataCalcMongo{}
// 	currentdisplayDataCalcMongo := displayDataCalcMongo{}
// 	prevdisplayDataCalcMongo := displayDataCalcMongo{}
// 	templateDataCalcMongo := displayDataCalcMongo{}
// 	q.C(meterMapping).Find(nil).Distinct("Building_Name", &buildingID)
// 	log.Println(buildingID)
// 	for _, each := range buildingID {
// 		q.C(meterMapping).Find(bson.M{"Building_Name": each}).Distinct("devID", &devID)
// 		thecount, err := q.C("am" + each.(string)).Find(nil).Count()
// 		if err == nil {

// 			if thecount > 0 {
// 				log.Println(displayDataCalcCollection + each.(string))
// 				log.Println(_NEC)

// 				for _, each2 := range devID {

// 					q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
// 					q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)
// 					q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&currentdisplayDataCalcMongo)
// 					q.C("am" + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 					log.Println(currentdisplayDataCalcMongo)
// 					if (currentdisplayDataCalcMongo.LastReportTime != time.Time{}) {

// 						for displayDataCalcMongos.LastReportTime.Before(templateDataCalcMongo.LastReportTime) {
// 							// cond, _ := colour.Red("DEST < SRC")
// 							// displayDataCalcMongos = displayDataCalcMongo{}
// 							log.Println(_AGG, each.(string), each2.(string), "DEST < SRC")

// 							// tick := time.Now()
// 							q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$gt": displayDataCalcMongos.LastReportTime}}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)
// 							q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").Limit(1).One(&prevdisplayDataCalcMongo)

// 							fmt.Println(displayDataCalcMongos.LastReportTime, currentdisplayDataCalcMongo.LastReportTime, prevdisplayDataCalcMongo.LastReportTime)
// 							currentdisplayDataCalcMongo.Usage = currentdisplayDataCalcMongo.PwrUsage - prevdisplayDataCalcMongo.PwrUsage
// 							// fmt.Println(currentdisplayDataCalcMongo.Usage)
// 							currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 							currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
// 							currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
// 							currentdisplayDataCalcMongo.CC = 4950
// 							currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 							currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
// 							currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
// 							currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID
// 							q.C("am"+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)

// 							currentdisplayDataCalcMongo = displayDataCalcMongo{}
// 							prevdisplayDataCalcMongo = displayDataCalcMongo{}
// 							q.C("am" + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 							fmt.Println(displayDataCalcMongos.LastReportTime)
// 						}
// 						mappingDevice = devMan{}
// 					}
// 					continue
// 				}

// 			} else {
// 				log.Println(_EMPTYDEST)
// 				displayDataCalcMongos = displayDataCalcMongo{}

// 				for _, each2 := range devID {
// 					q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
// 					// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)

// 					// q.C("am" + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 					q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)

// 					// fmt.Println(displayDataCollection + each.(string)
// 					// q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2}).Sort("-lastReportTime").Limit(1).One(&displayDataCalcMongos)
// 					// fmt.Println(displayDataMongo)
// 					// cond, _ := colour.Red("DEST < SRC")
// 					log.Println(_AGG, each.(string), each2.(string))
// 					fmt.Println(currentdisplayDataCalcMongo.LastReportTime)
// 					// tick := time.Now()
// 					// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": displayDataMongos.LastReportTime}).Sort("-lastReportTime").Limit(1).One(&currentdisplayDataCalcMongo)
// 					// fmt.Println(currentdisplayDataCalcMongo)

// 					q.C("am" + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").Limit(1).One(&prevdisplayDataCalcMongo)
// 					if (prevdisplayDataCalcMongo == displayDataCalcMongo{}) {
// 						currentdisplayDataCalcMongo.Usage = 0
// 					}
// 					currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 					currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
// 					currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
// 					currentdisplayDataCalcMongo.CC = 4950
// 					currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 					currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
// 					currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
// 					currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID

// 					fmt.Println(currentdisplayDataCalcMongo)
// 					q.C("am"+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)
// 					mappingDevice = devMan{}
// 					// displayDataMongos = displayDataMongo{}
// 					currentdisplayDataCalcMongo = displayDataCalcMongo{}
// 					prevdisplayDataCalcMongo = displayDataCalcMongo{}

// 				}

// 			}

// 		}
// 	}
// }

func (mycron LastStream) AggHour() {
	status := checkDBStatus()
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := aggHourStruct{}
		var containerdevMan []interface{}
		var containerdevID []interface{}
		explainacontainer := []aggHourStruct{}

		qu := mongo.DB(meter.DBName)
		qu.C(CdevMan).Find(nil).Distinct("Building_Name", &containerdevMan)

		// fmt.Println(containerdevMan)

		for _, one := range containerdevMan {
			qu.C(CdevMan).Find(bson.M{"Building_Name": one}).Distinct("devID", &containerdevID)
			// fmt.Println(containerdevID)
			for _, two := range containerdevID {

				fmt.Println("two" + two.(string))
				qu.C(hourCollection).Find(bson.M{"Device_ID": two}).Limit(1).Sort("-lastReportTime").One(&containerLastRecord)
				fmt.Println(containerLastRecord.DeviceID)
				// fmt.Println(containerLastRecord)
				if (aggHourStruct{}) == containerLastRecord {

					thepipe := pipeDeviceHourWhole(two.(string))
					qu.C(displayDataCalcCollection + one.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

					for _, each := range explainacontainer {
						fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
						tempTime := SetTimeStampForHour(each.LastReportTime)
						each.LastReportTime = tempTime
						each.PFLimit = 0.8

						qu.C(hourCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
					}
				} else {
					thepipe := pipeDeviceHour(containerLastRecord.LastReportTime, containerLastRecord.DeviceID)
					qu.C(displayDataCalcCollection + one.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
					fmt.Println(explainacontainer)
					for _, each := range explainacontainer {
						fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
						tempTime := SetTimeStampForHour(each.LastReportTime)
						each.LastReportTime = tempTime
						each.PFLimit = 0.8
						qu.C(hourCollection).Upsert(bson.M{"Device_ID": each.DeviceID, "lastReportTime": each.LastReportTime}, each)
					}
				}
				explainacontainer = []aggHourStruct{}
				containerLastRecord = aggHourStruct{}
			}
		}

		fmt.Println(_EOF)
		qu.Close()
	}
}

func (mycron LastStream) AggDay() {
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
func (mycron LastStream) AggMonth() {
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

// func (mycron LastStream) AggDay() {
// 	mongo := meter.GetMongo()
// 	// explaina := bson.M{}
// 	var containerdevMan []interface{}
// 	explainacontainer := []aggHourStruct{}
// 	thepipe := pipeDeviceHour()
// 	qu := mongo.DB(meter.DBName)
// 	qu.C(CdevMan).Find(nil).Distinct("Building_Name", &containerdevMan)
// 	// fmt.Println(containerdevMan)
// 	for _, j := range containerdevMan {
// 		qu.C(displayDataCalcCollection + j.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
// 		for _, k := range explainacontainer {

// 			for _, l := range explainacontainer {
// 				fmt.Println(l.LastReportTime)
// 				fmt.Println(l.LastReportTime.Hour())
// 				timeStampRepl := SetTimeStampForHour(l.LastReportTime)
// 				l.LastReportTime = timeStampRepl
// 				qu.C(hourCollection).Upsert(bson.M{"Device_ID": k.DeviceID, "lastReportTime": k.LastReportTime}, l)
// 			}
// 		}

// 		fmt.Println(explainacontainer)
// 	}
// }

// func (mycron LastStream) AggMonth() {
// 	mongo := meter.GetMongo()
// 	// explaina := bson.M{}
// 	var containerdevMan []interface{}
// 	explainacontainer := []aggHourStruct{}
// 	thepipe := pipeDeviceHour()
// 	qu := mongo.DB(meter.DBName)
// 	qu.C(CdevMan).Find(nil).Distinct("Building_Name", &containerdevMan)
// 	// fmt.Println(containerdevMan)
// 	for _, j := range containerdevMan {
// 		qu.C(displayDataCalcCollection + j.(string)).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
// 		for _, k := range explainacontainer {

// 			for _, l := range explainacontainer {
// 				fmt.Println(l.LastReportTime)
// 				fmt.Println(l.LastReportTime.Hour())
// 				timeStampRepl := SetTimeStampForHour(l.LastReportTime)
// 				l.LastReportTime = timeStampRepl
// 				qu.C(hourCollection).Upsert(bson.M{"Device_ID": k.DeviceID, "lastReportTime": k.LastReportTime}, l)
// 			}
// 		}

// 		fmt.Println(explainacontainer)
// 	}
// }

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
