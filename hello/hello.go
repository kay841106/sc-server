package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {

	session, err := mgo.Dial("140.118.70.136:27017")

	if err != nil {
		panic(err)
	}
	mux := mux.NewRouter()
	mux.Host("140.118.70.136:9112")
	mux.HandleFunc("/", middleware.HandleFunc(handler))
	// admindb := session.DB("admin")

	cred := mgo.Credential{
		Username:  "root",
		Password:  "123",
		Source:    "admin",
		Mechanism: "SCRAM-SHA-1",
	}

	ta := session.Login(&cred)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("enter.")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	fmt.Println("ta", ta)

	yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})
	// http.ListenAndServe(":8080", mux)

	num := []int{2, 3, 5}
	sum := 0
	for _, each := range num {
		sum += each

	}
	fmt.Println("asu", sum)
	// fmt.Println(er)
	args := os.Args[1:]
	configPath := args[1]
	filename, _ := filepath.Abs(configPath + "config.yml")

	fmt.Println(configPath)
	fmt.Println(filename)
	// resource.IniConf(filename)
	// resource.ConfPath = configPath
	// di, err := resource.GetDI()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(len(args))

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}
