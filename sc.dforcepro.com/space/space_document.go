package space

import (
	"fmt"
	"net/http"
	"strconv"

	"dforcepro.com/resource"
	"dforcepro.com/resource/db"
	"github.com/gorilla/mux"
)

type Doc interface {
	initApi(router *mux.Router)
}

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

var _di *resource.Di

func SetDI(c *resource.Di) {
	_di = c
}

func getMongo() db.Mongo {
	return _di.Mongodb
}

func GetMongo() db.Mongo {
	return _di.Mongodb
}

func _afterEndPoint(w http.ResponseWriter, req *http.Request) {

}

func checkBool(x string) bool {

	tmpVal, _ := strconv.ParseFloat(x, 32)
	fmt.Println(tmpVal)
	fmt.Println(x)
	if (tmpVal) > 5 {
		return true
	} else {
		return false
	}

}
