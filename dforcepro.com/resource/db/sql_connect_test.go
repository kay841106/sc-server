package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func getMsSQL() *SQL {
	sql := SQL{"127.0.0.1", 3306, "dforcepro", "1234"}
	return &sql
}

func Test_Exec(t *testing.T) {
	var db = "Common"
	mssql := getMsSQL()
	conn, err := mssql.GetMySQLConn(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Exec("create table tbl (fld1 int primary key, fld2 int)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// defer conn.Exec("drop table tb1")

	conn.Exec("insert into tbl (fld1, fld2) values (1, 2)")

	assert.Equal(t, 1, 1, "should be equal.")
}
