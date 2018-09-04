package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "dforcepro.com/api"
	"github.com/gorilla/mux"
)

func GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/space/state/id/{id}/name/{state}", Next: sca.deviceState, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/space/state/id/{id}/name/{state}/arg/{argbv}", Next: sca.deviceACState, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/space/status/id/{id}", Next: sca.deviceStatusRes, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/space/cam/posX/{pos1}/posY/{pos2}", Next: sca.camTurn, Method: "GET", Auth: false},
	}
}

func (sca SmartSpAPI) deviceState(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	deviceID := vars["id"]
	state := vars["state"]

	if deviceID == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if state == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// queryMap := util.GetQueryValue(req, []string{"start", "end"})
	// startTime := (*queryMap)["start"].(string)

	// endTime := (*queryMap)["end"].(string)

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=turn" + state

	//var result []interface{}
	//var err error

	response, err := http.Get(url)
	fmt.Println(response)
	//getMongo().DB(DBName).C(Coll).Find(bson.M{"GWID": stationID}).Limit(100).All(&result)
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
		// contents, err := ioutil.ReadAll(response.Body)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// var record powerCons

		// json.NewDecoder(response.Body).Decode(&record)
		// json.Marshal(record)
		// json.NewEncoder(w).Encode(&record)

	}

}
func (sca SmartSpAPI) camTurn(w http.ResponseWriter, req *http.Request) {
	log.Println("camTurn")
	vars := mux.Vars(req)
	pos1 := vars["pos1"]
	pos2 := vars["pos2"]

	if pos1 == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if pos2 == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := "http://admin:ntust27333141@140.118.19.197:7289/cgi/ptdc.cgi?command=set_relative_pos&posX=" + pos1 + "&posY=" + pos2

	response, err := http.Get(url)
	// fmt.Println(response)

	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
	}

}

func (sca SmartSpAPI) deviceACState(w http.ResponseWriter, req *http.Request) {
	log.Println("deviceACState")
	vars := mux.Vars(req)
	deviceID := vars["id"]
	state := vars["state"]
	argbv := vars["argbv"]

	if deviceID == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if state == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if argbv == "" || false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// queryMap := util.GetQueryValue(req, []string{"start", "end"})
	// startTime := (*queryMap)["start"].(string)

	// endTime := (*queryMap)["end"].(string)

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=" + state + "&arg1=" + argbv

	//var result []interface{}
	//var err error

	response, err := http.Get(url)
	// fmt.Println(response)
	//getMongo().DB(DBName).C(Coll).Find(bson.M{"GWID": stationID}).Limit(100).All(&result)
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
	}

}

func (sca SmartSpAPI) deviceStatusRes(w http.ResponseWriter, req *http.Request) {
	log.Println("deviceStatusRes")
	vars := mux.Vars(req)
	deviceID := vars["id"]

	if deviceID == "" || false {
		w.WriteHeader(http.StatusBadRequest)
	}

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/devices/" + deviceID

	if deviceID == "130" {

		response, err := http.Get(url)
		// fmt.Println(response)
		//getMongo().DB(DBName).C(Coll).Find(bson.M{"GWID": stationID}).Limit(100).All(&result)
		if err != nil {

			_di.Log.Err(err.Error())
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			defer response.Body.Close()
			// contents, err := ioutil.ReadAll(response.Body)
			// if err != nil {
			// 	fmt.Println(err)
			// }

			var tmprecord airCond

			json.NewDecoder(response.Body).Decode(&tmprecord)
			tmpVal := checkBool(tmprecord.Properties.UIACValue)
			// node = append(node, tmpVal)
			//w.Write([]byte(tmpVal))
			json.NewEncoder(w).Encode(checkState{State: tmpVal})
			return
		}

	}
	//var result []interface{}
	//var err error

	response, err := http.Get(url)
	// fmt.Println(response)
	//getMongo().DB(DBName).C(Coll).Find(bson.M{"GWID": stationID}).Limit(100).All(&result)
	if err != nil {

		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		// contents, err := ioutil.ReadAll(response.Body)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		var record deviceStatusRes

		json.NewDecoder(response.Body).Decode(&record)
		// json.Marshal(record)
		recordVal, _ := strconv.ParseBool(record.Properties.Value)
		json.NewEncoder(w).Encode(checkState{State: recordVal})
		fmt.Println(req.RemoteAddr)
		return
	}

}
