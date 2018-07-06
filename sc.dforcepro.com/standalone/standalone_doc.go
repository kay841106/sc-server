package standalone

import "time"

const (
	meterMapping            = "SC01_DeviceManager"
	rawDataAEMDRACollection = "SC01_RawData_AEMDRA_"
	rawDataCPMCollection    = "SC01_RawData_CPM72_"
	aemdra                  = "AEMDRA"
	cpm                     = "CPM"

	displayDataCalcCollection = "SC01_displayData_Calc_"
	displayDataCollection     = "SC01_displayData_"

	hourCollection  = "SC01_hour_All"
	dayCollection   = "SC01_day_All"
	monthCollection = "SC01_month_All"
	_EMPTYDEST      = "DESTINATION IS EMPTY"
	_AGG            = "AGGREGATION ->"
	_NEC            = "NON_EMPTY_COLL"
	_EOF            = "END_OF_FILE"
	//DBDispenser is database name
	DBDispenser      = "dispenser"
	RawDataDispenser = "Rawdata"
	//DBAirbox is database name
	DBAirbox              = "sc"
	MappingAirbox         = "Airbox_device_mapping"
	cAirbox               = "airbox"
	AirboxHourCollection  = "hour"
	AirboxDayCollection   = "day"
	AirboxMonthCollection = "month"

	DispenserHourCollection  = "Dispenser_hour"
	DispenserDayCollection   = "Dispenser_day"
	DispenserMonthCollection = "Dispenser_month"

	RawDataDispenserAPIUrl = "http://3egreenserver.ddns.net:9000/group/currentstatus/NTUST/NTUST_Dispenser"
)

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

type rawdataToDisplayData struct {
	PSum           float64   `json:"p_sum"`
	PfAvg          float64   `json:"pf_avg"`
	LastReportTime time.Time `json:"lastReportTime"`
	BlockID        string    `json:"blockId"`
	AeTot          float64   `json:"ae_tot"`
	DevID          string    `json:"devID"`
	GWID           string    `json:"GWID"`
}

//TESTING AIRBOX API POST TO OTHER API

type airboxTestingStruct struct {
	PM25 float64 `json:"PM2_5"`
}
