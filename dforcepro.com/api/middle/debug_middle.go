package middle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betacraft/yaag/middleware"
	"github.com/gorilla/mux"
)

type DebugMiddle bool

func (lm DebugMiddle) Enable() bool {
	return bool(lm)
}

func (lm DebugMiddle) GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			_di.Log.Debug("-------Debug Request-------")
			path, _ := mux.CurrentRoute(r).GetPathTemplate()
			path = fmt.Sprintf("%s,%s?%s", r.Method, path, r.URL.RawQuery)
			_di.Log.Debug("path: " + path)
			header, _ := json.Marshal(r.Header)
			_di.Log.Debug("header: " + string(header))
			b := middleware.ReadBody(r)
			out, _ := json.Marshal(b)
			_di.Log.Debug("body: " + string(out))
			_di.Log.Debug("-------End Debug Request-------")

			// _, err := getRedis(0).GetClient().Ping().Result()
			// if err != nil {
			// 	_di.Log.Err(err.Error())
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	w.Write([]byte(err.Error()))
			// 	return
			// }

			err := getMongo().Ping()
			if err != nil {
				getMongo().Refresh()
				_di.Log.Err(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
				w.Write([]byte(err.Error()))
				return
			}
			f(w, r)
		}
	}
}
