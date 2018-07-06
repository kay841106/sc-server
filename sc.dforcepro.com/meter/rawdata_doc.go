package meter

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"dforcepro.com/api"
	"gopkg.in/mgo.v2/bson"
)

type SCConfAPI bool

func (sca SCConfAPI) Enable() bool {
	return bool(sca)
}

func (sca SCConfAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/v1/rawdata", Next: sca.gEndpoint, Method: "GET", Auth: false},
	}
}
func (sca SCConfAPI) gEndpoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)

	queries := req.URL.Query()

	var results []interface{}
	mongo := getMongo()
	query := mongo.DB(DBName).C("SC01_RawData_AEMDRA_IB").Find(bson.M{})
	total, err := query.Count()
	var limit, page, totalPage = 100, 1, 1

	if err != nil {
		_di.Log.Err(err.Error())

	} else if total != 0 {
		limitStrAry, ok := queries["limit"]
		if ok {
			_limit, err := strconv.Atoi(limitStrAry[0])
			if err == nil {
				limit = _limit
			}
			if limit > 100 || limit < 1 {
				limit = 100
			}
		}

		pageStrAry, ok := queries["page"]
		if ok {
			_page, err := strconv.Atoi(pageStrAry[0])
			if err == nil {
				page = _page
			}
		}
		totalPage = int(math.Ceil(float64(total) / float64(limit)))

		if page > totalPage {
			page = totalPage
		} else if page < 1 {
			page = 1
		}
		query.Limit(limit).Skip(page - 1).All(&results)

	}
	responseJSON := onlyRes{&results}
	json.NewEncoder(w).Encode(responseJSON)
	_afterEndPoint(w, req)
}

func _beforeEndPoint(w http.ResponseWriter, req *http.Request) {
	// 檢查 DI 是否正確設置
	if _di == nil {
		log.Fatal(errors.New("DI doesn't set. Use func SetDI first"))
	}
}
