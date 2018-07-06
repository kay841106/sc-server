package cron

import "time"

import "dforcepro.com/resource"

var _di *resource.Di

func SetDI(c *resource.Di) {
	_di = c
}

type citys struct {
	City   string `json:"city,omitempty"`
	Cityid string `json:"cityid,omitempty"`
}

type towns struct {
	Town   string `json:"town,omitempty"`
	Cityid string `json:"cityid,omitempty"`
	Townid string `json:"townid,omitempty"`
}

type Rawdata struct {
	Rawdata Record `json:"records,omitempty"`
}

type Record struct {
	ContentDescription string      `json:"contentDescription,omitempty"`
	Locations          []Locations `json:"locations,omitempty"`
}

type Locations struct {
	DatasetDescription string     `json:"datasetDescription,omitempty"`
	Location           []Location `json:"location,omitempty"`
}

type Location struct {
	WeatherElement []WeatherElement `json:"weatherElement,omitempty"`
}

type WeatherElement struct {
	ElementName string `json:"elementName,omitempty"`
	Time        []Time `json:"time,omitempty"`
}

type WeatherGetS struct {
	ElementValue string    `json:"ElementValue,omitempty" bson:"ElementValue"`
	ElementName  string    `json:"WeatherElement,omitempty" bson:"WeatherElement"`
	Startime     time.Time `json:"startTime,omitempty" bson:"startTime"`
	EndTime      time.Time `json:"endTime,omitempty" bson:"endTime"`
}

type Time struct {
	Startime string `json:"startTime,omitempty"`
	EndTime  string `json:"endTime,omitempty"`
	// ElementValue string `json:"elementValue,omitempty"  bson:"elementValue"`
	ElementValue []ElVal `json:"elementValue,omitempty"`
}

type realdata struct {
	Town         string    `json:"town,omitempty" bson:"town"`
	ElementName  string    `json:"WeatherElement,omitempty" bson:"WeatherElement"`
	Startime     time.Time `json:"startTime,omitempty" bson:"startTime"`
	EndTime      time.Time `json:"endTime,omitempty" bson:"endTime"`
	ElementValue string    `json:"ElementValue,omitempty" bson:"ElementValue"`
}

type ElVal struct {
	Value    string `json:"value,omitempty"`
	Measures string `json:"measures,omitempty"`
}
