package dispenser

const (
	//DispenserDB for database name
	DispenserDB                = "dispenser"
	dispenserRawdataCollection = "Rawdata"
)

type dispenserData struct {
	Status         int     `json:"status"`
	Temp           int     `json:"temp"`
	Watts          float64 `json:"watts"`
	Current        int     `json:"current"`
	Devicenickname string  `json:"devicenickname"`
	Lastupdated    int64   `json:"lastupdated"`
	Address        string  `json:"address"`
}

type leostruct struct {
	Status         int     `json:"status"`
	Watts          float64 `json:"watts"`
	Devicenickname string  `json:"devicenickname"`
	Lastupdated    int64   `json:"lastupdated"`
	Address        string  `json:"address"`
}

type AllStatus struct {
	Status         int     `json:"status"`
	Temp           int     `json:"temp"`
	Watts          float64 `json:"watts"`
	Devicenickname string  `json:"devicenickname"`
	Lastupdated    int64   `json:"lastupdated"`
	Address        string  `json:"address"`
}

type Nina struct {
	UploadTime       string  `json:"UploadTime" bson:"UploadTime"`
	Watts            float32 `json:"Watts" bson:"Watts"`
	DeviceMacAddress string  `json:"DeviceMacAddress" bson:"DeviceMacAddress"`
	DeviceNickname   string  `json:"DeviceNickname" bson:"DeviceNickname"`
	Class            int     `json:"Class" bson:"Class"`
	TimeStamp        int     `json:"Timestamp" bson:"Timestamp"`
}

type dispenserboard struct {
	UploadTime  string `json:"UploadTime" bson:"UploadTime"`
	Device      string `json:"device" bson:"device"`
	Hottemp     string `json:"hotTemp" bson:"hotTemp"`
	Warmtemp    string `json:"warmTemp" bson:"warmTemp"`
	Coldtemp    string `json:"coldTemp" bson:"coldTemp"`
	TDS         string `json:"tds" bson:"tds"`
	Heating     int    `json:"heating" bson:"heating"`
	Cooling     int    `json:"cooling" bson:"cooling"`
	SavingPower int    `json:"savingpower" bson:"savingpower"`
	Sterilizing int    `json:"sterilizing" bson:"sterilizing"`
	Inputwater  int    `json:"inputwater" bson:"inputwater"`
	Waterlevel  int    `json:"waterlevel" bson:"waterlevel"`
	Hotoutput   int    `json:"hotoutput" bson:"hotoutput"`
	Warmoutput  int    `json:"warmoutput" bson:"warmoutput"`
	Coldeoutput int    `json:"coldoutput" bson:"coldoutput"`
}
