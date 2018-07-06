package alert

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"dforcepro.com/util"
	yaml "gopkg.in/yaml.v2"
)

const (
	WarnAlert       = 1 << iota
	ErrorAlert      = 1 << iota
	NotSupportAlert = 1 << iota
)

var ConfTpl *AlertConfTpl

type Alert struct {
	AlertType int
	Field     string
	Title     string
	Message   string
	Minute    float64
	Times     uint8
}

//小寫是給PUT /v1/device/alert用的
type FieldAlert struct {
	Field  string       `json:"field, "`  // 欄位名稱
	Enable bool         `json:"enable, "` // 是否啟用
	Title  string       `json:"title, "`  //	告警名稱
	Error  *AlertConfig `json:"error, "`  // 錯誤門檻設定
	Warn   *AlertConfig `json:"warn, "`   // 警告門檻設定
}

type AlertConfig struct {
	Upper      float64 `json:"upper, "`
	Lower      float64 `json:"lower, "`
	Minute     float64 `json:"minute, "`
	Times      uint8   `json:"times, "`
	MessageTpl string  `json:"messagetpl, "` // ex. The %s value is %.2f is %s then %.2f.
}

type AlertConfTpl map[string][]FieldAlert

func (ac *AlertConfTpl) GetConf(name string) *[]FieldAlert {
	conf, ok := (*ac)[name]
	if ok {
		return &conf
	}
	return nil
}

func (ac *Alert) ToKeyStr() string {
	return fmt.Sprintf("[%d]:%s", ac.AlertType, ac.Title)
}

func (ac *FieldAlert) Validate(inter interface{}) *Alert {
	s := reflect.ValueOf(inter).Elem()
	field := s.FieldByName(ac.Field)
	if !field.IsValid() {
		return nil
	}

	var value float64
	if "float64" != field.Type().Name() {
		return &Alert{AlertType: NotSupportAlert, Message: "not support not float64 value"}
	}
	value = field.Float()

	alert := validateValue(ac.Field, value, ac.Error, ErrorAlert)
	if alert != nil {
		alert.Title = ac.Title
		alert.Field = ac.Field
		return alert
	}

	alert = validateValue(ac.Field, value, ac.Warn, WarnAlert)
	if alert != nil {
		alert.Title = ac.Title
		alert.Field = ac.Field
		return alert
	}
	return nil
}

func validateValue(fieldName string, value float64, ac *AlertConfig, alertType int) *Alert {
	if value >= ac.Upper {
		message := fmt.Sprintf(ac.MessageTpl, fieldName, value, "upper", ac.Upper)
		return &Alert{AlertType: alertType, Message: message, Minute: ac.Minute, Times: ac.Times}
	}

	if value <= ac.Lower {
		message := fmt.Sprintf(ac.MessageTpl, fieldName, value, "lower", ac.Lower)
		return &Alert{AlertType: alertType, Message: message, Minute: ac.Minute, Times: ac.Times}
	}

	return nil
}

func InitTpl(ymlfile string) error {
	filename, _ := filepath.Abs(ymlfile)
	exist, _ := util.FileExists(filename)
	if !exist {
		return errors.New(fmt.Sprintf("file not exist: %s", filename))
	}
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}
	ac := AlertConfTpl{}
	err = yaml.Unmarshal(yamlFile, &ac)
	if err != nil {
		return err
	}
	ConfTpl = &ac
	return nil
}
