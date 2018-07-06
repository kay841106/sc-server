package dispenser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sc.dforcepro.com/meter"

	"dforcepro.com/api"
	"gopkg.in/mgo.v2/bson"
)

type Dispenser bool

func (sca Dispenser) Enable() bool {
	return bool(sca)
}

func (sca Dispenser) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/dispenser/{start}/{end}", Next: sca.getRawData, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/Short/{start}/{end}", Next: sca.getShortRawdata, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/address/{start}/{end}/{address}", Next: sca.getRawData_addr, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/Status/All/{start}/{end}", Next: sca.getAllstatus, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/device_id/{start}/{end}/{device_id}", Next: sca.getRawdata_deviceid, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/nina", Next: sca.nina, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/hezong/board/{device_id}", Next: sca.board, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/hezong/board/{device_id}/{start}/{end}", Next: sca.boardquery, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/report/power/{address}/{start}/{end}", Next: sca.reportPower, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/realtime/power/{address}", Next: sca.realPower, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/dispenser/ninaquery/{device}/{start}/{end}", Next: sca.ninaQuery, Method: "GET", Auth: false},
	}
}

func xxxpipeBuildAllLatest() []bson.M {
	pipeline := []bson.M{}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name": "$_id.Building_Name",
			"Gateway_ID":    "$_id.Gateway_ID",

			"lastReportTime": 1,
			"avg_Usage":      bson.M{"$divide": []interface{}{"$avg_Usage", 1000}},
			"avg_Demand":     bson.M{"$divide": []interface{}{"$avg_Demand", 1000}},
			"max_Demand":     bson.M{"$divide": []interface{}{"$max_Demand", 1000}},
			"min_Demand":     bson.M{"$divide": []interface{}{"$min_Demand", 1000}},
			"max_Usage":      bson.M{"$divide": []interface{}{"$max_Usage", 1000}},
			"min_Usage":      bson.M{"$divide": []interface{}{"$min_Usage", 1000}},

			"max_PF": 1,
			"min_PF": 1,
			"avg_PF": bson.M{"$abs": "$avg_PF"},
			"CC":     bson.M{"$multiply": []interface{}{"$CC", 1000}},
		},
	})

	return pipeline
}

//Get the device by the timestamp
//
//							/dispenser/{start}/{end}
//							           long     long
func (sca Dispenser) getRawData(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []dispenserData{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"temp":           1,
		"watts":          1,
		"current":        1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}}).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println("1st API ok")
}

//Get the short data struct to device by the timestamp
//
//							/dispenser/Short/{start}/{end}
//							          		 long     long
func (sca Dispenser) getShortRawdata(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []leostruct{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}}).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println("2nd API ok")
}

//Get the device by the address(Sort by depending)
//
//							/dispenser/{start}/{end}/{address}
//										long    long  string
func (sca Dispenser) getRawData_addr(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]
	addr := vars["address"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []leostruct{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}, "address": addr}).Sort("-lastupdated").Limit(1).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//mt.Println("3rd API OK")
}

//Get now the newest status
//
//							/dispenser/Status/All/{start}/{end}
//													long  long
func (sca Dispenser) getAllstatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []AllStatus{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"temp":           1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}, "status": 1}).Sort("-lastupdated").Limit(60).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println("4rd API ok")
}

//Get the device by the device_id(Sort by depending)
//
//							/dispenser/{start}/{end}/{device_id}
//										long    long  string
func (sca Dispenser) getRawdata_deviceid(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]
	d_id := vars["device_id"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []leostruct{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}, "devicenickname": d_id}).Sort("-lastupdated").Limit(1).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println("5rd API OK")
}

//nina api
func (sca Dispenser) nina(w http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// tmpstart := vars["start"]
	// tmpend := vars["end"]
	// d_id := vars["device_id"]

	// start, _ := strconv.ParseInt(tmpstart, 10, 64)
	// end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []Nina{}
	// selectData := bson.M{
	// 	"UploadTime":       1,
	// 	"Watts":            1,
	// 	"DeviceMacAddress": 1,
	// 	"DeviceNickname":   1,
	// 	"Class":            1,
	// 	"_id":              0,
	// }
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C("Nina").Find(bson.M{}).Sort("-UploadTime").Limit(20).All(&dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	err := json.NewEncoder(w).Encode(dispenserDatastruct)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("7th API OK")
}

func (sca Dispenser) board(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// tmpstart := vars["start"]
	// tmpend := vars["end"]
	deviceid := vars["device_id"]

	// start, _ := strconv.ParseInt(tmpstart, 10, 64)
	// end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []dispenserboard{}
	selectData := bson.M{
		"UploadTime":  1,
		"device":      1,
		"hotTemp":     1,
		"warmTemp":    1,
		"coldTemp":    1,
		"tds":         1,
		"heating":     1,
		"cooling":     1,
		"savingpower": 1,
		"sterilizing": 1,
		"inputwater":  1,
		"waterlevel":  1,
		"hotoutput":   1,
		"warmoutput":  1,
		"coldeoutput": 1,
		"_id":         0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C("Leo").Find(bson.M{"device": deviceid}).Sort("-UploadTime").Limit(1).Select(selectData).All(&dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	err := json.NewEncoder(w).Encode(dispenserDatastruct)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("8th API OK")
}

func (sca Dispenser) boardquery(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]
	deviceid := vars["device_id"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []dispenserboard{}
	selectData := bson.M{
		"UploadTime":  1,
		"device":      1,
		"hotTemp":     1,
		"warmTemp":    1,
		"coldTemp":    1,
		"tds":         1,
		"heating":     1,
		"cooling":     1,
		"savingpower": 1,
		"sterilizing": 1,
		"inputwater":  1,
		"waterlevel":  1,
		"hotoutput":   1,
		"warmoutput":  1,
		"coldeoutput": 1,
		"_id":         0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C("Leo").Find(bson.M{"Timestamp": bson.M{"$gte": start, "$lte": end}, "device": deviceid}).Sort("UploadTime").Select(selectData).All(&dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	err := json.NewEncoder(w).Encode(dispenserDatastruct)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("9th API OK")
}

func (sca Dispenser) reportPower(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]
	address := vars["address"]
	//fmt.Println(tmpstart, tmpend, address)
	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []leostruct{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"lastupdated": bson.M{"$gte": start, "$lte": end}, "address": address}).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	//fmt.Println("10th API OK")
}

func (sca Dispenser) realPower(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	address := vars["address"]
	//fmt.Println(tmpstart, tmpend, address)
	dispenserDatastruct := []leostruct{}
	selectData := bson.M{
		"address":        1,
		"status":         1,
		"watts":          1,
		"devicenickname": 1,
		"lastupdated":    1,
		"_id":            0,
	}
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C(dispenserRawdataCollection).Find(bson.M{"address": address}).Sort("-lastupdated").Limit(1).Select(selectData).All(&dispenserDatastruct)
	json.NewEncoder(w).Encode(dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	//fmt.Println("11th API OK")
}

func (sca Dispenser) ninaQuery(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tmpstart := vars["start"]
	tmpend := vars["end"]
	deviceid := vars["device"]

	start, _ := strconv.ParseInt(tmpstart, 10, 64)
	end, _ := strconv.ParseInt(tmpend, 10, 64)
	dispenserDatastruct := []Nina{}
	// selectData := bson.M{
	// 	"UploadTime":       1,
	// 	"Watts":            1,
	// 	"DeviceMacAddress": 1,
	// 	"DeviceNickname":   1,
	// 	"Class":            1,
	// 	"_id":              0,
	// }
	mongo := meter.GetMongo()
	mongo.DB(DispenserDB).C("Nina").Find(bson.M{"Timestamp": bson.M{"$gte": start, "$lte": end}, "DeviceNickname": deviceid}).All(&dispenserDatastruct)
	//fmt.Println(dispenserDatastruct)
	err := json.NewEncoder(w).Encode(dispenserDatastruct)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(tmpstart, tmpend, start, end, deviceid)
	//fmt.Println(dispenserDatastruct)
	//fmt.Println("12th API OK")
}
