package cron

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"dforcepro.com/cron"
	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/meter"
)

type Dispenser bool

//Enable for Dispenser
func (mycron Dispenser) Enable() bool {
	return bool(mycron)
}
func (mycron Dispenser) GetJobs() []cron.JobSpec {
	return []cron.JobSpec{
		cron.JobSpec{
			Spec: "*/30 * * * * *",
			Job:  mycron.dispenserGET,
		},
		// {
		// Spec: "*/30 * * * * *",
		// Job:  mycron.deviceHourDispenser,
		// },
	}
}

//Pipeline
func pipeDeviceHourWholeDispenser(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"address": devID,
			}}}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"address": "$address",

				"Hour": bson.M{"$hour": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},

				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				"Month": bson.M{"$month": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
				"Year":  bson.M{"$year": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
				"day":   bson.M{"$dayOfYear": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
			},
			"devicenickname": bson.M{"$last": "$devicenickname"},
			"temp":           bson.M{"$avg": "$temp"},
			"demandcount":    bson.M{"$avg": "$demandcount"},
			"watts":          bson.M{"$avg": "$watts"},
			// "wattsaverage":   bson.M{"$sum": bson.M{"$avg": "$watts"}},
			"current":     bson.M{"$avg": "$current"},
			"lastupdated": bson.M{"$last": "$lastupdated"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"address": "$_id.address",
			// "isoTime": bson.M{date"new Date($lastupdated*1000)",
			"devicenickname": 1,
			"temp":           1,
			"watts":          1,
			"current":        1,
			"lastupdated":    1,
			"demandcount":    1,
			"Hour":           "$_id.Hour",
			"timeStamp":      bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000000}}}},
			// "watts_average":
			// "timeStamp": time.Date("$_id.Year", "$_id.Month", 1, 0, 0, 0, 0, time.Local),

			// "Device_Type":      1,
			// "Building_Details": 1,

			// "CC": "$CC",

			// "avg_Usage":    1,
			// "avg_Demand":   1,
			// "max_Demand":   1,
			// "min_Demand":   1,
			// "max_Usage":    1,
			// "min_Usage":    1,
			// "avg_PF":       1,
			// "max_PF":       1,
			// "min_PF":       1,
			// "weather_Temp": 1,
		},
	})

	// pipeline = append(pipeline, bson.M{
	// 	"$sort": bson.M{
	// 		"lastReportTime": 1},
	// })

	return pipeline
}

func pipeDeviceHourDispenser(start int64, devID string) []bson.M {

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"lastupdated": bson.M{
					"$gt": start,
				}, "address": devID,
			},
		}}
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"Hour": bson.M{"$hour": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
			// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
			// "Month": bson.M{"$month": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
			"Year":    bson.M{"$year": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000000}}}}},
			"day":     bson.M{"$dayOfMonth": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000000}}}}},
			"address": "$address",
		},
	})

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"address": "$address",

				// "Hour": bson.M{"$hour": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000000}}}}},
				"Hour": "$Hour",
				// "Day":        bson.M{"$dayOfMonth": "$lastReportTime"},
				// "Month": bson.M{"$month": bson.M{"$add": []interface{}{time.Unix(0, 0), bson.M{"$multiply": []interface{}{"$lastupdated", 1000}}}}},
				"Year": "$Year",
				"day":  "$day",
			},
			"devicenickname": bson.M{"$last": "$devicenickname"},
			"temp":           bson.M{"$avg": "$temp"},
			"demandcount":    bson.M{"$avg": "$demandcount"},
			"watts":          bson.M{"$avg": "$watts"},
			// "wattsaverage":   bson.M{"$sum": bson.M{"$avg": "$watts"}},
			"current":     bson.M{"$avg": "$current"},
			"lastupdated": bson.M{"$last": "$lastupdated"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"address": "$_id.address",
			// "isoTime": bson.M{date"new Date($lastupdated*1000)",
			"devicenickname": 1,
			"temp":           1,
			"watts":          1,
			"current":        1,
			"lastupdated":    1,
			"demandcount":    1,
			"Hour":           "$_id.Hour",
			// "watts_average":

			// "Device_Type":      1,
			// "Building_Details": 1,

			// "CC": "$CC",

			// "avg_Usage":    1,
			// "avg_Demand":   1,
			// "max_Demand":   1,
			// "min_Demand":   1,
			// "max_Usage":    1,
			// "min_Usage":    1,
			// "avg_PF":       1,
			// "max_PF":       1,
			// "min_PF":       1,
			// "weather_Temp": 1,
		},
	})

	// pipeline = append(pipeline, bson.M{
	// 	"$sort": bson.M{
	// 		"lastReportTime": 1},
	// })

	return pipeline
}

//Function

//SetField function
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

func (s *CDispenserPOST) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func modifyString(s interface{}) (string, error) {
	erro := errors.New("Cannot convert type")
	switch i := s.(type) {
	case float64:
		j := float64(i)
		return strconv.FormatFloat(j, 'f', -1, 64), nil
	case float32:
		j := float64(i)
		return strconv.FormatFloat(j, 'f', -1, 32), nil
	case int64:
		j := int64(i)
		return strconv.Itoa(int(j)), nil
	case int32:
		j := int64(i)
		return strconv.Itoa(int(j)), nil
	case int:
		return strconv.Itoa(i), nil
	case uint64:
		j := uint64(i)
		return strconv.Itoa(int(j)), nil
	case uint32:
		j := uint64(i)
		return strconv.Itoa(int(j)), nil
	case uint:
		return strconv.Itoa(int(i)), nil
	case string:
		return i, nil
	case nil:
		return "null", nil
	default:
		return "", erro

	}

}

func modifyInt(s interface{}) (int, error) {
	erro := errors.New("Cannot convert type")
	switch i := s.(type) {
	case float64:
		j := float64(i)
		return int(j), nil
		// return strconv.FormatFloat(j, 'f', -1, 64), nil
	case float32:
		j := float64(i)
		return int(j), nil
	case int64:
		j := int64(i)
		return int(j), nil
	case int32:
		j := int64(i)
		return int(j), nil
	case int:
		return i, nil
	case uint64:
		j := uint64(i)
		return int(j), nil
	case uint32:
		j := uint64(i)
		return int(j), nil
	case uint:
		j := uint(i)
		return int(j), nil
	case bool:
		if i != false {
			return -1, nil
		}
		return 1, nil

	case string:
		j := string(i)

		k, err := strconv.Atoi(j)
		if err != nil {
			return 0, err
		}
		return k, nil

		// return i, nil
	case nil:
		return 0, nil
	default:

		return 0, erro

	}

}

func modifyFloat64(s interface{}) (float64, error) {
	erro := errors.New("Cannot convert type")
	switch i := s.(type) {
	case float64:
		return i, nil
		// return strconv.FormatFloat(j, 'f', -1, 64), nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		j := string(i)
		k, err := strconv.Atoi(j)
		if err != nil {
			return 0, err
		}
		return float64(k), nil

		// return i, nil
	case nil:
		return 0, nil
	default:

		return 0, erro

	}
}

func modifyInt64(s interface{}) (int64, error) {
	erro := errors.New("Cannot convert type")
	switch i := s.(type) {
	case float64:
		j := float64(i)
		return int64(j), nil
		// return strconv.FormatFloat(j, 'f', -1, 64), nil
	case float32:
		j := float64(i)
		return int64(j), nil
	case int64:
		return i, nil
	case int32:
		return int64(i), nil
	case int:
		return int64(i), nil
	case uint64:
		j := uint64(i)
		return int64(j), nil
	case uint32:
		j := uint64(i)
		return int64(j), nil
	case uint:
		j := uint(i)
		return int64(j), nil
	case string:
		j := string(i)
		if k, err := strconv.Atoi(j); err != nil {
			return int64(0), err
		} else {
			return int64(k), nil
		}
		// return i, nil
	case nil:
		return int64(0), nil
	default:

		return int64(0), erro

	}

}

func (mycron Dispenser) dispenserGET() {
	checkDBStatus()

	log.Println("RawDataDispenser STARTING")
	response, err := http.Get(RawDataDispenserAPIUrl)

	mongo := meter.GetMongo()
	q := mongo.DB(DBDispenser)

	if err != nil {
		log.Println(err)
		response.Body.Close()
		// //os.Exit(1)
	} else {

		// defer response.Body.Close()
		// fmt.Println(response.Body)
		container := CDispenserPOST{}
		tmpcontainer := tempResponseDispenser{}
		content, err := ioutil.ReadAll(response.Body)
		// fmt.Println(string(content))

		if err != nil {

			log.Println(err)
			// //os.Exit(1)
		}

		erro := json.Unmarshal(content, &tmpcontainer)
		if erro != nil {
			log.Println(erro)
			// //os.Exit(1)
		}

		err = meter.GetMongo().Ping()
		if err != nil {
			log.Println(err)
			// //os.Exit(1)
		}

		for _, each := range tmpcontainer.Data {

			starttime, err := modifyInt64(each["starttime"])
			if err != nil {
				log.Println(err)
				// //os.Exit(1)
			}
			lastupdated, err := modifyInt64(each["lastupdated"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			stoptime, err := modifyInt64(each["stoptime"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			voltages, err := modifyString(each["voltages"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			address, err := modifyString(each["address"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			alertCurrent, err := modifyString(each["alertCurrent"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			battery, err := modifyInt(each["battery"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			calibration, err := modifyString(each["calibration"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			current, err := modifyInt(each["current"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			currentac, err := modifyString(each["current_ac"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			currentbyte, err := modifyInt(each["current_byte"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			currenthealthValue, err := modifyInt(each["current_healthValue"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			demandCount, err := modifyInt(each["demandCount"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			demandWatts, err := modifyInt(each["demandWatts"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			devicename, err := modifyString(each["devicename"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			deviceNickname, err := modifyString(each["DeviceNickname"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			filterPer, err := modifyInt(each["filterPer"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			hz, err := modifyInt(each["hz"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			lastSendAlert, err := modifyInt(each["lastSendAlert"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			maxCurrentraw, err := modifyInt(each["maxCurrent_raw"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			number, err := modifyInt(each["number"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			status, err := modifyInt(each["status"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			stopCurrent, err := modifyInt(each["stopCurrent"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			temp, err := modifyInt(each["temp"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			typess, err := modifyString(each["type"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			username, err := modifyString(each["username"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			variantCurrent, err := modifyInt(each["variantCurrent"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			samplingLeng, err := modifyInt(each["samplingLeng"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}
			watts, err := modifyFloat64(each["watts"])
			if err != nil {
				log.Println(err)
				//os.Exit(1)
			}

			container.Address = address
			container.Starttime = starttime
			container.Stoptime = stoptime
			container.Voltages = voltages
			container.Lastupdated = lastupdated
			container.AlertCurrent = alertCurrent
			container.Battery = battery
			container.Calibration = calibration
			container.Current = current
			container.CurrentAc = currentac
			container.CurrentByte = currentbyte
			container.CurrentHealthValue = currenthealthValue
			container.DemandCount = demandCount
			container.DemandWatts = demandWatts
			container.Devicename = devicename
			container.DeviceNickname = deviceNickname
			container.FilterPer = filterPer
			container.Hz = hz
			container.LastSendAlert = lastSendAlert
			container.MaxCurrentRaw = maxCurrentraw
			container.Number = number
			container.Status = status
			container.StopCurrent = stopCurrent
			container.Temp = temp
			container.Type = typess
			container.Username = username
			container.VariantCurrent = variantCurrent
			container.SamplingLeng = samplingLeng
			container.Watts = watts

			fmt.Println(container)

			q.C(RawDataDispenser).Upsert(bson.M{"address": container.Address, "lastupdated": container.Lastupdated, "starttime": container.Starttime, "stoptime": container.Stoptime}, container)

		}
		log.Println("RawDataDispenser FINISH")
		q.Close()
		response.Body.Close()

	}

	//
}

func (sca Dispenser) deviceHourDispenser() {
	mongo := meter.GetMongo()
	containerArray := []CDispenserUnixTime{}

	var tmpDevID []interface{}
	tmpData := CDispenserGET{}
	container := CDispenserGET{}
	err := mongo.DB(DBDispenser).C(RawDataDispenser).Find(nil).Distinct("address", &tmpDevID)
	fmt.Println(tmpDevID)
	if err != nil {
		log.Println(err)
	}
	for _, j := range tmpDevID {
		mongo.DB(DBDispenser).C(DispenserHourCollection).Find(bson.M{"address": j.(string)}).Sort("-lastupdated").One(&tmpData)
		fmt.Println(tmpData)
		if (CDispenserGET{}) != tmpData {
			datetimePlus1Hour := tmpData.TimeStamp.Add(time.Hour * 1)
			timestampUnix := datetimePlus1Hour.Unix()
			thepipe := pipeDeviceHourDispenser(timestampUnix, tmpData.Address)
			mongo.DB(DBDispenser).C(RawDataDispenser).Pipe(thepipe).All(&containerArray)

			for _, k := range containerArray {
				// container.TimeStamp = time.Unix(k.Lastupdated, 0)
				// container.Address = k.Address
				// container.Current = k.Current
				// container.DemandCount = k.DemandCount
				// container.Status = k.Status
				// container.Temp = k.Temp
				// container.Watts = k.Watts
				// container.WattsAverage = k.WattsAverage

				container.TimeStamp = timeConvertToHour(time.Unix(0, k.Lastupdated*1000*1000))
				container.Address = k.Address
				container.Current = k.Current
				container.DemandCount = k.DemandCount
				container.DeviceName = k.DeviceName
				// container.Status = k.Status
				container.Temp = k.Temp
				container.Watts = k.Watts
				fmt.Println(container)

				// mongo.DB(DBDispenser).C(DispenserHourCollection).Upsert(bson.M{"address": container.Address, "timestamp": container.TimeStamp}, container)
			}
		} else {

			thepipe := pipeDeviceHourWholeDispenser(j.(string))
			fmt.Println("WHOLE", j.(string))
			mongo.DB(DBDispenser).C(RawDataDispenser).Pipe(thepipe).AllowDiskUse().All(&containerArray)
			// mongo.C(DBDispenser).Pipe(thepipe).AllowDiskUse().Explain(&explain)
			fmt.Println(containerArray)
			for _, k := range containerArray {
				// container.TimeStamp = time.Unix(0, k.Lastupdated*1000*1000)
				// container.Address = k.Address
				// container.Current = strconv.Itoa(k.Current)
				// container.DemandCount = strconv.Itoa(k.DemandCount)
				// container.Status = strconv.Itoa(k.Status)
				// container.Temp = strconv.Itoa(k.Temp)
				// container.Watts = strconv.FormatFloat(k.Watts, 'f', 0, 64)

				container.TimeStamp = timeConvertToHour(time.Unix(0, k.Lastupdated*1000*1000))
				container.Address = k.Address
				container.Current = k.Current
				container.DemandCount = k.DemandCount
				container.DeviceName = k.DeviceName
				// container.Status = k.Status
				container.Temp = k.Temp
				container.Watts = k.Watts
				fmt.Println(k, container.TimeStamp, k.Hour)
				// container.WattsAverage = strconv.Itoa(k.WattsAverage)

				// tempTime := SetTimeStampForHour(each.LastReportTime)
				// each.LastReportTime = tempTime
				// fmt.Println(each.DevID, each.LastReportTime)
				mongo.DB(DBDispenser).C(DispenserHourCollection).Upsert(bson.M{"address": container.Address, "timestamp": container.TimeStamp}, container)
			}
		}
		container = CDispenserGET{}
		tmpData = CDispenserGET{}
		containerArray = []CDispenserUnixTime{}
	}

	log.Println(_EOF)
}

// mongo.DB(DBDispenser).C(DispenserHourCollection).Find()
