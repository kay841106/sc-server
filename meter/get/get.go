package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"
	// "/Users/avbee/go/src/sc-server/meter/get/auth"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	// change due to high cpu using globalsign
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:10003"
	dbbackup = "140.118.122.103:27017"

	db             = "sc"
	c_lastreport   = "lastreport"
	c_devices      = "devices"
	c_gwtstat      = "gw_status"
	c_offlineChart = "offline_chart"

	c_hourly = "hour"
	c_daily  = "day"
	c_month  = "month"
)

var session *mgo.Session

func dbConnect() {
	var err error
	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbpublic, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 10,
	}

	session, err = mgo.DialWithInfo(dbInfo)

	if err != nil {
		os.Exit(1)
	}

}

type CPMSnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string        `json:"GW_ID" bson:"GW_ID"`
}

type Metrics struct {
	GET11 float64 `json:"ae_tot" bson:"ae_tot"`
	GET12 float64 `json:"p_sum" bson:"p_sum"`
	GET13 float64 `json:"pf_avg" bson:"pf_avg"`
}

type DSTwxTemplate struct {
	Name        string      `json:"name,omitempty" bson:"name"`
	Description string      `json:"description,omitempty" bson:"description"`
	BaseType    string      `json:"baseType,omitempty" bson:"baseType"`
	Ordinal     int         `json:"ordinal" bson:"ordinal"`
	Aspects     interface{} `json:"aspects,omitempty" bson:"aspects"`
}

type device struct {
	Rows []struct {
		DevID      int    `json:"DevID" bson:"DevID"`
		Floor      string `json:"Floor" bson:"Floor"`
		GWID       string `json:"GWID" bson:"GWID"`
		MGWID      string `json:"M_GWID" bson:"M_GWID"`
		MMAC       string `json:"M_MAC" bson:"M_MAC"`
		NUM        string `json:"NUM" bson:"NUM"`
		Place      string `json:"Place" bson:"Place"`
		Territory  string `json:"Territory" bson:"Territory"`
		Type       string `json:"Type" bson:"Type"`
		MACAddress string `json:"MAC_Address" bson:"MAC_Address"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			DevID      DSTwxTemplate `json:"DevID"`
			Floor      DSTwxTemplate `json:"Floor"`
			GWID       DSTwxTemplate `json:"GWID"`
			MGWID      DSTwxTemplate `json:"M_GWID"`
			MMAC       DSTwxTemplate `json:"M_MAC" `
			NUM        DSTwxTemplate `json:"NUM" `
			Place      DSTwxTemplate `json:"Place" `
			Territory  DSTwxTemplate `json:"Territory" `
			Type       DSTwxTemplate `json:"Type" `
			MACAddress DSTwxTemplate `json:"MACAddress" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type getlastreport struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string    `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string    `json:"GWID" bson:"GWID"`
	DevID         int       `json:"ID" bson:"ID"`
	Floor         string    `json:"Floor" bson:"Floor"`
	MGWID         string    `json:"M_GWID" bson:"M_GWID"`
	MMAC          string    `json:"M_MAC" bson:"M_MAC"`
	NUM           string    `json:"NUM" bson:"NUM"`
	Place         string    `json:"PLACE" bson:"PLACE"`
	Territory     string    `json:"TERRITORY" bson:"TERRITORY"`
	Type          string    `json:"TYPE" bson:"TYPE"`
	Mfloor        string    `json:"meter_floor" bson:"meter_floor"`
	Mplace        string    `json:"meter_place" bson:"meter_place"`
	Metrics       Metrics
}

type postlastreport struct {
	MACAddress string `json:"MAC_Address" bson:"MAC_Address"`
}

type postGWID struct {
	GWID string `json:"GW_ID" bson:"GW_ID"`
}

func gopostlastreport(w http.ResponseWriter, r *http.Request) {

	headercontainer := postlastreport{}

	json.NewDecoder(r.Body).Decode(&headercontainer)
	// fmt.Println(headercontainer)
	container := []getlastreport{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_lastreport)
	Mongo.Find(bson.M{"MAC_Address": headercontainer.MACAddress}).All(&container)
	json.NewEncoder(w).Encode(container)
	// fmt.Println(container)
}

func gopostlastreportAllBuilding(w http.ResponseWriter, r *http.Request) {

	buildingsID := []string{"aa:bb:02:03:01:01", "aa:bb:02:06:01:01", "aa:bb:02:06:01:02", "aa:bb:02:03:02:01",
		"aa:bb:02:03:03:01", "aa:bb:02:09:01:01", "aa:bb:02:10:01:01", "aa:bb:02:03:04:01", "aa:bb:02:14:01:02",
		"aa:bb:02:12:01:01", "aa:bb:02:15:01:01", "aa:bb:02:04:01:01", "aa:bb:02:16:01:01", "aa:bb:02:18:01:01", "aa:bb:02:18:01:02",
		"aa:bb:02:07:01:01"}

	container := []getlastreport{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_lastreport)
	err := Mongo.Find(bson.M{"MAC_Address": bson.M{"$in": buildingsID}}).All(&container)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(container)
	fmt.Println(buildingsID)
}

func gopostgwstat(w http.ResponseWriter, r *http.Request) {

	headercontainer := postGWID{}

	json.NewDecoder(r.Body).Decode(&headercontainer)
	// fmt.Println(headercontainer)
	container := []gwstat{}
	sess := session.Clone()

	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{"GW_ID": headercontainer.GWID}).All(&container)
	json.NewEncoder(w).Encode(container)
	// fmt.Println(container)
}

func gogetlastreport(w http.ResponseWriter, r *http.Request) {

	container := []getlastreport{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_lastreport)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	// fmt.Println(container)
}

type devices struct {
	MACAddress string `json:"MAC_Address" bson:"MAC_Address"`
	DevID      int    `json:"DevID" bson:"DevID"`
	Floor      string `json:"Floor" bson:"Floor"`
	GWID       string `json:"GWID" bson:"GWID"`
	MGWID      string `json:"M_GWID" bson:"M_GWID"`
	MMAC       string `json:"M_MAC" bson:"M_MAC"`
	NUM        string `json:"NUM" bson:"NUM"`
	Place      string `json:"Place" bson:"Place"`
	Territory  string `json:"Territory" bson:"Territory"`
	Type       string `json:"Type" bson:"Type"`
}

func gogetDevices(w http.ResponseWriter, r *http.Request) {

	container := []devices{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_devices)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	// fmt.Println(container)
}

type gwstat struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
	MGWID         string    `json:"M_GWID" bson:"M_GWID"`
	Place         string    `json:"Place" bson:"Place"`
}

type gwdata struct {
	GWID  string `json:"GWID" bson:"GWID"`
	MGWID string `json:"M_GWID" bson:"M_GWID"`
	// MMAC      string `json:"M_MAC" bson:"M_MAC"`
	// NUM       string `json:"NUM" bson:"NUM"`
	Place string `json:"Place" bson:"Place"`
	// Territory string `json:"Territory" bson:"Territory"`
	// Type      string `json:"Type" bson:"Type"`
}

//////////////////////////////////SPACE

type queryRes struct {
	Rows     *[]interface{} `json:"result,omitempty"`
	Total    int            `json:"total,omitempty"`
	AllPages int            `json:"allPages, omitempty"`
	Page     int            `json:"page,omitempty"`
	Limit    int            `json:"limit, omitempty"`
}
type onlyRes struct {
	Rows *[]interface{} `json:"result,omitempty"`
}

type deviceStatusRes struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	RoomID          int      `json:"roomID"`
	Type            string   `json:"type"`
	BaseType        string   `json:"baseType"`
	Enabled         bool     `json:"enabled"`
	Visible         bool     `json:"visible"`
	IsPlugin        bool     `json:"isPlugin"`
	ParentID        int      `json:"parentId"`
	RemoteGatewayID int      `json:"remoteGatewayId"`
	Interfaces      []string `json:"interfaces"`
	Properties      struct {
		Parameters []struct {
			ID                int `json:"id"`
			LastReportedValue int `json:"lastReportedValue"`
			LastSetValue      int `json:"lastSetValue"`
			Size              int `json:"size"`
			Value             int `json:"value"`
		} `json:"parameters"`
		PollingTimeSec         int    `json:"pollingTimeSec"`
		WakeUpTime             int    `json:"wakeUpTime"`
		ZwaveCompany           string `json:"zwaveCompany"`
		ZwaveInfo              string `json:"zwaveInfo"`
		ZwaveVersion           string `json:"zwaveVersion"`
		AlarmDelay             string `json:"alarmDelay"`
		AlarmExclude           string `json:"alarmExclude"`
		AlarmLevel             string `json:"alarmLevel"`
		AlarmTimeTimestamp     string `json:"alarmTimeTimestamp"`
		AlarmType              string `json:"alarmType"`
		ArmConditions          string `json:"armConditions"`
		ArmConfig              string `json:"armConfig"`
		ArmDelay               string `json:"armDelay"`
		ArmError               string `json:"armError"`
		ArmTimeTimestamp       string `json:"armTimeTimestamp"`
		Armed                  string `json:"armed"`
		BatteryLevel           string `json:"batteryLevel"`
		BatteryLowNotification string `json:"batteryLowNotification"`
		Configured             bool   `json:"configured"`
		Dead                   string `json:"dead"`
		DefInterval            string `json:"defInterval"`
		DeviceControlType      string `json:"deviceControlType"`
		DeviceIcon             string `json:"deviceIcon"`
		EmailNotificationID    string `json:"emailNotificationID"`
		EmailNotificationType  string `json:"emailNotificationType"`
		EndPointID             string `json:"endPointId"`
		FibaroAlarm            string `json:"fibaroAlarm"`
		FirmwareUpdate         string `json:"firmwareUpdate"`
		LastBreached           string `json:"lastBreached"`
		LiliOffCommand         string `json:"liliOffCommand"`
		LiliOnCommand          string `json:"liliOnCommand"`
		Log                    string `json:"log"`
		LogTemp                string `json:"logTemp"`
		Manufacturer           string `json:"manufacturer"`
		MarkAsDead             string `json:"markAsDead"`
		MaxInterval            string `json:"maxInterval"`
		MinInterval            string `json:"minInterval"`
		Model                  string `json:"model"`
		NodeID                 string `json:"nodeId"`
		ParametersTemplate     string `json:"parametersTemplate"`
		ProductInfo            string `json:"productInfo"`
		PushNotificationID     string `json:"pushNotificationID"`
		PushNotificationType   string `json:"pushNotificationType"`
		RemoteGatewayID        string `json:"remoteGatewayId"`
		SaveLogs               string `json:"saveLogs"`
		SerialNumber           string `json:"serialNumber"`
		SmsNotificationID      string `json:"smsNotificationID"`
		SmsNotificationType    string `json:"smsNotificationType"`
		StepInterval           string `json:"stepInterval"`
		Tamper                 string `json:"tamper"`
		UpdateVersion          string `json:"updateVersion"`
		UseTemplate            string `json:"useTemplate"`
		UserDescription        string `json:"userDescription"`
		Value                  string `json:"value"`
	} `json:"properties"`
	Actions struct {
		AbortUpdate       int `json:"abortUpdate"`
		ForceArm          int `json:"forceArm"`
		MeetArmConditions int `json:"meetArmConditions"`
		Reconfigure       int `json:"reconfigure"`
		RetryUpdate       int `json:"retryUpdate"`
		SetArmed          int `json:"setArmed"`
		SetInterval       int `json:"setInterval"`
		StartUpdate       int `json:"startUpdate"`
		UpdateFirmware    int `json:"updateFirmware"`
	} `json:"actions"`
	Created   int `json:"created"`
	Modified  int `json:"modified"`
	SortOrder int `json:"sortOrder"`
}

type airCond struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	RoomID     int    `json:"roomID"`
	Type       string `json:"type"`
	Visible    bool   `json:"visible"`
	Enabled    bool   `json:"enabled"`
	Properties struct {
		DeviceIcon  int    `json:"deviceIcon"`
		IP          string `json:"ip"`
		Port        int    `json:"port"`
		CurrentIcon string `json:"currentIcon"`
		Log         string `json:"log"`
		LogTemp     string `json:"logTemp"`
		MainLoop    string `json:"mainLoop"`
		UIACValue   string `json:"ui.AC.value"`
		Visible     string `json:"visible"`
		Rows        []struct {
			Type     string `json:"type"`
			Elements []struct {
				ID              int    `json:"id"`
				Lua             bool   `json:"lua"`
				WaitForResponse bool   `json:"waitForResponse"`
				Caption         string `json:"caption"`
				Name            string `json:"name"`
				Favourite       bool   `json:"favourite"`
				Main            bool   `json:"main"`
			} `json:"elements"`
		} `json:"rows"`
	} `json:"properties"`
	Actions struct {
		PressButton int `json:"pressButton"`
		SetSlider   int `json:"setSlider"`
	} `json:"actions"`
	Created   int `json:"created"`
	Modified  int `json:"modified"`
	SortOrder int `json:"sortOrder"`
}

type powerCons struct {
	ID  int `json:"id"`
	KWh int `json:"kWh"`
	W   int `json:"W"`
	Min int `json:"min"`
	Max int `json:"max"`
	Avg int `json:"avg"`
}

type checkState struct {
	State bool
}

func checkBool(x string) bool {

	tmpVal, _ := strconv.ParseFloat(x, 32)
	// fmt.Println(tmpVal)
	// fmt.Println(x)
	if (tmpVal) > 5 {
		return true
	} else {
		return false
	}

}

///////////////////////////////////

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// func unique2(intSlice []gwdata) []string {
// 	keys := make(map[string]bool)
// 	list := []string{}
// 	for _, entry := range intSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }

func gogetgwstat(w http.ResponseWriter, r *http.Request) {

	container := []gwstat{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{}).All(&container)
	json.NewEncoder(w).Encode(container)
	// fmt.Println(container)
}

type offlineChart struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int       `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MeterOffline  int       `json:"Meter_Offline" bson:"Meter_Offline"`
	GWOffline     int       `json:"GW_Offline" bson:"GW_Offline"`
}

type postTime struct {
	Start *string `json:"Start" bson:"Start"`
	Stop  *string `json:"Stop" bson:"Stop"`
}

func SetTimeStampForLastDay(theTime time.Time) time.Time {
	year, month, day := theTime.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
}

func UtcToLocal(theTime time.Time) time.Time {
	// year, month, day := theTime.Date()
	// h, m, s := theTime.Clock()

	// secondsEastOfUTC := int((8 * time.Hour).Seconds())
	// beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	return theTime.Add(time.Hour * 16)
}

func gogetgofflinechart(w http.ResponseWriter, r *http.Request) {

	headercontainer := postTime{}

	// container := []offlineChart{}
	json.NewDecoder(r.Body).Decode(&headercontainer)
	fmt.Print(headercontainer)
	if *headercontainer.Start != "" && *headercontainer.Stop != "" {

		start, e := time.ParseInLocation("2006-01-02", *headercontainer.Start, time.Local)
		stop, er := time.ParseInLocation("2006-01-02", *headercontainer.Stop, time.Local)

		if e != nil || er != nil {
			log.Println(e)
			log.Println(er)
		}

		sess := session.Clone()
		Mongo := sess.DB(db)
		container := []offlineChart{}
		var headercontainer2 []offlineChart
		// headercontainer3 := []offlineChart{}

		defer sess.Close()

		Mongo.C(c_offlineChart).Find(bson.M{"Timestamp": bson.M{"$gte": start, "$lte": SetTimeStampForLastDay(stop)}}).All(&container)
		for _, each := range container {
			each.Timestamp = each.Timestamp.Add(time.Hour * 8)
			headercontainer2 = append(headercontainer2, each)
		}
		// fmt.Println(container)
		json.NewEncoder(w).Encode(headercontainer2)
		log.Println("gogetgofflinechart")

	}
}

func gogetgwdetail(w http.ResponseWriter, r *http.Request) {

	container := []gwstat{}
	container2 := []gwdata{}
	container3 := []string{}
	container4 := []gwdata{}
	container5 := []gwdata{}
	sess := session.Clone()
	defer sess.Close()

	Mongo := sess.DB(db).C(c_gwtstat)
	Mongo.Find(bson.M{}).All(&container)
	sess.DB(db).C(c_devices).Find(bson.M{}).Distinct("GWID", &container2)

	for _, each := range container2 {
		// fmt.Print(each)
		for _, each2 := range container {

			if each2.GWID[0:7] != each.GWID[0:7] {
				// fmt.Println(each2.GWID)
				// continue
				container3 = append(container3, each.MGWID)
			}
		}

	}

	for _, each3 := range unique(container3) {
		sess.DB(db).C(c_devices).Find(bson.M{"M_GWID": each3}).Limit(1).All(&container4)
		for _, each := range container4 {
			container5 = append(container5, each)
		}

	}
	// fmt.Println(container5)
	json.NewEncoder(w).Encode(container5)
}

type getAgg struct {
	Timestamp  time.Time `json:"Timestamp" bson:"Timestamp"`
	MACAddress string    `json:"MAC_Address" bson:"MAC_Address"`
	GWID       string    `json:"GW_ID" bson:"GW_ID"`
	// Metrics    Metrics   `json:"Metrics" bson:"Metrics"`
	GET11 float64 `json:"pf_avg" bson:"pf_avg"`
	GET12 float64 `json:"ae_tot" bson:"ae_tot"`
	GET13 float64 `json:"p_sum" bson:"p_sum"`
}

type postAgg struct {
	Start      *string `json:"Start" bson:"Start"`
	Stop       *string `json:"Stop" bson:"Stop"`
	MACAddress *string `json:"MAC_Address" bson:"MAC_Address"`
}

func gopostqueryHourly(w http.ResponseWriter, r *http.Request) {

	headercontainer := postAgg{}

	json.NewDecoder(r.Body).Decode(&headercontainer)

	if *headercontainer.Start != "" && *headercontainer.Stop != "" {

		start, e := time.ParseInLocation("2006-01-02T15", *headercontainer.Start, time.Local)
		stop, er := time.ParseInLocation("2006-01-02T15", *headercontainer.Stop, time.Local)
		diff := stop.Sub(start)

		if e != nil || er != nil {
			log.Println(e)
			log.Println(er)
		}
		if diff <= time.Hour*24 {

			container := []getAgg{}
			sess := session.Clone()
			defer sess.Close()

			Mongo := sess.DB(db).C(c_hourly)
			Mongo.Find(bson.M{"MAC_Address": *headercontainer.MACAddress, "Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
			json.NewEncoder(w).Encode(container)
			fmt.Println(diff)
		} else {
			WARNING := "Time more than 24 hours"
			json.NewEncoder(w).Encode(WARNING)
		}
	}
}

func gopostqueryDaily(w http.ResponseWriter, r *http.Request) {

	headercontainer := postAgg{}

	json.NewDecoder(r.Body).Decode(&headercontainer)

	if *headercontainer.Start != "" && *headercontainer.Stop != "" {

		start, e := time.ParseInLocation("2006-01-02", *headercontainer.Start, time.Local)
		stop, er := time.ParseInLocation("2006-01-02", *headercontainer.Stop, time.Local)
		diff := stop.Sub(start)

		if e != nil || er != nil {
			log.Println(e)
			log.Println(er)
		}
		if diff <= time.Hour*21*31 {

			container := []getAgg{}
			sess := session.Clone()
			defer sess.Close()

			Mongo := sess.DB(db).C(c_daily)
			Mongo.Find(bson.M{"MAC_Address": *headercontainer.MACAddress, "Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
			json.NewEncoder(w).Encode(container)
			fmt.Println(diff)
		} else {
			WARNING := "Time more than 31 days"
			json.NewEncoder(w).Encode(WARNING)
		}
	}
}

func gopostqueryMonthly(w http.ResponseWriter, r *http.Request) {

	headercontainer := postAgg{}

	json.NewDecoder(r.Body).Decode(&headercontainer)

	if *headercontainer.Start != "" && *headercontainer.Stop != "" {

		start, e := time.ParseInLocation("2006-01", *headercontainer.Start, time.Local)
		stop, er := time.ParseInLocation("2006-01", *headercontainer.Stop, time.Local)
		diff := stop.Sub(start)

		if e != nil || er != nil {
			log.Println(e)
			log.Println(er)
		}
		if diff <= time.Hour*24*30*36 {

			container := []getAgg{}
			sess := session.Clone()
			defer sess.Close()

			Mongo := sess.DB(db).C(c_month)
			Mongo.Find(bson.M{"MAC_Address": *headercontainer.MACAddress, "Timestamp": bson.M{"$gte": start, "$lte": stop}}).All(&container)
			json.NewEncoder(w).Encode(container)
			fmt.Println(diff)
		} else {
			WARNING := "Time more than 36 months"
			json.NewEncoder(w).Encode(WARNING)
		}
	}
}

// `````
// SMART SPACE
// `````

func deviceState(w http.ResponseWriter, req *http.Request) {

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

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=turn" + state

	response, err := http.Get(url)
	// fmt.Println(response)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)

	}
}

func camTurn(w http.ResponseWriter, req *http.Request) {
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

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
	}

}

func deviceACState(w http.ResponseWriter, req *http.Request) {
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

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=" + state + "&arg1=" + argbv

	response, err := http.Get(url)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
	}

}

func thedeviceStatusRes(w http.ResponseWriter, req *http.Request) {
	log.Println("deviceStatusRes")
	vars := mux.Vars(req)
	deviceID := vars["id"]

	if deviceID == "" || false {
		w.WriteHeader(http.StatusBadRequest)
	}

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/devices/" + deviceID

	if deviceID == "130" {

		response, err := http.Get(url)

		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)

		} else {
			defer response.Body.Close()
			var tmprecord airCond

			json.NewDecoder(response.Body).Decode(&tmprecord)
			tmpVal := checkBool(tmprecord.Properties.UIACValue)

			json.NewEncoder(w).Encode(checkState{State: tmpVal})
			return
		}

	}

	response, err := http.Get(url)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	} else {
		defer response.Body.Close()
		var record deviceStatusRes

		json.NewDecoder(response.Body).Decode(&record)
		recordVal, _ := strconv.ParseBool(record.Properties.Value)
		json.NewEncoder(w).Encode(checkState{State: recordVal})
		// fmt.Println(req.RemoteAddr)
		return
	}

}

// ```
// AUTHORIZATION
// ```

var SigningKey = []byte("AddDevices")

type jwtSignation struct {
	jwt.StandardClaims
	Platform string `json:"platform,omitempty"`
	Pass     string `json:"pass,omitempty"`
}

type addDeviceList struct {
	Key               *string   `json:"APIKey" bson:"APIKey"`
	Timestamp_changes time.Time `json:"Timestamp_Changes" bson:"Timestamp_Changes"`
	Userchanges       *string   `json:"user" bson:"user"`
	DevID             *int      `json:"DevID" bson:"DevID"`
	Floor             *string   `json:"Floor" bson:"Floor"`
	GWID              *string   `json:"GWID" bson:"GWID"`
	MGWID             *string   `json:"M_GWID" bson:"M_GWID"`
	MMAC              *string   `json:"M_MAC" bson:"M_MAC"`
	NUM               *string   `json:"NUM" bson:"NUM"`
	Place             *string   `json:"Place" bson:"Place"`
	Territory         *string   `json:"Territory" bson:"Territory"`
	Type              *string   `json:"Type" bson:"Type"`
	MACAddress        *string   `json:"MAC_Address" bson:"MAC_Address"`
}

func genAuth(w http.ResponseWriter, req *http.Request) {
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

// type headercontainer struct {
// 	key *string `json:"key"`
// }

func validator(w http.ResponseWriter, r *http.Request) {
	log.Println("Validator")
	// vars := mux.Vars(req)
	// keyz := headercontainer{}
	container := addDeviceList{}
	json.NewDecoder(r.Body).Decode(&container)
	// keyzz := headercontainer

	mongo := session.Clone()
	container.Timestamp_changes = time.Now()
	fmt.Println(container)
	year, _, day := time.Now().Local().Date()
	yearDay := time.Now().Local().YearDay()
	// jwtcontainer := jwtSignation{}
	if container.Key != nil {
		token, err := jwt.ParseWithClaims(*container.Key, &jwtSignation{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return SigningKey, nil
		})

		if ve, lll := err.(*jwt.ValidationError); lll {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// http.Error(w, "ValidationError", http.StatusNotAcceptable)
				// json.NewEncoder(w).Encode("ValidationErrorMalformed")
				fmt.Println("ValidationErrorMalformed")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// http.Error(w, "ValidationErrorExpired", http.StatusNotAcceptable)
				// json.NewEncoder(w).Encode("ValidationErrorExpired")
				fmt.Println("ValidationErrorExpired")
			}
		} else {
			// if claims, ok := token.Claims.(*jwtSignation); ok && token.Valid {

			if claims, ok := token.Claims.(*jwtSignation); ok && token.Valid {

				claimsPointA := claims.StandardClaims.VerifyIssuer("APIUser", ok)
				claimsPointB := claims.StandardClaims.VerifyAudience(strconv.Itoa(year)+strconv.Itoa(yearDay)+strconv.Itoa(day), false)

				if claimsPointA != true {
					http.Error(w, "ValidationErrorIssuer", http.StatusNotAcceptable)
					fmt.Println("ValidationErrorIssuer")
				} else if claimsPointB != true {
					http.Error(w, "ValidationErrorExpired", http.StatusNotAcceptable)
					fmt.Println("ValidationErrorExpired")
				} else {
					json.NewDecoder(r.Body).Decode(&container)
					fmt.Println(strconv.Itoa(year)+strconv.Itoa(yearDay)+strconv.Itoa(day), claimsPointB)

					if container.DevID == nil {
						http.Error(w, "BuildingDetails null", http.StatusBadRequest)
						fmt.Println("BuildingDetails null")
					} else if container.Floor == nil {
						http.Error(w, "BuildingName null", http.StatusBadRequest)
						fmt.Println("BuildingName null")
					} else if container.GWID == nil {
						http.Error(w, "DeviceBrand null", http.StatusBadRequest)
						fmt.Println("DeviceBrand null")
					} else if container.MACAddress == nil {
						http.Error(w, "DeviceDetails null", http.StatusBadRequest)
						fmt.Println("DeviceDetails null")
					} else if container.MGWID == nil {
						http.Error(w, "devID null", http.StatusBadRequest)
						fmt.Println("devID null")
					} else if container.MMAC == nil {
						http.Error(w, "DeviceInfo null", http.StatusBadRequest)
						fmt.Println("DeviceInfo null")
					} else if container.NUM == nil {
						http.Error(w, "DeviceName null", http.StatusBadRequest)
						fmt.Println("DeviceName null")
					} else if container.Place == nil {
						http.Error(w, "DeviceType null", http.StatusBadRequest)
						fmt.Println("DeviceType null")
					} else if container.Territory == nil {
						http.Error(w, "Floor null", http.StatusBadRequest)
						fmt.Println("Floor null")
					} else if container.Type == nil {
						http.Error(w, "GWID null", http.StatusBadRequest)
						fmt.Println("GWID null")
					} else {
						fmt.Println("BDetails :", container.MACAddress)

						erro := mongo.DB(db).C(c_devices).Update(bson.M{"MAC_Address": container.MACAddress}, bson.M{"$set": container})
						w.WriteHeader(http.StatusOK)
						fmt.Println("Device Added : ", container.MACAddress)
						json.NewEncoder(w).Encode(container)

						if erro != nil {

							w.WriteHeader(http.StatusBadRequest)
							log.Println(container.MACAddress, container, err)
						}
					}

				}

			}
		}
	}
	// http.Error(w, "no data", http.StatusBadRequest)
}

func AddDev(w http.ResponseWriter, req *http.Request) {

	log.Println("AddDev")

	// tokenString := req.Header.Get("authorization")
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("Unexpected sigining method: %v", token.Header["alg"])
	// 	}

	// 	return _CLIENTSEC, nil
	// })
	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	container := addDeviceList{}

	json.NewDecoder(req.Body).Decode(&container)

	mongo := session.Clone()
	err := session.DB(db).C(c_devices).Update(bson.M{"MAC_Address": container.MACAddress}, bson.M{"$set": container})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(container.MACAddress, container, err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(container)

	mongo.Close()
	// 	fmt.Println(req.RemoteAddr, claims)
	// } else {
	// 	fmt.Println(err)
	// }
}

func meterDetail(w http.ResponseWriter, r *http.Request) {

	zeros := 0
	aspect := bson.M{}

	sess := session.Clone()
	vars := mux.Vars(r)
	deviceID := vars["macID"]
	fmt.Println(deviceID)

	shit := "Please fill the MAC ID"
	container := device{}

	mongo := sess.DB(db).C(c_devices)

	if deviceID == "" {
		json.NewEncoder(w).Encode(shit)

	}
	err := mongo.Find(bson.M{"MAC_Address": deviceID}).All(&container.Rows)
	// err := mongo.Find(bson.M{"GWID": bson.M{"$in": "/^meter_03/"}}).All(&container.Rows)
	fmt.Println(err)

	container.Datashape.FieldDefinitions.DevID.Ordinal = zeros
	container.Datashape.FieldDefinitions.DevID.BaseType = `NUMBER`
	container.Datashape.FieldDefinitions.DevID.Aspects = aspect
	container.Datashape.FieldDefinitions.DevID.Name = `DevID`
	container.Datashape.FieldDefinitions.DevID.Description = `DevID`

	container.Datashape.FieldDefinitions.Floor.Ordinal = zeros
	container.Datashape.FieldDefinitions.Floor.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Floor.Aspects = aspect
	container.Datashape.FieldDefinitions.Floor.Name = `Floor`
	container.Datashape.FieldDefinitions.Floor.Description = `Floor`

	container.Datashape.FieldDefinitions.GWID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GWID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.GWID.Aspects = aspect
	container.Datashape.FieldDefinitions.GWID.Name = `GWID`
	container.Datashape.FieldDefinitions.GWID.Ordinal = zeros
	container.Datashape.FieldDefinitions.GWID.Description = `GWID`

	container.Datashape.FieldDefinitions.MGWID.Ordinal = zeros
	container.Datashape.FieldDefinitions.MGWID.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MGWID.Aspects = aspect
	container.Datashape.FieldDefinitions.MGWID.Name = `M_GWID`
	container.Datashape.FieldDefinitions.MGWID.Description = `M_GWID`

	container.Datashape.FieldDefinitions.MMAC.Ordinal = zeros
	container.Datashape.FieldDefinitions.MMAC.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MMAC.Aspects = aspect
	container.Datashape.FieldDefinitions.MMAC.Name = `M_MAC`
	container.Datashape.FieldDefinitions.MMAC.Description = `M_MAC`

	container.Datashape.FieldDefinitions.NUM.Ordinal = zeros
	container.Datashape.FieldDefinitions.NUM.BaseType = `STRING`
	container.Datashape.FieldDefinitions.NUM.Aspects = aspect
	container.Datashape.FieldDefinitions.NUM.Name = `NUM`
	container.Datashape.FieldDefinitions.NUM.Description = `NUM`

	container.Datashape.FieldDefinitions.Place.Ordinal = zeros
	container.Datashape.FieldDefinitions.Place.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Place.Aspects = aspect
	container.Datashape.FieldDefinitions.Place.Name = `Place`
	container.Datashape.FieldDefinitions.Place.Description = `Place`

	container.Datashape.FieldDefinitions.Territory.Ordinal = zeros
	container.Datashape.FieldDefinitions.Territory.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Territory.Aspects = aspect
	container.Datashape.FieldDefinitions.Territory.Name = `Territory`
	container.Datashape.FieldDefinitions.Territory.Description = `Territory`

	container.Datashape.FieldDefinitions.Type.Ordinal = zeros
	container.Datashape.FieldDefinitions.Type.BaseType = `STRING`
	container.Datashape.FieldDefinitions.Type.Aspects = aspect
	container.Datashape.FieldDefinitions.Type.Name = `Type`
	container.Datashape.FieldDefinitions.Type.Description = `Type`

	container.Datashape.FieldDefinitions.MACAddress.Ordinal = zeros
	container.Datashape.FieldDefinitions.MACAddress.BaseType = `STRING`
	container.Datashape.FieldDefinitions.MACAddress.Aspects = aspect
	container.Datashape.FieldDefinitions.MACAddress.Name = `MAC_Address`
	container.Datashape.FieldDefinitions.MACAddress.Description = `MAC_Address`

	json.NewEncoder(w).Encode(container)

}

func test(w http.ResponseWriter, r *http.Request) {

	// zeros := 0
	// aspect := bson.M{}

	sess := session.Clone()
	// vars := mux.Vars(r)
	// deviceID := vars["macID"]
	// fmt.Println(deviceID)

	// shit := "Please fill the MAC ID"
	// container := device{}
	var amam []interface{}
	mongo := sess.DB(db).C(c_devices)

	// if deviceID == "" {
	// 	json.NewEncoder(w).Encode(shit)

	// }
	// err := mongo.Find(bson.M{"MAC_Address": deviceID}).All(&container.Rows)
	// err := mongo.Find(bson.M{}).All(&amam)
	err := mongo.Find(bson.M{"GWID": bson.M{"$regex": bson.RegEx{`^meter_05`, "i"}}}).All(&amam)
	json.NewEncoder(w).Encode(amam)
	fmt.Println(err)

}
func checkDBStatus() bool {
	err := session.Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		session.Close()
		time.Sleep(5 * time.Second)
		session.Refresh()
	}
	fmt.Println("DB GOOD")
	return true
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if checkDBStatus(); true {
		// A very simple health .
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		sess := session.Clone()
		aa := sess.LiveServers()
		fmt.Println(aa)

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		io.WriteString(w, `{"alive": true}`)
	}
}

// ```
// MAIN
// ```

func main() {

	dbConnect()
	// auth.GenAuth()
	router := mux.NewRouter()

	// HEALTHCHECK
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	// additional API for query all buildings data
	router.HandleFunc("/meter/lastreport/allbuilding", gopostlastreportAllBuilding).Methods("GET")
	router.HandleFunc("/meter/offlinechart", gogetgofflinechart).Methods("POST")
	router.HandleFunc("/meter/meterdetail/{macID}", meterDetail).Methods("GET")

	//Authorization
	router.HandleFunc("/meter/auth/{user}", genAuth).Methods("GET")
	router.HandleFunc("/meter/devices", validator).Methods("POST")

	// All API
	router.HandleFunc("/meter/lastreport", gopostlastreport).Methods("POST")
	router.HandleFunc("/meter/gwstat", gopostgwstat).Methods("POST")

	router.HandleFunc("/meter/hourly", gopostqueryHourly).Methods("POST")
	router.HandleFunc("/meter/daily", gopostqueryDaily).Methods("POST")
	router.HandleFunc("/meter/monthly", gopostqueryMonthly).Methods("POST")

	router.HandleFunc("/meter/lastreport", gogetlastreport).Methods("GET")
	router.HandleFunc("/meter/devices", gogetDevices).Methods("GET")
	router.HandleFunc("/meter/gwstat", gogetgwstat).Methods("GET")

	router.HandleFunc("/space/state/id/{id}/name/{state}", deviceState).Methods("GET")
	router.HandleFunc("/space/state/id/{id}/name/{state}/arg/{argbv}", deviceACState).Methods("GET")
	router.HandleFunc("/space/status/id/{id}", thedeviceStatusRes).Methods("GET")
	router.HandleFunc("/space/cam/posX/{pos1}/posY/{pos2}", camTurn).Methods("GET")

	router.HandleFunc("/test", test).Methods("GET")

	log.Println(http.ListenAndServe(":8081", router))

}
