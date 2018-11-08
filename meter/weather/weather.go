package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"

	// change due to high cpu using globalsign
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dblocal  = "172.16.0.132:27017"
	dbpublic = "140.118.70.136:10003"

	dbbackup = "140.118.122.103:27017"

	weatherCollection = "weather"
	localTimeZoneExt  = "+08:00"
	c_weather         = "weather"
	db                = "sc"
)

var bannr = `
Program name : weather

maintainer   : avbee.lab@gmail.com

Date         : 09/06/2018

`

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

func db_connect() {

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(dbbackup, ",", -1),
		Database: "admin",
		Username: "dontask",
		Password: "idontknow",
		Timeout:  time.Second * 2,
	}
	session, _ = mgo.DialWithInfo(dbInfo)
}

func getObjectIDTwoArg(GWID string, macID string, timestamp int64) bson.ObjectId {
	var b [12]byte
	var sum [8]byte
	var c [4]byte
	// timestamp := time.Unix(LastReportTime, 0)
	// binary.BigEndian.PutUint32(b[:], uint32(timestamp))
	binary.BigEndian.PutUint32(c[:], uint32(timestamp))

	did := sum[:]
	gid := sum[:]

	hw := md5.New()
	hw.Write([]byte(GWID))
	copy(did, hw.Sum(nil))
	hw.Write([]byte(macID))
	copy(gid, hw.Sum(nil))
	// b[0] = c[:1]
	b[0] = c[0]
	b[1] = c[1]
	b[2] = c[2]
	b[3] = c[3]
	b[4] = did[4]
	b[5] = did[5]
	b[6] = did[6]
	b[7] = did[7]
	b[8] = gid[4]
	b[9] = gid[5]
	b[10] = gid[6]
	b[11] = gid[7]

	dst := hex.EncodeToString(b[:])
	theid := bson.ObjectIdHex(dst)

	fmt.Println(theid, uint32(timestamp))
	return theid

}

func getweather() {
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
		finaldata.Town = townID

		for _, weather := range a.Rawdata.Locations[0].Location[0].WeatherElement {
			if weather.ElementName == "T" {

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

			}
			finaldata.ElementName = weather.ElementName
		}

		fmt.Println("finish")

	}

}

func main() {
	db_connect()
	fmt.Println(bannr)
	c := cron.New()

	c.AddFunc("@daily", getweather)

	c.Start()
	select {}
}
