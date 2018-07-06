package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"dforcepro.com/cron"
	"gopkg.in/mgo.v2/bson"
	alert "sc.dforcepro.com/alert"
	"sc.dforcepro.com/meter"
)

const (
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

type LastStream bool

func (mycron LastStream) Enable() bool {
	return bool(mycron)
}

func (mycron LastStream) GetJobs() []cron.JobSpec {
	return []cron.JobSpec{
		cron.JobSpec{
			Spec: "*/30 * * * * *",
			Job:  mycron.updateStream,
		},
		{
			Spec: "*/30 * * * * *",
			Job:  mycron.updateError,
		},
		{
			// Spec: "*/30 * * * * *",
			Spec: "0 1 * * * *",
			Job:  mycron.AggHour,
		},
		{
			Spec: "0 1 1 * * *",
			Job:  mycron.AggDay,
		},
		{
			Spec: "0 1 1 1 * *",
			Job:  mycron.AggMonth,
		},
		// {
		// 	Spec: "*/30 * * * * *",
		// 	Job:  mycron.displayDataCalc,
		// },

		{
			Spec: "*/45 * * * * *",
			Job:  mycron.updateStreamAll,
		},
	}
}
func (mycron LastStream) updateStreamAll() {
	status := checkDBStatus()
	if status == true {
		Mongo := meter.GetMongo()
		// Mongo.DB(meter.DBName).C(streamAllCollection).RemoveAll(bson.M{"lastReportTime": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		container := []meter.DisplayDataElement2nd{}
		container_devman := []meter.DeviceManagerS{}
		devID := []string{}
		Mongo.DB(meter.DBName).C(CdevMan).Find(nil).All(&container_devman)
		for _, j := range container_devman {
			if strings.Contains(j.DeviceInfo, "- Total") || strings.Contains(j.DeviceInfo, "-total") {
				devID = append(devID, j.DeviceID)
				fmt.Println("UPDATE STREAM ALL", j.DeviceInfo)
			}
			fmt.Println("NOT STREAM ALL", j.DeviceInfo)
		}
		// devID := [9]string{"33000509b52f0201", "33000509b52f5a01", "33000509b52f3501", "33000509b53b1801", "33000509b52f1001", "33000509b53b4901", "33000509b53b4902", "33000509b53b7901", "33000509b52f2101", "33000509b53b0501"}

		Mongo.DB(meter.DBName).C(streamCollection).Find(bson.M{"Device_ID": bson.M{"$in": devID}}).All(&container)

		for _, j := range container {

			Mongo.DB(meter.DBName).C(streamAllCollection).Upsert(bson.M{"Device_ID": j.DeviceID}, j)
			//reset struct
			j = meter.DisplayDataElement2nd{}
		}

		fmt.Println(container)
		Mongo.Close()

		// container:=
	}
}
func (mycron LastStream) updateStream() {
	status := checkDBStatus()
	if status == true {
		fmt.Println("hello")
		Mongo := meter.GetMongo()
		Mongo.DB(meter.DBName).C(streamCollection).RemoveAll(bson.M{"lastReportTime": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		container := []meter.DisplayDataElement2nd{}
		containerSingle := meter.DisplayDataElement2nd{}
		devMan := []meter.DeviceManagerS{}
		weatherS := []WeatherGetS{}
		// cont2 := alert.AlertValue{}
		// tmpContainer := []meter.DisplayDataElement2nd{}
		// tmpContainerB := meter.DisplayDataElement{}
		// containerNull := []meter.DisplayDataElement2nd{}
		var tmpDevID []string

		var buildContainer []interface{}

		Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{}).Distinct("Building_Name", &buildContainer)

		for _, element := range buildContainer {
			Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)

			for _, elementID := range tmpDevID {

				Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)

				count, _ := Mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Count()

				if count != 0 {
					Mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Limit(1).Sort("-lastReportTime").All(&container)
					// fmt.Println(container)

					for _, j := range container {
						// Mongo.DB(meter.DBName).C(weatherCollection).Find(bson.M{"WeatherElement": "T", "startTime": bson.M{"$lte": j.LastReportTime}, "endTime": bson.M{"$gte": j.LastReportTime}}).All(&weatherS)
						Mongo.DB(meter.DBName).C(weatherCollection).Find(bson.M{"WeatherElement": "T", "startTime": bson.M{"$lte": j.LastReportTime}, "endTime": bson.M{"$gte": j.LastReportTime}}).Limit(1).All(&weatherS)
						fmt.Println(weatherS)
						for _, k := range weatherS {
							tmpValWeather, _ := strconv.Atoi(k.ElementValue)
							containerSingle.WeatherTemp = tmpValWeather
						}
						// times := 0
						containerSingle.LastReportTime = j.LastReportTime
						containerSingle.AvgPF = j.AvgPF
						containerSingle.BuildingDetails = j.BuildingDetails
						containerSingle.BuildingName = j.BuildingName
						containerSingle.DeviceID = elementID
						containerSingle.GatewayID = j.GatewayID
						containerSingle.PwrDemand = j.PwrDemand / 1000 //KW
						containerSingle.PwrUsage = j.Usage
						// containerSingle.PwrUsage = j.PwrUsage / 1000   //KW

						// append(cont, cont2)
						Mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, containerSingle)
						containerSingle = meter.DisplayDataElement2nd{}
						// mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, j)
						if j.BuildingDetails != "" {
							fmt.Println(j.BuildingName)
						}

					}

				} else {
					Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"devID": elementID}).All(&devMan)
					// fmt.Println(devMan)
					for _, detailDevice := range devMan {

						containerSingle.AvgPF = 0
						containerSingle.BuildingDetails = detailDevice.BuildingDetails
						containerSingle.BuildingName = detailDevice.BuildingName
						containerSingle.DeviceID = detailDevice.DeviceID
						containerSingle.GatewayID = detailDevice.GatewayID
						containerSingle.PwrDemand = 0
						containerSingle.PwrUsage = 0
						containerSingle.LastReportTime = time.Time{}
						Mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, containerSingle)
						containerSingle = meter.DisplayDataElement2nd{}

					}

				}

			}

			fmt.Println("finish")
			Mongo.Close()
		}

	}
}

func (mycron LastStream) updateError() {
	status := checkDBStatus()
	if status == true {
		fmt.Println("hello")

		mongo := meter.GetMongo()

		container := []meter.DisplayDataElement2nd{}
		mongo.DB(meter.DBName).C(statusCollection).RemoveAll(bson.M{"lastReportTime": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		// cont := []alert.AlertValue{}
		cont2 := alert.AlertValue{}
		// tmpContainer := meter.DisplayData{}
		// tmpContainerB := meter.DisplayDataElement{}
		// containerNull := []meter.DisplayDataElement2nd{}
		containerNull := []meter.DeviceManagerS{}
		var tmpDevID []string

		Mongo := meter.GetMongo()

		var buildContainer []interface{}

		Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{}).Distinct("Building_Name", &buildContainer)

		for _, element := range buildContainer {
			Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)
			// fmt.Println(tmpDevID)
			for _, elementID := range tmpDevID {

				// mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)

				count, _ := mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Count()

				if count != 0 {

					mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Limit(1).Sort("-lastReportTime").All(&container)
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
						cont2.LastReportTime = j.LastReportTime
						cont2.DeviceDetails.BuildingDetails = j.BuildingDetails
						cont2.DeviceDetails.BuildingName = j.BuildingName
						cont2.DeviceDetails.GatewayID = j.GatewayID

						// append(cont, cont2)
						mongo.DB(meter.DBName).C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
						cont2 = alert.AlertValue{}

						// mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, j)
						if j.BuildingDetails == "" {
							fmt.Println(j.BuildingName)
						}

					}

				} else {
					mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"devID": elementID}).All(&containerNull)
					for _, detailDevice := range containerNull {
						fmt.Println(detailDevice)
						cont2.DeviceID = elementID
						cont2.OnlineStatus = false
						cont2.DeviceDetails.BuildingDetails = detailDevice.BuildingDetails
						cont2.DeviceDetails.BuildingName = detailDevice.BuildingName
						cont2.DeviceDetails.GatewayID = detailDevice.GatewayID
						cont2.DownTime.HoursD = 99999
						cont2.DownTime.MinutesD = 99999
						cont2.DownTime.SecondD = 99999
						cont2.LastReportTime = time.Time{}

						mongo.DB(meter.DBName).C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
						cont2 = alert.AlertValue{}

					}

				}

			}

			fmt.Println("finish")
			mongo.Close()
		}

	}
}

func (mycron LastStream) detailStatus() {
	status := checkDBStatus()
	if status == true {
		fmt.Println("detailStatus")

		mongo := meter.GetMongo()

		container := []meter.DisplayDataElement2nd{}
		mongo.DB(meter.DBName).C(statusCollection).RemoveAll(bson.M{"lastReportTime": bson.M{"$lt": time.Now().Add(time.Duration(-1) * time.Hour)}})
		// cont := []alert.AlertValue{}
		cont2 := alert.AlertValue{}
		// tmpContainer := meter.DisplayData{}
		// tmpContainerB := meter.DisplayDataElement{}
		// containerNull := []meter.DisplayDataElement2nd{}
		containerNull := []meter.DeviceManagerS{}
		var tmpDevID []string

		Mongo := meter.GetMongo()

		var buildContainer []interface{}

		Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{}).Distinct("Building_Name", &buildContainer)

		for _, element := range buildContainer {
			Mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)
			// fmt.Println(tmpDevID)
			for _, elementID := range tmpDevID {

				// mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"Building_Name": element}).Distinct("devID", &tmpDevID)

				count, _ := mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Count()

				if count != 0 {

					mongo.DB(meter.DBName).C(coll + element.(string)).Find(bson.M{"Device_ID": elementID}).Limit(1).Sort("-lastReportTime").All(&container)
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
						cont2.LastReportTime = j.LastReportTime
						cont2.DeviceDetails.BuildingDetails = j.BuildingDetails
						cont2.DeviceDetails.BuildingName = j.BuildingName
						cont2.DeviceDetails.GatewayID = j.GatewayID

						// append(cont, cont2)
						mongo.DB(meter.DBName).C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
						cont2 = alert.AlertValue{}

						// mongo.DB(meter.DBName).C(streamCollection).Upsert(bson.M{"Device_ID": elementID}, j)
						if j.BuildingDetails == "" {
							fmt.Println(j.BuildingName)
						}

					}

				} else {
					mongo.DB(meter.DBName).C(meter.DeviceManager).Find(bson.M{"devID": elementID}).All(&containerNull)
					for _, detailDevice := range containerNull {
						fmt.Println(detailDevice)
						cont2.DeviceID = elementID
						cont2.OnlineStatus = false
						cont2.DeviceDetails.BuildingDetails = detailDevice.BuildingDetails
						cont2.DeviceDetails.BuildingName = detailDevice.BuildingName
						cont2.DeviceDetails.GatewayID = detailDevice.GatewayID
						cont2.DownTime.HoursD = 99999
						cont2.DownTime.MinutesD = 99999
						cont2.DownTime.SecondD = 99999
						cont2.LastReportTime = time.Time{}

						mongo.DB(meter.DBName).C(statusCollection).Upsert(bson.M{"Device_ID": elementID}, cont2)
						cont2 = alert.AlertValue{}

					}

				}

			}

			fmt.Println("finish")
			mongo.Close()
		}

	}
}
