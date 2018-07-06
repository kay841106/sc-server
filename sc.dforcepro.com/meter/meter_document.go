package meter

import (
	"crypto/md5"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"dforcepro.com/resource"
	"dforcepro.com/resource/db"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

const (
	_CLIENTSEC = "Y7WfGYtOHGBjMMigZ6QrcvveYuNDEgepBuBpYJr2lCB-UYdRRTFe5swVQW8iLh5a"
)

var SigningKey = []byte("AddDevices")

type jwtSignation struct {
	jwt.StandardClaims
	Platform string `json:"platform,omitempty"`
	Pass     string `json:"pass,omitempty"`
}
type Doc interface {
	initApi(router *mux.Router)
}

type queryRes struct {
	Rows     *[]interface{} `json:"result,omitempty"`
	Total    int            `json:"total,omitempty"`
	AllPages int            `json:"allPages,omitempty"`
	Page     int            `json:"page,omitempty"`
	Limit    int            `json:"limit,omitempty"`
}
type onlyRes struct {
	Rows *[]interface{} `json:"result,omitempty"`
}
type rawAEMDRA struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	QA             float64       `json:"qa,omitempty" bson:"qa"`
	UCA            float64       `json:"uca,omitempty" bson:"uca"`
	PSum           float64       `json:"p_sum,omitempty" bson:"p_sum"`
	PB             float64       `json:"pb,omitempty" bson:"pb"`
	PC             float64       `json:"pc,omitempty" bson:"pc"`
	PA             float64       `json:"pa,omitempty" bson:"pa"`
	PFAvg          float64       `json:"pf_avg,omitempty" bson:"pf_avg"`
	PFA            float64       `json:"pfa,omitempty" bson:"pfa"`
	PFC            float64       `json:"pfc,omitempty" bson:"pfc"`
	PFB            float64       `json:"pfb,omitempty" bson:"pfb"`
	QC             float64       `json:"qc,omitempty" bson:"qc"`
	QSum           float64       `json:"q_sum,omitempty" bson:"q_sum"`
	UAvg           float64       `json:"u_avg,omitempty" bson:"u_avg"`
	AEA            float64       `json:"aea,omitempty" bson:"aea"`
	ULNAvg         float64       `json:"uln_avg,omitempty" bson:"uln_avg"`
	AEC            float64       `json:"aec,omitempty" bson:"aec"`
	QB             float64       `json:"qb,omitempty" bson:"qb"`
	UA             float64       `json:"ua,omitempty" bson:"ua"`
	UC             float64       `json:"uc,omitempty" bson:"uc"`
	UB             float64       `json:"ub,omitempty" bson:"ub"`
	LastReportTime time.Time     `json:"lastReportTime,omitempty" bson:"lastReportTime"`
	BlockID        string        `json:"blockId,omitempty" bson:"blockId"`
	IAvg           float64       `json:"i_avg,omitempty" bson:"i_avg"`
	IA             float64       `json:"ia,omitempty" bson:"ia"`
	IC             float64       `json:"ic,omitempty" bson:"ic"`
	UAB            float64       `json:"uab,omitempty" bson:"uab"`
	AETot          float64       `json:"ae_tot,omitempty" bson:"ae_tot"`
	IB             float64       `json:"ib,omitempty" bson:"ib"`
	RETot          float64       `json:"re_tot,omitempty" bson:"re_tot"`
	DevID          string        `json:"devID,omitempty" bson:"devID"`
	GWID           string        `json:"GWID,omitempty" bson:"GWID"`
	AEB            float64       `json:"aeb,omitempty" bson:"aeb"`
	REB            float64       `json:"reb,omitempty" bson:"reb"`
	REC            float64       `json:"rec,omitempty" bson:"rec"`
	REA            float64       `json:"rea,omitempty" bson:"rea"`
	Wire           float64       `json:"Wire,omitempty" bson:"wire"`
	UBC            float64       `json:"ubc,omitempty" bson:"ubc"`
	Freq           float64       `json:"freq,omitempty" bson:"freq"`
	SSum           float64       `json:"s_sum,omitempty" bson:"s_sum"`
	SC             float64       `json:"sc,omitempty" bson:"sc"`
	SB             float64       `json:"sb,omitempty" bson:"sb"`
	SA             float64       `json:"sa,omitempty" bson:"sa"`
}

type aggAllToday struct {
	Rows []struct {
		GatewayID      string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime time.Time `json:"lastReportTime" bson:"lastReportTime"`
		CC             float32   `json:"CC" bson:"CC"`
		PwrUsage       float64   `json:"avg_Usage" bson:"avg_Usage"`
		BuildingName   string    `json:"Building_Name" bson:"Building_Name"`
		PwrDemand      float64   `json:"Pwr_Demand" bson:"avg_Demand"`
		MaxDemand      float64   `json:"max_Demand" bson:"max_Demand"`
		MinDemand      float64   `json:"min_Demand" bson:"min_Demand"`
		MinUsage       float64   `json:"min_Usage" bson:"min_Usage"`
		MaxUsage       float64   `json:"max_Usage" bson:"max_Usage"`
		Weather        int8      `json:"weather_Temp" bson:"weather"`
		AvgPF          float64   `json:"PF" bson:"avg_PF"`
		MaxPF          float64   `json:"max_PF" bson:"max_PF"`
		MinPF          float64   `json:"min_PF" bson:"min_PF"`
	}
	Datashape struct {
		FieldDefinitions struct {
			GatewayID      DSTwxTemplate `json:"GatewayID"`
			LastReportTime DSTwxTemplate `json:"lastReportTime"`
			CC             DSTwxTemplate `json:"CC"`
			PwrUsage       DSTwxTemplate `json:"PwrUsage"`
			BuildingName   DSTwxTemplate `json:"Building_Name" `
			PwrDemand      DSTwxTemplate `json:"PwrDemand" `
			MaxDemand      DSTwxTemplate `json:"MaxDemand" `
			MinDemand      DSTwxTemplate `json:"MinDemand" `
			MinUsage       DSTwxTemplate `json:"MinUsage" `
			MaxUsage       DSTwxTemplate `json:"MaxUsage" `
			Weather        DSTwxTemplate `json:"Weather" `
			AvgPF          DSTwxTemplate `json:"AvgPF" `
			MaxPF          DSTwxTemplate `json:"MaxPF"`
			MinPF          DSTwxTemplate `json:"MinPF"`
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type aggAllNow struct {
	Rows []struct {
		GatewayID      string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime time.Time `json:"lastReportTime" bson:"lastReportTime"`
		CC             float32   `json:"CC" bson:"CC"`
		PwrUsage       float64   `json:"avg_Usage" bson:"avg_Usage"`
		BuildingName   string    `json:"Building_Name" bson:"Building_Name"`
		PwrDemand      float64   `json:"Pwr_Demand" bson:"avg_Demand"`
		MaxDemand      float64   `json:"max_Demand" bson:"max_Demand"`
		MinDemand      float64   `json:"min_Demand" bson:"min_Demand"`
		MinUsage       float64   `json:"min_Usage" bson:"min_Usage"`
		MaxUsage       float64   `json:"max_Usage" bson:"max_Usage"`
		AvgPF          float32   `json:"PF" bson:"avg_PF"`
		MaxPF          float32   `json:"max_PF" bson:"max_PF"`
		MinPF          float32   `json:"min_PF" bson:"min_PF"`
	}
	Datashape struct {
		FieldDefinitions struct {
			GatewayID      DSTwxTemplate `json:"GatewayID"`
			LastReportTime DSTwxTemplate `json:"lastReportTime"`
			CC             DSTwxTemplate `json:"CC"`
			PwrUsage       DSTwxTemplate `json:"PwrUsage"`
			BuildingName   DSTwxTemplate `json:"Building_Name" `
			PwrDemand      DSTwxTemplate `json:"PwrDemand" `
			MaxDemand      DSTwxTemplate `json:"MaxDemand" `
			MinDemand      DSTwxTemplate `json:"MinDemand" `
			MinUsage       DSTwxTemplate `json:"MinUsage" `
			MaxUsage       DSTwxTemplate `json:"MaxUsage" `
			AvgPF          DSTwxTemplate `json:"AvgPF" `
			MaxPF          DSTwxTemplate `json:"MaxPF"`
			MinPF          DSTwxTemplate `json:"MinPF"`
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}
type strAllHourOnTimeWODS struct {
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
	PF              float64   `json:"avg_PF" bson:"avg_PF"`
	MinPF           float64   `json:"min_PF" bson:"min_PF"`
	PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
	PwrUsage        float64   `json:"Pwr_Usage" bson:"avg_Usage"`
	MaxPF           float64   `json:"max_PF" bson:"max_PF"`
	MinDemand       float64   `json:"min_Demand" bson:"min_Demand"`
	MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"`
	MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`
	CC              float64   `json:"CC" bson:"CC"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
}

type strAllHourOnTime struct {
	Rows []struct {
		GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
		MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"`
		BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
		BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
		PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
		PF              float64   `json:"avg_PF" bson:"avg_PF"`
		MinPF           float64   `json:"min_PF" bson:"min_PF"`
		PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
		PwrUsage        float64   `json:"total_Usage" bson:"total_Usage"`
		MaxPF           float64   `json:"max_PF" bson:"max_PF"`
		MinDemand       float64   `json:"min_Demand" bson:"min_Demand"`
		MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"`
		MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`
		CC              float64   `json:"CC" bson:"CC"`
		WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			GatewayID       DSTwxTemplate `json:"GatewayID"`
			LastReportTime  DSTwxTemplate `json:"lastReportTime"`
			CC              DSTwxTemplate `json:"CC"`
			PwrUsage        DSTwxTemplate `json:"total_Usage"`
			BuildingName    DSTwxTemplate `json:"Building_Name" `
			BuildingDetails DSTwxTemplate `json:"Building_Details" `
			PwrDemand       DSTwxTemplate `json:"PwrDemand" `
			MaxDemand       DSTwxTemplate `json:"MaxDemand" `
			MinDemand       DSTwxTemplate `json:"MinDemand" `
			MinUsage        DSTwxTemplate `json:"MinUsage" `
			MaxUsage        DSTwxTemplate `json:"MaxUsage" `
			PF              DSTwxTemplate `json:"PF" `
			PFLimit         DSTwxTemplate `json:"PFLimit" `
			MaxPF           DSTwxTemplate `json:"MaxPF"`
			MinPF           DSTwxTemplate `json:"MinPF"`
			WeatherTemp     DSTwxTemplate `json:"weather_Temp"`
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type strAllDayMonthOnTime struct {
	Rows []struct {
		GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
		MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"`
		BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
		BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
		PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
		PF              float64   `json:"avg_PF" bson:"avg_PF"`
		MinPF           float64   `json:"min_PF" bson:"min_PF"`
		PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
		PwrUsage        float64   `json:"avg_Usage" bson:"avg_Usage"`
		TotalUsage      float64   `json:"total_Usage" bson:"total_Usage"`
		MaxPF           float64   `json:"max_PF" bson:"max_PF"`
		MinDemand       float64   `json:"min_Demand" bson:"min_Demand"`
		MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"`
		MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`
		CC              float64   `json:"CC" bson:"CC"`
		WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`

		// PrevmaxUsage    float64 `json:"Prev_max_Usage" bson:"Prev_max_Usage"`
		// PrevminUsage    float64 `json:"Prev_min_Usage" bson:"Prev_min_Usage"`
		// PrevPwrUsage    float64 `json:"Prev_Pwr_Usage" bson:"Prev_Pwr_Usage"`
		// PrevPwrDemand   float64 `json:"Prev_Pwr_Demand" bson:"Prev_Pwr_Demand"`
		// PrevmaxDemand   float64 `json:"Prev_max_Demand" bson:"Prev_max_Demand"`
		// PrevminDemand   float64 `json:"Prev_min_Demand" bson:"Prev_min_Demand"`
		// PrevmaxPF       float64 `json:"Prev_max_PF" bson:"Prev_max_PF"`
		// PrevminPF       float64 `json:"Prev_min_PF" bson:"Prev_min_PF"`
		// PrevPF          float64 `json:"Prev_PF" bson:"Prev_PF"`
		// PrevWeatherTemp int     `json:"Prev_weather_Temp" bson:"Prev_weather_Temp"`
	} `json:"rows"`

	Datashape struct {
		FieldDefinitions struct {
			GatewayID       DSTwxTemplate `json:"GatewayID"`
			LastReportTime  DSTwxTemplate `json:"lastReportTime"`
			CC              DSTwxTemplate `json:"CC"`
			PwrUsage        DSTwxTemplate `json:"avg_Usage"`
			BuildingName    DSTwxTemplate `json:"Building_Name" `
			BuildingDetails DSTwxTemplate `json:"Building_Details" `
			PwrDemand       DSTwxTemplate `json:"avg_Demand" `
			MaxDemand       DSTwxTemplate `json:"max_Demand" `
			MinDemand       DSTwxTemplate `json:"min_Demand" `
			MinUsage        DSTwxTemplate `json:"min_Usage" `
			MaxUsage        DSTwxTemplate `json:"max_Usage" `
			PF              DSTwxTemplate `json:"avg_PF" `
			PFLimit         DSTwxTemplate `json:"PF_Limit" `
			MaxPF           DSTwxTemplate `json:"max_PF"`
			MinPF           DSTwxTemplate `json:"min_PF"`
			WeatherTemp     DSTwxTemplate `json:"weather_Temp"`

			// PrevmaxUsage DSTwxTemplate `json:"Prev_max_Usage" `
			// PrevminUsage DSTwxTemplate `json:"Prev_min_Usage" `
			// PrevPwrUsage DSTwxTemplate `json:"Prev_Pwr_Usage" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type AggAllOnTime struct {
	Rows []struct {
		GatewayID      string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime time.Time `json:"lastReportTime" bson:"lastReportTime"`
		CC             float32   `json:"CC" bson:"CC"`
		PwrUsage       float64   `json:"Pwr_Usage" bson:"avg_Usage"`
		BuildingName   string    `json:"Building_Name" bson:"Building_Name"`
		PwrDemand      float64   `json:"Pwr_Demand" bson:"avg_Demand"`
		AvgPF          float32   `json:"PF" bson:"PF"`
		Weather        int8      `json:"weather_Temp" bson:"weather"`
		Usage          float32   `json:"Usage" bson:"Usage"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			GatewayID      DSTwxTemplate `json:"GatewayID"`
			LastReportTime DSTwxTemplate `json:"lastReportTime"`
			CC             DSTwxTemplate `json:"CC"`
			PwrUsage       DSTwxTemplate `json:"Pwr_Usage"`
			Usage          DSTwxTemplate `json:"Usage"`
			BuildingName   DSTwxTemplate `json:"Building_Name" `
			PwrDemand      DSTwxTemplate `json:"Pwr_Demand" `
			AvgPF          DSTwxTemplate `json:"PF" `
			Weather        DSTwxTemplate `json:"weather" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type DisplayData struct {
	Rows []struct {
		GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
		LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
		PwrUsage        float64   `json:"Pwr_Usage" bson:"Pwr_Usage"`
		BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
		BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
		DevID           string    `json:"Device_ID" bson:"Device_ID"`
		PwrDemand       float64   `json:"Pwr_Demand" bson:"Pwr_Demand"`
		PF              float32   `json:"PF" bson:"PF"`
		Usage           float32   `json:"Usage" bson:"Usage"`
		Weather         int8      `json:"weather_Temp" bson:"weather_Temp"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			GatewayID       DSTwxTemplate `json:"Gateway_ID"`
			LastReportTime  DSTwxTemplate `json:"lastReportTime"`
			PwrUsage        DSTwxTemplate `json:"Pwr_Usage"`
			BuildingName    DSTwxTemplate `json:"Building_Name" `
			BuildingDetails DSTwxTemplate `json:"Building_Details" `
			PwrDemand       DSTwxTemplate `json:"Pwr_Demand" `
			PF              DSTwxTemplate `json:"PF" `
			Usage           DSTwxTemplate `json:"Usage" `
			DevID           DSTwxTemplate `json:"Device_ID" `
			Weather         DSTwxTemplate `json:"weather_Temp" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type DisplayDataElement struct {
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
}
type DeviceManagerS struct {
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	DeviceID        string    `json:"devID" bson:"devID"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	GatewayID       string    `json:"GWID" bson:"GWID"`
	DeviceInfo      string    `json:"Device_Info" bson:"Device_Info"`
}

type ListDevices struct {
	Rows []struct {
		LastReportTime  time.Time `json:"Time_Added" bson:"Time_Added"`
		DeviceID        string    `json:"Device_ID" bson:"devID"`
		BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
		BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
		GatewayID       string    `json:"Gateway_ID" bson:"GWID"`
		DeviceInfo      string    `json:"Device_Info" bson:"Device_Info"`
		DeviceBrand     string    `json:"Device_Brand" bson:"Device_Brand"`
		DeviceDetails   string    `json:"Device_Details" bson:"Device_Details"`
		DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
		Floor           string    `json:"Floor" bson:"Floor"`
		Facility        string    `json:"Facility" bson:"Facility"`
		DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	} `json:"rows"`
	Datashape struct {
		FieldDefinitions struct {
			GatewayID       DSTwxTemplate `json:"Gateway_ID"`
			DeviceInfo      DSTwxTemplate `json:"Device_Info"`
			DeviceID        DSTwxTemplate `json:"Device_ID"`
			LastReportTime  DSTwxTemplate `json:"Time_Added"`
			BuildingName    DSTwxTemplate `json:"Building_Name" `
			BuildingDetails DSTwxTemplate `json:"Building_Details" `
			DeviceType      DSTwxTemplate `json:"Device_Type" `
			Floor           DSTwxTemplate `json:"Floor" `
			DeviceName      DSTwxTemplate `json:"Device_Name" `
			DeviceDetails   DSTwxTemplate `json:"Device_Details" `
			DeviceBrand     DSTwxTemplate `json:"Device_Brand" `
			Facility        DSTwxTemplate `json:"Facility" `
		} `json:"fieldDefinitions"`
	} `json:"dataShape"`
}

type DisplayDataElement2nd struct {
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	PwrUsage        float64   `json:"Pwr_Usage" bson:"Pwr_Usage"`
	PwrDemand       float64   `json:"Pwr_Demand" bson:"Pwr_Demand"`
	AvgPF           float64   `json:"PF" bson:"PF"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	Usage           float64   `json:"Usage" bson:"Usage"`
}

type autoGenerated struct {
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	CC              float64   `json:"CC" bson:"CC"`
	PwrUsage        float64   `json:"Pwr_Usage" bson:"avg_Usage"`
	Floor           string    `json:"Floor" bson:"Floor"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	PwrDemand       float64   `json:"Pwr_Demand" bson:"avg_Demand"`
	QuarterPost     int       `json:"QuarterPost" bson:"QuarterPost"`
	Facility        string    `json:"Facility" bson:"Facility"`
	DeviceDetails   string    `json:"Device_Details" bson:"Device_Details"`
	AvgPF           float64   `json:"PF" bson:"avg_PF"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
}
type AggHourStruct struct {
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	CC              float64   `json:"CC" bson:"CC"`
	PwrUsage        float64   `json:"avg_Usage" bson:"avg_Usage"`
	Floor           string    `json:"Floor" bson:"Floor"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	WeatherTemp     int       `json:"weather_Temp" bson:"weather_Temp"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
	AvgPF           float64   `json:"avg_PF" bson:"avg_PF"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"`
	PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
	MinPF           float64   `json:"min_PF" bson:"min_PF"`
	MaxPF           float64   `json:"max_PF" bson:"max_PF"`
	MinDemand       float64   `json:"min_Demand" bson:"min_Demand"`
	MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"`
	MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`
}
type AggDayStruct struct {
	GatewayID       string    `json:"Gateway_ID" bson:"Gateway_ID"`
	LastReportTime  time.Time `json:"lastReportTime" bson:"lastReportTime"`
	CC              float64   `json:"CC" bson:"CC"`
	PwrUsage        float64   `json:"avg_Usage" bson:"avg_Usage"`
	Floor           string    `json:"Floor" bson:"Floor"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	TotalUsage      float64   `json:"total_Usage" bson:"total_Usage"`
	PwrDemand       float64   `json:"avg_Demand" bson:"avg_Demand"`
	AvgPF           float64   `json:"PF" bson:"avg_PF"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	DeviceID        string    `json:"Device_ID" bson:"Device_ID"`
	MaxUsage        float64   `json:"max_Usage" bson:"max_Usage"`
	MinPF           float64   `json:"min_PF" bson:"min_PF"`
	MaxPF           float64   `json:"max_PF" bson:"max_PF"`
	MinDemand       float64   `json:"min_Demand" bson:"min_Demand"`
	MaxDemand       float64   `json:"max_Demand" bson:"max_Demand"`
	MinUsage        float64   `json:"min_Usage" bson:"min_Usage"`
	PFLimit         float64   `json:"PF_Limit" bson:"PF_Limit"`
}

type node struct {
	Rows []autoGenerated `json:"result,omitempty"`
}

type rawCPM struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	PFA            float64       `json:"pfa,omitempty" bson:"pfa"`
	PFC            float64       `json:"pfc,omitempty" bson:"pfc"`
	PFB            float64       `json:"pfb,omitempty" bson:"pfb"`
	IAvgTHD        float64       `json:"iavg_thd,omitempty" bson:"iavg_thd"`
	LastReportTime time.Time     `json:"lastReportTime,omitempty" bson:"lastReportTime"`
	IAvg           float64       `json:"i_avg,omitempty" bson:"i_avg"`
	UCA            float64       `json:"uca,omitempty" bson:"uca"`
	UBC            float64       `json:"ubc,omitempty" bson:"ubc"`
	IA             float64       `json:"ia,omitempty" bson:"ia"`
	Freq           float64       `json:"freq,omitempty" bson:"freq"`
	UAB            float64       `json:"uab,omitempty" bson:"uab"`
	IB             float64       `json:"ib,omitempty" bson:"ib"`
	UAvg           float64       `json:"u_avg,omitempty" bson:"u_avg"`
	AETot          float64       `json:"ae_tot,omitempty" bson:"ae_tot"`
	Wire           float64       `json:"wire,omitempty" bson:"wire"`
	PSum           float64       `json:"p_sum,omitempty" bson:"p_sum"`
	PC             float64       `json:"pc,omitempty" bson:"pc"`
	IC             float64       `json:"ic,omitempty" bson:"ic"`
	SSum           float64       `json:"s_sum,omitempty" bson:"s_sum"`
	DevID          string        `json:"devID,omitempty" bson:"devID"`
	GWID           string        `json:"GWID,omitempty" bson:"GWID"`
	PB             float64       `json:"pb,omitempty" bson:"pb"`
	ULNAvg         float64       `json:"uln_avg,omitempty" bson:"uln_avg"`
	UAvgTHD        float64       `json:"uavg_thd,omitempty" bson:"uavg_thd"`
	PA             float64       `json:"pa,omitempty" bson:"pa"`
	SC             float64       `json:"sc,omitempty" bson:"sc"`
	SB             float64       `json:"sb,omitempty" bson:"sb"`
	SA             float64       `json:"sa,omitempty" bson:"sa"`
	UA             float64       `json:"ua,omitempty" bson:"ua"`
	PFAvg          float64       `json:"pf_avg,omitempty" bson:"pf_avg"`
	UC             float64       `json:"uc,omitempty" bson:"uc"`
	UB             float64       `json:"ub,omitempty" bson:"ub"`
}
type powerCons struct {
	ID  int `json:"id"`
	KWh int `json:"kWh"`
	W   int `json:"W"`
	Min int `json:"min"`
	Max int `json:"max"`
	Avg int `json:"avg"`
}

type RegisterDevID struct {
	GatewayID       string    `json:"GWID" bson:"GWID"`
	DeviceBrand     string    `json:"Device_Brand" bson:"Device_Brand"`
	DeviceID        string    `json:"devID" bson:"devID"`
	DeviceDetails   string    `json:"Device_Details" bson:"Device_Details"`
	DeviceName      string    `json:"Device_Name" bson:"Device_Name"`
	DeviceInfo      string    `json:"Device_Info" bson:"Device_Info"`
	DeviceType      string    `json:"Device_Type" bson:"Device_Type"`
	Floor           string    `json:"Floor" bson:"Floor"`
	Facility        string    `json:"Facility" bson:"Facility"`
	BuildingName    string    `json:"Building_Name" bson:"Building_Name"`
	BuildingDetails string    `json:"Building_Details" bson:"Building_Details"`
	TimeAdded       time.Time `json:"Time_Added" bson:"Time_Added"`
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

func (w *rawAEMDRA) GenObjectId() {
	if bson.ObjectId("") == w.ID {

		w.ID = getObjectIDTwoArg(w.GWID, w.DevID, w.LastReportTime)
	}
}

func (w *rawCPM) GenObjectId() {
	if bson.ObjectId("") == w.ID {

		w.ID = getObjectIDTwoArg(w.GWID, w.DevID, w.LastReportTime)
	}
}

func getObjectIDTwoArg(GWID string, DevID string, timestamp time.Time) bson.ObjectId {
	var b [12]byte
	var sum [8]byte

	// timestamp := time.Unix(LastReportTime, 0)
	binary.BigEndian.PutUint32(b[:], uint32(timestamp.Unix()))

	did := sum[:]
	gid := sum[:]

	hw := md5.New()
	hw.Write([]byte(GWID))
	copy(did, hw.Sum(nil))
	hw.Write([]byte(DevID))
	copy(gid, hw.Sum(nil))

	b[4] = did[1]
	b[5] = did[2]
	b[6] = did[3]
	b[7] = did[4]
	b[8] = gid[5]
	b[9] = gid[6]
	b[10] = gid[7]
	b[11] = gid[8]

	return bson.ObjectId(b[:])

}

func GetObjectIDOneArg(DevID string, timestamp time.Time) bson.ObjectId {
	var b [12]byte
	var sum [9]byte

	// timestamp := time.Unix(LastReportTime, 0)
	binary.BigEndian.PutUint32(b[:], uint32(timestamp.Unix()))

	// did := sum[:]
	gid := sum[:]

	hw := md5.New()
	// hw.Write([]byte(GWID))
	// copy(did, hw.Sum(nil))
	hw.Write([]byte(DevID))
	copy(gid, hw.Sum(nil))

	b[4] = gid[1]
	b[5] = gid[2]
	b[6] = gid[3]
	b[7] = gid[4]
	b[8] = gid[5]
	b[9] = gid[6]
	b[10] = gid[7]
	b[11] = gid[8]

	return bson.ObjectId(b[:])

}

type ResDoc struct {
	UploadData      string   `json:"Upload_Data,omitempty" bson:"Upload_Data"`     // 資料上傳
	ConfigFlag      string   `json:"Config_Flag,omitempty" bson:"Config_Flag"`     // 設定檔更新
	IPURL           []string `json:"IP_URL,omitempty" bson:"IP_URL"`               // 資料中心_IP_1 ~IP_5
	DateTime        int64    `json:"DateTime,omitempty" bson:"DateTime"`           // 今天日期時間
	SendRate        int      `json:"Send_Rate,omitempty" bson:"Send_Rate"`         // 傳送頻率
	GainRate        int      `json:"Gain_Rate,omitempty" bson:"Gain_Rate"`         // 資料採集頻率
	Resend          string   `json:"Resend,omitempty" bson:"Resend"`               // 資料重送
	BackupTime      int      `json:"Backup_Time,omitempty" bson:"Backup_Time"`     // 資料儲放時間
	MACAddr         string   `json:"MAC_Address,omitempty" bson:"MAC_Address"`     // MAC_Address
	StationID       string   `json:"Station_ID,omitempty" bson:"Station_ID"`       // Station_ID
	ResendTimeStart int64    `json:"Resend_time_S,omitempty" bson:"Resend_time_S"` // 要求特定時間重送資料
	ResendTimeEnd   int64    `json:"Resend_time_E,omitempty" bson:"Resend_time_E"` // 要求特定時間重送資料
	Command         int      `json:"Command,omitempty" bson:"Command"`             // 重置累積功率
	Version         string   `json:"GW_Ver,omitempty" bson:"GW_Ver"`               // 軔體版本確認
}

type DSTwxTemplate struct {
	Name        string      `json:"name,omitempty" bson:"name"`
	Description string      `json:"description,omitempty" bson:"description"`
	BaseType    string      `json:"baseType,omitempty" bson:"baseType"`
	Ordinal     int         `json:"ordinal" bson:"ordinal"`
	Aspects     interface{} `json:"aspects,omitempty" bson:"aspects"`
}

type DSTwx struct {
	FieldDefinitions FieldDefinitions `json:"fieldDefinitions"`
}

type FieldDefinitions struct {
	LastReportTime  DSTwxTemplate `json:"lastReportTime"`
	DevID           DSTwxTemplate `json:"Device_ID"`
	BuildingName    DSTwxTemplate `json:"Building_Name" `
	BuildingDetails DSTwxTemplate `json:"Building_Details" `
	GatewayID       DSTwxTemplate `json:"Gateway_ID" `
}

type genIDStruct struct {
	AccessToken string `json:"access_token" bson:"access_token"`
	Expire      int    `json:"expires_in" bson:"expires_in"`
	Type        string `json:"token_type" bson:"token_type"`
}
