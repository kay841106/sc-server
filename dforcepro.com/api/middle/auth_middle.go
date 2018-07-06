package middle

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"dforcepro.com/resource/db"
	"dforcepro.com/util"
	"github.com/gorilla/mux"
)

type AuthMiddle bool

const (
	PenddingMinute = 24 * 60 //閒置自動登出時間，單位分鐘
)

var (
	authMap map[string]bool = make(map[string]bool)
)

func AddAuthPath(path string, auth bool) {
	authMap[path] = auth
}

func IsAuth(path string, method string) bool {
	key := fmt.Sprintf("%s:%s", path, method)
	auth, ok := authMap[key]

	if ok {
		return auth
	}
	return false
}

func (am AuthMiddle) Enable() bool {
	return bool(am)
}

func (am AuthMiddle) GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// one time scope setup area for middleware
		return func(w http.ResponseWriter, r *http.Request) {
			// TODO: remove thie back door
			admin := r.Header.Get("dforcegod")
			// ... pre handler functionality
			path, err := mux.CurrentRoute(r).GetPathTemplate()
			_di.Log.Debug("curentRout: " + path)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
			isAuth := IsAuth(path, r.Method)
			_di.Log.Debug(fmt.Sprintf("Path is Auth: %t", isAuth))
			if isAuth && admin != "god" {
				isLogin, err := checkLogin(r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(err.Error()))
					return
				}
				if !isLogin {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				r.Header.Set("isLogin", "true")
				userMap := *(util.DecodeToken(r))

				if userMap == nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("taken error"))
					return
				}

				_di.Log.Debug(fmt.Sprintf("%v", userMap))
				sysCode, _ := util.GetSysCode(r)
				tokenSys := userMap["system"].(string)
				if sysCode != tokenSys {
					w.WriteHeader(http.StatusUnauthorized)
					_di.Log.Debug(fmt.Sprintf("Toke system is %s. Header system is %s", tokenSys, sysCode))
					w.Write([]byte("taken error"))
					return
				}

				r.Header.Set("AuthID", userMap["id"].(string))
				r.Header.Set("AuthAccount", userMap["account"].(string))
				r.Header.Set("AuthName", userMap["name"].(string))
				r.Header.Set("AuthGroup", userMap["group"].(string))
			} else if admin == "god" {
				r.Header.Set("AuthID", "")
				r.Header.Set("AuthAccount", "god")
				r.Header.Set("AuthName", "god")
				r.Header.Set("AuthGroup", "admin")
			}

			f(w, r)

			setToken := r.Header.Get("SET_TOKEN")

			if setToken == "" {
				return
			}
			redisClient := getRedis(db.RedisDB_Token).GetClient()
			device := util.GetClientKey(r)
			err = redisClient.SAdd(setToken, device).Err()

			if err != nil {
				_di.Log.Err(err.Error())
			}

			if PenddingMinute > 0 {
				redisClient.Expire(setToken, PenddingMinute*time.Minute)
			}
		}
	}
}

func checkLogin(req *http.Request) (bool, error) {
	token := req.Header.Get("token")
	if token == "" {
		return false, nil
	}
	redisClient := getRedis(db.RedisDB_Token).GetClient()
	count := redisClient.SCard(token).Val()
	device := util.GetClientKey(req)
	if count == 0 {
		// 從未登入過
		return false, nil
	}

	isLogin := redisClient.SIsMember(token, device).Val()
	if isLogin {
		// 同一台裝置登入
		if PenddingMinute > 0 {
			redisClient.Expire(token, time.Minute*PenddingMinute)
		}
		return true, nil
	}
	return false, errors.New("User multi-device login")
}
