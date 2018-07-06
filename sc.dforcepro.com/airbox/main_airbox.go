package airbox

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"dforcepro.com/api"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/meter"
)

//SCAirbox API
type SCAirbox bool

//Enable SCAirbox bool
func (sca SCAirbox) Enable() bool {
	return bool(sca)
}

//GetAPIs for SCAirbox
func (sca SCAirbox) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/space/stream", Next: sca.SpaceAirboxStream, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/space/oneday/{devid}", Next: sca.SpaceAirboxOneDay, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/airbox/mapping", Next: sca.AirboxMapping, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/airbox/stream", Next: sca.AirboxStream, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/airbox/hour/{devid}/{start}/{end}", Next: sca.GetTotalHourAirbox, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/airbox/day/{devid}/{start}/{end}", Next: sca.GetTotalDayAirbox, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/airbox/month/{devid}/{start}/{end}", Next: sca.GetTotalMonthAirbox, Method: "GET", Auth: false},
	}
}

func (sca SCAirbox) AirboxMapping(w http.ResponseWriter, req *http.Request) {
	zeros := 0
	aspect := bson.M{}
	container := MappingData{}

	Mongo := meter.GetMongo()

	err := Mongo.DB(DBName).C(AirboxMapping).Find(nil).All(&container.Rows)

	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.DevID.Description = `DevID`
	container.Datashape.FieldDefinitions.DevID.Name = `DevID`
	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros

	container.Datashape.FieldDefinitions.Location.Aspects = aspect
	container.Datashape.FieldDefinitions.Location.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Location.Description = `Location`
	container.Datashape.FieldDefinitions.Location.Name = `Location`
	container.Datashape.FieldDefinitions.Location.Ordinal = zeros

	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(container)
	}

}

//SpaceAirboxStream for API
func (sca SCAirbox) SpaceAirboxStream(w http.ResponseWriter, req *http.Request) {
	zeros := 0
	aspect := bson.M{}
	container := ScheckStatus{}
	ConstspaceAirboxdevID := "781463DA0149AD7C"
	Mongo := meter.GetMongo()

	err := Mongo.DB(DBName).C(AirboxStream).Find(bson.M{"Device_ID": ConstspaceAirboxdevID}).All(&container.Rows)

	container.Datashape.FieldDefinitions.CO.Aspects = aspect
	container.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO.Description = `CO`
	container.Datashape.FieldDefinitions.CO.Name = `CO`
	container.Datashape.FieldDefinitions.CO.Ordinal = zeros

	container.Datashape.FieldDefinitions.CO2.Aspects = aspect
	container.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO2.Description = `CO2`
	container.Datashape.FieldDefinitions.CO2.Name = `CO2`
	container.Datashape.FieldDefinitions.CO2.Ordinal = zeros

	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros

	container.Datashape.FieldDefinitions.Humidity.Aspects = aspect
	container.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

	container.Datashape.FieldDefinitions.Noise.Aspects = aspect
	container.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Noise.Description = `Noise`
	container.Datashape.FieldDefinitions.Noise.Name = `Noise`
	container.Datashape.FieldDefinitions.Noise.Ordinal = zeros

	container.Datashape.FieldDefinitions.PM25.Aspects = aspect
	container.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Ordinal = zeros

	container.Datashape.FieldDefinitions.Temp.Aspects = aspect
	container.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Temp.Description = `Temp`
	container.Datashape.FieldDefinitions.Temp.Name = `Temp`
	container.Datashape.FieldDefinitions.Temp.Ordinal = zeros

	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(container)
	}

}

//AirboxStream for API
func (sca SCAirbox) AirboxStream(w http.ResponseWriter, req *http.Request) {

	zeros := 0
	aspect := bson.M{}
	container := ScheckStatus{}

	Mongo := meter.GetMongo()

	err := Mongo.DB(DBName).C(AirboxStream).Find(nil).All(&container.Rows)

	// err := Mongo.DB(DBName).C(AirboxStream).Find(bson.M{"Device_ID": ConstspaceAirboxdevID}).All(&container.Rows)

	container.Datashape.FieldDefinitions.CO.Aspects = aspect
	container.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO.Description = `CO`
	container.Datashape.FieldDefinitions.CO.Name = `CO`
	container.Datashape.FieldDefinitions.CO.Ordinal = zeros

	container.Datashape.FieldDefinitions.CO2.Aspects = aspect
	container.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO2.Description = `CO2`
	container.Datashape.FieldDefinitions.CO2.Name = `CO2`
	container.Datashape.FieldDefinitions.CO2.Ordinal = zeros

	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros

	container.Datashape.FieldDefinitions.Humidity.Aspects = aspect
	container.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

	container.Datashape.FieldDefinitions.Noise.Aspects = aspect
	container.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Noise.Description = `Noise`
	container.Datashape.FieldDefinitions.Noise.Name = `Noise`
	container.Datashape.FieldDefinitions.Noise.Ordinal = zeros

	container.Datashape.FieldDefinitions.PM25.Aspects = aspect
	container.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Ordinal = zeros

	container.Datashape.FieldDefinitions.Temp.Aspects = aspect
	container.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Temp.Description = `Temp`
	container.Datashape.FieldDefinitions.Temp.Name = `Temp`
	container.Datashape.FieldDefinitions.Temp.Ordinal = zeros

	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(container)
	}

}
func (sca SCAirbox) SpaceAirboxOneDay(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	devid := vars["devid"]

	aspect := bson.M{}
	zeros := 0

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	mongo := meter.GetMongo()

	container := ScheckStatus{}

	// pipeline:=[]bson.M{}
	// pipeline=append(pipeline, bson.M{
	// 	"$group":bson.M{
	// 		"$_id":bson.M{

	// 		}
	// 	}
	// })
	mongo.DB(DBName).C(AirboxRaw).Find(bson.M{"Device_ID": devid, "Upload_Time": bson.M{"$gte": start, "$lte": end}}).All(&container.Rows)

	container.Datashape.FieldDefinitions.CO.Aspects = aspect
	container.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO.Description = `CO`
	container.Datashape.FieldDefinitions.CO.Name = `CO`
	container.Datashape.FieldDefinitions.CO.Ordinal = zeros

	container.Datashape.FieldDefinitions.Location.Aspects = aspect
	container.Datashape.FieldDefinitions.Location.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Location.Description = `Location`
	container.Datashape.FieldDefinitions.Location.Name = `Location`
	container.Datashape.FieldDefinitions.Location.Ordinal = zeros

	container.Datashape.FieldDefinitions.CO2.Aspects = aspect
	container.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.CO2.Description = `CO2`
	container.Datashape.FieldDefinitions.CO2.Name = `CO2`
	container.Datashape.FieldDefinitions.CO2.Ordinal = zeros

	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros

	container.Datashape.FieldDefinitions.Humidity.Aspects = aspect
	container.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
	container.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

	container.Datashape.FieldDefinitions.Noise.Aspects = aspect
	container.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Noise.Description = `Noise`
	container.Datashape.FieldDefinitions.Noise.Name = `Noise`
	container.Datashape.FieldDefinitions.Noise.Ordinal = zeros

	container.Datashape.FieldDefinitions.PM25.Aspects = aspect
	container.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
	container.Datashape.FieldDefinitions.PM25.Ordinal = zeros

	container.Datashape.FieldDefinitions.Temp.Aspects = aspect
	container.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.Temp.Description = `Temp`
	container.Datashape.FieldDefinitions.Temp.Name = `Temp`
	container.Datashape.FieldDefinitions.Temp.Ordinal = zeros

	json.NewEncoder(w).Encode(container)

	fmt.Println(req.RemoteAddr)
	req.Body.Close()
}

func (sca SCAirbox) GetTotalHourAirbox(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// bid := vars["bid"]
	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	// var amp []interface{}
	containerz := ScheckStatus{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02T15", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02T15", end, time.Local)

		fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := meter.GetMongo()
		qu := mongo.DB(DBName).C(AirboxHour)
		qu.Find(bson.M{"Device_ID": devid, "Upload_Time": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)
		// fmt.Println(devid, tmpstart, container.Rows, err)

		// qu.Find(bson.M{"Device_ID": devid, "Upload_Time": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)

		containerz.Datashape.FieldDefinitions.CO.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO.Description = `CO`
		containerz.Datashape.FieldDefinitions.CO.Name = `CO`
		containerz.Datashape.FieldDefinitions.CO.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Location.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Location.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.Location.Description = `Location`
		containerz.Datashape.FieldDefinitions.Location.Name = `Location`
		containerz.Datashape.FieldDefinitions.Location.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.CO2.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO2.Description = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Name = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.DevID.Aspects = aspect
		containerz.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Humidity.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		containerz.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		containerz.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Noise.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Noise.Description = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Name = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.PM25.Aspects = aspect
		containerz.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Temp.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Temp.Description = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Name = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Ordinal = zeros

		json.NewEncoder(w).Encode(containerz)

		fmt.Println(req.RemoteAddr)
		req.Body.Close()
	}
}

// func (sca SCBuildAPI) postRawAEMDRA(w http.ResponseWriter, req *http.Request) {

// 	var devID string
// 	var GWID string

// 	_beforeEndPoint(w, req)

// 	pkg := rawAEMDRA{}

// 	_ = json.NewDecoder(req.Body).Decode(&pkg)

// 	mongo := getMongo()

// 	Resfrom := ResDoc{}

// 	devID = pkg.DevID

// 	v := reflect.ValueOf(pkg)
// 	value := make([]interface{}, v.NumField())
// 	for i := 0; i < v.NumField(); i++ {
// 		value[i] = v.Field(i).Interface()
// 	}
// 	pkg.GenObjectId
// }

func (sca SCAirbox) GetTotalDayAirbox(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	containerz := ScheckStatus{}
	// containerForDevManager := []DeviceManagerS{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02", end, time.Local)

		fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := meter.GetMongo()
		// mongo.DB(DBName).C(DeviceManager).Find(bson.M{"devID": devid}).One(&containerForDevManager)

		qu := mongo.DB(DBName).C(AirboxDay)
		qu.Find(bson.M{"Device_ID": devid, "Upload_Time": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)

		// qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)
		containerz.Datashape.FieldDefinitions.CO.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO.Description = `CO`
		containerz.Datashape.FieldDefinitions.CO.Name = `CO`
		containerz.Datashape.FieldDefinitions.CO.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Location.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Location.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.Location.Description = `Location`
		containerz.Datashape.FieldDefinitions.Location.Name = `Location`
		containerz.Datashape.FieldDefinitions.Location.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.CO2.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO2.Description = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Name = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.DevID.Aspects = aspect
		containerz.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Humidity.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		containerz.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		containerz.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Noise.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Noise.Description = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Name = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.PM25.Aspects = aspect
		containerz.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Temp.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Temp.Description = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Name = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Ordinal = zeros

		json.NewEncoder(w).Encode(&containerz)

		fmt.Println(req.RemoteAddr)
		req.Body.Close()

	}

}

func (sca SCAirbox) GetTotalMonthAirbox(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	containerz := ScheckStatus{}
	// containerForDevManager := DeviceManagerS{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01", end, time.Local)

		fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := meter.GetMongo()
		// mongo.DB(DBName).C(DeviceManager).Find(bson.M{"devID": devid}).All(&containerForDevManager)

		qu := mongo.DB(DBName).C(AirboxMonth)
		qu.Find(bson.M{"Device_ID": devid, "Upload_Time": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)

		// qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)
		containerz.Datashape.FieldDefinitions.CO.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO.Description = `CO`
		containerz.Datashape.FieldDefinitions.CO.Name = `CO`
		containerz.Datashape.FieldDefinitions.CO.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.CO2.Aspects = aspect
		containerz.Datashape.FieldDefinitions.CO2.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.CO2.Description = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Name = `CO2`
		containerz.Datashape.FieldDefinitions.CO2.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Location.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Location.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.Location.Description = `Location`
		containerz.Datashape.FieldDefinitions.Location.Name = `Location`
		containerz.Datashape.FieldDefinitions.Location.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.DevID.Aspects = aspect
		containerz.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
		containerz.Datashape.FieldDefinitions.DevID.Description = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
		containerz.Datashape.FieldDefinitions.DevID.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Humidity.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Humidity.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Humidity.Description = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Name = `Humidity`
		containerz.Datashape.FieldDefinitions.Humidity.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		containerz.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		containerz.Datashape.FieldDefinitions.LastReportTime.Description = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Name = `Upload_Time`
		containerz.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Noise.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Noise.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Noise.Description = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Name = `Noise`
		containerz.Datashape.FieldDefinitions.Noise.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.PM25.Aspects = aspect
		containerz.Datashape.FieldDefinitions.PM25.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.PM25.Description = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Name = `PM2_5`
		containerz.Datashape.FieldDefinitions.PM25.Ordinal = zeros

		containerz.Datashape.FieldDefinitions.Temp.Aspects = aspect
		containerz.Datashape.FieldDefinitions.Temp.BaseType = `NUMBER`
		containerz.Datashape.FieldDefinitions.Temp.Description = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Name = `Temp`
		containerz.Datashape.FieldDefinitions.Temp.Ordinal = zeros

		json.NewEncoder(w).Encode(&containerz)

		fmt.Println(req.RemoteAddr)
		req.Body.Close()
	}
}
