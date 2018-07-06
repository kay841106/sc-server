package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"dforcepro.com/api"
	"dforcepro.com/util"
	"github.com/gorilla/mux"
)

const (
	RegionTable = "Region"
)

type RegionAPI bool

func (r RegionAPI) createEndpoint(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	code := vars["code"]
	region := Region{}
	_ = json.NewDecoder(req.Body).Decode(&region)
	region.ParentCode = code
	region.Enable = true

	db, err := _di.SQL.GetMySQLConn(Database)
	if err != nil {
		_di.Log.Err(err.Error())
		return
	}
	defer db.Close()
	_, ok := region.Insert(db)

	if ok {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r RegionAPI) getEndpoint(w http.ResponseWriter, req *http.Request) {

	queries := req.URL.Query()
	fmt.Println(queries)
	ids, ok := queries["ids"]
	if !ok {
		vars := mux.Vars(req)
		code := vars["code"]
		if "" == code {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ids = append(ids, code)
	}
	db, err := _di.SQL.GetMySQLConn(Database)
	if err != nil {
		_di.Log.Err(err.Error())
		return
	}
	defer db.Close()

	result := GetSubRegion(db, ids...)

	if len(*result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("not found"))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func (r RegionAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/country", Next: r.createEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/country/{code:[a-z]+}/city", Next: r.createEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/city/{code:[a-z]+}/town", Next: r.createEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/city/{code:[a-z]+}/town", Next: r.getEndpoint, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/city", Next: r.getEndpoint, Method: "GET", Auth: false},
	}

}

func (r RegionAPI) Enable() bool {
	return bool(r)
}

type Region struct {
	ID         int
	Name       string `json:"name,omitempty"`
	Sort       int16
	Code       string `json:"code,omitempty"`
	Enable     bool
	ParentCode string
}

type RegionResponse struct {
	Town_ID   string
	Town_Name string
}

func (r Region) Insert(db *sql.DB) (int64, bool) {
	var result sql.Result
	var err error
	if db == nil {
		_di.Log.Err("db is nil.")
		return 0, false
	}
	sqlTpl := "INSERT INTO %s (`name`,`code`,`sort`,`enable`,`parent_code`) %s"
	if r.ParentCode == "" {
		var subQuery = fmt.Sprintf("SELECT ?,?,count(*)+1,?,? FROM %s WHERE `enable` = 1 AND `parent_code` IS NULL", RegionTable)
		fmt.Println(fmt.Sprintf(sqlTpl, RegionTable, subQuery), r.Name, r.Code, r.Enable, nil)
		result, err = db.Exec(fmt.Sprintf(sqlTpl, RegionTable, subQuery), r.Name, r.Code, r.Enable, nil)
	} else {
		var subQuery = fmt.Sprintf("SELECT ?,?,count(*)+1,?,? FROM %s WHERE `enable` = 1 AND `parent_code` = '%s'", RegionTable, r.ParentCode)
		result, err = db.Exec(fmt.Sprintf(sqlTpl, RegionTable, subQuery), r.Name, r.Code, r.Enable, r.ParentCode)
	}
	fmt.Println("insert")
	if err != nil {
		_di.Log.Err(err.Error())
		return 0, false
	}
	id, err := result.LastInsertId()
	if err != nil {
		_di.Log.Err(err.Error())
	}
	_di.Log.Info(string(id))
	return id, true
}

func GetSubRegion(db *sql.DB, codes ...string) *[]interface{} {
	fmt.Println("Query")
	if db == nil {
		_di.Log.Err("db is nil.")
		return nil
	}
	condition := util.JoinStrWithQuotation(util.SymbolComma, util.SymbolSingleQuotation, codes...)
	sqlTpl := "SELECT `code`, `name` FROM `Region` WHERE `enable` = 1 AND `parent_code` in (%s)"
	querySQL := fmt.Sprintf(sqlTpl, condition)
	rows, err := db.Query(querySQL)
	if err != nil {
		_di.Log.Debug(err.Error())
		return nil
	}
	defer rows.Close()
	var results []interface{}
	for rows.Next() {

		data := RegionResponse{}
		if err := rows.Scan(&data.Town_ID, &data.Town_Name); err != nil {
			_di.Log.Err(err.Error())
			return &results
		}
		results = append(results, data)
	}
	return &results
}
