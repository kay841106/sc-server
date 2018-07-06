package util

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetClientKey(req *http.Request) string {
	deviceID := req.Header.Get("X-Device-ID")
	if deviceID != "" {
		log.Debug(deviceID)
		return MD5(deviceID)
	}

	real := req.Header.Get("X-Real-IP")
	if real == "" {
		real = req.RemoteAddr
	}
	log.Debug(real)
	log.Debug(fmt.Sprintf("add: %s, agent: %s", real, req.UserAgent()))
	return MD5(real + req.UserAgent())
}

func IsLogin(req *http.Request) bool {
	isLoginStr := req.Header.Get("isLogin")
	isLogin, err := strconv.ParseBool(isLoginStr)
	if err != nil {
		return false
	}
	return isLogin
}

func DecodeToken(req *http.Request) *map[string]interface{} {
	token := req.Header.Get("Token")
	if token == "" {
		return nil
	}
	mapSerialize, err := DecodeMap(token)
	if err != nil {
		return nil
	}
	return mapSerialize
}

func GetQueryValue(req *http.Request, keys []string) *map[string]interface{} {
	queries := req.URL.Query()
	result := make(map[string]interface{})

	for _, key := range keys {
		value, ok := queries[key]
		if !ok {
			// if key not exist. use empty string
			result[key] = ""
			continue
		}
		if len(value) == 1 {
			result[key] = value[0]
		} else {
			result[key] = value
		}
	}
	return &result
}

func GetPostValue(req *http.Request, defaultEmpty bool, keys []string) (*map[string]interface{}, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for _, key := range keys {
		if vs := req.PostForm[key]; len(vs) > 0 {
			result[key] = vs[0]
		} else if defaultEmpty {
			result[key] = ""
		}
	}
	return &result, nil
}

type RequestFile struct {
	ReqFile   multipart.File
	ReqHeader *multipart.FileHeader
}

func GetMutiFormPostValue(req *http.Request, fileKeys []string, valueKeys []string) (map[string]RequestFile, map[string]interface{}, error) {
	req.ParseMultipartForm(32 << 20)

	fileMap := make(map[string]RequestFile)
	for _, fk := range fileKeys {
		file, handler, err := req.FormFile(fk)
		if err != nil {
			for _, value := range fileMap {
				defer value.ReqFile.Close()
			}
			return nil, nil, err
		}
		fileMap[fk] = RequestFile{file, handler}
	}

	valueMap := make(map[string]interface{})
	for _, vk := range valueKeys {
		valueMap[vk] = req.FormValue(vk)
	}
	return fileMap, valueMap, nil
}

func GetPathVars(req *http.Request, keys []string) map[string]interface{} {
	vars := mux.Vars(req)
	if len(vars) == 0 {
		return nil
	}
	valueMap := make(map[string]interface{})
	for _, key := range keys {
		if v, ok := vars[key]; ok {
			valueMap[key] = v
		} else {
			valueMap[key] = ""
		}
	}
	if len(valueMap) == 0 {
		return nil
	}
	return valueMap
}

var (
	systemMap = map[string]string{
		"f91c0edc018cccab7e524c099990550d": "lzw",
		"174de676895fbb5239d3a12b95a3a0fb": "ytz",
		"0901dfcac280f58e6527c5502ddd075f": "ytz-web",
	}
)

func GetSysCode(req *http.Request) (string, bool) {
	system := req.Header.Get("System")
	sysCode, ok := systemMap[system]
	if !ok {
		return "", false
	}
	return sysCode, true
}
