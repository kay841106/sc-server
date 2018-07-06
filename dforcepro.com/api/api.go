package api

import (
	"fmt"
	"net/http"

	"dforcepro.com/api/middle"
	"github.com/gorilla/mux"
)

const (
	StatusRemote = 106
)

type APIconf struct {
	Router      *mux.Router
	MiddleWares *[]middle.Middleware
}

type APIHandler struct {
	Path   string
	Next   func(http.ResponseWriter, *http.Request)
	Method string
	Auth   bool
}

type API interface {
	GetAPIs() *[]*APIHandler
	Enable() bool
}

func InitAPI(conf *APIconf, apis ...API) {
	for _, myapi := range apis {
		if myapi.Enable() {
			addHandler(conf, myapi.GetAPIs())
		}
	}
}

func addHandler(conf *APIconf, apiHandlers *[]*APIHandler) {
	router := conf.Router
	for _, handler := range *apiHandlers {
		middle.AddAuthPath(fmt.Sprintf("%s:%s", handler.Path, handler.Method), handler.Auth)
		router.HandleFunc(handler.Path, middle.BuildChain(handler.Next, *conf.MiddleWares...)).Methods(handler.Method)
	}
}
