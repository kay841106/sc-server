package alert

import (
	"time"
)

const (
	day_Duration = 24
)

type MeterALert struct {
	AlertValue *[]AlertValue `json:"Alert_Value"`
}

type AirboxAlert struct {
	AirboxAlertValue *[]AirboxAlertValue `json:"Alert_Value"`
}

// func NewMeterAlert(devID string) (*MeterALert, error) {
// 	conf := alert.ConfTpl.GetConf("MeterDoc")
// 	if conf == nil {
// 		return nil, errors.New("not exist: MeterDoc")
// 	}
// 	return &MeterALert{ID: devID, AlertValue: conf}, nil
// }

// func FindMeterAlertConfByID(devID string) (*MeterALert, error) {
// 	meterAlert = MeterALert{}
// 	// err := getMongo.
// }

func CompareTimeStamp(timeStamp time.Time) (time.Duration, time.Duration, time.Duration, bool) {
	var status bool
	// now := time.Now()
	diff := time.Since(timeStamp)

	if diff.Hours() > day_Duration {
		status = false
	} else {
		status = true
	}
	h := diff / time.Hour
	diff = diff % time.Hour
	m := diff / time.Minute
	diff = diff % time.Minute
	s := diff / time.Second
	return h, m, s, status
}
