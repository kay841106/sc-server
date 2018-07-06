package alert

import "time"

type AlertValue struct {
	DeviceID string `json:"Device_ID" bson:"Device_ID"`

	OnlineStatus   bool      `json:"Online_Status"`
	LastReportTime time.Time `json:"lastReportTime"`
	DeviceDetails  DeviceDetails
	DownTime       struct {
		HoursD   time.Duration `json:"HoursD"`
		MinutesD time.Duration `json:"MinutesD"`
		SecondD  time.Duration `json:"SecondD"`
	} `json:"DownTime"`
}
type AlertValueThingworx struct {
	Rows []struct {
		DeviceID string `json:"Device_ID" bson:"Device_ID"`

		OnlineStatus   bool      `json:"Online_Status"`
		LastReportTime time.Time `json:"lastReportTime"`
		DeviceDetails  DeviceDetails
		DownTime       struct {
			HoursD   time.Duration `json:"HoursD"`
			MinutesD time.Duration `json:"MinutesD"`
			SecondD  time.Duration `json:"SecondD"`
		} `json:"DownTime"`
	} `json:"rows"`
}

type DeviceDetails struct {
	BuildingName    string `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string `json:"Building_Details" bson:"Building_Details"`
	GatewayID       string `json:"Gateway_ID" bson:"Gateway_ID"`
}

// type DeviceDown struct {
// 	Month          int8
// 	Times          int16
// 	DownTime       time.Duration `json:"DownTime"`
// 	LastReportTime time.Time     `json:"lastReportTime"`
// }

type AirboxAlertValue struct {
	DeviceID string `json:"Device_ID" bson:"Device_ID"`

	OnlineStatus bool      `json:"Online_Status"`
	UploadTime   time.Time `json:"Upload_Time"`
	Location     string    `json:"Location"`
	DownTime     struct {
		HoursD   time.Duration `json:"HoursD"`
		MinutesD time.Duration `json:"MinutesD"`
		SecondD  time.Duration `json:"SecondD"`
	} `json:"DownTime"`
}
type AirboxAlertValueThingworx struct {
	Rows []struct {
		DeviceID string `json:"Device_ID" bson:"Device_ID"`

		OnlineStatus bool      `json:"Online_Status"`
		UploadTime   time.Time `json:"Upload_Time"`
		Location     string    `json:"Location"`
		DownTime     struct {
			HoursD   time.Duration `json:"HoursD"`
			MinutesD time.Duration `json:"MinutesD"`
			SecondD  time.Duration `json:"SecondD"`
		} `json:"DownTime"`
	} `json:"rows"`
}
