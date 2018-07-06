package middle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dforcepro.com/util"
	"github.com/betacraft/yaag/middleware"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type LogMiddle bool

func (lm LogMiddle) Enable() bool {
	return bool(lm)
}

func (lm LogMiddle) GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_di.Log.Debug("-------Log-------")
			token := r.Header.Get("Token")
			if token != "" {
				system, ok := util.GetSysCode(r)
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Must set System in the header."))
					return
				}

				name := r.Header.Get("AuthName")
				account := r.Header.Get("AuthAccount")
				_di.Log.Debug(name + "-" + account)
				path, _ := mux.CurrentRoute(r).GetPathTemplate()
				path = fmt.Sprintf("%s,%s?%s", r.Method, path, r.URL.RawQuery)
				_di.Log.Debug(path)
				header, _ := json.Marshal(r.Header)
				_di.Log.Debug(string(header))
				if r.Method == "GET" {
					er := lm.InsertLog(name, account, system, string(header), "", path)
					if er != nil {
						_di.Log.Debug(er.Error())
					}
				} else {

					b := middleware.ReadBody(r)
					out, _ := json.Marshal(b)
					_di.Log.Debug(string(out))
					er := lm.InsertLog(name, account, system, string(header), string(out), path)
					if er != nil {
						_di.Log.Debug(er.Error())
					}
				}
			}
			_di.Log.Debug("-------End Log-------")
			f(w, r)
		}
	}
}

func (lm LogMiddle) InsertLog(user string, account string, sys string, header string, body string, path string) error {
	return getMongo().DB("Common").C("Log").Insert(bson.M{
		"name":      user,
		"account":   account,
		"sys":       sys,
		"header":    header,
		"path":      path,
		"body":      body,
		"timestamp": int32(time.Now().Unix()),
	})

}
