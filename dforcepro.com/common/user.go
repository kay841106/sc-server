package common

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"dforcepro.com/api"
	"dforcepro.com/resource/db"
	"dforcepro.com/util"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserAPI bool

const (
	UserC           = "User"
	PublicKeyField  = "publicKey"
	PrivateKeyField = "privateKey"
	FormatPEM       = "pem"
	FormatPKIX      = "pkix"
)

func decodeData(cryptoData string, req *http.Request) (*bytes.Buffer, error) {
	key := util.GetClientKey(req)
	redisClient := getRedis(db.RedisDB_LoginKey).GetClient()

	bytes, err := redisClient.HGet(key, PrivateKeyField).Bytes()
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(bytes)
	if err != nil {
		return nil, err
	}
	return util.DecodeString(cryptoData, privateKey)
}

func (ua UserAPI) getKeyEndpoint(w http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()
	var format string
	typeStr, ok := queries["type"]
	if ok {
		format = typeStr[0]
	} else {
		format = FormatPKIX
	}

	key := util.GetClientKey(req)
	redisClient := getRedis(db.RedisDB_LoginKey).GetClient()
	pkBytes, err := redisClient.HGet(key, PublicKeyField+format).Bytes()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(pkBytes)
		return
	}

	var size = 1024
	sizeStr, ok := queries["size"]
	if ok {
		size, err = strconv.Atoi(sizeStr[0])
		if err != nil {
			size = 1024
		}
	}
	var serverPrivateKey *rsa.PrivateKey
	prBytes, err := redisClient.HGet(key, PrivateKeyField).Bytes()
	if err != nil {
		serverPrivateKey, err = rsa.GenerateKey(rand.Reader, size)
	} else {
		serverPrivateKey, err = x509.ParsePKCS1PrivateKey(prBytes)
	}

	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(serverPrivateKey)
	var buf *bytes.Buffer
	if format == FormatPKIX {
		buf, err = util.EncodePublicKey(serverPrivateKey.Public())
	} else if format == FormatPEM {
		buf, err = util.EncodePublicKeyPem(serverPrivateKey.Public())
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not support type."))
		return
	}

	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = redisClient.HSet(key, PublicKeyField+format, buf.Bytes()).Err()
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_di.Log.Info(fmt.Sprintf("insert redis: %s - %s", key, PublicKeyField))
	err = redisClient.HSet(key, PrivateKeyField, privateBytes).Err()
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_di.Log.Info(fmt.Sprintf("insert redis: %s - %s", key, PrivateKeyField))
	// 設定redis key expiretime
	redisClient.Expire(key, 5*time.Minute)

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (ua UserAPI) getTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	queryMap := util.GetQueryValue(req, []string{"cryptoToken"})
	crytoToken, ok := (*queryMap)["cryptoToken"]
	if !ok || crytoToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	crytoTokenStr := crytoToken.(string)
	_di.Log.Debug(crytoTokenStr)
	tokenByte, err := decodeData(crytoTokenStr, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(tokenByte.Bytes())
}

func (ua UserAPI) loginEndpoint(w http.ResponseWriter, req *http.Request) {
	// 判斷使用者是否有token (末實作)
	// 從資料庫取出使用者帳號資料
	// 從Redis取出private key
	// 解密碼出來做MD5比對使用者密碼
	// 成功 - 產生token 將使用者資訊加密存入redis方便取用
	// 寫入redis
	isLogin := util.IsLogin(req)
	if isLogin {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	sysCode, ok := util.GetSysCode(req)
	if !ok {
		_di.Log.Debug("system not set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	account := req.PostForm.Get("user")
	cryptoPwd := req.PostForm.Get("cryptoPwd")
	publicKeyStr := req.PostForm.Get("publicKey")

	if len(account) == 0 || len(cryptoPwd) == 0 || len(publicKeyStr) == 0 {
		_di.Log.Debug("postform not complete")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	plaintPwdBuf, err := decodeData(cryptoPwd, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	userPwd := util.MD5(plaintPwdBuf.String())

	user := User{}
	err = (&user).FindByAccAndPwd(account, userPwd, sysCode)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("account or password error."))
		_di.Log.Debug(fmt.Sprintf("Login with: %s - %s - %s - %s",
			account, userPwd, sysCode, plaintPwdBuf.String()))
		_di.Log.Debug(err.Error())
		return
	}

	group := user.GetGroup(sysCode)
	if group == "" || user.Enable == false {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Permission denied"))
		return
	}

	userSerialize := map[string]interface{}{
		"id":      user.ID.Hex(),
		"account": user.Account,
		"name":    user.Name,
		"group":   group,
		"system":  sysCode,
		"s":       time.Now().Unix(),
	}
	mapSerialize := map[string]interface{}(userSerialize)
	token, err := util.EncodeMap(&mapSerialize)
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Token generate fail"))
		return
	}
	_di.Log.Debug(fmt.Sprintf("Login token: %s", token))

	// set header for middleware to set toke
	req.Header.Add("SET_TOKEN", token)

	// publicKey size must be 2048 to encode token
	if publicKeyStr == "dforcepro" {
		w.Write([]byte(token))
		return
	}

	publicKey, err := util.DecodePublicKey(publicKeyStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	encodeToken, err := util.EncodeString(token, publicKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Encode token error"))
		_di.Log.Debug(err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(encodeToken.Bytes())

}

func (ua UserAPI) logoutEndpoint(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	redisClient := getRedis(db.RedisDB_Token).GetClient()
	err := redisClient.Del(token).Err()
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (ua UserAPI) createEndpoint(w http.ResponseWriter, req *http.Request) {
	// 從Header中取出要登入的系統
	// 解析JSON to User struct
	// 解密pwd & 加密為MD5
	// 存入MongoDB
	sysCode, ok := util.GetSysCode(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDoc := User{}
	_ = json.NewDecoder(req.Body).Decode(&userDoc)
	userDoc.Enable = true
	userDoc.MustChangePwd = false
	userDoc.ID = bson.NewObjectId()
	plaintPwdBuf, err := decodeData(userDoc.Pwd, req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	userDoc.Pwd = util.MD5(plaintPwdBuf.String())

	if len(userDoc.Systems) > 0 {
		for _, sys := range userDoc.Systems {
			userSys := _UserSys{
				SysCode:    sys,
				Group:      userDoc.Group,
				Status:     true,
				CreateTime: time.Now().Unix(),
			}
			userDoc.Sys = append(userDoc.Sys, userSys)
		}
	} else {
		userSys := _UserSys{
			SysCode:    sysCode,
			Group:      userDoc.Group,
			Status:     true,
			CreateTime: time.Now().Unix(),
		}
		userDoc.Sys = append(userDoc.Sys, userSys)
	}

	mongo := getMongo()
	queryValues := util.GetQueryValue(req, []string{"type"})
	if createType := (*queryValues)["type"].(string); createType == "sub" {
		authID := req.Header.Get("AuthID")
		if !bson.IsObjectIdHex(authID) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("AuthID error"))
			return
		}
		parentUser := bson.ObjectIdHex(authID)
		userDoc.ParentUser = &parentUser
		subAccCount, err := mongo.DB(Database).C(UserC).Find(bson.M{"parentuser": parentUser}).Count()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		authAccount := req.Header.Get("AuthAccount")
		serialNumStr, err := util.IntToFixStrLen(subAccCount+1, 3)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("sub account is too much"))
			return
		}
		userDoc.Account = authAccount + serialNumStr
	}

	err = mongo.DB(Database).C(UserC).Insert(userDoc)

	if err != nil {
		_di.Log.WriteFile(fmt.Sprintf("input/%s/%s", UserC, userDoc.ID.Hex()), toJSONByte(userDoc))
		_di.Log.Err(err.Error())
		// 回傳錯誤息
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can not insert data."))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("success"))
	}
}

func (ua UserAPI) updateEndpoint(w http.ResponseWriter, req *http.Request) {
	pathVar := util.GetPathVars(req, []string{"ID"})
	userID := ""
	if len(pathVar) == 0 {
		userID = req.Header.Get("AuthID")
	} else {
		userID = pathVar["ID"].(string)
	}
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user id error."))
		return
	}

	userSet := bson.M{}
	_ = json.NewDecoder(req.Body).Decode(&userSet)
	allowField := []string{"company", "name", "description", "phone", "address", "telephone"}
	for key := range userSet {
		if !util.IsStrInList(key, allowField...) {
			delete(userSet, key)
		}
	}

	err := getMongo().DB(Database).C(UserC).UpdateId(userID, bson.M{"$set": userSet})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("success"))
	}
}

func (ua UserAPI) getSelfEndpoint(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("AuthID")

	userResult := User{}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(userID)}},
		bson.M{"$lookup": bson.M{
			"from":         UserC,
			"localField":   "_id",
			"foreignField": "parentuser",
			"as":           "subusers"},
		},
	}
	pipe := getMongo().DB(Database).C(UserC).Pipe(pipeline)
	err := pipe.One(&userResult)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResult)
}

func (ua UserAPI) getEncryptPwdEndpoint(w http.ResponseWriter, req *http.Request) {
	key := util.GetClientKey(req)
	redisClient := getRedis(db.RedisDB_LoginKey).GetClient()

	bytes, err := redisClient.HGet(key, PublicKeyField+FormatPKIX).Bytes()
	if err != nil {
		_di.Log.Err(fmt.Sprintf("key: %s. redis msg: %s", key, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	publicKey, err := util.DecodePublicKey(string(bytes))
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(req)
	pwd := vars["pwd"]
	ecodePwdBuf, err := util.EncodeString(pwd, publicKey)
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(ecodePwdBuf.Bytes())
	}
}

func (ua UserAPI) searchEndpoint(w http.ResponseWriter, req *http.Request) {
	limit, page := 100, 1
	paramMap := util.GetQueryValue(req, []string{"limit", "page"})
	validate := map[string][]string{
		"limit": []string{"Numeric"},
		"page":  []string{"Numeric"},
	}
	ok, errMap := util.CheckParam(*paramMap, validate)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMap)
		return
	}
	if value, ok := (*paramMap)["limit"]; ok {
		limit, _ = strconv.Atoi(value.(string))
	}
	if value, ok := (*paramMap)["page"]; ok {
		page, _ = strconv.Atoi(value.(string))
	}

	sysCode, ok := util.GetSysCode(req)
	if !ok {
		_di.Log.Debug("system not set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &User{}
	mogonQuery := user.Find(bson.M{
		"sys": bson.M{"$elemMatch": bson.M{"syscode": sysCode}},
	})

	result, err := util.MongoPagination(mogonQuery, limit, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var userSearchList []UserSearch
	iter := mogonQuery.Iter()
	for iter.Next(user) {
		userSearch := &UserSearch{}
		userSearch.convert(user, sysCode)
		userSearchList = append(userSearchList, *userSearch)
	}
	if err := iter.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	result.Raws = userSearchList
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (ua UserAPI) getOneEndpoint(w http.ResponseWriter, req *http.Request) {
	vars := util.GetPathVars(req, []string{"ID"})
	if vars == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Path ID not set."))
		return
	}
	syscode, ok := util.GetSysCode(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("header system not set."))
		return
	}
	userID := vars["ID"].(string)
	result := User{}
	err := getMongo().DB(Database).C(UserC).
		FindId(userID).
		One(&result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	us := &UserSearch{}
	us.convert(&result, syscode)
	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(us)
}

func (ua UserAPI) changePwdEndpoint(w http.ResponseWriter, req *http.Request) {
	pathVar := util.GetPathVars(req, []string{"ID"})
	userID := ""
	if len(pathVar) == 0 {
		userID = req.Header.Get("AuthID")
	} else {
		userID = pathVar["ID"].(string)
	}
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user id error."))
		return
	}

	userDoc := &User{}
	err := userDoc.FindByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not find user."))
		return
	}

	postVal, err := util.GetPostValue(req, false, []string{"pwd"})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if len(*postVal) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("must post pwd."))
		return
	}
	cryptoPwd := (*postVal)["pwd"].(string)

	plaintPwdBuf, err := decodeData(cryptoPwd, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	md5pwd := util.MD5(plaintPwdBuf.String())
	err = userDoc.ChangePwd(md5pwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("success"))
}

func (ua UserAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/v1/key", Next: ua.getKeyEndpoint, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/v1/token", Next: ua.getTokenEndpoint, Method: "GET", Auth: false},
		&api.APIHandler{Path: "/v1/login", Next: ua.loginEndpoint, Method: "POST", Auth: false},
		&api.APIHandler{Path: "/v1/logout", Next: ua.logoutEndpoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/user/_search", Next: ua.searchEndpoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/user/pwd", Next: ua.changePwdEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/v1/user/{ID}/pwd", Next: ua.changePwdEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/v1/user/log", Next: ua.getLogEndPoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/user/{ID}", Next: ua.getOneEndpoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/user/{ID}", Next: ua.updateEndpoint, Method: "PUT", Auth: true},
		&api.APIHandler{Path: "/v1/user", Next: ua.createEndpoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/v1/user", Next: ua.updateEndpoint, Method: "PUT", Auth: true},
		&api.APIHandler{Path: "/v1/user", Next: ua.getSelfEndpoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/pwd/{pwd}", Next: ua.getEncryptPwdEndpoint, Method: "GET", Auth: false},
	}
}

func (ua UserAPI) Enable() bool {
	return bool(ua)
}

type UserSearch struct {
	ID          string `json:"ID,omitempty"`
	Company     string `json:"Company,omitempty"`
	Account     string `json:"Account,omitempty"`
	Name        string `json:"Name,omitempty"`
	Email       string `json:"Email,omitempty"`
	Phone       string `json:"Phone"`
	Telephone   string `json:"Telephone"`
	Address     string `json:"Address"`
	Group       string `json:"Group,omitempty"`
	Description string `json:"Description"`
	Enable      bool   `json:"Enable,omitempty"`
}

func (us *UserSearch) convert(u *User, sysCode string) {
	us.ID = u.ID.Hex()
	us.Name = u.Name
	us.Account = u.Account
	us.Company = u.Company
	us.Email = u.Email
	us.Phone = u.Phone
	us.Description = u.Description
	us.Enable = u.Enable
	us.Group = u.GetGroup(sysCode)
	us.Address = u.Address
	us.Telephone = u.Telephone
}

type User struct {
	ID            bson.ObjectId     `json:"ID,omitempty" bson:"_id"`
	Company       string            `json:"Company,omitempty"`
	Name          string            `json:"Name,omitempty"`
	Account       string            `json:"Account,omitempty"`
	Email         string            `json:"Email,omitempty"`
	Phone         string            `json:"Phone,omitempty"`
	Address       string            `json:"Address,omitempty"`
	Telephone     string            `json:"Telephone,omitempty"`
	Description   string            `json:"Description,omitempty"`
	Pwd           string            `json:"Pwd,omitempty"`
	PublicKey     string            `json:"PublicKey,omitempty" bson:"-"`
	PushToken     map[string]string `json:"PushToken,omitempty"`
	Group         string            `json:"Group,omitempty" bson:"-"`
	Systems       []string          `json:"System,omitempty" bson:"-"`
	Sys           []_UserSys        `json:"-"`
	ParentUser    *bson.ObjectId    `json:"-"`
	SubUsers      []User            `json:"SubUsers,omitempty"`
	Enable        bool
	MustChangePwd bool
}

type SimpleUser struct {
	ID   bson.ObjectId `json:"ID,omitempty"`
	Name string        `json:"Name,omitempty"`
}

func (u *User) ToSimple() *SimpleUser {
	return &SimpleUser{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (u *User) Insert(us ..._UserSys) error {
	return nil
}

func (u *User) ChangePwd(newpwd string) error {
	return getMongo().DB(Database).C(UserC).UpdateId(u.ID.Hex(), bson.M{"$set": bson.M{"pwd": newpwd}})
}

func (u *User) Find(query bson.M) *mgo.Query {
	return getMongo().DB(Database).C(UserC).Find(query)
}

func (u *User) FindByID(id string) error {
	return getMongo().DB(Database).C(UserC).FindId(id).One(u)
}

func (u *User) FindByAccount(acc string, sysCode string) error {
	return getMongo().DB(Database).C(UserC).Find(
		bson.M{
			"account": acc,
			"sys":     bson.M{"$elemMatch": bson.M{"syscode": sysCode}},
		},
	).One(u)
}

func (u *User) FindByAccAndPwd(account string, pwd string, sysCode string) error {
	return getMongo().DB(Database).C(UserC).Find(
		bson.M{
			"account": account,
			"pwd":     pwd,
			"sys":     bson.M{"$elemMatch": bson.M{"syscode": sysCode}},
		},
	).One(u)
}

func (u *User) GetGroup(sysCode string) string {
	for _, val := range u.Sys {
		if val.SysCode == sysCode {
			return val.Group
		}
	}
	return ""
}

type _UserSys struct {
	SysCode    string `json:"SysCode,omitempty"`
	Group      string `json:"Group,omitempty"` // 紀錄權限
	Status     bool   `json:"Status,omitempty"`
	CreateTime int64  `json:"CreateTime,omitempty"`
	UpdateTime int64  `json:"UpdateTime,omitempty"`
}

func (ua UserAPI) getLogEndPoint(w http.ResponseWriter, req *http.Request) {

	queries := req.URL.Query()
	s, ok := queries["start"]
	e, ok2 := queries["end"]
	if !ok {
		s = append(s, "2000-01-01")
	}
	if !ok2 {
		e = append(e, "2100-01-01")
	}
	shortForm := "2006-01-02"
	t1, _ := time.Parse(shortForm, s[0])
	t2, _ := time.Parse(shortForm, e[0])
	var op []OperateLog

	system, ok := util.GetSysCode(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must set System in the header."))
	}

	mongo := getMongo()

	query := mongo.DB("Common").C("Log").Find(bson.M{
		"sys": system,
		"timestamp": bson.M{
			"$gte": t1.Unix(),
			"$lt":  t2.Unix(),
		},
	})

	total, err := query.Count()
	if err != nil {
		return
	}
	if total == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	query.All(&op)
	result := make([]interface{}, 0)
	for _, element := range op {
		result = append(result, element)
	}
	w.Header().Set("Content-Type", "application/json")
	responseJSON := queryResNoPage{&result, total}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseJSON)

}

type OperateLog struct {
	Path      string `json:"path"`
	Body      string `json:"body"`
	Name      string `json:"name"`
	Account   string `json:"account"`
	Header    string `json:"header"`
	Timestamp int32  `json:"timestamp"`
}

type queryResNoPage struct {
	Raws  *[]interface{} `json:"raw,omitempty"`
	Total int
}
