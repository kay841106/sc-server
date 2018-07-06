package alert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type AlertTest struct {
	CPURate    float64
	MacAddress string
}

func Test_Alert_Validate(t *testing.T) {
	ea := AlertConfig{Upper: 20.2, Lower: 2.4, Minute: time.Minute.Minutes(), Times: 5, MessageTpl: "The %s value is %.2f is %s then %.2f"}
	wa := AlertConfig{Upper: 15.2, Lower: 4.1, Minute: time.Duration(30).Minutes(), Times: 5, MessageTpl: "The %s value is %.2f is %s then %.2f"}
	fa := FieldAlert{Field: "CPURate", Enable: true, Error: &ea, Warn: &wa}

	alertTest := AlertTest{CPURate: 3.3, MacAddress: "aaddd"}
	alert := fa.Validate(&alertTest)
	assert.NotNil(t, alert, "must not be nil")
	assert.Equal(t, WarnAlert, alert.AlertType, "must be WarnAlert")

	alertTest.CPURate = 21.0
	alert = fa.Validate(&alertTest)
	assert.NotNil(t, alert, "must not be nil")
	assert.Equal(t, ErrorAlert, alert.AlertType, "must be ErrorAlert")

	alertTest.CPURate = 9.2
	alert = fa.Validate(&alertTest)
	assert.Nil(t, alert, "have no alert")

	fa.Field = "MacAddress"
	alert = fa.Validate(&alertTest)
	assert.NotNil(t, alert, "must not be nil")
	assert.Equal(t, NotSupportAlert, alert.AlertType, "must be NotSupportAlert")
}

func Test_GetConf(t *testing.T) {
	InitTpl("alert_conf_tpl.yml")
	docs := ConfTpl.GetConf("WaterDoc")
	fieldAlert := (*docs)[0]
	assert.Equal(t, 2, len(*docs))
	assert.Equal(t, float64(2), fieldAlert.Error.Minute)
}
