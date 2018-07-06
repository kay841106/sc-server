package search

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type HourlyDoc struct {
	ID    bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Hour  uint8         `json:"Hour,omitempty" bson:"hour"`
	Day   uint8         `json:"Day,omitempty" bson:"day"`
	Month uint8         `json:"Month,omitempty" bson:"month"`
	Year  uint32        `json:"Year,omitempty" bson:"year"`

	LastReportTime time.Time `json:"LastReportTime, omitempty" bson:"lastReportTime"`

	Floor           string `json:"Floor,omitempty" bson:"Floor"`
	BuildingName    string `json:"BuildingName,omitempty" bson:"Building_Name"`
	BuildingDetails string `json:"BuildingDetails,omitempty" bson:"Building_Details"`
	DeviceName      string `json:"DeviceName,omitempty" bson:"Device_Name"`
	DeviceDetails   string `json:"DeviceDetails,omitempty" bson:"Building_Details"`
	DeviceID        string `json:"DeviceID,omitempty" bson:"Device_ID"`

	MaxPWRDemand float64 `json:"MaxPWRDemand,omitempty" bson:"max_Pwr_Demand"`
	WeatherTemp  float32 `json:"WeatherTemp,omitempty" bson:"weather_temp"`
	MaxPWRUsage  float64 `json:"MaxPWRUsage,omitempty" bson:"max_Pwr_Usage"`
	MaxPF        float64 `json:"MaxPF,omitempty" bson:"max_PF"`
	TotalUsage   float64 `json:"TotalUsage,omitempty" bson:"total_Usage"`
	CC           float64 `json:"CC,omitempty" bson:"CC"`
	AvgDemand    float64 `json:"AvgDemand,omitempty" bson:"avg_Demand"`
	PFLimit      float64 `json:"PFLimit,omitempty" bson:"PF_Limit"`
	AvgUsage     float64 `json:"AvgUsage,omitempty" bson:"avg_Usage"`
	AvgPF        float64 `json:"AvgPF,omitempty" bson:"avg_PF"`
}

type HourlyCompDoc struct {
	ID    bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Hour  uint8         `json:"Hour,omitempty" bson:"hour"`
	Day   uint8         `json:"Day,omitempty" bson:"day"`
	Month uint8         `json:"Month,omitempty" bson:"month"`
	Year  uint32        `json:"Year,omitempty" bson:"year"`

	LastReportTime time.Time `json:"LastReportTime, omitempty" bson:"lastReportTime"`

	Floor           string `json:"Floor,omitempty" bson:"Floor"`
	BuildingName    string `json:"BuildingName,omitempty" bson:"Building_Name"`
	BuildingDetails string `json:"BuildingDetails,omitempty" bson:"Building_Details"`
	DeviceName      string `json:"DeviceName,omitempty" bson:"Device_Name"`
	DeviceDetails   string `json:"DeviceDetails,omitempty" bson:"Building_Details"`
	DeviceID        string `json:"DeviceID,omitempty" bson:"Device_ID"`

	MaxPWRDemand float64 `json:"MaxPWRDemand,omitempty" bson:"max_Pwr_Demand"`
	WeatherTemp  float32 `json:"WeatherTemp,omitempty" bson:"weather_temp"`
	MaxPWRUsage  float64 `json:"MaxPWRUsage,omitempty" bson:"max_Pwr_Usage"`
	MaxPF        float64 `json:"MaxPF,omitempty" bson:"max_PF"`
	TotalUsage   float64 `json:"TotalUsage,omitempty" bson:"total_Usage"`
	CC           float64 `json:"CC,omitempty" bson:"CC"`
	AvgDemand    float64 `json:"AvgDemand,omitempty" bson:"avg_Demand"`
	PFLimit      float64 `json:"PFLimit,omitempty" bson:"PF_Limit"`
	AvgUsage     float64 `json:"AvgUsage,omitempty" bson:"avg_Usage"`
	AvgPF        float64 `json:"AvgPF,omitempty" bson:"avg_PF"`

	PrevUsage      float64 `json:"PrevUsage,omitempty" bson:"prev_Usage"`
	PrevPF         float64 `json:"PrevPF,omitempty" bson:"prev_PF"`
	PrevDemand     float64 `json:"PrevDemand,omitempty" bson:"prev_Demand"`
	PrevTotalUsage float64 `json:"PrevTotalUsage,omitempty" bson:"prev_total_Usage"`
}

func (sca ESearch) storeHour(w http.ResponseWriter, req *http.Request) {
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.2.10:9201"))
	if err != nil {
		// Handle error
	}
}
