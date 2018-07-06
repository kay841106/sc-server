package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AppStat struct {
	Version    string         `json:"version"`
	QueueMax   int64          `json:"queue_max"`
	QueueUsage int64          `json:"queue_usage"`
	Total      int64          `json:"total_count"`
	Ios        StatPushResult `json:"ios,omitempty"`
	Android    StatPushResult `json:"android,omitempty"`
}
type StatPushResult struct {
	Success int `json:"push_success"`
	Error   int `json:"push_error"`
}
type PushResult struct {
	Counts  int       `json:"counts"`
	Log     []PushLog `json:"logs"`
	Success string    `json:"success"`
}
type PushDataArray struct {
	Notifications []PushData `json:"notifications"`
}
type PushLog struct {
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Token    string `json:"token"`
	Message  string `json:"message"`
	Err      string `json:"error"`
}
type PushData struct {
	Tokens   []string    `json:"tokens"`
	PType    int64       `json:"platform"`
	Priority string      `json:"priority"`
	Data     DataContent `json:"data"`
}
type DataContent struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

const (
	address      = "http://localhost:8089/"
	pushRoute    = "api/push"
	appStatRoute = "api/stat/app"
)

func GoRushStatApp() (*AppStat, error) {
	// Set up a connection to the server.

	resp, err := http.Get(fmt.Sprintf("%s%s", address, appStatRoute))
	body := AppStat{}
	if err != nil {
		fmt.Println(err)
	} else {

		json.NewDecoder(resp.Body).Decode(&body)
		return &body, nil
	}
	defer resp.Body.Close()
	return &body, err
}
func GoRushPush(token []string, title string, message string, platform int64, priority string) (*PushResult, error) {
	// Set up a connection to the server.
	u := PushDataArray{}
	push := PushData{token, platform, priority, DataContent{title, message}}
	u.Notifications = append(u.Notifications, push)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	resp, _ := http.Post(fmt.Sprintf("%s%s", address, pushRoute), "application/json; charset=utf-8", b)
	body := PushResult{}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&body)
	return &body, nil

}
