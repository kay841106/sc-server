package meter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"dforcepro.com/api"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/alert"
)

const (
	streamCollection    = "SC01_Stream"
	statusCollection    = "SC01_Status"
	streamAllCollection = "SC01_Stream_All"
	//DeviceManager Collection
	DeviceManager = "SC01_DeviceManager"

	collMonthTotal = "SC01_month_All"

	coll = "SC01_displayData_Calc_"

	collAll = "SC01_displayData_All"

	collAllCampus    = "SC01_displayData_Calc_main_substation"
	collHour         = "SC01_hour_All"
	collHourTotal    = "SC01_hour_"
	lookupHourOnTime = "SC01_hour_All"

	collDayTotal    = "SC01_day_All"
	collDay         = "SC01_day_All"
	lookupDayOnTime = "SC01_day_All"

	collMonth         = "SC01_month_All"
	lookupMonthOnTime = "SC01_month_All"

	_EmptyStr = ""

	//RSASecret for Validation
	RSASecret = "Y7WfGYtOHGBjMMigZ6QrcvveYuNDEgepBuBpYJr2lCB-UYdRRTFe5swVQW8iLh5a"
)

//SCBuildAPI Building
type SCBuildAPI bool

//Enable the API
func (sca SCBuildAPI) Enable() bool {
	return bool(sca)
}

//GetAPIs router
func (sca SCBuildAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{

		&api.APIHandler{Path: "/mappingmeter/{devid}", Next: sca.mappingMeter, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/Status/{b_id}", Next: sca.checkLastAllDev, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/agg/All/today", Next: sca.aggBuildingAllToday, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/agg/All/{start}/{end}", Next: sca.aggBuildingAllOnTime, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/explain/", Next: sca.queryExplain, Method: "GET", Auth: false},
		// &api.APIHandler{Path: "/agg/All/latest", Next: sca.aggBuildingAllLatest, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/Status/Dev/All", Next: sca.checkStatus, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/Stream", Next: sca.DataStream, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/Stream/All", Next: sca.DataStreamAll, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/agg/explain/{b_id}/{start}/{end}", Next: sca.aggBuildingAllExplain, Method: "GET", Auth: false},
		// &api.APIHandler{Path: "/agg/hour/{start}/{end}/{bid}", Next: sca.aggAllHourOnTime, Method: "GET", Auth: false},
		// &api.APIHandler{Path: "/agg/day/{start}/{end}/{bid}", Next: sca.aggAllDayOnTime, Method: "GET", Auth: false},
		// &api.APIHandler{Path: "/agg/month/{start}/{end}/{bid}", Next: sca.aggAllMonthOnTime, Method: "GET", Auth: false},

		&api.APIHandler{Path: "/agg/total/hour/{devid}/{start}/{end}", Next: sca.GetTotalHour, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/agg/total/day/{devid}/{start}/{end}", Next: sca.GetTotalDay, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/agg/total/month/{devid}/{start}/{end}", Next: sca.GetTotalMonth, Method: "GET", Auth: false},
		// &api.APIHandler{Path: "/agg/total/month/{bid}/{devid}/{start}/{end}", Next: sca.aggAllMonthOnTime, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/add/{key}", Next: sca.Validator, Method: "POST", Auth: false},
		&api.APIHandler{Path: "/genID/{user}", Next: sca.GenAuth, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/listDevices", Next: sca.ListDevices, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/removeDevice", Next: sca.RemoveDev, Method: "POST", Auth: false},
		&api.APIHandler{Path: "/hand", Next: sca.testaest, Method: "POST", Auth: false},
		// &api.APIHandler{Path: "/add/{GatewayID}/{DeviceBrand}/{DeviceID}/{DeviceDetails}/{DeviceName}/{DeviceInfo}/{DeviceType}/{Floor}/{BuildingName}/{BuildingDetails}", Next: sca.AddDev, Method: "POST", Auth: false},
	}

}

///TIME
func tdayStart() time.Time {
	t := time.Now().UTC()
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func tdayEnd() time.Time {
	t := time.Now().UTC()
	fmt.Println(t)
	year, month, day := t.Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, time.Local)

}

func tdayCust(H int, M int, S int) time.Time {
	t := time.Now()
	year, month, day := t.Date()
	return time.Date(year, month, day, H, M, S, 0, t.Location())
}

///PIPE
func pipeBuildAllOnTime(start string, end string) []bson.M {
	var pipeline []bson.M

	if start != "" && end != "" {
		tmpstart, e := time.Parse("2006-01-02", start)
		tmpend, er := time.Parse("2006-01-02", end)
		fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {

			_di.Log.Err(er.Error())
		}

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"lastReportTime": bson.M{
					"$gte": tmpstart,
					"$lte": tmpend,
				}},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name": "$Building_Name",
			"Gateway_ID":    "$Gateway_ID",

			"lastReportTime": "$lastReportTime",
			// KW
			"avg_Usage":  bson.M{"$divide": []interface{}{"$Pwr_Usage", 1000}},
			"avg_Demand": bson.M{"$divide": []interface{}{"$Pwr_Demand", 1000}},
			"PF":         "$PF",
			"CC":         bson.M{"$multiply": []interface{}{"$CC", 1000}},
			"weather":    "$weather_Temp",
		},
	})

	return pipeline
}

func pipeAllHourOnTime(start string, end string, bid string) []bson.M {
	var pipeline []bson.M

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02T15", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02T15", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{

				"$and": []bson.M{
					bson.M{
						"lastReportTime": bson.M{
							"$gte": tmpstart,
							"$lte": tmpend},
					},
					bson.M{"Building_Name": bid},
				},
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name":    "$Building_Name",
			"Building_Details": "$Building_Details",
			"Gateway_ID":       "$Gateway_ID",

			"lastReportTime": "$lastReportTime",
			// KW
			"max_Usage":    bson.M{"$divide": []interface{}{"$max_Usage", 1000}},
			"min_Usage":    bson.M{"$divide": []interface{}{"$min_Usage", 1000}},
			"Pwr_Usage":    bson.M{"$divide": []interface{}{"$Pwr_Usage", 1000}},
			"Pwr_Demand":   bson.M{"$divide": []interface{}{"$Pwr_Demand", 1000}},
			"max_Demand":   bson.M{"$divide": []interface{}{"$max_Demand", 1000}},
			"min_Demand":   bson.M{"$divide": []interface{}{"$min_Demand", 1000}},
			"max_PF":       "$max_PF",
			"min_PF":       "$min_PF",
			"PF":           "$PF",
			"PF_Limit":     "$PF_Limit",
			"CC":           bson.M{"$multiply": []interface{}{"$CC", 1000}},
			"weather_Temp": "$weather_Temp",
		},
	})

	return pipeline
}

func pipeAllDayOnTime(start string, end string, bid string) []bson.M {
	var pipeline []bson.M

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{

				"$and": []bson.M{
					bson.M{
						"lastReportTime": bson.M{
							"$gte": tmpstart,
							"$lte": tmpend},
					},
					bson.M{"Building_Name": bid},
				},
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name":    "$Building_Name",
			"Building_Details": "$Building_Details",
			"Gateway_ID":       "$Gateway_ID",

			"lastReportTime": "$lastReportTime",
			// KW
			"max_Usage":    bson.M{"$divide": []interface{}{"$max_Usage", 1000}},
			"min_Usage":    bson.M{"$divide": []interface{}{"$min_Usage", 1000}},
			"Pwr_Usage":    bson.M{"$divide": []interface{}{"$Pwr_Usage", 1000}},
			"Total_Usage":  bson.M{"$divide": []interface{}{bson.M{"$sum": "$Pwr_Usage"}, 1000}},
			"Pwr_Demand":   bson.M{"$divide": []interface{}{"$Pwr_Demand", 1000}},
			"max_Demand":   bson.M{"$divide": []interface{}{"$max_Demand", 1000}},
			"min_Demand":   bson.M{"$divide": []interface{}{"$min_Demand", 1000}},
			"max_PF":       "$max_PF",
			"min_PF":       "$min_PF",
			"PF":           "$PF",
			"PF_Limit":     "$PF_Limit",
			"CC":           bson.M{"$multiply": []interface{}{"$CC", 1000}},
			"weather_Temp": "$weather_Temp",
		},
	})

	return pipeline
}

func pipeAllMonthOnTime(start string, end string, bid string) []bson.M {
	var pipheline []bson.M

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		pipheline = append(pipheline, bson.M{
			"$match": bson.M{

				"$and": []bson.M{
					bson.M{
						"lastReportTime": bson.M{
							"$gte": tmpstart,
							"$lte": tmpend},
					},
					bson.M{"Building_Name": bid},
				},
			},
		})
	}

	return pipheline
}

func pipeBuildAllToday() []bson.M {

	start := tdayStart()
	end := tdayEnd()
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"lastReportTime": bson.M{
					"$gt": start,
					"$lt": end,
				}, "Device_ID": "3000509b52f3501"},
		}}

	// fmt.Println(start)
	// fmt.Println(end)

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Building_Name": "$Building_Name",
				"Hour":          bson.M{"$hour": "$lastReportTime"},
				"Gateway_ID":    "$Gateway_ID",
			},
			"lastReportTime": bson.M{"$last": "$lastReportTime"},
			"max_Demand":     bson.M{"$max": "$Pwr_Demand"},
			"min_Demand":     bson.M{"$min": "$Pwr_Demand"},
			"max_Usage":      bson.M{"$max": "$Pwr_Usage"},
			"min_Usage":      bson.M{"$min": "$Pwr_Usage"},
			"avg_PF":         bson.M{"$avg": bson.M{"$abs": "$PF"}},
			"max_PF":         bson.M{"$max": bson.M{"$abs": "$PF"}},
			"min_PF":         bson.M{"$min": bson.M{"$abs": "$PF"}},
			"avg_Usage":      bson.M{"$avg": "$Pwr_Usage"},
			"avg_Demand":     bson.M{"$avg": "$Pwr_Demand"},
			"total_Usage":    bson.M{"$sum": "$Pwr_Usage"},
			"CC":             bson.M{"$avg": "$CC"},
			"weather":        bson.M{"$avg": "$weather_Temp"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Building_Name": "$_id.Building_Name",
			"Gateway_ID":    "$_id.Gateway_ID",

			"CC":             bson.M{"$multiply": []interface{}{"$CC", 1000}},
			"lastReportTime": 1,
			"total_Usage":    bson.M{"$divide": []interface{}{"$total_Usage", 1000}},
			"avg_Usage":      bson.M{"$divide": []interface{}{"$avg_Usage", 1000}},
			"avg_Demand":     bson.M{"$divide": []interface{}{"$avg_Demand", 1000}},
			"max_Demand":     bson.M{"$divide": []interface{}{"$max_Demand", 1000}},
			"min_Demand":     bson.M{"$divide": []interface{}{"$min_Demand", 1000}},
			"max_Usage":      bson.M{"$divide": []interface{}{"$max_Usage", 1000}},
			"min_Usage":      bson.M{"$divide": []interface{}{"$min_Usage", 1000}},

			"max_PF":  1,
			"min_PF":  1,
			"weather": 1,
			"avg_PF":  bson.M{"$abs": "$avg_PF"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": 1},
	})

	return pipeline
}

func pipeBuildAllLatest() []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$sort": bson.M{
				"lastReportTime": 1}},
	}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Building_Name": "$Building_Name",
				"Gateway_ID":    "$Gateway_ID",
			},
			"lastReportTime": bson.M{"$last": "$lastReportTime"},
			"avg_Usage":      bson.M{"$last": "$Pwr_Usage"},
			"avg_Demand":     bson.M{"$last": "$Pwr_Demand"},
			"max_Demand":     bson.M{"$max": "$Pwr_Demand"},
			"min_Demand":     bson.M{"$min": "$Pwr_Demand"},
			"max_Usage":      bson.M{"$max": "$Pwr_Usage"},
			"min_Usage":      bson.M{"$min": "$Pwr_Usage"},
			"avg_PF":         bson.M{"$last": "$PF"},
			"max_PF":         bson.M{"$max": bson.M{"$abs": "$PF"}},
			"min_PF":         bson.M{"$min": bson.M{"$abs": "$PF"}},
			"CC":             bson.M{"$avg": "$CC"},
		},
	})
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

func pipeDevCheckStatus(devID string) []bson.M {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"Device_ID": devID,
			},
		}}

	pipeline = []bson.M{
		bson.M{
			"$sort": bson.M{
				"lastReportTime": 1}},
	}

	pipeline = append(pipeline, bson.M{

		"$group": bson.M{
			"_id": bson.M{
				"Device_ID":  "$Device_ID",
				"Gateway_ID": "$Gateway_ID",
			},
			"lastReportTime":   bson.M{"$last": "$lastReportTime"},
			"Building_Details": bson.M{"$last": "$Building_Details"},
			"Building_Name":    bson.M{"$last": "$Building_Name"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,

			"Device_ID":        "$_id.Device_ID",
			"Gateway_ID":       "$_id.Gateway_ID",
			"lastReportTime":   1,
			"Building_Details": 1,
			"Building_Name":    1,
		},
	})

	return pipeline
}

func pipeDevIDArray(devIDArray []string, sortBy string) []bson.M {
	// pipeline := []bson.M{
	// 	bson.M{"$match": bson.M{"Device_ID": bson.M{"$in": devIDArray}}},
	// }

	pipeline := []bson.M{
		bson.M{
			"$sort": bson.M{
				"lastReportTime": 1},
		}}

	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"Device_ID": "$Device_ID",
			},
			"lastReportTime": bson.M{"$last": "$lastReportTime"},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id":            0,
			"Device_ID":      "$_id.Device_ID",
			"lastReportTime": "$lastReportTime"},
	})

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"lastReportTime": -1},
	})
	return pipeline
}

///MASHAL

func (sca SCBuildAPI) mappingMeter(w http.ResponseWriter, req *http.Request) {
	log.Println("mappingMeter")
	vars := mux.Vars(req)
	devID := vars["devid"]
	_beforeEndPoint(w, req)

	tmpContainer := ListDevices{}

	var explain interface{}
	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()

	mongo.DB(DBName).C(DeviceManager).Find(bson.M{"devID": devID}).All(&tmpContainer.Rows)
	mongo.DB(DBName).C(DeviceManager).Find(nil).Explain(&explain)
	// fmt.Println(explain)

	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Name = `Time_Added`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Description = `Time_Added`

	tmpContainer.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	tmpContainer.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Name = `Device_Brand`
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Description = `Device_Brand`

	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Name = `Device_Details`
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Description = `Device_Details`

	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Name = `Device_Info`
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Description = `Device_Info`

	tmpContainer.Datashape.FieldDefinitions.DeviceName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Name = `Device_Name`
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Description = `Device_Name`

	tmpContainer.Datashape.FieldDefinitions.DeviceID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Description = `Device_ID`

	tmpContainer.Datashape.FieldDefinitions.DeviceType.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceType.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Name = `Device_Type`
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Description = `Device_Type`

	tmpContainer.Datashape.FieldDefinitions.Floor.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Floor.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.Floor.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Floor.Name = `Floor`
	tmpContainer.Datashape.FieldDefinitions.Floor.Description = `Floor`

	tmpContainer.Datashape.FieldDefinitions.Facility.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Facility.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.Facility.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Facility.Name = `Facility`
	tmpContainer.Datashape.FieldDefinitions.Facility.Description = `Facility`

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	defer req.Body.Close()
}
func (sca SCBuildAPI) checkLastAllDev(w http.ResponseWriter, req *http.Request) {
	log.Println("checkLastAllDev")
	_beforeEndPoint(w, req)

	vars := mux.Vars(req)
	buildingName := vars["b_id"]
	container := DisplayData{}

	var result []string
	var err error
	aspect := bson.M{}
	zeros := 0

	mongo := getMongo()
	err = mongo.DB(DBName).C("SC01_DeviceManager").Find(bson.M{"Building_Name": buildingName}).Distinct("devID", &result)

	thepipe := pipeDevIDArray(result, "lastReportTime")

	pipe := mongo.DB(DBName).C(coll + buildingName).Pipe(thepipe).AllowDiskUse().All(&container.Rows)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		log.Println(currTime.String() + coll + buildingName)
	}

	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	container.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	container.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros
	container.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	container.Datashape.FieldDefinitions.DevID.Description = `Device_ID`

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	mongo.Close()
	// fmt.Println(req.RemoteAddr)

}

func (sca SCBuildAPI) aggBuildingAllLatest(w http.ResponseWriter, req *http.Request) {
	log.Println("aggBuildingAllLatest")
	_beforeEndPoint(w, req)

	container := aggAllNow{}

	// var result []string
	var err error
	aspect := bson.M{}
	zeros := 0

	mongo := getMongo()
	thepipe := pipeBuildAllLatest()
	pipe := mongo.DB(DBName).C(collAll).Pipe(thepipe).AllowDiskUse().All(&container.Rows)

	container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	container.Datashape.FieldDefinitions.CC.Ordinal = zeros
	container.Datashape.FieldDefinitions.CC.BaseType = `STRING`
	container.Datashape.FieldDefinitions.CC.Aspects = aspect
	container.Datashape.FieldDefinitions.CC.Name = `CC`
	container.Datashape.FieldDefinitions.CC.Description = `CC`

	container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrUsage.Name = `PwrUsage`
	container.Datashape.FieldDefinitions.PwrUsage.Description = `PwrUsage`

	container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrDemand.Name = `PwrDemand`
	container.Datashape.FieldDefinitions.PwrDemand.Description = `PwrDemand`

	container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxDemand.Name = `MaxDemand`
	container.Datashape.FieldDefinitions.MaxDemand.Description = `MaxDemand`

	container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinDemand.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MinDemand.Name = `MinDemand`
	container.Datashape.FieldDefinitions.MinDemand.Description = `MinDemand`

	container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MinUsage.Name = `MinUsage`
	container.Datashape.FieldDefinitions.MinUsage.Description = `MinUsage`

	container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxUsage.Name = `MaxUsage`
	container.Datashape.FieldDefinitions.MaxUsage.Description = `MaxUsage`

	container.Datashape.FieldDefinitions.AvgPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.AvgPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.AvgPF.Aspects = aspect
	container.Datashape.FieldDefinitions.AvgPF.Name = `AvgPF`
	container.Datashape.FieldDefinitions.AvgPF.Description = `AvgPF`

	container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxPF.Name = `MaxPF`
	container.Datashape.FieldDefinitions.MaxPF.Description = `MaxPF`

	container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MinPF.Name = `MinPF`
	container.Datashape.FieldDefinitions.MinPF.Description = `MinPF`

	if pipe != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		log.Println(currTime.String() + collAll)
	}

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

//DataStream provides newest data for each device ID
func (sca SCBuildAPI) DataStream(w http.ResponseWriter, req *http.Request) {

	log.Println("DataStream")
	_beforeEndPoint(w, req)

	tmpContainer := DisplayData{}

	var explain interface{}
	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()

	mongo.DB(DBName).C(streamCollection).Find(nil).All(&tmpContainer.Rows)
	mongo.DB(DBName).C(streamCollection).Find(nil).Explain(&explain)
	// fmt.Println(explain)

	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	tmpContainer.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	tmpContainer.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Name = `PwrDemand`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Description = `PwrDemand`

	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Name = `PwrUsage`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Description = `PwrUsage`

	tmpContainer.Datashape.FieldDefinitions.Weather.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Weather.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.Weather.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Weather.Name = `Weather`
	tmpContainer.Datashape.FieldDefinitions.Weather.Description = `Weather`

	tmpContainer.Datashape.FieldDefinitions.PF.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PF.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PF.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PF.Name = `PF`
	tmpContainer.Datashape.FieldDefinitions.PF.Description = `PF`

	tmpContainer.Datashape.FieldDefinitions.DevID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DevID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.DevID.Description = `Device_ID`

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

//DataStreamAll provides newest data for each building
func (sca SCBuildAPI) DataStreamAll(w http.ResponseWriter, req *http.Request) {
	log.Println("DataStreamAll")
	_beforeEndPoint(w, req)

	tmpContainer := DisplayData{}

	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()

	mongo.DB(DBName).C(streamAllCollection).Find(aspect).All(&tmpContainer.Rows)
	// fmt.Println(tmpContainer)
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	tmpContainer.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	tmpContainer.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Name = `Pwr_Demand`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Description = `PwrDemand`

	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Name = `Pwr_Usage`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Description = `PwrUsage`

	tmpContainer.Datashape.FieldDefinitions.Weather.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Weather.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.Weather.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Weather.Name = `weather_Temp`
	tmpContainer.Datashape.FieldDefinitions.Weather.Description = `Weather`

	tmpContainer.Datashape.FieldDefinitions.PF.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PF.BaseType = `NUMBER`
	tmpContainer.Datashape.FieldDefinitions.PF.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PF.Name = `PF`
	tmpContainer.Datashape.FieldDefinitions.PF.Description = `PF`

	tmpContainer.Datashape.FieldDefinitions.DevID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DevID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DevID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DevID.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.DevID.Description = `Device_ID`

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

func (sca SCBuildAPI) checkStatus(w http.ResponseWriter, req *http.Request) {
	log.Println("checkStatus")
	_beforeEndPoint(w, req)

	tmpContainer := alert.AlertValueThingworx{}

	// var tmpDevID []string

	// zeros := 0
	// aspect := bson.M{}

	mongo := getMongo()
	var explain interface{}
	// var buildContainer []interface{}

	mongo.DB(DBName).C(statusCollection).Find(bson.M{}).All(&tmpContainer.Rows)
	mongo.DB(DBName).C(statusCollection).Find(nil).Explain(&explain)
	// fmt.Println(explain)

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()

}

func (sca SCBuildAPI) aggBuildingAllToday(w http.ResponseWriter, req *http.Request) {
	log.Println("aggBuildingAllToday")
	_beforeEndPoint(w, req)

	container := aggAllToday{}

	// var result []string
	var err error

	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()
	thepipe := pipeBuildAllToday()
	pipe := mongo.DB(DBName).C(collAllCampus).Pipe(thepipe).AllowDiskUse().All(&container.Rows)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		log.Println(currTime.String() + collAll)
	}

	container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	container.Datashape.FieldDefinitions.CC.Ordinal = zeros
	container.Datashape.FieldDefinitions.CC.BaseType = `STRING`
	container.Datashape.FieldDefinitions.CC.Aspects = aspect
	container.Datashape.FieldDefinitions.CC.Name = `CC`
	container.Datashape.FieldDefinitions.CC.Description = `CC`

	container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrUsage.Name = `PwrUsage`
	container.Datashape.FieldDefinitions.PwrUsage.Description = `PwrUsage`

	container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrDemand.Name = `PwrDemand`
	container.Datashape.FieldDefinitions.PwrDemand.Description = `PwrDemand`

	container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxDemand.Name = `MaxDemand`
	container.Datashape.FieldDefinitions.MaxDemand.Description = `MaxDemand`

	container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinDemand.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MinDemand.Name = `MinDemand`
	container.Datashape.FieldDefinitions.MinDemand.Description = `MinDemand`

	container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MinUsage.Name = `MinUsage`
	container.Datashape.FieldDefinitions.MinUsage.Description = `MinUsage`

	container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxUsage.Name = `MaxUsage`
	container.Datashape.FieldDefinitions.MaxUsage.Description = `MaxUsage`

	container.Datashape.FieldDefinitions.AvgPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.AvgPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.AvgPF.Aspects = aspect
	container.Datashape.FieldDefinitions.AvgPF.Name = `AvgPF`
	container.Datashape.FieldDefinitions.AvgPF.Description = `AvgPF`

	container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxPF.Name = `MaxPF`
	container.Datashape.FieldDefinitions.MaxPF.Description = `MaxPF`

	container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MinPF.Name = `MinPF`
	container.Datashape.FieldDefinitions.MinPF.Description = `MinPF`

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

func (sca SCBuildAPI) aggBuildingAllOnTime(w http.ResponseWriter, req *http.Request) {
	log.Println("aggBuildingAllOnTime")
	_beforeEndPoint(w, req)
	vars := mux.Vars(req)
	start := vars["start"]
	end := vars["end"]
	tmpContainer := AggAllOnTime{}

	// var result []string
	var err error

	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()
	thepipe := pipeBuildAllOnTime(start, end)
	pipe := mongo.DB(DBName).C(collAll).Pipe(thepipe).AllowDiskUse().All(&tmpContainer.Rows)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		log.Println(currTime.String() + collAll)
	}
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	tmpContainer.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	tmpContainer.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	tmpContainer.Datashape.FieldDefinitions.CC.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.CC.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.CC.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.CC.Name = `Building_Details`
	tmpContainer.Datashape.FieldDefinitions.CC.Description = `Building_Details`

	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.PwrUsage.Description = `Device_ID`

	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.PwrDemand.Description = `Device_ID`

	tmpContainer.Datashape.FieldDefinitions.AvgPF.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.AvgPF.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.AvgPF.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.AvgPF.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.AvgPF.Description = `Device_ID`

	tmpContainer.Datashape.FieldDefinitions.Weather.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Weather.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.Weather.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Weather.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.Weather.Description = `Device_ID`

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

func (sca SCBuildAPI) aggAllHourOnTime(w http.ResponseWriter, req *http.Request) {
	log.Println("aggAllHourOnTime")
	_beforeEndPoint(w, req)
	vars := mux.Vars(req)
	start := vars["start"]
	end := vars["end"]
	bid := vars["bid"]

	container := strAllHourOnTime{}

	// var result []string
	var err error

	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()
	thepipe := pipeAllHourOnTime(start, end, bid)
	pipe := mongo.DB(DBName).C(collHour).Pipe(thepipe).AllowDiskUse().All(&container.Rows)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		fmt.Println(currTime.String() + collAll)
	}

	container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
	container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

	container.Datashape.FieldDefinitions.CC.Ordinal = zeros
	container.Datashape.FieldDefinitions.CC.BaseType = `STRING`
	container.Datashape.FieldDefinitions.CC.Aspects = aspect
	container.Datashape.FieldDefinitions.CC.Name = `CC`
	container.Datashape.FieldDefinitions.CC.Description = `CC`

	container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrUsage.Name = `PwrUsage`
	container.Datashape.FieldDefinitions.PwrUsage.Description = `PwrUsage`

	container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	container.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	container.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	container.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	container.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	container.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.PwrDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.PwrDemand.Name = `PwrDemand`
	container.Datashape.FieldDefinitions.PwrDemand.Description = `PwrDemand`

	container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxDemand.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxDemand.Name = `MaxDemand`
	container.Datashape.FieldDefinitions.MaxDemand.Description = `MaxDemand`

	container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinDemand.BaseType = `DATETIME`
	container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
	container.Datashape.FieldDefinitions.MinDemand.Name = `MinDemand`
	container.Datashape.FieldDefinitions.MinDemand.Description = `MinDemand`

	container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MinUsage.Name = `MinUsage`
	container.Datashape.FieldDefinitions.MinUsage.Description = `MinUsage`

	container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxUsage.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxUsage.Name = `MaxUsage`
	container.Datashape.FieldDefinitions.MaxUsage.Description = `MaxUsage`

	container.Datashape.FieldDefinitions.PF.Ordinal = zeros
	container.Datashape.FieldDefinitions.PF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PF.Aspects = aspect
	container.Datashape.FieldDefinitions.PF.Name = `PF`
	container.Datashape.FieldDefinitions.PF.Description = `PF`

	container.Datashape.FieldDefinitions.PFLimit.Ordinal = zeros
	container.Datashape.FieldDefinitions.PFLimit.BaseType = `STRING`
	container.Datashape.FieldDefinitions.PFLimit.Aspects = aspect
	container.Datashape.FieldDefinitions.PFLimit.Name = `PFLimit`
	container.Datashape.FieldDefinitions.PFLimit.Description = `PFLimit`

	container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MaxPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MaxPF.Name = `MaxPF`
	container.Datashape.FieldDefinitions.MaxPF.Description = `MaxPF`

	container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
	container.Datashape.FieldDefinitions.MinPF.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
	container.Datashape.FieldDefinitions.MinPF.Name = `MinPF`
	container.Datashape.FieldDefinitions.MinPF.Description = `MinPF`

	container.Datashape.FieldDefinitions.WeatherTemp.Ordinal = zeros
	container.Datashape.FieldDefinitions.WeatherTemp.BaseType = `STRING`
	container.Datashape.FieldDefinitions.WeatherTemp.Aspects = aspect
	container.Datashape.FieldDefinitions.WeatherTemp.Name = `weather_Temp`
	container.Datashape.FieldDefinitions.WeatherTemp.Description = `weather_Temp`

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

func (sca SCBuildAPI) aggAllDayOnTime(w http.ResponseWriter, req *http.Request) {
	log.Println("aggAllDayOnTime")
	_beforeEndPoint(w, req)
	vars := mux.Vars(req)
	start := vars["start"]
	end := vars["end"]
	bid := vars["bid"]

	container := strAllDayMonthOnTime{}

	// var result []string
	var err error

	mongo := getMongo()
	thepipe := pipeAllDayOnTime(start, end, bid)
	pipe := mongo.DB(DBName).C(collDay).Pipe(thepipe).AllowDiskUse().All(&container.Rows)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		log.Println(currTime.String() + collAll)
	}

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	fmt.Println(req.RemoteAddr)
	mongo.Close()

}

func (sca SCBuildAPI) aggAllMonthOnTime(w http.ResponseWriter, req *http.Request) {
	log.Println("aggAllMonthOnTime")
	_beforeEndPoint(w, req)
	vars := mux.Vars(req)
	start := vars["start"]
	end := vars["end"]
	bid := vars["bid"]

	// container := strAllDayMonthOnTime{}

	var container []interface{}
	// var result []string
	var err error

	mongo := getMongo()
	thepipe := pipeAllMonthOnTime(start, end, bid)
	pipe := mongo.DB(DBName).C(collMonth).Pipe(thepipe).AllowDiskUse().All(&container)

	if pipe != nil {

	}

	// fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		fmt.Println(currTime.String() + collAll)
	}

	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	mongo.Close()
	// fmt.Println(req.RemoteAddr)

}

func (sca SCBuildAPI) aggBuildingAllExplain(w http.ResponseWriter, req *http.Request) {
	log.Println("aggBuildingAllExplain")
	_beforeEndPoint(w, req)

	vars := mux.Vars(req)
	start := vars["start"]
	end := vars["end"]
	bid := vars["bid"]

	var result []string
	var err error

	mongo := getMongo()

	thepipe := pipeAllMonthOnTime(start, end, bid)
	explainErr := bson.M{}

	pipe := mongo.DB(DBName).C(collMonth).Pipe(thepipe).AllowDiskUse().Explain(explainErr)

	if pipe != nil {

	}

	fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		fmt.Println(currTime.String())
	}

	json.NewEncoder(w).Encode(explainErr)
	_afterEndPoint(w, req)
	fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()

}

// This function allows to know details of process during mongodb query.
func (sca SCBuildAPI) queryExplain(w http.ResponseWriter, req *http.Request) {
	_beforeEndPoint(w, req)

	// vars := mux.Vars(req)
	// buildingName := vars["b_id"]

	var result []string
	// var results []interface{}
	var err error

	mongo := getMongo()
	// err = mongo.DB(DBName).C("SC01_DeviceManager").Find(bson.M{"Building_Name": buildingName}).Distinct("devID", &result)

	// thepipe := PipeDevIDArray(result, "lastReportTime")
	thepipe := pipeBuildAllLatest()
	explainErr := bson.M{}

	pipe := mongo.DB(DBName).C(collAll).Pipe(thepipe)

	if pipe != nil {

	}
	fmt.Println(explainErr)

	fmt.Println(result)

	if err != nil {

		_di.Log.Err(err.Error())
	} else {
		currTime := time.Now()

		fmt.Println(currTime.String() + collAll)
	}

	json.NewEncoder(w).Encode(&explainErr)
	_afterEndPoint(w, req)
	fmt.Println(req.RemoteAddr)
	req.Body.Close()
	mongo.Close()
}

func (sca SCBuildAPI) GetTotalHour(w http.ResponseWriter, req *http.Request) {
	log.Println("GetTotalHour")
	vars := mux.Vars(req)
	// bid := vars["bid"]
	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	// var amp []interface{}
	container := strAllHourOnTime{}
	containerz := strAllHourOnTime{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02T15", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02T15", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := getMongo()
		qu := mongo.DB(DBName).C(collHour)
		qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)
		// fmt.Println(devid, tmpstart, container.Rows, err)
		for _, j := range containerz.Rows {
			//Convert to KW
			j.MaxDemand = j.MaxDemand / 1000
			j.MinDemand = j.MinDemand / 1000
			j.PwrDemand = j.PwrDemand / 1000
			// j.PwrUsage = j.PwrUsage / 1000
			// j.MinUsage = j.MinUsage / 1000
			// j.MaxUsage = j.MaxUsage / 1000
			container.Rows = append(container.Rows, j)
		}
		// qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)
		container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
		container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
		container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
		container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
		container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

		container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
		container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
		container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

		container.Datashape.FieldDefinitions.CC.Ordinal = zeros
		container.Datashape.FieldDefinitions.CC.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.CC.Aspects = aspect
		container.Datashape.FieldDefinitions.CC.Name = `CC`
		container.Datashape.FieldDefinitions.CC.Description = `CC`

		container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrUsage.Name = `total_Usage`
		container.Datashape.FieldDefinitions.PwrUsage.Description = `total_Usage`

		container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
		container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

		container.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
		container.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

		container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrDemand.Name = `avg_Demand`
		container.Datashape.FieldDefinitions.PwrDemand.Description = `avg_Demand`

		container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxDemand.Name = `max_Demand`
		container.Datashape.FieldDefinitions.MaxDemand.Description = `max_Demand`

		container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MinDemand.Name = `min_Demand`
		container.Datashape.FieldDefinitions.MinDemand.Description = `min_Demand`

		container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MinUsage.Name = `min_Usage`
		container.Datashape.FieldDefinitions.MinUsage.Description = `min_Usage`

		container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxUsage.Name = `max_Usage`
		container.Datashape.FieldDefinitions.MaxUsage.Description = `max_Usage`

		container.Datashape.FieldDefinitions.PF.Ordinal = zeros
		container.Datashape.FieldDefinitions.PF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PF.Aspects = aspect
		container.Datashape.FieldDefinitions.PF.Name = `avg_PF`
		container.Datashape.FieldDefinitions.PF.Description = `avg_PF`

		container.Datashape.FieldDefinitions.PFLimit.Ordinal = zeros
		container.Datashape.FieldDefinitions.PFLimit.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PFLimit.Aspects = aspect
		container.Datashape.FieldDefinitions.PFLimit.Name = `PF_Limit`
		container.Datashape.FieldDefinitions.PFLimit.Description = `PF_Limit`

		container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxPF.Name = `max_PF`
		container.Datashape.FieldDefinitions.MaxPF.Description = `max_PF`

		container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MinPF.Name = `min_PF`
		container.Datashape.FieldDefinitions.MinPF.Description = `min_PF`

		container.Datashape.FieldDefinitions.WeatherTemp.Ordinal = zeros
		container.Datashape.FieldDefinitions.WeatherTemp.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.WeatherTemp.Aspects = aspect
		container.Datashape.FieldDefinitions.WeatherTemp.Name = `weather_Temp`
		container.Datashape.FieldDefinitions.WeatherTemp.Description = `weather_Temp`

		json.NewEncoder(w).Encode(container)
		_afterEndPoint(w, req)
		// fmt.Println(req.RemoteAddr)
		req.Body.Close()
		mongo.Close()
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

func (sca SCBuildAPI) GetTotalDay(w http.ResponseWriter, req *http.Request) {
	log.Println("GetTotalDay")
	vars := mux.Vars(req)

	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	container := strAllDayMonthOnTime{}
	containerz := strAllDayMonthOnTime{}

	// containerForDevManager := []DeviceManagerS{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01-02", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01-02", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := getMongo()
		// mongo.DB(DBName).C(DeviceManager).Find(bson.M{"devID": devid}).One(&containerForDevManager)

		qu := mongo.DB(DBName).C(collDayTotal)
		qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)

		for _, j := range containerz.Rows {
			//Convert to KW
			j.MaxDemand = j.MaxDemand / 1000
			j.MinDemand = j.MinDemand / 1000
			j.PwrDemand = j.PwrDemand / 1000
			// j.PwrUsage = j.PwrUsage / 1000
			// j.MinUsage = j.MinUsage / 1000
			// j.MaxUsage = j.MaxUsage / 1000
			container.Rows = append(container.Rows, j)
		}
		// qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)
		container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
		container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
		container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
		container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
		container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

		container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
		container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
		container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

		container.Datashape.FieldDefinitions.CC.Ordinal = zeros
		container.Datashape.FieldDefinitions.CC.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.CC.Aspects = aspect
		container.Datashape.FieldDefinitions.CC.Name = `CC`
		container.Datashape.FieldDefinitions.CC.Description = `CC`

		container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrUsage.Name = `avg_Usage`
		container.Datashape.FieldDefinitions.PwrUsage.Description = `avg_Usage`

		container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
		container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

		container.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
		container.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

		container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrDemand.Name = `avg_Demand`
		container.Datashape.FieldDefinitions.PwrDemand.Description = `avg_Demand`

		container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxDemand.Name = `max_Demand`
		container.Datashape.FieldDefinitions.MaxDemand.Description = `max_Demand`

		container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MinDemand.Name = `min_Demand`
		container.Datashape.FieldDefinitions.MinDemand.Description = `min_Demand`

		container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MinUsage.Name = `min_Usage`
		container.Datashape.FieldDefinitions.MinUsage.Description = `min_Usage`

		container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxUsage.Name = `max_Usage`
		container.Datashape.FieldDefinitions.MaxUsage.Description = `max_Usage`

		container.Datashape.FieldDefinitions.PF.Ordinal = zeros
		container.Datashape.FieldDefinitions.PF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PF.Aspects = aspect
		container.Datashape.FieldDefinitions.PF.Name = `avg_PF`
		container.Datashape.FieldDefinitions.PF.Description = `avg_PF`

		container.Datashape.FieldDefinitions.PFLimit.Ordinal = zeros
		container.Datashape.FieldDefinitions.PFLimit.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PFLimit.Aspects = aspect
		container.Datashape.FieldDefinitions.PFLimit.Name = `PF_Limit`
		container.Datashape.FieldDefinitions.PFLimit.Description = `PF_Limit`

		container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxPF.Name = `max_PF`
		container.Datashape.FieldDefinitions.MaxPF.Description = `max_PF`

		container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MinPF.Name = `min_PF`
		container.Datashape.FieldDefinitions.MinPF.Description = `min_PF`

		container.Datashape.FieldDefinitions.WeatherTemp.Ordinal = zeros
		container.Datashape.FieldDefinitions.WeatherTemp.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.WeatherTemp.Aspects = aspect
		container.Datashape.FieldDefinitions.WeatherTemp.Name = `weather_Temp`
		container.Datashape.FieldDefinitions.WeatherTemp.Description = `weather_Temp`

		json.NewEncoder(w).Encode(&container)
		_afterEndPoint(w, req)
		// fmt.Println(req.RemoteAddr)
		req.Body.Close()
		mongo.Close()

	}

}

func (sca SCBuildAPI) GetTotalMonth(w http.ResponseWriter, req *http.Request) {
	log.Println("GetTotalMonth")
	vars := mux.Vars(req)

	devid := vars["devid"]
	start := vars["start"]
	end := vars["end"]

	zeros := 0
	aspect := bson.M{}

	container := strAllDayMonthOnTime{}
	containerz := strAllDayMonthOnTime{}

	// containerForDevManager := DeviceManagerS{}

	if start != "" && end != "" {
		tmpstart, e := time.ParseInLocation("2006-01", start, time.Local)
		tmpend, er := time.ParseInLocation("2006-01", end, time.Local)

		// fmt.Println(tmpend)

		if e != nil {
			_di.Log.Err(e.Error())
		}

		if er != nil {
			_di.Log.Err(er.Error())
		}

		mongo := getMongo()
		// mongo.DB(DBName).C(DeviceManager).Find(bson.M{"devID": devid}).All(&containerForDevManager)

		qu := mongo.DB(DBName).C(collMonthTotal)
		qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).All(&containerz.Rows)
		for _, j := range containerz.Rows {
			//Convert to KW
			j.MaxDemand = j.MaxDemand / 1000
			j.MinDemand = j.MinDemand / 1000
			j.PwrDemand = j.PwrDemand / 1000
			// j.PwrUsage = j.PwrUsage / 1000
			// j.MinUsage = j.MinUsage / 1000
			// j.MaxUsage = j.MaxUsage / 1000
			container.Rows = append(container.Rows, j)
		}
		// qu.Find(bson.M{"Device_ID": devid, "lastReportTime": bson.M{"$gte": tmpstart, "$lte": tmpend}}).Explain(&amp)
		container.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
		container.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
		container.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
		container.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
		container.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

		container.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
		container.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
		container.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
		container.Datashape.FieldDefinitions.LastReportTime.Name = `lastReportTime`
		container.Datashape.FieldDefinitions.LastReportTime.Description = `lastReportTime`

		container.Datashape.FieldDefinitions.CC.Ordinal = zeros
		container.Datashape.FieldDefinitions.CC.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.CC.Aspects = aspect
		container.Datashape.FieldDefinitions.CC.Name = `CC`
		container.Datashape.FieldDefinitions.CC.Description = `CC`

		container.Datashape.FieldDefinitions.PwrUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrUsage.Name = `avg_Usage`
		container.Datashape.FieldDefinitions.PwrUsage.Description = `avg_Usage`

		container.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
		container.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

		container.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
		container.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
		container.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
		container.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
		container.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

		container.Datashape.FieldDefinitions.PwrDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.PwrDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PwrDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.PwrDemand.Name = `avg_Demand`
		container.Datashape.FieldDefinitions.PwrDemand.Description = `avg_Demand`

		container.Datashape.FieldDefinitions.MaxDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxDemand.Name = `max_Demand`
		container.Datashape.FieldDefinitions.MaxDemand.Description = `max_Demand`

		container.Datashape.FieldDefinitions.MinDemand.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinDemand.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinDemand.Aspects = aspect
		container.Datashape.FieldDefinitions.MinDemand.Name = `min_Demand`
		container.Datashape.FieldDefinitions.MinDemand.Description = `min_Demand`

		container.Datashape.FieldDefinitions.MinUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MinUsage.Name = `min_Usage`
		container.Datashape.FieldDefinitions.MinUsage.Description = `min_Usage`

		container.Datashape.FieldDefinitions.MaxUsage.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxUsage.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxUsage.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxUsage.Name = `max_Usage`
		container.Datashape.FieldDefinitions.MaxUsage.Description = `max_Usage`

		container.Datashape.FieldDefinitions.PF.Ordinal = zeros
		container.Datashape.FieldDefinitions.PF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PF.Aspects = aspect
		container.Datashape.FieldDefinitions.PF.Name = `avg_PF`
		container.Datashape.FieldDefinitions.PF.Description = `avg_PF`

		container.Datashape.FieldDefinitions.PFLimit.Ordinal = zeros
		container.Datashape.FieldDefinitions.PFLimit.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.PFLimit.Aspects = aspect
		container.Datashape.FieldDefinitions.PFLimit.Name = `PF_Limit`
		container.Datashape.FieldDefinitions.PFLimit.Description = `PF_Limit`

		container.Datashape.FieldDefinitions.MaxPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MaxPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MaxPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MaxPF.Name = `max_PF`
		container.Datashape.FieldDefinitions.MaxPF.Description = `max_PF`

		container.Datashape.FieldDefinitions.MinPF.Ordinal = zeros
		container.Datashape.FieldDefinitions.MinPF.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.MinPF.Aspects = aspect
		container.Datashape.FieldDefinitions.MinPF.Name = `min_PF`
		container.Datashape.FieldDefinitions.MinPF.Description = `min_PF`

		container.Datashape.FieldDefinitions.WeatherTemp.Ordinal = zeros
		container.Datashape.FieldDefinitions.WeatherTemp.BaseType = `NUMBER`
		container.Datashape.FieldDefinitions.WeatherTemp.Aspects = aspect
		container.Datashape.FieldDefinitions.WeatherTemp.Name = `weather_Temp`
		container.Datashape.FieldDefinitions.WeatherTemp.Description = `weather_Temp`

		json.NewEncoder(w).Encode(&container)
		_afterEndPoint(w, req)
		// fmt.Println(req.RemoteAddr)
		req.Body.Close()
		mongo.Close()
	}

}

func (sca SCBuildAPI) testaest(w http.ResponseWriter, req *http.Request) {
	// aa := `{

	type aaa struct {
		Address            string  `json:"address"`
		Voltages           string  `json:"voltages,omitempty"`
		Hz                 int     `json:"hz,omitempty"`
		DemandCount        int     `json:"demandCount,omitempty"`
		DemandWatts        int     `json:"demandWatts,omitempty"`
		Current            int     `json:"current,omitempty"`
		Battery            int     `json:"battery,omitempty"`
		Devicename         string  `json:"devicename,omitempty"`
		Temp               int     `json:"temp,omitempty"`
		CurrentAc          string  `json:"current_ac,omitempty"`
		Watts              float64 `json:"watts,omitempty"`
		Username           string  `json:"username,omitempty"`
		AlertCurrent       string  `json:"alertCurrent,omitempty"`
		LastSendAlert      int     `json:"lastSendAlert,omitempty"`
		Calibration        string  `json:"calibration,omitempty"`
		Lastupdated        int64   `json:"lastupdated,omitempty"`
		Number             int     `json:"number,omitempty"`
		Status             int     `json:"status" bson:"status"`
		Starttime          int64   `json:"starttime,omitempty"`
		Stoptime           int64   `json:"stoptime,omitempty"`
		StopCurrent        int     `json:"stopCurrent,omitempty"`
		MaxCurrentRaw      int     `json:"maxCurrent_raw,omitempty"`
		CurrentByte        int     `json:"current_byte,omitempty"`
		VariantCurrent     int     `json:"variantCurrent,omitempty"`
		SamplingLeng       int     `json:"samplingLeng,omitempty"`
		FilterPer          int     `json:"filterPer,omitempty"`
		Type               string  `json:"type,omitempty"`
		CurrentHealthValue int     `json:"current_healthValue,omitempty"`
		DeviceNickname     string  `json:"DeviceNickname,omitempty"`
	}

	aa := aaa{}
	json.NewDecoder(req.Body).Decode(&aa)

	// 	"address":             "98:07:2D:0C:45:68",
	// 	"voltages":            "110",
	// 	"hz":                  0,
	// 	"demandCount":         1,
	// 	"demandWatts":         0,
	// 	"current":             0,
	// 	"battery":             81,
	// 	"devicename":          "",
	// 	"temp":                27,
	// 	"current_ac":          "111.2",
	// 	"watts":               12.2,
	// 	"username":            "NTUST",
	// 	"alertCurrent":        "1",
	// 	"lastSendAlert":       0,
	// 	"calibration":         "0",
	// 	"lastupdated":         1526278351511,
	// 	"number":              0,
	// 	"status":              -1,
	// 	"starttime":           0,
	// 	"stoptime":            1526278339831,
	// 	"stopCurrent":         100,
	// 	"maxCurrent_raw":      -1,
	// 	"current_byte":        -1,
	// 	"variantCurrent":      -1,
	// 	"samplingLeng":        -1,
	// 	"filterPer":           1,
	// 	"type":                "CM",
	// 	"current_healthValue": -1,
	// 	"DeviceNickname":      "國際大樓11F"
	// },`
	fmt.Println(aa)
}

func (sca SCBuildAPI) GenAuth(w http.ResponseWriter, req *http.Request) {
	log.Println("GenAuth")
	vars := mux.Vars(req)
	theuser := vars["user"]

	year, _, day := time.Now().Date()
	yearDay := time.Now().YearDay()

	container := jwtSignation{
		jwt.StandardClaims{
			Issuer:   theuser,
			Audience: strconv.Itoa(year) + strconv.Itoa(yearDay) + strconv.Itoa(day),
		},
		"thingworx",
		"123",
	}
	// 	jwt.StandardClaims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, container)

	tokenString, _ := token.SignedString(SigningKey)
	w.Write([]byte(tokenString))
	// json.NewEncoder(w).Encode(tokenString)
	// fmt.Println(SigningKey)
}

func (sca SCBuildAPI) Validator(w http.ResponseWriter, req *http.Request) {
	log.Println("Validator")
	vars := mux.Vars(req)
	keyz := vars["key"]
	container := RegisterDevID{}
	mongo := getMongo()

	year, _, day := time.Now().Local().Date()
	yearDay := time.Now().Local().YearDay()
	// jwtcontainer := jwtSignation{}
	token, err := jwt.ParseWithClaims(keyz, &jwtSignation{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return SigningKey, nil
	})

	if ve, lll := err.(*jwt.ValidationError); lll {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			http.Error(w, "ValidationError", http.StatusNotAcceptable)
			fmt.Println("ValidationErrorMalformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			http.Error(w, "ValidationErrorExpired", http.StatusNotAcceptable)
			fmt.Println("ValidationErrorExpired")
		}
	} else {
		// if claims, ok := token.Claims.(*jwtSignation); ok && token.Valid {

		if claims, ok := token.Claims.(*jwtSignation); ok && token.Valid {

			claimsPointA := claims.StandardClaims.VerifyIssuer("aaeon", ok)
			claimsPointB := claims.StandardClaims.VerifyAudience(strconv.Itoa(year)+strconv.Itoa(yearDay)+strconv.Itoa(day), false)

			if claimsPointA != true {
				http.Error(w, "ValidationErrorIssuer", http.StatusNotAcceptable)
				fmt.Println("ValidationErrorIssuer")
			} else if claimsPointB != true {
				http.Error(w, "ValidationErrorExpired", http.StatusNotAcceptable)
				fmt.Println("ValidationErrorExpired")
			} else {
				json.NewDecoder(req.Body).Decode(&container)
				fmt.Println(strconv.Itoa(year)+strconv.Itoa(yearDay)+strconv.Itoa(day), claimsPointB)

				if container.BuildingDetails == "" || container.BuildingDetails == "undefined" {
					http.Error(w, "BuildingDetails null", http.StatusBadRequest)
					fmt.Println("BuildingDetails null")
				} else if container.BuildingName == "" || container.BuildingName == "undefined" {
					http.Error(w, "BuildingName null", http.StatusBadRequest)
					fmt.Println("BuildingName null")
				} else if container.BuildingDetails == "" || container.BuildingDetails == "undefined" {
					http.Error(w, "DeviceBrand null", http.StatusBadRequest)
					fmt.Println("DeviceBrand null")
				} else if container.DeviceDetails == "" || container.DeviceDetails == "undefined" {
					http.Error(w, "DeviceDetails null", http.StatusBadRequest)
					fmt.Println("DeviceDetails null")
				} else if container.DeviceID == "" || container.DeviceID == "undefined" {
					http.Error(w, "devID null", http.StatusBadRequest)
					fmt.Println("devID null")
				} else if container.DeviceInfo == "" || container.DeviceInfo == "undefined" {
					http.Error(w, "DeviceInfo null", http.StatusBadRequest)
					fmt.Println("DeviceInfo null")
				} else if container.DeviceName == "" || container.DeviceName == "undefined" {
					http.Error(w, "DeviceName null", http.StatusBadRequest)
					fmt.Println("DeviceName null")
				} else if container.DeviceType == "" || container.DeviceType == "undefined" {
					http.Error(w, "DeviceType null", http.StatusBadRequest)
					fmt.Println("DeviceType null")
				} else if container.Floor == "" || container.Floor == "undefined" {
					http.Error(w, "Floor null", http.StatusBadRequest)
					fmt.Println("Floor null")
				} else if container.GatewayID == "" || container.GatewayID == "undefined" {
					http.Error(w, "GWID null", http.StatusBadRequest)
					fmt.Println("GWID null")
				} else {
					fmt.Println("BDetails :", container.BuildingDetails)
					container.TimeAdded = time.Now()
					erro := mongo.DB(DBName).C(DeviceManager).Upsert(bson.M{"devID": container.DeviceID}, container)
					w.WriteHeader(http.StatusOK)
					fmt.Println("Device Added : ", container.DeviceID)
					json.NewEncoder(w).Encode(container)
					_afterEndPoint(w, req)
					if erro != nil {

						w.WriteHeader(http.StatusBadRequest)
						log.Println(container.DeviceID, container, err)
					}
				}

			}

		}
	}
}

func (sca SCBuildAPI) AddDev(w http.ResponseWriter, req *http.Request) {
	log.Println("AddDev")
	_beforeEndPoint(w, req)

	// tokenString := req.Header.Get("authorization")
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("Unexpected sigining method: %v", token.Header["alg"])
	// 	}

	// 	return _CLIENTSEC, nil
	// })
	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	container := RegisterDevID{}

	json.NewDecoder(req.Body).Decode(&container)

	mongo := getMongo()
	err := mongo.DB(DBName).C(DeviceManager).Upsert(bson.M{"devID": container.DeviceID}, container)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(container.DeviceID, container, err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(container)
	_afterEndPoint(w, req)
	mongo.Close()
	// 	fmt.Println(req.RemoteAddr, claims)
	// } else {
	// 	fmt.Println(err)
	// }
}

func (sca SCBuildAPI) ListDevices(w http.ResponseWriter, req *http.Request) {
	log.Println("ListDevices")
	_beforeEndPoint(w, req)

	tmpContainer := ListDevices{}

	var explain interface{}
	zeros := 0
	aspect := bson.M{}

	mongo := getMongo()

	mongo.DB(DBName).C(DeviceManager).Find(nil).All(&tmpContainer.Rows)
	mongo.DB(DBName).C(DeviceManager).Find(nil).Explain(&explain)
	// fmt.Println(explain)

	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.BaseType = `DATETIME`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Name = `Time_Added`
	tmpContainer.Datashape.FieldDefinitions.LastReportTime.Description = `Time_Added`

	tmpContainer.Datashape.FieldDefinitions.GatewayID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.GatewayID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Name = `Gateway_ID`
	tmpContainer.Datashape.FieldDefinitions.GatewayID.Description = `Gateway_ID`

	tmpContainer.Datashape.FieldDefinitions.BuildingName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Name = `Building_Name`
	tmpContainer.Datashape.FieldDefinitions.BuildingName.Description = `Building_Name`

	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Name = `Building_Details`
	tmpContainer.Datashape.FieldDefinitions.BuildingDetails.Description = `Building_Details`

	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Name = `Device_Brand`
	tmpContainer.Datashape.FieldDefinitions.DeviceBrand.Description = `Device_Brand`

	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Name = `Device_Details`
	tmpContainer.Datashape.FieldDefinitions.DeviceDetails.Description = `Device_Details`

	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Name = `Device_Info`
	tmpContainer.Datashape.FieldDefinitions.DeviceInfo.Description = `Device_Info`

	tmpContainer.Datashape.FieldDefinitions.DeviceName.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceName.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Name = `Device_Name`
	tmpContainer.Datashape.FieldDefinitions.DeviceName.Description = `Device_Name`

	tmpContainer.Datashape.FieldDefinitions.DeviceID.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceID.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Name = `Device_ID`
	tmpContainer.Datashape.FieldDefinitions.DeviceID.Description = `Device_ID`

	tmpContainer.Datashape.FieldDefinitions.DeviceType.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.DeviceType.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Name = `Device_Type`
	tmpContainer.Datashape.FieldDefinitions.DeviceType.Description = `Device_Type`

	tmpContainer.Datashape.FieldDefinitions.Floor.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Floor.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.Floor.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Floor.Name = `Floor`
	tmpContainer.Datashape.FieldDefinitions.Floor.Description = `Floor`

	tmpContainer.Datashape.FieldDefinitions.Facility.Ordinal = zeros
	tmpContainer.Datashape.FieldDefinitions.Facility.BaseType = `STRING`
	tmpContainer.Datashape.FieldDefinitions.Facility.Aspects = aspect
	tmpContainer.Datashape.FieldDefinitions.Facility.Name = `Facility`
	tmpContainer.Datashape.FieldDefinitions.Facility.Description = `Facility`

	json.NewEncoder(w).Encode(tmpContainer)
	_afterEndPoint(w, req)
	// fmt.Println(req.RemoteAddr)
	defer req.Body.Close()
	mongo.Close()

}

func (sca SCBuildAPI) RemoveDev(w http.ResponseWriter, req *http.Request) {
	log.Println("RemoveDev")
	_beforeEndPoint(w, req)

	// tokenString := req.Header.Get("authorization")
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("Unexpected sigining method: %v", token.Header["alg"])
	// 	}

	// 	return _CLIENTSEC, nil
	// })
	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	container := RegisterDevID{}

	json.NewDecoder(req.Body).Decode(&container)
	// fmt.Println(container)
	mongo := getMongo()
	err := mongo.DB(DBName).C(DeviceManager).Remove(bson.M{"devID": container.DeviceID, "GWID": container.GatewayID})

	// err := mongo.DB(DBName).C(DeviceManager).Upsert(bson.M{"devID": container.DeviceID}, container)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)

		http.Error(w, container.DeviceID+" "+container.GatewayID+" not found", http.StatusNotFound)
		fmt.Println(container.DeviceID, container, err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(container)
	}

	_afterEndPoint(w, req)
	mongo.Close()
	// 	fmt.Println(req.RemoteAddr, claims)
	// } else {
	// 	fmt.Println(err)
	// }
}

//GW Status
