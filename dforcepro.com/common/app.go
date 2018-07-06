package common

import (
	"net/http"

	"dforcepro.com/api"
	"dforcepro.com/util"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type AppAPI bool

const (
	AppC = "App"
)

type App struct {
	ID       bson.ObjectId `json:"-,omitempty" bson:"_id"`
	System   string        `json:"system,omitempty"`
	Type     string        `json:"type,omitempty"`
	Version  string        `json:"version,omitempty"`
	Filename string        `json:"filename,omitempty"`
}

func (aa AppAPI) getVersionEndpoint(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	vType := vars["type"]
	if !util.IsStrInList(vType, "ios", "android") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sysCode, ok := util.GetSysCode(req)
	if !ok {
		_di.Log.Debug("system not set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app := App{}
	mongo := getMongo()
	mongo.DB(Database).C(AppC).Find(bson.M{"system": sysCode, "type": vType}).One(&app)

	if app.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte("can not find")))
		return
	}

	w.Write([]byte(app.Version))
}

func (aa AppAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/v1/app_version/{type}", Next: aa.getVersionEndpoint, Method: "GET", Auth: false},
	}
}

func (aa AppAPI) Enable() bool {
	return bool(aa)
}

///app_version/{type}
