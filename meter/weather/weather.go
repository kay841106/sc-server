package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

const (
	weatherCollection = "weather"
	localTimeZoneExt  = "+08:00"
	c_weather         = "weather"
	db                = "sc"
)

type citys struct {
	City   string `json:"city,omitempty"`
	Cityid string `json:"cityid,omitempty"`
}

type towns struct {
	Town   string `json:"town,omitempty"`
	Cityid string `json:"cityid,omitempty"`
	Townid string `json:"townid,omitempty"`
}

type rawdata struct {
	Rawdata record `json:"records,omitempty"`
}

type record struct {
	ContentDescription string      `json:"contentDescription,omitempty"`
	Locations          []locations `json:"locations,omitempty"`
}

type locations struct {
	DatasetDescription string     `json:"datasetDescription,omitempty"`
	Location           []location `json:"location,omitempty"`
}

type location struct {
	WeatherElement []weatherElement `json:"weatherElement,omitempty"`
}

type weatherElement struct {
	ElementName string    `json:"elementName,omitempty"`
	Thetime     []thetime `json:"time,omitempty"`
}

type weatherGetS struct {
	ElementValue string    `json:"ElementValue,omitempty" bson:"ElementValue"`
	ElementName  string    `json:"WeatherElement,omitempty" bson:"WeatherElement"`
	Startime     time.Time `json:"startTime,omitempty" bson:"startTime"`
	EndTime      time.Time `json:"endTime,omitempty" bson:"endTime"`
}

type thetime struct {
	Startime string `json:"startTime,omitempty"`
	EndTime  string `json:"endTime,omitempty"`
	// ElementValue string `json:"elementValue,omitempty"  bson:"elementValue"`
	ElementValue []elVal `json:"elementValue,omitempty"`
}

type realdata struct {
	Town         string    `json:"town,omitempty" bson:"town"`
	ElementName  string    `json:"WeatherElement,omitempty" bson:"WeatherElement"`
	Startime     time.Time `json:"startTime,omitempty" bson:"startTime"`
	EndTime      time.Time `json:"endTime,omitempty" bson:"endTime"`
	ElementValue string    `json:"ElementValue,omitempty" bson:"ElementValue"`
}

type elVal struct {
	Value    string `json:"value,omitempty"`
	Measures string `json:"measures,omitempty"`
}

var session *mgo.Session

func init() {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN("140.118.70.136:10003", ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 2,
	}
	session, _ = mgo.DialWithInfo(dbInfo)
}

func getweather() {
	// mark := make(map[string]string)
	sess := session.Clone()
	fmt.Println("hello")

	cityID := "F-D0047-063"
	townID := "大安區"
	fmt.Println("start")
	response, err := http.Get("http://opendata.cwb.gov.tw/api/v1/rest/datastore/" + cityID + "?locationName=" + townID + "&elementName=PoP,T,Wx,RH,WeatherDescription&sort=time&Authorization=CWB-2FA1D452-8CE2-4EDC-BCDD-B550B36061E1")
	fmt.Println(response, "end")
	if err == nil {
		defer response.Body.Close()

		a := rawdata{}
		//json.Unmarshal(contents, &a)
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		// fmt.Printf("%s\n", contents)
		err = json.Unmarshal(contents, &a)
		fmt.Println(string(contents))

		finaldata := realdata{}
		// finaldata.City = id
		finaldata.Town = townID

		// str := a.Rawdata.Locations[0].Location[0].WeatherElement[0].Time[0].Startime
		// str = strings.Replace(str, " ", "T", -1) + ".371Z"
		// t, _ := time.Parse(time.RFC3339, str)
		// finaldata.Startime = t

		// strt := a.Rawdata.Locations[0].Location[0].WeatherElement[0].Time[0].EndTime
		// strt = strings.Replace(strt, " ", "T", -1) + ".371Z"
		// t2, _ := time.Parse(time.RFC3339, strt)
		// finaldata.EndTime = t2
		// fmt.Println(a.Rawdata.Locations[0].Location[0].WeatherElement[0].Time[0].ElementValue)
		for _, weather := range a.Rawdata.Locations[0].Location[0].WeatherElement {
			if weather.ElementName == "T" {
				// fmt.Println(weather.Time[0].ElementValue[0])
				for _, j := range weather.Thetime {
					fmt.Println(j.EndTime)
					tick, _ := time.Parse(time.RFC3339, strings.Replace(j.Startime, " ", "T", -1)+localTimeZoneExt)
					tock, _ := time.Parse(time.RFC3339, strings.Replace(j.EndTime, " ", "T", -1)+localTimeZoneExt)
					finaldata.Startime = tick
					finaldata.EndTime = tock

					for _, k := range j.ElementValue {

						finaldata.ElementValue = k.Value
						finaldata.ElementName = weather.ElementName

					}
					fmt.Println(finaldata)

					sess.DB(db).C(c_weather).Upsert(bson.M{"weatherElement": finaldata.ElementName, "startTime": finaldata.Startime, "endTime": finaldata.EndTime}, finaldata)
				}

				// finaldata.ElementName = weather.ElementName
				// finaldata.ElementValue = weather.Time[i].ElementValue[i].Value j.ElementValue

			}
			finaldata.ElementName = weather.ElementName

			// finaldata.ElementValue = weather.Time[0].ElementValue.Value
			// fmt.Println("TEST" + weather.Time[0].ElementValue[0].Value)
			//finaldata.WeatherElementreal = append(finaldata.WeatherElementreal, struc)

			//c.Insert(finaldata)

			// c.Upsert(bson.M{"town": finaldata.Town, "weatherElement": finaldata.ElementName, "startTime": finaldata.Startime, "endTime": finaldata.EndTime}, finaldata)
		}
		// fmt.Println(a.Rawdata.Locations[0].Location[0].WeatherElement)
		//c.Update(bson.M{"town": finaldata.Town}, bson.M{"$set": finaldata})
		fmt.Println("finish")

	}

}

func main() {
	c := cron.New()

	c.AddFunc("0 1 * * * *", getweather)

	c.Start()
	select {}
}
