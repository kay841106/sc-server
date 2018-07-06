package middle

import (
	"net/http"

	"dforcepro.com/resource"
	"dforcepro.com/resource/db"
)

type middle interface {
	Enable() bool
	GetMiddleWare() func(f http.HandlerFunc) http.HandlerFunc
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

var (
	_di *resource.Di
)

func SetDI(c *resource.Di) {
	_di = c
}

// buildChain builds the middlware chain recursively, functions are first class
func BuildChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](BuildChain(f, m[1:len(m)]...))
}

func GetMiddlewares(middles ...middle) *[]Middleware {
	var middlewares []Middleware
	for _, m := range middles {
		if m.Enable() {
			middlewares = append(middlewares, m.GetMiddleWare())
		}
	}
	return &middlewares
}

func getRedis(redisdb int) *db.Redis {
	return (&_di.Redis).DB(redisdb)
}

func getMongo() db.Mongo {
	return _di.Mongodb
}
