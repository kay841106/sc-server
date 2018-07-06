package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"dforcepro.com/api"
	"dforcepro.com/api/middle"
	"dforcepro.com/common"
	"dforcepro.com/cron"
	"dforcepro.com/resource"

	scair "sc.dforcepro.com/airbox"
	sccron "sc.dforcepro.com/cron"
	scdispenser "sc.dforcepro.com/dispenser"
	scmet "sc.dforcepro.com/meter"
	scspa "sc.dforcepro.com/space"
	scstalone "sc.dforcepro.com/standalone"

	"github.com/betacraft/yaag/yaag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	args := os.Args[1:]
	if count := len(args); count == 0 {
		fmt.Println("parameters error.")
		return
	}
	configPath := args[1]

	filename, _ := filepath.Abs(configPath + "config.yml")
	resource.IniConf(filename)
	resource.ConfPath = configPath
	di, err := resource.GetDI()
	if err != nil {
		panic(err)
	}
	defer resource.Close()

	di.Log.StartLog()
	// doc.SetDI(di)
	scmet.SetDI(di)
	scspa.SetDI(di)
	scair.SetDI(di)
	// scstalone.SetDI(di)
	// scflux.SetDI(di)

	// util.InitValidator()
	// iniTaskServer(di)
	switch usage := args[0]; usage {
	case "api":
		runAPI(di)
	case "cron":
		runCron(di)
	case "standalone":
		runStandAlone(di)
	default:
		fmt.Println("parameters error. only api or job")
	}

}

func runAPI(di *resource.Di) {

	router := mux.NewRouter()
	yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})

	middle.SetDI(di)
	middlewares := middle.GetMiddlewares(middle.DebugMiddle(true), middle.GenDocMiddle(false))

	// apiConf := &api.APIconf{Router: router, MiddleWares: middlewares}
	// middleConf := di.APIConf.Middle
	// middle.SetDI(di)
	// middlewares := middle.GetMiddlewares(middle.DebugMiddle(true), middle.GenDocMiddle(true), middle.AuthMiddle(true))
	// middlewares := middle.GetMiddlewares(
	// 	middle.DebugMiddle(middleConf.Debug),
	// 	middle.GenDocMiddle(middleConf.GenDoc),
	// )
	apiConf := &api.APIconf{Router: router, MiddleWares: middlewares}

	common.SetDI(di)

	// api.InitAPI(apiConf, scmet.SCConfAPI(true), scmet.SCBuildAPI(true), scspa.SmartSpAPI(true), scair.SCAirbox(true), scflux.Influx(true))
	api.InitAPI(apiConf, scmet.SCConfAPI(true), scmet.SCBuildAPI(true), scspa.SmartSpAPI(true), scair.SCAirbox(true), scdispenser.Dispenser(true))
	di.Log.Debug("Start API. Port: " + di.APIConf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", di.APIConf.Port), router))
}

func runCron(di *resource.Di) {
	sccron.SetDI(di)
	// cron.StartCronJob(sccron.LastStream(true))
	cron.StartCronJob(sccron.Weather(true), sccron.LastStream(false), sccron.Airbox(true), sccron.Dispenser(true))
	// cron.StartCronJob(sccron.LastStream(true))

}

func runStandAlone(di *resource.Di) {
	// middlewares := middle.GetMiddlewares(middle.DebugMiddle(true), middle.GenDocMiddle(false))
	// go scstalone.DisplayDataCalc()
	scstalone.AirboxStreamAll()
	// cron.StartCronJob(sccron.LastStream(true))
	// cron.StartCronJob(sccron.Weather(false), sccron.LastStream(true), sccron.Airbox(false), sccron.Dispenser(false))
	// cron.StartCronJob(sccron.LastStream(true))

}
