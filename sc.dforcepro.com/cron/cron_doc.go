package cron

import (
	"log"
	"time"

	"sc.dforcepro.com/meter"
)

const (
	RawDataDispenserAPIUrl = "http://3egreenserver.ddns.net:9000/group/currentstatus/NTUST/NTUST_Dispenser"
	//DBDispenser is database name
	DBDispenser      = "dispenser"
	RawDataDispenser = "Rawdata"
	//DBAirbox is database name
	DBAirbox      = "sc"
	MappingAirbox = "Airbox_device_mapping"
	cAirbox       = "airbox"
	AirboxStatus  = "AirboxStatus"

	AirboxHourCollection  = "hour"
	AirboxDayCollection   = "day"
	AirboxMonthCollection = "month"

	DispenserHourCollection  = "Dispenser_hour"
	DispenserDayCollection   = "Dispenser_day"
	DispenserMonthCollection = "Dispenser_month"
)

func checkDBStatus() bool {
	err := meter.GetMongo().Ping()
	for err != nil {
		log.Println("Connection to DB is down, restarting ....")
		meter.GetMongo().Close()
		time.Sleep(5 * time.Second)
		meter.GetMongo().Refresh()
	}
	return true
}

func timeConvertToMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}
func timeConvertToDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
func timeConvertToHour(t time.Time) time.Time {
	year, month, day := t.Date()
	hour := t.Hour()
	return time.Date(year, month, day, hour, 0, 0, 0, time.Local)
}

//METER

type detailStatusStruct struct {
	GWID     string `json:"GWID" bson:"GWID"`
	GWStatus bool   `json:"GWStatus" bson:"GWStatus"`
	Device   []struct {
		ID             string    `json:"ID" bson:"ID"`
		Status         bool      `json:"status" bson:"status"`
		LastReportTime time.Time `json:"lastReportTime" bson:"lastReportTime"`
		Metric         struct {
			Usage  float64 `json:"usage" bson:"usage"`
			PF     float64 `json:"PF" bson:"PF"`
			Demand float64 `json:"demand" bson:"demand"`
		} `json:"Metric"`
		DownTime struct {
			HoursD   time.Duration `json:"HoursD"`
			MinutesD time.Duration `json:"MinutesD"`
			SecondD  time.Duration `json:"SecondD"`
		} `json:"DownTime"`
	} `json:"Device"`
}

type displayDataMongo struct {
	DeviceID       string    `json:"Device_ID" bson:"Device_ID"`
	PwrDemand      int       `json:"Pwr_Demand" bson:"Pwr_Demand"`
	LastReportTime time.Time `json:"lastReportTime" bson:"lastReportTime"`
	WeatherTemp    int       `json:"weather_Temp" bson:"weather_Temp"`
	PF             float64   `json:"PF" bson:"PF"`
	PwrUsage       float64   `json:"Pwr_Usage" bson:"Pwr_Usage"`
	DeviceType     string    `json:"Device_Type" bson:"Device_Type"`
	BuildingName   string    `json:"Building_Name" bson:"Building_Name"`
	GatewayID      string    `json:"Gateway_ID" bson:"Gateway_ID"`
}

type displayDataCalcMongo struct {
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	QuarterPost     int       `json:"QuarterPost" bson:"QuarterPost"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	CC              float64   `json:"CC" bson:"CC"`
	Floor           string    `json:"Floor" bson:"Floor"`
	PwrUsage        float64   `json:"Pwr_Usage" bson:"Pwr_Usage"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	PwrDemand       int       `json:"Pwr_Demand" bson:"Pwr_Demand"`
	PF              float64   `json:"PF" bson:"PF"`
	DeviceDetails   string    `json:"Device_Details" bson:"Device_Details"`
	Usage           float64   `json:"Usage" bson:"Usage"`
}
type devMan struct {
	GatewayID       string `json:"GWID" bson:"GWID"`
	DeviceBrand     string `json:"Device_Brand" bson:"Device_Brand"`
	DeviceID        string `json:"devID" bson:"devID"`
	DeviceDetails   string `json:"Device_Details" bson:"Device_Details"`
	DeviceName      string `json:"Device_Name" bson:"Device_Name"`
	DeviceInfo      string `json:"Device_Info" bson:"Device_Info"`
	DeviceType      string `json:"Device_Type" bson:"Device_Type"`
	Floor           string `json:"Floor" bson:"Floor"`
	BuildingName    string `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string `json:"Building_Details" bson:"Building_Details"`
}

type mappingGetBuildingAndDevice struct {
	DeviceID     string `json:"devID" bson:"devID"`
	BuildingName string `json:"Building_Name" bson:"Building_Name"`
}
type aggHourStruct struct {
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	CC              float64   `json:"CC" bson:"CC"`
	Usage           float64   `json:"total_Usage" bson:"total_Usage"` // KWh
	Floor           string    `json:"Floor" bson:"Floor"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
	AvgPF           float64   `json:"avg_PF" bson:"avg_PF"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"` //KWh
	PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
	MinPF           float64   `json:"min_PF" bson:"min_PF"`
	MaxPF           float64   `json:"max_PF" bson:"max_PF"`
	MinDemand       float64   `json:"min_Demand" bson:"min_Demand"` //W
	MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"` //W
	MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`   //KWh
}
type displayDataElement2nd struct {
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	PwrUsage        float64   `json:"Usage" bson:"Usage"`
	PwrDemand       float64   `json:"Pwr_Demand" bson:"Pwr_Demand"`
	AvgPF           float64   `json:"PF" bson:"PF"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
}

//Dispenser
type CDispenssssserGET struct {
	Address      string    `json:"address"`
	DemandCount  string    `json:"demandCount"`
	Status       string    `json:"status"`
	Current      string    `json:"current"`
	Watts        string    `json:"watts"`
	WattsAverage string    `json:"watts_Average"`
	Temp         string    `json:"temp"`
	TimeStamp    time.Time `json:"timeStamp"`
}

type CDispenserGET struct {
	Address     string `json:"address"`
	DemandCount int    `json:"demandCount"`
	// Status       int       `json:"status"`
	Current int     `json:"current"`
	Watts   float64 `json:"watts"`
	// WattsAverage float64   `json:"watts_Average"`
	Temp       int       `json:"temp"`
	TimeStamp  time.Time `json:"timeStamp"`
	DeviceName string    `json:"devicenickname" bson:"devicenickname"`
	Hour       int       `json:"Hour" bson:"Hour"`
}
type CDispenserUnixTime struct {
	Address      string    `json:"address"`
	DemandCount  int       `json:"demandCount"`
	Status       int       `json:"status"`
	Current      int       `json:"current"`
	Watts        float64   `json:"watts"`
	WattsAverage float64   `json:"watts_Average"`
	Temp         int       `json:"temp"`
	TimeStamp    time.Time `json:"timeStamp"`
	Lastupdated  int64     `json:"lastupdated"`
	DeviceName   string    `json:"devicenickname" bson:"devicenickname"`
	Hour         int       `json:"Hour" bson:"Hour"`
}

//Airbox
type CAirbox struct {
	PM25           int       `json:"PM2_5" bson:"PM2_5"`
	LastReportTime time.Time `json:"Upload_Time" bson:"Upload_Time"`
	CO2            int       `json:"CO2" bson:"CO2"`
	CO             int       `json:"CO" bson:"CO"`
	Noise          int       `json:"Noise" bson:"Noise"`
	Temp           float64   `json:"Temp" bson:"Temp"`
	Humidity       int       `json:"Humidity" bson:"Humidity"`
	DevID          string    `json:"Device_ID" bson:"Device_ID"`
}

type AirboxDeviceMapping struct {
	DevID            string    `json:"DevID" bson:"DevID"`
	Location         string    `json:"Location" bson:"Location"`
	InstallationDate time.Time `json:"InstallationDate" bson:"InstallationDate"`
}

type airboxHttpReq struct {
	Timestamp     string  `json:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MACAddress    string  `json:"MAC_Address"`
	GWID          string  `json:"GW_ID"`
	CPURate       float64 `json:"CPU_rate"`
	StorageRate   int     `json:"Storage_rate"`
	GET11         float64 `json:"GET_1_1"`
	GET12         float64 `json:"GET_1_2"`
}

//Dispenser POST
type CDispenserPOST struct {
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

type tempResponseDispenser struct {
	Result bool                     `json:"result"`
	Data   []map[string]interface{} `json:"data"`
	// Data []struct {
	// 	Address            string
	// 	Voltages           string  `json:"voltages,omitempty"`
	// 	Hz                 int     `json:"hz,omitempty"`
	// 	DemandCount        int     `json:"demandCount,omitempty"`
	// 	DemandWatts        int     `json:"demandWatts,omitempty"`
	// 	Current            int     `json:"current,omitempty"`
	// 	Battery            int     `json:"battery,omitempty"`
	// 	Devicename         string  `json:"devicename,omitempty"`
	// 	Temp               int     `json:"temp,omitempty"`
	// 	CurrentAc          string  `json:"current_ac,omitempty"`
	// 	Watts              float64 `json:"watts,omitempty"`
	// 	Username           string  `json:"username,omitempty"`
	// 	AlertCurrent       string  `json:"alertCurrent,omitempty"`
	// 	LastSendAlert      int     `json:"lastSendAlert,omitempty"`
	// 	Calibration        string  `json:"calibration,omitempty"`
	// 	Lastupdated        int64   `json:"lastupdated,omitempty"`
	// 	Number             int     `json:"number,omitempty"`
	// 	Status             int     `json:"status" bson:"status"`
	// 	Starttime          int64   `json:"starttime,omitempty"`
	// 	Stoptime           int64   `json:"stoptime,omitempty"`
	// 	StopCurrent        int     `json:"stopCurrent,omitempty"`
	// 	MaxCurrentRaw      int     `json:"maxCurrent_raw,omitempty"`
	// 	CurrentByte        int     `json:"current_byte,omitempty"`
	// 	VariantCurrent     int     `json:"variantCurrent,omitempty"`
	// 	SamplingLeng       int     `json:"samplingLeng,omitempty"`
	// 	FilterPer          int     `json:"filterPer,omitempty"`
	// 	Type               string  `json:"type,omitempty"`
	// 	CurrentHealthValue int     `json:"current_healthValue,omitempty"`
	// 	DeviceNickname     string  `json:"DeviceNickname,omitempty"`
	// }
	// `json:"data" bson:"data"`
}
