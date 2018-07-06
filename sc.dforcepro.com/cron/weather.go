package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"dforcepro.com/cron"
	"gopkg.in/mgo.v2/bson"
	"sc.dforcepro.com/meter"
)

const (
	weatherCollection = "SC01_Weather"
	localTimeZoneExt  = "+08:00"
)

type Weather bool

func (mycron Weather) Enable() bool {
	return bool(mycron)
}

func (mycron Weather) GetJobs() []cron.JobSpec {
	return []cron.JobSpec{
		cron.JobSpec{
			Spec: "0 1 1 * * *",
			// Spec: "*/30 * * * * *",
			Job: mycron.getweather,
		},
	}
}

func (mycron Weather) getweather() {
	// mark := make(map[string]string)
	fmt.Println("hello")
	// mark["A"] = "F-D0047-063"
	// mark["B"] = "F-D0047-075"
	// mark["C"] = "F-D0047-051"
	// mark["D"] = "F-D0047-079"
	// mark["E"] = "F-D0047-067"
	// mark["F"] = "F-D0047-071"
	// mark["G"] = "F-D0047-003"
	// mark["H"] = "F-D0047-007"
	// mark["I"] = "F-D0047-057"
	// mark["J"] = "F-D0047-009"
	// mark["K"] = "F-D0047-015"
	// mark["M"] = "F-D0047-023"
	// mark["N"] = "F-D0047-017"
	// mark["O"] = "F-D0047-055"
	// mark["P"] = "F-D0047-027"
	// mark["Q"] = "F-D0047-029"
	// mark["T"] = "F-D0047-033"
	// mark["U"] = "F-D0047-043"
	// mark["V"] = "F-D0047-037"
	// mark["W"] = "F-D0047-087"
	// mark["X"] = "F-D0047-047"
	// mark["Z"] = "F-D0047-083"

	// fmt.Println("getweather")
	// // cityarray := []citys{}

	mongo := meter.GetMongo()
	// c := mongo.DB(meter.DBName).C(weatherCollection)
	// c.Find(bson.M{}).All(&cityarray)
	/*cityarray := []citys{}

	c := session.DB(zbtDb).C("City")

	c.Find(bson.M{}).All(&cityarray)
	fmt.Println(len(cityarray))*/
	// for _, manycity := range cityarray {
	// fmt.Println("start")
	// id := manycity.Cityid
	// townsarray := []towns{}
	// c := mongo.DB(meter.DBName).C("Town")
	// fmt.Println(manycity.Cityid)
	// err := c.Find(bson.M{"cityid": bson.M{"$in": []string{manycity.Cityid}}}).All(&townsarray)
	// fmt.Println(err)
	// fmt.Println(len(townsarray))
	// for _, manytown := range townsarray {

	cityID := "F-D0047-063"
	townID := "大安區"
	fmt.Println("start")
	response, err := http.Get("http://opendata.cwb.gov.tw/api/v1/rest/datastore/" + cityID + "?locationName=" + townID + "&elementName=PoP,T,Wx,RH,WeatherDescription&sort=time&Authorization=CWB-2FA1D452-8CE2-4EDC-BCDD-B550B36061E1")
	fmt.Println("end")
	if err != nil {
		defer response.Body.Close()

		a := Rawdata{}
		//json.Unmarshal(contents, &a)
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		// fmt.Printf("%s\n", contents)
		err = json.Unmarshal(contents, &a)
		// fmt.Println(string(contents))

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
				for _, j := range weather.Time {
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

					mongo.DB(meter.DBName).C(weatherCollection).Upsert(bson.M{"weatherElement": finaldata.ElementName, "startTime": finaldata.Startime, "endTime": finaldata.EndTime}, finaldata)
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
		mongo.Close()
	}

}
