package middle

import (
	"net/http"

	"dforcepro.com/util"
)

type AccessMiddle bool

func (am AccessMiddle) Enable() bool {
	return bool(am)
}

func (am AccessMiddle) GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// one time scope setup area for middleware
		return func(w http.ResponseWriter, r *http.Request) {
			// TODO: remove thie back door
			queryValues := util.GetQueryValue(r, []string{"token", "system"})
			token := (*queryValues)["token"].(string)
			system := (*queryValues)["system"].(string)
			if token != "" {
				r.Header.Add("Token", token)
			}
			if system != "" {
				r.Header.Add("System", system)
			}
			f(w, r)
		}
	}
}
