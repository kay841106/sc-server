//  @author  Fransiscus Bimo
//  @version 1.0, 06/22/18
// Add: {
// 	controlWoodhouse: agent server -> woodhouse gateway to control utilities
// }
//  AVBEE

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

const (
	mac              = "aa:bb:05:01:01:01"
	gwID             = "space_01"
	spaceAirboxdevID = "781463DA0149AD7C"
	contentType      = "application/json"
)

// struct to push monitorWoodhouse
type agentPOST struct {
	Timestamp     string  `json:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MACAddress    string  `json:"MAC_Address"`
	GWID          string  `json:"GW_ID"`
	CPURate       float64 `json:"CPU_rate"`
	StorageRate   int     `json:"Storage_rate"`
	GET11         bool    `json:"GET_1_1"`
	GET12         bool    `json:"GET_1_2"`
	GET13         bool    `json:"GET_1_3"`
	GET14         bool    `json:"GET_1_4"`
	GET15         bool    `json:"GET_1_5"`
	GET16         bool    `json:"GET_1_6"`
	GET17         bool    `json:"GET_1_7"`
	GET18         bool    `json:"GET_1_8"`
	GET19         bool    `json:"GET_1_9"`
	GET110        bool    `json:"GET_1_10"`

	SET11  int `json:"SET_1_1_0"`
	SET12  int `json:"SET_1_2_0"`
	SET13  int `json:"SET_1_3_0"`
	SET14  int `json:"SET_1_4_0"`
	SET15  int `json:"SET_1_5_0"`
	SET16  int `json:"SET_1_6_0"`
	SET17  int `json:"SET_1_7_0"`
	SET18  int `json:"SET_1_8_0"`
	SET19  int `json:"SET_1_9_0"`
	SET110 int `json:"SET_1_10_0"`

	// Disable
	// SET111        float64 `json:"GET_1_11"`
	// GET112        float64 `json:"GET_1_12"`
	// GET113        float64 `json:"GET_1_13"`
	// GET114        float64 `json:"GET_1_14"`
}

// struct to push controlWoodhouse
type agent2POST struct {
	MACAddress string `json:"MAC_Address" bson:"MAC_Address"`
	GWID       string `json:"GW_ID" bson:"GW_ID"`
}

// type JsonStruct struct {
// 	Data json.RawMessage
// }

//struct to store JSON response rawdata
type rawResponse struct {
	Properties struct {
		Value string `json:"value"`
	} `json:"properties"`
}

//struct to store Air conditioner JSON response rawdata
type rawResponseAirCond struct {
	Properties struct {
		UIACValue string `json:"ui.AC.value"`
	} `json:"properties"`
}

//airbox inside wooden house (Disable)
// type airboxRaw struct {
// 	PM25           float64   `json:"PM2_5" bson:"PM2_5"`
// 	LastReportTime time.Time `json:"Upload_Time" bson:"Upload_Time"`
// 	CO2            float64   `json:"CO2" bson:"CO2"`
// 	CO             float64   `json:"CO" bson:"CO"`
// 	Noise          float64   `json:"Noise" bson:"Noise"`
// 	Temp           float64   `json:"Temp" bson:"Temp"`
// 	Humidity       float64   `json:"Humidity" bson:"Humidity"`
// 	DevID          string    `json:"Device_ID" bson:"Device_ID"`
// }

//rule to check status
func checkBool(x string) bool {

	tmpVal, _ := strconv.ParseFloat(x, 32)
	// fmt.Println(tmpVal)
	// fmt.Println(x)
	if (tmpVal) > 5 {
		return true
	}
	return false

}

func checkState(x int) string {
	if x == 1 {
		return "On"
	}
	return "Off"
}

func checkStateAC(x int) string {
	if x == 1 {
		return "1"
	}
	return "2"
}

func convertBool(x bool) int {
	if x == true {
		return 1
	}
	return 0
}

func reverseState(x int) int {
	if x == 1 {
		return 0
	}
	return 1
}

func pushCmd(deviceID string, state string) string {
	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=turn" + state
	fmt.Println(url)
	req, err := http.Get(url)

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(req.Status, req.Header)
	body, _ := ioutil.ReadAll(req.Body)

	defer req.Body.Close()
	return string(body)
}

func pushCmdAC(deviceID string, argbv string) string {
	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/callAction?deviceID=" + deviceID + "&name=pressButton&arg1=" + argbv
	req, err := http.Get(url)
	fmt.Println(url)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(req.Status, req.Header)

	body, _ := ioutil.ReadAll(req.Body)

	defer req.Body.Close()
	return string(body)
}

func collectDataValue(deviceID string) bool {

	url := "http://m10513020@gapps.ntust.edu.tw:Ntust27333141@140.118.19.197:7288/api/devices/" + deviceID

	//get from url
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	//deviceID belongs to air conditioner
	if deviceID == "130" {

		container := rawResponseAirCond{}

		json.NewDecoder(response.Body).Decode(&container)
		defer response.Body.Close()
		dataval := checkBool(container.Properties.UIACValue)

		return dataval

	}

	container := rawResponse{}

	json.NewDecoder(response.Body).Decode(&container)
	defer response.Body.Close()

	dataval, _ := strconv.ParseBool(container.Properties.Value)

	return dataval
}

// Disable
// func getAirboxData() *airboxRaw {
// 	container := airboxRaw{}
// 	dbInfo := &mgo.DialInfo{
// 		Addrs:    strings.SplitN("140.118.70.136:27017", ",", -1),
// 		Database: "admin",
// 		Username: "root",
// 		Password: "123",
// 	}
// 	mongo, _ := mgo.DialWithInfo(dbInfo)
// 	// fmt.Println(err)
// 	mongo.DB("sc").C("airbox").Find(bson.M{"Device_ID": spaceAirboxdevID, "Upload_Time": bson.M{"$gt": time.Now().Add(time.Duration(-1) * time.Hour)}}).Sort("-Upload_Time").Limit(1).One(&container)

// 	if (container == airboxRaw{}) {
// 		container = airboxRaw{}
// 	}
// 	return &container
// }

// log confirmation
func confirm(x int, i int) string {
	return ("Success : " + strconv.Itoa(x) + " State : " + strconv.Itoa(i))
}

func monitorWoodhouse() {
	fmt.Println("Function :", "monitorWoodhouse")

	var indexdata []bool

	timestamp := time.Now()
	thetime := timestamp.Format("2006-01-02 15:04:05")
	timeunix := timestamp.Unix()
	// fmt.Println(timestamp)
	// container := agentPOST{}

	deviceID := []int{102, 95, 20, 27, 42, 51, 56, 61, 66, 130}

	for _, i := range deviceID {
		indexdata = append(indexdata, collectDataValue(strconv.Itoa(i)))
		// indexdata[h] = collectDataValue(strconv.Itoa(i))
		// fmt.Println(h, i)
	}

	// Disable
	// airboxcontainer := getAirboxData()

	container := agentPOST{
		Timestamp:     thetime,
		TimestampUnix: timeunix,
		MACAddress:    mac,
		GWID:          gwID,
		CPURate:       1.0,
		StorageRate:   1,
		GET11:         indexdata[0],
		GET12:         indexdata[1],
		GET13:         indexdata[2],
		GET14:         indexdata[3],
		GET15:         indexdata[4],
		GET16:         indexdata[5],
		GET17:         indexdata[6],
		GET18:         indexdata[7],
		GET19:         indexdata[8],
		GET110:        indexdata[9],

		SET11:  convertBool(indexdata[0]),
		SET12:  convertBool(indexdata[1]),
		SET13:  convertBool(indexdata[2]),
		SET14:  convertBool(indexdata[3]),
		SET15:  convertBool(indexdata[4]),
		SET16:  convertBool(indexdata[5]),
		SET17:  convertBool(indexdata[6]),
		SET18:  convertBool(indexdata[7]),
		SET19:  convertBool(indexdata[8]),
		SET110: convertBool(indexdata[9]),

		// Disable
		// GET111:        airboxcontainer.CO2,
		// GET112:        airboxcontainer.Humidity,
		// GET113:        airboxcontainer.PM25,
		// GET114:        airboxcontainer.Temp,
	}

	b := new(bytes.Buffer)

	json.NewEncoder(b).Encode(container)

	req, err := http.NewRequest("POST", "https://beta2-api.dforcepro.com/gateway/v1/rawdata", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Data :", req.Body)

	theclient := &http.Client{}
	resp, erro := theclient.Do(req)
	if err != nil {
		log.Println(erro)
	}

	fmt.Println(resp.Status, resp.Header)
	// fmt.Println(container)
	// res, _ := http.Post("https://beta2-api.dforcepro.com/gateway/v1/rawdata", contentType, b)
	body, _ := ioutil.ReadAll(resp.Body)
	// return string(body)
	defer resp.Body.Close()
	fmt.Println(string(body))
}

func controlWoodhouse() {
	fmt.Println("Function :", "controlWoodhouse")

	var indexdata []bool

	container := agent2POST{
		MACAddress: mac,
		GWID:       gwID,
	}

	deviceID := []int{102, 95, 20, 27, 42, 51, 56, 61, 66, 130}
	docID := []int{0, 0, 20, 27, 42, 51, 56, 61, 66, 76} // ID based on Doc
	multiconditionID := []int{5, 7, 38, 40}              //On, Off

	for _, i := range deviceID {
		indexdata = append(indexdata, collectDataValue(strconv.Itoa(i)))
	}

	b := new(bytes.Buffer)

	json.NewEncoder(b).Encode(container)

	req, err := http.NewRequest("POST", "https://beta2-api.dforcepro.com/gateway/v1/command", b)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
	}

	theclient := &http.Client{}

	resp, erro := theclient.Do(req)
	if err != nil {
		log.Println(erro)
	}

	bodys, _ := ioutil.ReadAll(resp.Body)

	if string(bodys) == "null\n" {
		return
	}
	fmt.Println(string(bodys))
	/////////////////
	var indexCmd []string
	var vi map[string]int

	// DEBUGGING
	// jsonstr, _ := json.Marshal(bson.M{"SET_1_1_0": 0, "SET_1_3_0": 0, "SET_1_7_0": 0, "SET_1_9_0": 0})
	json.Unmarshal(bodys, &vi)

	for each := range vi {
		indexCmd = strings.Split(each, "_")
		fmt.Println()
		n, _ := strconv.Atoi(indexCmd[2])
		n = n - 1
		if convertBool(indexdata[n]) != vi[each] {
			switch deviceID[n] {
			case 130:
				pushCmdAC(strconv.Itoa(docID[9]), checkStateAC(vi[each]))
				confirm(n, vi[each])
			case 102:
				if vi[each] == 1 {
					pushCmd(strconv.Itoa(multiconditionID[0]), "On")
					confirm(n, vi[each])
					continue
				} else {
					// fmt.Println(multiconditionID[1])
					pushCmd(strconv.Itoa(multiconditionID[1]), "On")
					confirm(n, vi[each])
				}
			case 95:
				if vi[each] == 1 {
					pushCmd(strconv.Itoa(multiconditionID[2]), "On")
					confirm(n, vi[each])

				} else {
					// fmt.Println(strconv.Itoa(multiconditionID[3]))
					pushCmd(strconv.Itoa(multiconditionID[3]), "On")
					confirm(n, vi[each])
				}
			default:
				// fmt.Println(strconv.Itoa(docID[i]))
				pushCmd(strconv.Itoa(docID[n]), checkState(vi[each]))
				confirm(n, vi[each])

			}
		} else {
			switch deviceID[n] {
			case 130:
				pushCmdAC(strconv.Itoa(docID[9]), checkStateAC(reverseState(vi[each])))
				confirm(n, reverseState(vi[each]))
			case 102:
				if reverseState(vi[each]) == 1 {
					pushCmd(strconv.Itoa(multiconditionID[0]), "On")
					confirm(n, reverseState(vi[each]))
					continue
				} else {
					// fmt.Println(multiconditionID[1])
					pushCmd(strconv.Itoa(multiconditionID[1]), "On")
					confirm(n, reverseState(vi[each]))
				}
			case 95:
				if reverseState(vi[each]) == 1 {
					pushCmd(strconv.Itoa(multiconditionID[2]), "On")
					confirm(n, reverseState(vi[each]))

				} else {
					// fmt.Println(strconv.Itoa(multiconditionID[3]))
					pushCmd(strconv.Itoa(multiconditionID[3]), "On")
					confirm(n, reverseState(vi[each]))
				}
			default:
				// fmt.Println(strconv.Itoa(docID[i]))
				pushCmd(strconv.Itoa(docID[n]), checkState(reverseState(vi[each])))
				confirm(n, reverseState(vi[each]))

			}
		}

	}
	defer resp.Body.Close()

}

func main() {

	c := cron.New()

	c.AddFunc("*/10 * * * * *", monitorWoodhouse)
	c.AddFunc("*/10 * * * * *", controlWoodhouse)

	c.Start()
	select {}

}
