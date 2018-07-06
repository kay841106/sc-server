package common

import (
	"fmt"
	"testing"

	"dforcepro.com/resource"
	"dforcepro.com/resource/db"
	"dforcepro.com/resource/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func Test_Insert(t *testing.T) {
	mssql := db.SQL{Host: "127.0.0.1", Port: 3306, User: "dforcepro", Pass: "1234"}
	log := logger.Logger{Path: "./", Duration: "minute", DebugMode: true}
	log.StartLog()
	_di = &resource.Di{SQL: mssql, Log: log}
	var id1, id2, id3 int64
	region := Region{Name: "taiwan", Sort: 1, Code: "tw", Enable: true}
	dbconn, err := _di.SQL.GetMySQLConn(Database)
	assert.Nil(t, err)
	defer dbconn.Close()
	id1, ok := region.Insert(dbconn)
	assert.True(t, ok)

	region = Region{Name: "taipei", Sort: 1, Code: "tp", Enable: true, ParentCode: "tw"}
	id2, ok = region.Insert(dbconn)
	assert.True(t, ok)
	region = Region{Name: "信義區", Sort: 1, Code: "102", Enable: true, ParentCode: "tp"}
	id3, ok = region.Insert(dbconn)
	assert.True(t, ok)
	fmt.Println(id1, id2, id3)
}

func Test_GetSubRegion(t *testing.T) {
	mssql := db.SQL{Host: "127.0.0.1", Port: 3306, User: "dforcepro", Pass: "1234"}
	dbconn, err := mssql.GetMySQLConn(Database)
	assert.Nil(t, err)
	defer dbconn.Close()
	fmt.Println("aaa")
	results := GetSubRegion(dbconn, "tw")
	fmt.Printf("%v", results)
}
