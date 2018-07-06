package cron

import (
	"fmt"
	"log"
	"time"

	"dforcepro.com/cron"
	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/airbox"
	alert "sc.dforcepro.com/alert"
	"sc.dforcepro.com/meter"
)

//Airbox cron
type Airbox bool

//Enable Airbox cron
func (mycron Airbox) Enable() bool {
	return bool(mycron)
}

//GetJobs cron
func (mycron Airbox) GetJobs() []cron.JobSpec {
	return []cron.JobSpec{
		cron.JobSpec{
			// 	Spec: "* * * * * *",
			// 	Job:  mycron.airboxStreamAll,
			// },
			// {
			// Spec: "*/30 * * * * *",
			Spec: "0 0 1 * * *",
			Job:  mycron.AggHourAirbox,
		},
		{
			Spec: "0 0 1 1 * *",
			// Spec: "*/30 * * * * *",
			Job: mycron.AggDayAirbox,
		},

		{
			Spec: "0 0 1 1 1 *",
			// Spec: "*/30 * * * * *",
			Job: mycron.AggMonthAirbox,
		},
	}
}

func pipeDeviceHourAirbox(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",
				"hour":      bson.M{"$hour": "$Upload_Time"},
				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month":      bson.M{"$month": "$lastReportTime"},
				"year": bson.M{"$year": "$Upload_Time"},
				"day":  bson.M{"$dayOfYear": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}

func pipeDeviceHourTimeAirbox(start time.Time, devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
				"Upload_Time": bson.M{
					"$gte": start,
				},
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",
				"hour":      bson.M{"$hour": "$Upload_Time"},
				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month":      bson.M{"$month": "$lastReportTime"},
				"year": bson.M{"$year": "$Upload_Time"},
				"day":  bson.M{"$dayOfYear": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}

func pipeDeviceDayAirbox(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",
				//"hour":      bson.M{"$hour": "$Upload_Time"},
				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month":      bson.M{"$month": "$lastReportTime"},
				"year": bson.M{"$year": "$Upload_Time"},
				"day":  bson.M{"$dayOfYear": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}
func pipeDeviceDayTimeAirbox(start time.Time, devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
				"Upload_Time": bson.M{
					"$gte": start,
				},
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",

				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month":      bson.M{"$month": "$lastReportTime"},
				"year": bson.M{"$year": "$Upload_Time"},
				"day":  bson.M{"$dayOfYear": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}

func pipeDeviceMonthAirbox(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",

				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				"Month": bson.M{"$month": "$lastReportTime"},
				"year":  bson.M{"$year": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}
func pipeDeviceMonthTimeAirbox(start time.Time, devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
				"Upload_Time": bson.M{
					"$gte": start,
				},
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",

				"Month": bson.M{"$month": "$lastReportTime"},
				"year":  bson.M{"$year": "$Upload_Time"},
			},
			"Upload_Time": bson.M{"$last": "$Upload_Time"},
			"PM2_5":       bson.M{"$avg": "$PM2_5"},
			"CO":          bson.M{"$avg": "$CO"},
			"CO2":         bson.M{"$avg": "$CO2"},
			"Noise":       bson.M{"$avg": "$Noise"},
			"Temp":        bson.M{"$avg": "$Temp"},
			"Humidity":    bson.M{"$avg": "$Humidity"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID": "$_id.Device_ID",

			"Upload_Time": 1,
			"PM2_5":       1,
			"CO":          1,
			"CO2":         1,
			"Noise":       1,
			"Temp":        1,
			"Humidity":    1,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"Upload_Time": 1},
	})

	return pipeline
}

func (mycron Airbox) airboxStreamAll() {
	fmt.Println("Start")
	mappingTemporaryData := airbox.MappingTemporaryData{}
	status := checkDBStatus()
	if status == true {
		Mongo := meter.GetMongo()
		// Mongo.DB(meter.DBName).C(streamAllCollection).RemoveAll(bson.M{"lastReportTime": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		container := airbox.CroncheckStatus{}
		// containerArray := []airbox.ScheckStatus{}
		fmt.Println("pass")
		devID := []string{}

		Mongo.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{}).Distinct("DevID", &devID)
		if devID == nil {
			log.Panic()
		}

		for _, j := range devID {
			fmt.Println(j)
			Mongo.DB(airbox.DBName).C(airbox.AirboxRaw).Find(bson.M{"Device_ID": j}).Sort("-Upload_Time").One(&container)
			Mongo.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": j}).One(&mappingTemporaryData)
			fmt.Println(container)
			container.Location = mappingTemporaryData.Location
			Mongo.DB(airbox.DBName).C(airbox.AirboxStream).Upsert(bson.M{"Device_ID": j}, container)
			mappingTemporaryData = airbox.MappingTemporaryData{}
			container = airbox.CroncheckStatus{}

			// Mongo.DB(meter.DBName).C(streamCollection).Find(bson.M{"Device_ID": bson.M{"$in": devID}}).All(&container)

			// for _, j := range container {

			// 	Mongo.DB(meter.DBName).C(streamAllCollection).Upsert(bson.M{"Device_ID": j.DeviceID}, j)
			// 	//reset struct
			// 	j = meter.DisplayDataElement2nd{}
			// }

			fmt.Println(container)

		}
	}
}

//AggHourAirbox aggregate hourly data
func (mycron Airbox) AggHourAirbox() {
	status := checkDBStatus()
	mappingTemporaryData := airbox.MappingTemporaryData{}
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := airbox.CroncheckStatus{}
		var containerdevID []interface{}
		explainacontainer := []airbox.CroncheckStatus{}

		qu := mongo.DB(DBAirbox)
		qu.C(cAirbox).Find(nil).Distinct("Device_ID", &containerdevID)
		// fmt.Println(containerdevMan)

		for _, two := range containerdevID {

			fmt.Println("two" + two.(string))
			qu.C(AirboxHourCollection).Find(bson.M{"Device_ID": two.(string)}).Sort("-Upload_Time").One(&containerLastRecord)
			fmt.Println(containerLastRecord)
			// fmt.Println(containerLastRecord)
			if (airbox.CroncheckStatus{}) == containerLastRecord {

				thepipe := pipeDeviceHourAirbox(two.(string))
				qu.C(cAirbox).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

				for _, each := range explainacontainer {

					tempTime := SetTimeStampForHour(each.LastReportTime)
					each.LastReportTime = tempTime
					// fmt.Println(each.DevID, each.LastReportTime)
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxHourCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
				}
			} else {
				thepipe := pipeDeviceHourTimeAirbox(containerLastRecord.LastReportTime, containerLastRecord.DevID)
				qu.C(cAirbox).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				// fmt.Println(explainacontainer)
				for _, each := range explainacontainer {

					tempTime := SetTimeStampForHour(each.LastReportTime)
					each.LastReportTime = tempTime
					fmt.Println(each.DevID, each.LastReportTime)
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxHourCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
				}
			}
			explainacontainer = []airbox.CroncheckStatus{}
			containerLastRecord = airbox.CroncheckStatus{}
		}

		log.Println(_EOF)
		qu.Close()
	}
}

//AggDayAirbox aggregate daily data
func (mycron Airbox) AggDayAirbox() {
	mappingTemporaryData := airbox.MappingTemporaryData{}
	status := checkDBStatus()
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := airbox.CroncheckStatus{}
		var containerdevID []interface{}
		explainacontainer := []airbox.CroncheckStatus{}

		qu := mongo.DB(DBAirbox)
		qu.C(cAirbox).Find(nil).Distinct("Device_ID", &containerdevID)

		// fmt.Println(containerdevMan)
		for _, two := range containerdevID {

			fmt.Println(two)
			qu.C(AirboxDayCollection).Find(bson.M{"Device_ID": two.(string)}).Sort("-Upload_Time").One(&containerLastRecord)
			// fmt.Println("devID_query = ", containerLastRecord.DeviceID, -two)

			if (airbox.CroncheckStatus{}) == containerLastRecord {

				thepipe := pipeDeviceDayAirbox(two.(string))
				qu.C(AirboxHourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForDay(each.LastReportTime)
					each.LastReportTime = tempTime
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxDayCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
					// each = meter.AggDayStruct{}
				}
				explainacontainer = []airbox.CroncheckStatus{}

			} else {

				thepipe := pipeDeviceDayTimeAirbox(containerLastRecord.LastReportTime, containerLastRecord.DevID)
				qu.C(AirboxHourCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					// fmt.Println(each.BuildingName, each.DeviceID, each.LastReportTime)
					tempTime := SetTimeStampForDay(each.LastReportTime)
					each.LastReportTime = tempTime
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxDayCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
					// each = meter.AggDayStruct{}
				}

				explainacontainer = []airbox.CroncheckStatus{}

			}

			containerLastRecord = airbox.CroncheckStatus{}
		}
		log.Println("AggDayAirbox", _EOF)
		qu.Close()
	}
}

//AggMonthAirbox aggregate monthly data
func (mycron Airbox) AggMonthAirbox() {
	mappingTemporaryData := airbox.MappingTemporaryData{}
	status := checkDBStatus()
	if status == true {
		mongo := meter.GetMongo()
		containerLastRecord := airbox.CroncheckStatus{}
		var containerdevMan []interface{}
		explainacontainer := []airbox.CroncheckStatus{}

		qu := mongo.DB(DBAirbox)
		qu.C(cAirbox).Find(nil).Distinct("Device_ID", &containerdevMan)

		// fmt.Println(containerdevMan)

		for _, two := range containerdevMan {

			fmt.Println(two)
			qu.C(AirboxMonthCollection).Find(bson.M{"Device_ID": two.(string)}).Sort("-Upload_Time").One(&containerLastRecord)
			fmt.Println(containerLastRecord.DevID)
			if (airbox.CroncheckStatus{}) == containerLastRecord {

				thepipe := pipeDeviceMonthAirbox(two.(string))
				qu.C(AirboxDayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)

				for _, each := range explainacontainer {
					fmt.Println(each.DevID, each.LastReportTime)
					tempTime := SetTimeStampForMonth(each.LastReportTime)
					fmt.Println(tempTime)
					each.LastReportTime = tempTime
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxMonthCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
				}
			} else {
				thepipe := pipeDeviceMonthTimeAirbox(containerLastRecord.LastReportTime, containerLastRecord.DevID)
				qu.C(AirboxDayCollection).Pipe(thepipe).AllowDiskUse().All(&explainacontainer)
				fmt.Println(explainacontainer)
				for _, each := range explainacontainer {
					fmt.Println(each.DevID, each.LastReportTime)
					tempTime := SetTimeStampForMonth(each.LastReportTime)

					each.LastReportTime = tempTime
					qu.DB(airbox.DBName).C(MappingAirbox).Find(bson.M{"DevID": each.DevID}).One(&mappingTemporaryData)
					each.Location = mappingTemporaryData.Location
					qu.C(AirboxMonthCollection).Upsert(bson.M{"Device_ID": each.DevID, "Upload_Time": each.LastReportTime}, each)
					mappingTemporaryData = airbox.MappingTemporaryData{}
				}
				explainacontainer = []airbox.CroncheckStatus{}

			}

			containerLastRecord = airbox.CroncheckStatus{}

		}

		log.Println("AggMonthAirbox", _EOF)
		qu.Close()
	}
}

func (mycron LastStream) AirboxupdateError() {
	status := checkDBStatus()
	if status == true {
		fmt.Println("hello")

		mongo := meter.GetMongo()

		container := []airbox.AirboxPayload{}
		q := mongo.DB(DBAirbox)
		q.C(AirboxStatus).RemoveAll(bson.M{"Upload_Time": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		// cont := []alert.AlertValue{}
		cont2 := alert.AirboxAlertValue{}
		// tmpContainer := meter.DisplayData{}
		// tmpContainerB := meter.DisplayDataElement{}
		// containerNull := []meter.DisplayDataElement2nd{}
		// containerNull := []AirboxDeviceMapping{}
		var tmpDevID []string

		// Mongo := meter.GetMongo()

		// var buildContainer []interface{}

		q.C(MappingAirbox).Find(bson.M{}).Distinct("Device_ID", &tmpDevID)

		// fmt.Println(tmpDevID)
		for _, elementID := range tmpDevID {

			// mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)
			q.C(cAirbox).Find(bson.M{"Device_ID": elementID}).Limit(1).Sort("-Upload_Time").All(&container)
			// fmt.Println(container)

			for _, j := range container {

				// times := 0
				ho, mo, se, status := alert.CompareTimeStamp(j.LastReportTime)
				// tmpContainer.Rows = append(tmpContainer.Rows, j)
				cont2.DeviceID = j.DeviceID
				cont2.OnlineStatus = status
				cont2.DownTime.HoursD = ho
				cont2.DownTime.MinutesD = mo
				cont2.DownTime.SecondD = se
				cont2.UploadTime = j.LastReportTime

				// append(cont, cont2)
				q.C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
				cont2 = alert.AirboxAlertValue{}

				// mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, j)

			}

			// } else {
			// 	q.C(meter.DeviceManager).Find(bson.M{"devID": elementID}).All(&containerNull)
			// 	for _, detailDevice := range containerNull {
			// 		fmt.Println(detailDevice)
			// 		cont2.DeviceID = elementID
			// 		cont2.OnlineStatus = false
			// 		cont2.DeviceDetails.BuildingDetails = detailDevice.BuildingDetails
			// 		cont2.DeviceDetails.BuildingName = detailDevice.BuildingName
			// 		cont2.DeviceDetails.GatewayID = detailDevice.GatewayID
			// 		cont2.DownTime.HoursD = 99999
			// 		cont2.DownTime.MinutesD = 99999
			// 		cont2.DownTime.SecondD = 99999
			// 		cont2.LastReportTime = time.Time{}

			// 		q.C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
			// 		cont2 = alert.AlertValue{}

			// 	}

			// }

			fmt.Println("finish")
		}

	}
}
