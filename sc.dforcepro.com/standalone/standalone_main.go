package standalone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/airbox"
	"sc.dforcepro.com/meter"
)

// //SCStandAlone Building
// type SCStandAlone bool

// //Enable the API
// func (sca SCStandAlone) Enable() bool {
// 	return bool(sca)
// }

func checkDBStatus() bool {
	err := meter.GetMongo().Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		time.Sleep(5 * time.Second)
		meter.GetMongo().Refresh()
		err = meter.GetMongo().Ping()
	}
	return true
}

func DisplayDataCalc() {
	for {

		checkDBStatus()
		//connect to DB
		var buildingID []interface{}
		var devID []interface{}

		// var thecount int
		mappingDevice := devMan{}
		mongo := meter.GetMongo()
		q := mongo.DB(meter.DBName)
		// displayDataMongos := displayDataMongo{}
		displayDataCalcMongos := &displayDataCalcMongo{}
		currentdisplayDataCalcMongo := &displayDataCalcMongo{}
		prevdisplayDataCalcMongo := &displayDataCalcMongo{}
		templateDataCalcMongo := &displayDataCalcMongo{}
		q.C(meterMapping).Find(nil).Distinct("Building_Name", &buildingID)
		log.Println(buildingID)
		for _, each := range buildingID {
			q.C(meterMapping).Find(bson.M{"Building_Name": each}).Distinct("devID", &devID)
			thecount, err := q.C(displayDataCalcCollection + each.(string)).Find(nil).Count()
			if err == nil {

				if thecount > 0 {
					log.Println(displayDataCalcCollection + each.(string))
					log.Println(_NEC)

					for _, each2 := range devID {

						q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)
						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&currentdisplayDataCalcMongo)
						q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
						// log.Println(currentdisplayDataCalcMongo)
						if (currentdisplayDataCalcMongo.LastReportTime != time.Time{}) {

							for displayDataCalcMongos.LastReportTime.Before(templateDataCalcMongo.LastReportTime) {

								// cond, _ := colour.Red("DEST < SRC")
								// displayDataCalcMongos = displayDataCalcMongo{}
								log.Println(_AGG, each.(string), each2.(string), "DEST < SRC")

								// tick := time.Now()
								q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$gt": displayDataCalcMongos.LastReportTime}}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)
								q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").One(&prevdisplayDataCalcMongo)

								fmt.Println(displayDataCalcMongos.LastReportTime, currentdisplayDataCalcMongo.LastReportTime, "prev", prevdisplayDataCalcMongo.LastReportTime)
								currentdisplayDataCalcMongo.Usage = currentdisplayDataCalcMongo.PwrUsage - prevdisplayDataCalcMongo.PwrUsage
								// fmt.Println(currentdisplayDataCalcMongo.Usage)
								currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
								currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
								currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
								currentdisplayDataCalcMongo.CC = 4950
								currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
								currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
								currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
								currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID
								q.C(displayDataCalcCollection+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)

								// currentdisplayDataCalcMongo = displayDataCalcMongo{}
								// prevdisplayDataCalcMongo = displayDataCalcMongo{}
								q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
								fmt.Println(displayDataCalcMongos.LastReportTime)
							}
							mappingDevice = devMan{}
							// displayDataCalcMongos = displayDataCalcMongo{}
							// currentdisplayDataCalcMongo = displayDataCalcMongo{}
						}
						continue
					}

				} else {
					log.Println(_EMPTYDEST)
					// displayDataCalcMongos = displayDataCalcMongo{}

					for _, each2 := range devID {
						q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
						// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)

						// q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)

						// fmt.Println(displayDataCollection + each.(string)
						// q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2}).Sort("-lastReportTime").Limit(1).One(&displayDataCalcMongos)
						// fmt.Println(displayDataMongo)
						// cond, _ := colour.Red("DEST < SRC")
						log.Println(_AGG, each.(string), each2.(string))
						fmt.Println(currentdisplayDataCalcMongo.LastReportTime)
						// tick := time.Now()
						// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": displayDataMongos.LastReportTime}).Sort("-lastReportTime").Limit(1).One(&currentdisplayDataCalcMongo)
						// fmt.Println(currentdisplayDataCalcMongo)

						q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").Limit(1).One(&prevdisplayDataCalcMongo)
						if &prevdisplayDataCalcMongo == nil {
							currentdisplayDataCalcMongo.Usage = 0
						}
						currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
						currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
						currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
						currentdisplayDataCalcMongo.CC = 4950
						currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
						currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
						currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
						currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID

						fmt.Println(currentdisplayDataCalcMongo)
						q.C(displayDataCalcCollection+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)
						err := q.C(displayDataCalcCollection+each.(string)).EnsureIndexKey("lastReportTime", "Device_ID")
						err1 := q.C(displayDataCalcCollection + each.(string)).EnsureIndexKey("lastReportTime")
						err2 := q.C(displayDataCalcCollection + each.(string)).EnsureIndexKey("Device_ID")

						if err != nil {
							log.Println(err)
						} else if err1 != nil {
							log.Println(err)
						} else if err2 != nil {
							log.Println(err)

						}
						mappingDevice = devMan{}
						// displayDataMongos = displayDataMongo{}
						// currentdisplayDataCalcMongo = displayDataCalcMongo{}
						// prevdisplayDataCalcMongo = displayDataCalcMongo{}

					}

				}

			}
		}
		mongo.Close()
		time.Sleep(45 * time.Second)

	}

}

// func DisplayDataCalc2() {

// 	for {
// 		//connect to DB
// 		var buildingID []interface{}
// 		var devID []interface{}

// 		// var thecount int
// 		mappingDevice := devMan{}
// 		mongo := meter.GetMongo()
// 		q := mongo.DB(meter.DBName)
// 		// displayDataMongos := displayDataMongo{}
// 		displayDataCalcMongos := &displayDataCalcMongo{}
// 		currentdisplayDataCalcMongo := &displayDataCalcMongo{}
// 		prevdisplayDataCalcMongo := &displayDataCalcMongo{}
// 		templateDataCalcMongo := &displayDataCalcMongo{}
// 		q.C(meterMapping).Find(nil).Distinct("Building_Name", &buildingID)
// 		log.Println(buildingID)
// 		for _, each := range buildingID {
// 			q.C(meterMapping).Find(bson.M{"Building_Name": each}).Distinct("devID", &devID)
// 			thecount, err := q.C(displayDataCalcCollection + each.(string)).Find(nil).Count()
// 			if err == nil {

// 				if thecount > 0 {
// 					log.Println(displayDataCalcCollection + each.(string))
// 					log.Println(_NEC)

// 					for _, each2 := range devID {

// 						q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
// 						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)
// 						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&currentdisplayDataCalcMongo)
// 						q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 						// log.Println(each2.(string))
// 						if (currentdisplayDataCalcMongo.LastReportTime != time.Time{}) {

// 							for displayDataCalcMongos.LastReportTime.Before(templateDataCalcMongo.LastReportTime) {
// 								displayDataCalcMongos = displayDataCalcMongo{}
// 								currentdisplayDataCalcMongo = displayDataCalcMongo{}
// 								// cond, _ := colour.Red("DEST < SRC")
// 								// displayDataCalcMongos = displayDataCalcMongo{}
// 								log.Println(_AGG, each.(string), each2.(string), "DEST < SRC")

// 								// tick := time.Now()
// 								q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$gt": displayDataCalcMongos.LastReportTime}}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)
// 								q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").Limit(1).One(&prevdisplayDataCalcMongo)

// 								// fmt.Println(displayDataCalcMongos.LastReportTime, currentdisplayDataCalcMongo.LastReportTime, prevdisplayDataCalcMongo.LastReportTime)
// 								currentdisplayDataCalcMongo.Usage = currentdisplayDataCalcMongo.PwrUsage - prevdisplayDataCalcMongo.PwrUsage
// 								// fmt.Println(currentdisplayDataCalcMongo.Usage)
// 								currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 								currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
// 								currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
// 								currentdisplayDataCalcMongo.CC = 4950
// 								currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 								currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
// 								currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
// 								currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID
// 								q.C(displayDataCalcCollection+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)

// 								currentdisplayDataCalcMongo = displayDataCalcMongo{}
// 								prevdisplayDataCalcMongo = displayDataCalcMongo{}
// 								q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 								// fmt.Println(displayDataCalcMongos.LastReportTime)
// 							}
// 							mappingDevice = devMan{}

// 						}
// 						continue
// 					}

// 				} else {
// 					log.Println(_EMPTYDEST)
// 					displayDataCalcMongos = displayDataCalcMongo{}

// 					for _, each2 := range devID {
// 						q.C(meterMapping).Find(bson.M{"devID": each2}).One(&mappingDevice)
// 						// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&templateDataCalcMongo)

// 						// q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("-lastReportTime").One(&displayDataCalcMongos)
// 						q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string)}).Sort("lastReportTime").One(&currentdisplayDataCalcMongo)

// 						// fmt.Println(displayDataCollection + each.(string)
// 						// q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2}).Sort("-lastReportTime").Limit(1).One(&displayDataCalcMongos)
// 						// fmt.Println(displayDataMongo)
// 						// cond, _ := colour.Red("DEST < SRC")
// 						log.Println(_AGG, each.(string), each2.(string))
// 						fmt.Println(currentdisplayDataCalcMongo.LastReportTime)
// 						// tick := time.Now()
// 						// q.C(displayDataCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": displayDataMongos.LastReportTime}).Sort("-lastReportTime").Limit(1).One(&currentdisplayDataCalcMongo)
// 						// fmt.Println(currentdisplayDataCalcMongo)

// 						q.C(displayDataCalcCollection + each.(string)).Find(bson.M{"Device_ID": each2.(string), "lastReportTime": bson.M{"$lt": currentdisplayDataCalcMongo.LastReportTime}}).Sort("-lastReportTime").Limit(1).One(&prevdisplayDataCalcMongo)
// 						if (prevdisplayDataCalcMongo == displayDataCalcMongo{}) {
// 							currentdisplayDataCalcMongo.Usage = 0
// 						}
// 						currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 						currentdisplayDataCalcMongo.BuildingName = mappingDevice.BuildingName
// 						currentdisplayDataCalcMongo.BuildingDetails = mappingDevice.BuildingDetails
// 						currentdisplayDataCalcMongo.CC = 4950
// 						currentdisplayDataCalcMongo.DeviceDetails = mappingDevice.DeviceDetails
// 						currentdisplayDataCalcMongo.DeviceName = mappingDevice.DeviceName
// 						currentdisplayDataCalcMongo.DeviceType = mappingDevice.DeviceType
// 						currentdisplayDataCalcMongo.Floor = mappingDevice.GatewayID

// 						// fmt.Println(currentdisplayDataCalcMongo)
// 						q.C(displayDataCalcCollection+each.(string)).Upsert(bson.M{"Device_ID": currentdisplayDataCalcMongo.DeviceID, "lastReportTime": currentdisplayDataCalcMongo.LastReportTime}, currentdisplayDataCalcMongo)
// 						err := q.C(displayDataCalcCollection+each.(string)).EnsureIndexKey("lastReportTime", "Device_ID")
// 						err1 := q.C(displayDataCalcCollection + each.(string)).EnsureIndexKey("lastReportTime")
// 						err2 := q.C(displayDataCalcCollection + each.(string)).EnsureIndexKey("Device_ID")

// 						if err != nil {
// 							log.Println(err)
// 						} else if err1 != nil {
// 							log.Println(err)
// 						} else if err2 != nil {
// 							log.Println(err)

// 						}
// 						mappingDevice = devMan{}
// 						// displayDataMongos = displayDataMongo{}
// 						currentdisplayDataCalcMongo = displayDataCalcMongo{}
// 						prevdisplayDataCalcMongo = displayDataCalcMongo{}

// 					}

// 				}

// 			}
// 		}
// 		q.Close()
// 		time.Sleep(30 * time.Second)

// 	}

// }

func airboxTestingPOST(pm25 float64) {
	url := ""
	log.Println("URL:>", url)

	container := airboxTestingStruct{
		// PM25: r.Intn(50),
	}

	// var jsonStr = []byte(container)
	jsonStr, _ := json.Marshal(container)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "testingToPOST")
	req.Header.Set("Content-Type", "application/json")

}

// func aaaa() {
// 	//variable declaration
// 	var temporarybuilding []string
// 	var temporarydevID []string

// 	//struct declaration
// 	aemdraRawdataToDisplayData := rawdataToDisplayData{}
// 	cpmRawdataToDisplayData := rawdataToDisplayData{}

// 	// aemdraRawdataToDisplayDataArray := []rawdataToDisplayData{}
// 	cpmRawdataToDisplayDataArray := []rawdataToDisplayData{}

// 	//initial setup
// 	mongo := meter.GetMongo()
// 	q := mongo.DB(meter.DBName)
// 	coll, _ := q.CollectionNames()

// 	//distinct BuildingName
// 	q.C(meterMapping).Find(nil).Distinct("Building_Name", &temporarybuilding)

// 	//
// 	for _, i := range coll {
// 		for _, each := range temporarybuilding {

// 			// //check string contgain rawdata or not
// 			if contain := strings.Contains(rawDataAEMDRACollection+each, i); contain {

// 				// // //get device id
// 				q.C(meterMapping).Find(bson.M{"Building_Name": each, "Device_Type": aemdra}).Distinct("devID", &temporarydevID)

// 				for _, j := range temporarydevID {

// 					// // //querysource data
// 					q.C(rawDataAEMDRACollection + each).Find(bson.M{"devID": j}).Sort("-lastReportTime").One(&aemdraRawdataToDisplayData)

// 					//debug time
// 					fmt.Println(aemdraRawdataToDisplayData.LastReportTime)

// 					// // //querysource data
// 					q.C(displayDataCollection + each).Find(bson.M{"Device_ID": j})

// 				}
// 			} else if contain := strings.Contains(rawDataCPMCollection+each, i); contain {

// 				// // //get device id
// 				q.C(meterMapping).Find(bson.M{"Building_Name": each, "Device_Type": cpm}).Distinct("devID", &temporarydevID)

// 				for _, j := range temporarydevID {

// 					// // //querysource data
// 					q.C(rawDataCPMCollection + each).Find(bson.M{"devID": j}).Sort("-lastReportTime").One(&cpmRawdataToDisplayData)

// 					// // //querysource data
// 					q.C(displayDataCollection + each).Find(bson.M{"Device_ID": j}).One(&cpmRawdataToDisplayDataArray)

// 				}
// 			}
// 		}

// 		//compare if its timestamp not None

// 		//if destination < source do aggregate

// 		// //pipeline

// 		//aggregate weather inside

// 		//compare destination is None but source not None
// 		// // it means new Collection

// 		// //pipeline

// 	}
// }

func AirboxStreamAll() {
	for {
		fmt.Println("Start")
		mappingTemporaryData := airbox.MappingTemporaryData{}
		status := checkDBStatus()
		if status == true {
			Mongo := meter.GetMongo()

			container := airbox.CroncheckStatus{}

			fmt.Println("pass")
			devID := []string{}
			q := Mongo.DB(airbox.DBName)
			q.C(MappingAirbox).Find(bson.M{}).Distinct("DevID", &devID)
			if devID == nil {
				log.Panic()
			}

			for _, j := range devID {
				fmt.Println(j)
				q.C(airbox.AirboxRaw).Find(bson.M{"Device_ID": j}).Sort("-Upload_Time").One(&container)
				q.C(MappingAirbox).Find(bson.M{"DevID": j}).One(&mappingTemporaryData)
				fmt.Println(container)
				container.Location = mappingTemporaryData.Location
				q.C(airbox.AirboxStream).Upsert(bson.M{"Device_ID": j}, container)
				mappingTemporaryData = airbox.MappingTemporaryData{}
				container = airbox.CroncheckStatus{}

				fmt.Println(container)

			}
			Mongo.Close()
		}
		time.Sleep(30 * time.Second)
	}
}
