package airbox

import (
	"time"

	"dforcepro.com/resource"
	"sc.dforcepro.com/meter"
)

const (
	//DBName for Airbox
	DBName = "sc"
	//AirboxRaw for Airbox
	AirboxRaw = "airbox"
	//AirboxStream for Airbox
	AirboxStream = "Stream"
	//AirboxHour for Airbox
	AirboxHour = "hour"
	//AirboxDay for Airbox
	AirboxDay = "day"
	//AirboxMonth for Airbox
	AirboxMonth = "month"
	//AirboxMapping for Airbox
	AirboxMapping = "Airbox_device_mapping"
)

//CroncheckStatus Struct for checkStatus func
type CroncheckStatus struct {
	PM25           int       `json:"PM2_5" bson:"PM2_5"`
	LastReportTime time.Time `json:"Upload_Time" bson:"Upload_Time"`
	CO2            int       `json:"CO2" bson:"CO2"`
	CO             int       `json:"CO" bson:"CO"`
	Noise          int       `json:"Noise" bson:"Noise"`
	Temp           float64   `json:"Temp" bson:"Temp"`
	Humidity       int       `json:"Humidity" bson:"Humidity"`
	DevID          string    `json:"Device_ID" bson:"Device_ID"`
	Location       string    `json:"Location" bson:"Location"`
}

type AirboxPayload struct {
	DeviceID       string    `json:"Device_ID" bson:"Device_ID"`
	PM25           int       `json:"PM2_5" bson:"PM2_5"`
	LastReportTime time.Time `json:"Upload_Time" bson:"Upload_Time"`
	CO2            int       `json:"CO2" bson:"CO2"`
	CO             int       `json:"CO" bson:"CO"`
	Noise          int       `json:"Noise" bson:"Noise"`
	Temp           float64   `json:"Temp" bson:"Temp"`
	Humidity       int       `json:"Humidity" bson:"Humidity"`
	DevID          string    `json:"Device_ID" bson:"Device_ID"`
	Location       string    `json:"Location" bson:"Location"`
}

//MappingTemporaryData for Airbox
type MappingTemporaryData struct {
	DevID    string `json:"DevID" bson:"DevID"`
	Location string `json:"Location" bson:"Location"`
}

//MappingData for Airbox
type MappingData struct {
	Rows []struct {
		DevID    string `json:"DevID" bson:"DevID"`
		Location string `json:"Location" bson:"Location"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			Location meter.DSTwxTemplate `json:"Location"`
			DevID    meter.DSTwxTemplate `json:"DevID"`
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

//ScheckStatus Struct for checkStatus func
type ScheckStatus struct {
	Rows []struct {
		PM25           int       `json:"PM2_5" bson:"PM2_5"`
		LastReportTime time.Time `json:"Upload_Time" bson:"Upload_Time"`
		CO2            int       `json:"CO2" bson:"CO2"`
		CO             int       `json:"CO" bson:"CO"`
		Noise          int       `json:"Noise" bson:"Noise"`
		Temp           float64   `json:"Temp" bson:"Temp"`
		Humidity       int       `json:"Humidity" bson:"Humidity"`
		DevID          string    `json:"Device_ID" bson:"Device_ID"`
		Location       string    `json:"Location" bson:"Location"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			PM25           meter.DSTwxTemplate `json:"PM25"`
			DevID          meter.DSTwxTemplate `json:"Device_ID"`
			LastReportTime meter.DSTwxTemplate `json:"Upload_Time"`
			CO2            meter.DSTwxTemplate `json:"CO2" `
			CO             meter.DSTwxTemplate `json:"CO" `
			Noise          meter.DSTwxTemplate `json:"Noise" `
			Temp           meter.DSTwxTemplate `json:"Temp" `
			Humidity       meter.DSTwxTemplate `json:"Humidity" `
			Location       meter.DSTwxTemplate `json:"Location" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

// Airbox Test Payload 06/08/2018

// var _di *resource.Di

// func getMongo() db.Mongo {
// 	return _di.Mongodb
// }

var _di *resource.Di

func SetDI(c *resource.Di) {
	_di = c
}
