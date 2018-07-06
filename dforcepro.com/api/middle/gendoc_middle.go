package middle

import (
	"net/http"

	"github.com/betacraft/yaag/middleware"
)

type GenDocMiddle bool

func (gdm GenDocMiddle) Enable() bool {
	return bool(gdm)
}

func (gdm GenDocMiddle) GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// one time scope setup area for middleware
		return middleware.HandleFunc(f)
	}
}
