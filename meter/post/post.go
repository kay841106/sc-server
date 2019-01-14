package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"strings"
	"time"

	colors "github.com/fatih/color"
	"github.com/globalsign/mgo/bson"

	// mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	"github.com/globalsign/mgo"

	"github.com/gorilla/mux"
)

const (
	db        = "sc"
	db_airbox = "airbox"
	dblocal   = "172.16.0.132:27017"
	dbpublic  = "140.118.70.136:10003"
	dbleoass  = "140.118.123.95:27017"
	dbbackup  = "140.118.122.103:27017"
	// c            = "testing",
	c_lastreport = "lastreport"
	c_aemdra     = "aemdra"
	c_cpm        = "cpm"
	c_gw_status  = "gw_status"
	c_devices    = "devices"

	c_airboxraw = "airbox_raw"

	insertSuccess = "OK"
)

var bannr = `
Program name : POST API

maintainer   : avbee.lab@gmail.com

Date         : 21/10/2018

`

type session struct {
	theSess *mgo.Session
}

func (s *session) startSession() *session {
	return &session{s.theSess.Clone()}
}

// var session *mgo.Session

type dataInHourCPM struct {
	datetime   time.Time
	MACAddress string
	GWID       string
	hours      []CPMSnd
}

type dataInHourAEMDRA struct {
	Datetime   time.Time            `bson:"Timestamp"`
	MACAddress string               `bson:"MACAddress"`
	GWID       string               `bson:"GWID"`
	Hours      []dataInMinuteAEMDRA `bson:"hours"`
}

type dataInMinuteAEMDRA struct {
	Datetime time.Time   `bson:"Timestamp"`
	minutes  []AEMDRASnd `bson:"minutes"`
}

type CPMSnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string        `json:"GW_ID" bson:"GW_ID"`
	CPURate       float64       `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   int           `json:"Storage_rate" bson:"Storage_rate"`
	GET11         float64       `json:"ae_tot" bson:"ae_tot"`
	GET12         float64       `json:"wire" bson:"wire"`
	GET13         float64       `json:"freq" bson:"freq"`
	GET14         float64       `json:"ua" bson:"ua"`
	GET15         float64       `json:"ub" bson:"ub"`
	GET16         float64       `json:"uc" bson:"uc"`
	GET17         float64       `json:"u_avg" bson:"u_avg"`
	GET18         float64       `json:"uab" bson:"uab"`
	GET19         float64       `json:"ubc" bson:"ubc"`
	GET110        float64       `json:"uca" bson:"uca"`
	GET111        float64       `json:"uln_avg" bson:"uln_avg"`
	GET112        float64       `json:"ia" bson:"ia"`
	GET113        float64       `json:"ib" bson:"ib"`
	GET114        float64       `json:"ic" bson:"ic"`
	GET115        float64       `json:"i_avg" bson:"i_avg"`
	GET116        float64       `json:"pa" bson:"pa"`
	GET117        float64       `json:"pb" bson:"pb"`
	GET118        float64       `json:"pc" bson:"pc"`
	GET119        float64       `json:"p_sum" bson:"p_sum"`
	GET120        float64       `json:"sa" bson:"sa"`
	GET121        float64       `json:"sb" bson:"sb"`
	GET122        float64       `json:"sc" bson:"sc"`
	GET123        float64       `json:"s_sum" bson:"s_sum"`
	GET124        float64       `json:"pfa" bson:"pfa"`
	GET125        float64       `json:"pfb" bson:"pfb"`
	GET126        float64       `json:"pfc" bson:"pfc"`
	GET127        float64       `json:"pf_avg" bson:"pf_avg"`
	GET128        float64       `json:"uavg_thd" bson:"avg_thd"`
	GET129        float64       `json:"iavg_thd" bson:"iavg_thd"`
}

type CPMRcv struct {
	Timestamp     string  `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MACAddress    string  `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string  `json:"GW_ID" bson:"GW_ID"`
	CPURate       float64 `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   int     `json:"Storage_rate" bson:"Storage_rate"`
	GET11         float64 `json:"GET_1_1" bson:"GET_1_1"`
	GET12         float64 `json:"GET_1_2" bson:"GET_1_2"`
	GET13         float64 `json:"GET_1_3" bson:"GET_1_3"`
	GET14         float64 `json:"GET_1_4" bson:"GET_1_4"`
	GET15         float64 `json:"GET_1_5" bson:"GET_1_5"`
	GET16         float64 `json:"GET_1_6" bson:"GET_1_6"`
	GET17         float64 `json:"GET_1_7" bson:"GET_1_7"`
	GET18         float64 `json:"GET_1_8" bson:"GET_1_8"`
	GET19         float64 `json:"GET_1_9" bson:"GET_1_9"`
	GET110        float64 `json:"GET_1_10" bson:"GET_1_10"`
	GET111        float64 `json:"GET_1_11" bson:"GET_1_11"`
	GET112        float64 `json:"GET_1_12" bson:"GET_1_12"`
	GET113        float64 `json:"GET_1_13" bson:"GET_1_13"`
	GET114        float64 `json:"GET_1_14" bson:"GET_1_14"`
	GET115        float64 `json:"GET_1_15" bson:"GET_1_15"`
	GET116        float64 `json:"GET_1_16" bson:"GET_1_16"`
	GET117        float64 `json:"GET_1_17" bson:"GET_1_17"`
	GET118        float64 `json:"GET_1_18" bson:"GET_1_18"`
	GET119        float64 `json:"GET_1_19" bson:"GET_1_19"`
	GET120        float64 `json:"GET_1_20" bson:"GET_1_20"`
	GET121        float64 `json:"GET_1_21" bson:"GET_1_21"`
	GET122        float64 `json:"GET_1_22" bson:"GET_1_22"`
	GET123        float64 `json:"GET_1_23" bson:"GET_1_23"`
	GET124        float64 `json:"GET_1_24" bson:"GET_1_24"`
	GET125        float64 `json:"GET_1_25" bson:"GET_1_25"`
	GET126        float64 `json:"GET_1_26" bson:"GET_1_26"`
	GET127        float64 `json:"GET_1_27" bson:"GET_1_27"`
	GET128        float64 `json:"GET_1_28" bson:"GET_1_28"`
	GET129        float64 `json:"GET_1_29" bson:"GET_1_29"`
}

type AEMDRASnd struct {
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	Timestamp     time.Time     `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64         `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	MACAddress    string        `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string        `json:"GW_ID" bson:"GW_ID"`
	CPURate       float64       `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   int           `json:"Storage_rate" bson:"Storage_rate"`
	GET11         float64       `json:"ae_tot" bson:"ae_tot"`
	GET12         float64       `json:"wire" bson:"wire"`
	GET13         float64       `json:"freq" bson:"freq"`
	GET14         float64       `json:"ua" bson:"ua"`
	GET15         float64       `json:"ub" bson:"ub"`
	GET16         float64       `json:"uc" bson:"uc"`
	GET17         float64       `json:"u_avg" bson:"u_avg"`
	GET18         float64       `json:"uab" bson:"uab"`
	GET19         float64       `json:"ubc" bson:"ubc"`
	GET110        float64       `json:"uca" bson:"uca"`
	GET111        float64       `json:"uln_avg" bson:"uln_avg"`
	GET112        float64       `json:"ia" bson:"ia"`
	GET113        float64       `json:"ib" bson:"ib"`
	GET114        float64       `json:"ic" bson:"ic"`
	GET115        float64       `json:"i_avg" bson:"i_avg"`
	GET116        float64       `json:"pa" bson:"pa"`
	GET117        float64       `json:"pb" bson:"pb"`
	GET118        float64       `json:"pc" bson:"pc"`
	GET119        float64       `json:"p_sum" bson:"p_sum"`
	GET120        float64       `json:"qa" bson:"qa"`
	GET121        float64       `json:"qb" bson:"qb"`
	GET122        float64       `json:"qc" bson:"qc"`
	GET123        float64       `json:"q_sum" bson:"q_sum"`
	GET124        float64       `json:"sa" bson:"sa"`
	GET125        float64       `json:"sb" bson:"sb"`
	GET126        float64       `json:"sc" bson:"sc"`
	GET127        float64       `json:"s_sum" bson:"s_sum"`
	GET128        float64       `json:"pfa" bson:"pfa"`
	GET129        float64       `json:"pfb" bson:"pfb"`
	GET130        float64       `json:"pfc" bson:"pfc"`
	GET131        float64       `json:"pf_avg" bson:"pf_avg"`
	GET132        float64       `json:"aea" bson:"aea"`
	GET133        float64       `json:"aeb" bson:"aeb"`
	GET134        float64       `json:"aec" bson:"aec"`
	GET135        float64       `json:"rea" bson:"rea"`
	GET136        float64       `json:"reb" bson:"reb"`
	GET137        float64       `json:"rec" bson:"rec"`
	GET138        float64       `json:"re_tot" bson:"re_tot"`
}

type AEMDRARcv struct {
	Timestamp     string  `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MACAddress    string  `json:"MAC_Address" bson:"MAC_Address"`
	GWID          string  `json:"GW_ID" bson:"GW_ID"`
	CPURate       float64 `json:"CPU_rate" bson:"CPU_rate"`
	StorageRate   int     `json:"Storage_rate" bson:"Storage_rate"`
	GET11         float64 `json:"GET_1_1" bson:"GET_1_1"`
	GET12         float64 `json:"GET_1_2" bson:"GET_1_2"`
	GET13         float64 `json:"GET_1_3" bson:"GET_1_3"`
	GET14         float64 `json:"GET_1_4" bson:"GET_1_4"`
	GET15         float64 `json:"GET_1_5" bson:"GET_1_5"`
	GET16         float64 `json:"GET_1_6" bson:"GET_1_6"`
	GET17         float64 `json:"GET_1_7" bson:"GET_1_7"`
	GET18         float64 `json:"GET_1_8" bson:"GET_1_8"`
	GET19         float64 `json:"GET_1_9" bson:"GET_1_9"`
	GET110        float64 `json:"GET_1_10" bson:"GET_1_10"`
	GET111        float64 `json:"GET_1_11" bson:"GET_1_11"`
	GET112        float64 `json:"GET_1_12" bson:"GET_1_12"`
	GET113        float64 `json:"GET_1_13" bson:"GET_1_13"`
	GET114        float64 `json:"GET_1_14" bson:"GET_1_14"`
	GET115        float64 `json:"GET_1_15" bson:"GET_1_15"`
	GET116        float64 `json:"GET_1_16" bson:"GET_1_16"`
	GET117        float64 `json:"GET_1_17" bson:"GET_1_17"`
	GET118        float64 `json:"GET_1_18" bson:"GET_1_18"`
	GET119        float64 `json:"GET_1_19" bson:"GET_1_19"`
	GET120        float64 `json:"GET_1_20" bson:"GET_1_20"`
	GET121        float64 `json:"GET_1_21" bson:"GET_1_21"`
	GET122        float64 `json:"GET_1_22" bson:"GET_1_22"`
	GET123        float64 `json:"GET_1_23" bson:"GET_1_23"`
	GET124        float64 `json:"GET_1_24" bson:"GET_1_24"`
	GET125        float64 `json:"GET_1_25" bson:"GET_1_25"`
	GET126        float64 `json:"GET_1_26" bson:"GET_1_26"`
	GET127        float64 `json:"GET_1_27" bson:"GET_1_27"`
	GET128        float64 `json:"GET_1_28" bson:"GET_1_28"`
	GET129        float64 `json:"GET_1_29" bson:"GET_1_29"`
	GET130        float64 `json:"GET_1_30" bson:"GET_1_30"`
	GET131        float64 `json:"GET_1_31" bson:"GET_1_31"`
	GET132        float64 `json:"GET_1_32" bson:"GET_1_32"`
	GET133        float64 `json:"GET_1_33" bson:"GET_1_33"`
	GET134        float64 `json:"GET_1_34" bson:"GET_1_34"`
	GET135        float64 `json:"GET_1_35" bson:"GET_1_35"`
	GET136        float64 `json:"GET_1_36" bson:"GET_1_36"`
	GET137        float64 `json:"GET_1_37" bson:"GET_1_37"`
	GET138        float64 `json:"GET_1_38" bson:"GET_1_38"`
}

type aemdraMetrics struct {
	GET11  float64 `json:"ae_tot" bson:"ae_tot"`
	GET12  float64 `json:"wire" bson:"wire"`
	GET13  float64 `json:"freq" bson:"freq"`
	GET14  float64 `json:"ua" bson:"ua"`
	GET15  float64 `json:"ub" bson:"ub"`
	GET16  float64 `json:"uc" bson:"uc"`
	GET17  float64 `json:"u_avg" bson:"u_avg"`
	GET18  float64 `json:"uab" bson:"uab"`
	GET19  float64 `json:"ubc" bson:"ubc"`
	GET110 float64 `json:"uca" bson:"uca"`
	GET111 float64 `json:"uln_avg" bson:"uln_avg"`
	GET112 float64 `json:"ia" bson:"ia"`
	GET113 float64 `json:"ib" bson:"ib"`
	GET114 float64 `json:"ic" bson:"ic"`
	GET115 float64 `json:"i_avg" bson:"i_avg"`
	GET116 float64 `json:"pa" bson:"pa"`
	GET117 float64 `json:"pb" bson:"pb"`
	GET118 float64 `json:"pc" bson:"pc"`
	GET119 float64 `json:"p_sum" bson:"p_sum"`
	GET120 float64 `json:"qa" bson:"qa"`
	GET121 float64 `json:"qb" bson:"qb"`
	GET122 float64 `json:"qc" bson:"qc"`
	GET123 float64 `json:"q_sum" bson:"q_sum"`
	GET124 float64 `json:"sa" bson:"sa"`
	GET125 float64 `json:"sb" bson:"sb"`
	GET126 float64 `json:"sc" bson:"sc"`
	GET127 float64 `json:"s_sum" bson:"s_sum"`
	GET128 float64 `json:"pfa" bson:"pfa"`
	GET129 float64 `json:"pfb" bson:"pfb"`
	GET130 float64 `json:"pfc" bson:"pfc"`
	GET131 float64 `json:"pf_avg" bson:"pf_avg"`
	GET132 float64 `json:"aea" bson:"aea"`
	GET133 float64 `json:"aeb" bson:"aeb"`
	GET134 float64 `json:"aec" bson:"aec"`
	GET135 float64 `json:"rea" bson:"rea"`
	GET136 float64 `json:"reb" bson:"reb"`
	GET137 float64 `json:"rec" bson:"rec"`
	GET138 float64 `json:"re_tot" bson:"re_tot"`
}
type gwstat struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
}

type cpmMetrics struct {
	GET11  float64 `json:"ae_tot" bson:"ae_tot"`
	GET12  float64 `json:"wire" bson:"wire"`
	GET13  float64 `json:"freq" bson:"freq"`
	GET14  float64 `json:"ua" bson:"ua"`
	GET15  float64 `json:"ub" bson:"ub"`
	GET16  float64 `json:"uc" bson:"uc"`
	GET17  float64 `json:"u_avg" bson:"u_avg"`
	GET18  float64 `json:"uab" bson:"uab"`
	GET19  float64 `json:"ubc" bson:"ubc"`
	GET110 float64 `json:"uca" bson:"uca"`
	GET111 float64 `json:"uln_avg" bson:"uln_avg"`
	GET112 float64 `json:"ia" bson:"ia"`
	GET113 float64 `json:"ib" bson:"ib"`
	GET114 float64 `json:"ic" bson:"ic"`
	GET115 float64 `json:"i_avg" bson:"i_avg"`
	GET116 float64 `json:"pa" bson:"pa"`
	GET117 float64 `json:"pb" bson:"pb"`
	GET118 float64 `json:"pc" bson:"pc"`
	GET119 float64 `json:"p_sum" bson:"p_sum"`
	GET120 float64 `json:"sa" bson:"sa"`
	GET121 float64 `json:"sb" bson:"sb"`
	GET122 float64 `json:"sc" bson:"sc"`
	GET123 float64 `json:"s_sum" bson:"s_sum"`
	GET124 float64 `json:"pfa" bson:"pfa"`
	GET125 float64 `json:"pfb" bson:"pfb"`
	GET126 float64 `json:"pfc" bson:"pfc"`
	GET127 float64 `json:"pf_avg" bson:"pf_avg"`
	GET128 float64 `json:"uavg_thd" bson:"avg_thd"`
	GET129 float64 `json:"iavg_thd" bson:"iavg_thd"`
}

type cpmlastreport struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
	// Status        bool      `json:"Status" bson:"Status"`
	MACAddress string     `json:"MAC_Address" bson:"MAC_Address"`
	Metrics    cpmMetrics `json:"metrics" bson:"metrics"`
}

type aemdralastreport struct {
	Timestamp     time.Time `json:"Timestamp" bson:"Timestamp"`
	TimestampUnix int64     `json:"Timestamp_Unix" bson:"Timestamp_Unix"`
	GWID          string    `json:"GW_ID" bson:"GW_ID"`
	// Status        bool      `json:"Status" bson:"Status"`
	MACAddress string        `json:"MAC_Address" bson:"MAC_Address"`
	Metrics    aemdraMetrics `json:"metrics" bson:"metrics"`
}

type structmetrics struct {
	getMtrcs float64 `json:"ae_tot" bson:"ae_tot"`
}

type getmetrics struct {
	Metrics structmetrics `json:"metrics" bson:"metrics"`
}

func (s *session) GWAuth(gwid string) bool {

	sess := s.startSession().theSess

	slice := []string{}

	Mongo := sess.DB(db)
	defer sess.Close()
	Mongo.C(c_devices).Find(bson.M{}).Distinct("GWID", &slice)

	m := make(map[string]bool)
	for i := 0; i < len(slice); i++ {
		if slice[i] != "" {
			m[slice[i][0:8]] = true
			// fmt.Println(m[slice[i][0:8]])
		}

		// fmt.Println(slice)
	}

	// // Error because node mac is not registered, will ignore >8 char
	// if _, ok := m[gwid]; ok {
	// 	return true
	// }
	// return false
	if gwid != "" {

		if _, ok := m[gwid[0:8]]; ok {
			fmt.Print(gwid)
			return true
		}

	}
	return false
}

func (s *session) MACAuth(macid string) bool {

	slice := []string{}
	// var i interface{}
	// sess := session.Clone()
	sess := s.startSession().theSess

	Mongo := sess.DB(db)
	defer sess.Close()
	Mongo.C(c_devices).Find(bson.M{}).Distinct("MAC_Address", &slice)
	// fmt.Println(i)
	m := make(map[string]bool)
	for i := 0; i < len(slice); i++ {
		m[slice[i]] = true
	}
	// for _,a:=range
	// if contains()

	if _, ok := m[macid]; ok {
		return true
	}

	return false
}

func checkTime(x int64) bool {
	switch time.Unix(x, 0).Local().Hour() == time.Now().Local().Hour() {
	case false:
		return false
	default:
		return true
	}
}

func syncTime(x int64) time.Time {

	if checkTime(x) != true {
		log.Println("##Incorrect Time Format##")

		// Previously -1 day in UTC format. Data still +8 format
		t := time.Unix(x, 0).UTC()

		year, month, day := t.Date()
		hour, min, sec := t.Clock()
		return time.Date(year, month, day+1, hour, min, sec, 0, time.UTC)

	}
	t := time.Unix(x, 0)
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)

}

func convertHour(t int64) time.Time {
	times := time.Unix(t, 0).UTC()
	year, month, day := times.Date()
	hour, _, _ := times.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

func statuscheck(thetime int64) bool {
	if thetime-time.Now().Unix() < 3600 {
		return true
	}
	return false
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

	// fmt.Println(theid, uint32(timestamp))
	return theid

}

// func recalcUnix(unixtime int64) int64 {
// 	return unixtime - (8 * time.Hour.Nanoseconds() / int64(time.Second))
// }

func (s *session) aemdraPost(w http.ResponseWriter, r *http.Request) {
	i := 60

	sess := s.startSession().theSess
	err := sess.Ping()

	// even without Ping, code can run well
	for err != nil {
		err = sess.Ping()
		log.Println("Waiting for DB to Online .... ", i)
		time.Sleep(time.Second * 5)
		i--
		if i <= 0 {
			os.Exit(1)
		}

	}

	// sess := session.Clone()
	// defer sess.Close()

	Mongo := sess.DB(db) // mongo := client.DB(db).C(c)

	defer sess.Close()
	container := AEMDRARcv{}

	json.NewDecoder(r.Body).Decode(&container)
	// defer r.Body.Close()

	containerSnd := &AEMDRASnd{

		ID:            getObjectIDTwoArg(container.GWID, container.MACAddress, container.TimestampUnix),
		TimestampUnix: container.TimestampUnix,
		Timestamp:     time.Unix(container.TimestampUnix, 0).UTC(),
		//
		MACAddress: container.MACAddress,
		GWID:       container.GWID,

		//
		CPURate:     container.CPURate,
		StorageRate: container.StorageRate,
		GET11:       container.GET11,
		GET12:       container.GET12,
		GET13:       container.GET13,
		GET14:       container.GET14,
		GET15:       container.GET15,
		GET16:       container.GET16,
		GET17:       container.GET17,
		GET18:       container.GET18,
		GET19:       container.GET19,
		GET110:      container.GET110,
		GET111:      container.GET111,
		GET112:      container.GET112,
		GET113:      container.GET113,
		GET114:      container.GET114,
		GET115:      container.GET115,
		GET116:      container.GET116,
		GET117:      container.GET117,
		GET118:      container.GET118,
		GET119:      container.GET119,
		GET120:      container.GET120,
		GET121:      container.GET121,
		GET122:      container.GET122,
		GET123:      container.GET123,
		GET124:      container.GET124,
		GET125:      container.GET125,
		GET126:      container.GET126,
		GET127:      container.GET127,
		GET128:      container.GET128,
		GET129:      container.GET129,
		GET130:      container.GET130,
		GET131:      container.GET131,
		GET132:      container.GET132,
		GET133:      container.GET133,
		GET134:      container.GET134,
		GET135:      container.GET135,
		GET136:      container.GET136,
		GET137:      container.GET137,
		GET138:      container.GET138,
	}

	lastreportmetrics := aemdraMetrics{
		GET11:  container.GET11,
		GET12:  container.GET12,
		GET13:  container.GET13,
		GET14:  container.GET14,
		GET15:  container.GET15,
		GET16:  container.GET16,
		GET17:  container.GET17,
		GET18:  container.GET18,
		GET19:  container.GET19,
		GET110: container.GET110,
		GET111: container.GET111,
		GET112: container.GET112,
		GET113: container.GET113,
		GET114: container.GET114,
		GET115: container.GET115,
		GET116: container.GET116,
		GET117: container.GET117,
		GET118: container.GET118,
		GET119: container.GET119,
		GET120: container.GET120,
		GET121: container.GET121,
		GET122: container.GET122,
		GET123: container.GET123,
		GET124: container.GET124,
		GET125: container.GET125,
		GET126: container.GET126,
		GET127: container.GET127,
		GET128: container.GET128,
		GET129: container.GET129,
		GET130: container.GET130,
		GET131: container.GET131,
		GET132: container.GET132,
		GET133: container.GET133,
		GET134: container.GET134,
		GET135: container.GET135,
		GET136: container.GET136,
		GET137: container.GET137,
		GET138: container.GET138,
	}
	// fmt.Println("init:", time.Unix(container.TimestampUnix, 0).UTC(), "crc:", time.Unix(recalcUnix(container.TimestampUnix), 0).UTC())
	if s.GWAuth(containerSnd.GWID) == true {
		// fmt.Print(containerSnd.MACAddress)
		GWStatuscontainer := gwstat{
			Timestamp:     time.Now().UTC(),
			TimestampUnix: time.Now().Unix(),
			GWID:          containerSnd.GWID[0:8],
			// Status:        statuscheck(containerSnd.TimestampUnix),
		}
		// update gwstatus
		err = Mongo.C(c_gw_status).Update(bson.M{"GW_ID": containerSnd.GWID[0:8]}, bson.M{"$set": GWStatuscontainer})
		if err != nil {
			mylog(container.GWID, container.MACAddress, "")
			fmt.Println(err)

			return
		}

		if s.MACAuth(containerSnd.MACAddress) == true {

			// Mongo.C(c_lastreport).Upsert(bson.M{"MAC_Address": containerSnd.MACAddress}, containerSnd)

			// Mongo.C(c_gw_status).Upsert(bson.M{"GW_ID": containerSnd.GWID}, GWStatuscontainer)
			// update cpm rawdata
			err := Mongo.C(c_aemdra).Insert(containerSnd)
			if err != nil {
				mylog(container.GWID, container.MACAddress, err.Error())

				// json.NewEncoder(w).Encode(err)

			} else {
				log.Println(containerSnd.MACAddress, insertSuccess)
			}

			json.NewEncoder(w).Encode(containerSnd)

			// containerSnd.ID = bson.NewObjectId()
			// update lastreport

			getthemetrics := getmetrics{}

			Lastreportcontainer := aemdralastreport{
				Timestamp:     time.Unix(containerSnd.TimestampUnix, 0).UTC(),
				TimestampUnix: containerSnd.TimestampUnix,
				GWID:          containerSnd.GWID[0:8],
				// Status:        statuscheck(containerSnd.TimestampUnix),
				MACAddress: containerSnd.MACAddress,
				Metrics:    lastreportmetrics,
			}

			err = Mongo.C(c_lastreport).Find(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}).One(&getthemetrics)
			if err != nil {
				// fmt.Println("false")
				json.NewEncoder(w).Encode(err)
			}

			Lastreportcontainer.Metrics.GET11 = Lastreportcontainer.Metrics.GET11 - getthemetrics.Metrics.getMtrcs
			err = Mongo.C(c_lastreport).Update(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}, bson.M{"$set": Lastreportcontainer})
			if err != nil {
				mylog(container.GWID, container.MACAddress, err.Error())

			}
			// defer sess.Close()

			// DB BACKUP
			// sess2 := db_connect2()
			// Mongo2 := sess2.DB(db)
			// err = Mongo2.C(c_aemdra).Insert(containerSnd)
			// err = Mongo2.C(c_gw_status).Update(bson.M{"GW_ID": containerSnd.GWID[0:8]}, bson.M{"$set": GWStatuscontainer})
			// err = Mongo2.C(c_lastreport).Update(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}, bson.M{"$set": Lastreportcontainer})
			// json.NewEncoder(w).Encode(containerSnd)
		} else {

			mylog(container.GWID, container.MACAddress, "MAC cannot found!")
			json.NewEncoder(w).Encode(container.MACAddress + "  MAC cannot found!")
		}

	} else {
		mylog(container.GWID, container.MACAddress, "MAC cannot found!")
		json.NewEncoder(w).Encode(container.GWID + "  GW cannot found!")

	}

}

func (s *session) cpmPost(w http.ResponseWriter, r *http.Request) {
	i := 60
	// bodyBytes, err2 := ioutil.ReadAll(r.Body)
	// if err2 != nil {
	// 	fmt.Print(err2)
	// } else {
	// 	fmt.Print(string(bodyBytes))
	// }
	sess := s.startSession().theSess
	err := sess.Ping()

	for err != nil {
		err = sess.Ping()
		log.Println("Waiting for DB to Online .... ", i)
		time.Sleep(time.Second * 5)
		i--
		if i <= 0 {
			os.Exit(1)
		}

	}

	Mongo := sess.DB(db)
	defer sess.Close()
	container := &CPMRcv{}

	json.NewDecoder(r.Body).Decode(&container)
	// defer r.Body.Close()

	containerSnd := &CPMSnd{

		ID:            getObjectIDTwoArg(container.GWID, container.MACAddress, container.TimestampUnix),
		TimestampUnix: container.TimestampUnix,
		Timestamp:     time.Unix(container.TimestampUnix, 0).UTC(),
		//
		MACAddress: container.MACAddress,
		GWID:       container.GWID,

		//
		CPURate:     container.CPURate,
		StorageRate: container.StorageRate,
		GET11:       container.GET11,
		GET12:       container.GET12,
		GET13:       container.GET13,
		GET14:       container.GET14,
		GET15:       container.GET15,
		GET16:       container.GET16,
		GET17:       container.GET17,
		GET18:       container.GET18,
		GET19:       container.GET19,
		GET110:      container.GET110,
		GET111:      container.GET111,
		GET112:      container.GET112,
		GET113:      container.GET113,
		GET114:      container.GET114,
		GET115:      container.GET115,
		GET116:      container.GET116,
		GET117:      container.GET117,
		GET118:      container.GET118,
		GET119:      container.GET119,
		GET120:      container.GET120,
		GET121:      container.GET121,
		GET122:      container.GET122,
		GET123:      container.GET123,
		GET124:      container.GET124,
		GET125:      container.GET125,
		GET126:      container.GET126,
		GET127:      container.GET127,
		GET128:      container.GET128,
		GET129:      container.GET129,
	}

	lastreportmetrics := cpmMetrics{
		GET11:  container.GET11,
		GET12:  container.GET12,
		GET13:  container.GET13,
		GET14:  container.GET14,
		GET15:  container.GET15,
		GET16:  container.GET16,
		GET17:  container.GET17,
		GET18:  container.GET18,
		GET19:  container.GET19,
		GET110: container.GET110,
		GET111: container.GET111,
		GET112: container.GET112,
		GET113: container.GET113,
		GET114: container.GET114,
		GET115: container.GET115,
		GET116: container.GET116,
		GET117: container.GET117,
		GET118: container.GET118,
		GET119: container.GET119,
		GET120: container.GET120,
		GET121: container.GET121,
		GET122: container.GET122,
		GET123: container.GET123,
		GET124: container.GET124,
		GET125: container.GET125,
		GET126: container.GET126,
		GET127: container.GET127,
		GET128: container.GET128,
		GET129: container.GET129,
	}

	// fmt.Println(containerSnd)
	if s.GWAuth(containerSnd.GWID) == true {

		GWStatuscontainer := gwstat{
			Timestamp:     time.Now().UTC(),
			TimestampUnix: time.Now().Unix(),
			GWID:          containerSnd.GWID[0:8],
		}

		err = Mongo.C(c_gw_status).Update(bson.M{"GW_ID": containerSnd.GWID[0:8]}, bson.M{"$set": GWStatuscontainer})

		if err != nil {
			mylog(container.GWID, container.MACAddress, err.Error())
			return
		}

		if s.MACAuth(containerSnd.MACAddress) == true {
			// update cpm rawdata
			err := Mongo.C(c_cpm).Insert(containerSnd)

			if err != nil {
				mylog(container.GWID, container.MACAddress, err.Error())

				// json.NewEncoder(w).Encode(err)
			} else {
				log.Println(containerSnd.MACAddress, insertSuccess)
			}
			json.NewEncoder(w).Encode(containerSnd)

			// declare the struct
			getthemetrics := getmetrics{}

			// update lastreport
			Lastreportcontainer := cpmlastreport{
				Timestamp:     time.Unix(containerSnd.TimestampUnix, 0).UTC(),
				TimestampUnix: containerSnd.TimestampUnix,
				GWID:          containerSnd.GWID[0:8],

				MACAddress: containerSnd.MACAddress,
				Metrics:    lastreportmetrics,
			}

			// get previous metric from db
			err = Mongo.C(c_lastreport).Find(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}).One(&getthemetrics)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			}

			// get delta metrics and set to the container
			Lastreportcontainer.Metrics.GET11 = Lastreportcontainer.Metrics.GET11 - getthemetrics.Metrics.getMtrcs
			err = Mongo.C(c_lastreport).Update(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}, bson.M{"$set": Lastreportcontainer})
			if err != nil {
				mylog(container.GWID, container.MACAddress, "  ")
				fmt.Println(err)

			}
			// defer sess.Close()

			// DB BACKUP
			// sess2 := db_connect2()
			// Mongo2 := sess2.DB(db)
			// err = Mongo2.C(c_cpm).Insert(containerSnd)
			// err = Mongo2.C(c_gw_status).Update(bson.M{"GW_ID": containerSnd.GWID[0:8]}, bson.M{"$set": GWStatuscontainer})
			// err = Mongo2.C(c_lastreport).Update(bson.M{"MAC_Address": Lastreportcontainer.MACAddress}, bson.M{"$set": Lastreportcontainer})
			// json.NewEncoder(w).Encode(containerSnd)
		} else {
			mylog(container.GWID, container.MACAddress, "MAC cannot found!")
			json.NewEncoder(w).Encode(container.MACAddress + " MAC cannot found!")
		}

	} else {
		mylog(container.GWID, container.MACAddress, "MAC cannot found!")
		json.NewEncoder(w).Encode(container.GWID + " GW cannot found!")

	}

}

func db_connect() *mgo.Session {

	inpEnv := getenvvar()

	dbInfo := &mgo.DialInfo{
		Addrs:    strings.SplitN(inpEnv.Mongo, ",", -1),
		Database: inpEnv.Database,
		Username: inpEnv.Username,
		Password: inpEnv.Password,
		Timeout:  time.Second * 10,
	}

	sess, err := mgo.DialWithInfo(dbInfo)
	if err != nil {
		os.Exit(1)
	}
	return sess
}

func mylog(x string, y string, z string) {
	thelog := "GW: " + x + "|  MAC: " + y + "|  STATUS: " + colors.RedString(z)
	log.Println(thelog)

}

type ENV struct {
	Port     string
	Mongo    string
	Database string
	Username string
	Password string
}

// var mongo_url string = Getenv().Mongo
// var database string = Getenv().Database
// var username string = Getenv().Username
// var password string = Getenv().Password

func getenvvar() ENV {

	theenv := ENV{}
	theenv.Mongo = os.Getenv("MONGO_URL")
	if theenv.Mongo == "" {
		theenv.Mongo = "172.16.0.132:27017"
	}
	theenv.Database = os.Getenv("DB_AUTH")
	if theenv.Database == "" {
		theenv.Database = "admin"
	}
	theenv.Username = os.Getenv("USERNAME")
	if theenv.Username == "" {
		theenv.Username = "dontask"
	}
	theenv.Password = os.Getenv("PASS")
	if theenv.Password == "" {
		theenv.Password = "idontknow"
	}

	return theenv
}

func main() {
	fmt.Print(bannr)
	sess := db_connect()

	v := session{sess}

	router := mux.NewRouter()
	router.HandleFunc("/meter/aemdra", v.aemdraPost).Methods("POST")
	router.HandleFunc("/meter/cpm", v.cpmPost).Methods("POST")

	// srv := &http.Server{
	// 	Handler: router,
	// 	Addr:    ":8080",
	// 	// Good practice: enforce timeouts for servers you create!
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }
	// http2.ConfigureServer(srv, nil)
	log.Fatal(http.ListenAndServe(":8080", router))

	// log.Println(srv.ListenAndServe())

}
