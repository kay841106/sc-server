package resource

import (
	"errors"
	"io/ioutil"

	"dforcepro.com/resource/db"
	"dforcepro.com/resource/logger"
	"dforcepro.com/resource/storage"
	amqpConf "github.com/RichardKnop/machinery/v1/config"
	yaml "gopkg.in/yaml.v2"
)

// 做為讀寫分離用，依參數回傳對應的Connect，預留做之後的擴充
const (
	READ  = iota
	WRITE = iota
)

type Di struct {
	Mongodb      db.Mongo            `yaml:"mongodb,omitempty"`
	SQL          db.SQL              `yaml:"sql,omitempty"`
	Elastic      db.Elastic          `yaml:"elastic,omitempty"`
	Redis        db.Redis            `yaml:"redis,omitempty"`
	Log          logger.Logger       `yaml:"log,omitempty"`
	APIConf      APIConf             `yaml:"api,omitempty"`
	Rabbitmq     amqpConf.Config     `yaml:"rabbitmq,omitempty"`
	FileStorage  storage.FileStorage `yaml:"fileStorage,omitempty"`
	ImageStorage storage.FileStorage `yaml:"imageStorage,omitempty"`
}

type APIConf struct {
	Port   string `yaml:"port,omitempty"`
	Middle struct {
		GenDoc bool `yaml:"gen_doc,omitempty"`
		Auth   bool `yaml:"auth,omitempty"`
		Log    bool `yaml:"log,omitempty"`
		Access bool `yaml:"access,omitempty"`
		Debug  bool `yaml:"debug,omitempty"`
	} `yaml:"middle,omitempty"`
}

type _JobRedis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	DB   string `yaml:"db"`
}

var (
	_DiConf  Di
	ConfPath string
)

// 初始化設定檔，讀YAML檔
func IniConf(path string) *Di {
	yamlFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &_DiConf)
	if err != nil {
		panic(err)
	}

	return &_DiConf
}

// 取得Conf
func GetDI() (*Di, error) {
	if (Di{}) != _DiConf {
		return &_DiConf, nil
	}
	return nil, errors.New("Conf doesn't ini. Use func IniConf first")
}

func Close() {
	_DiConf.Mongodb.Close()
	_DiConf.Redis.Close()
}
