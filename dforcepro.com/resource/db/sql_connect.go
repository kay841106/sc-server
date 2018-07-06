package db

import (
	"database/sql"
	"fmt"
)

type SQL struct {
	Host string `yaml:"host"`
	Port int16  `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

var _conn *sql.DB

func (s SQL) GetMySQLConn(db string) (*sql.DB, error) {
	if _conn != nil {
		return _conn, nil
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true", s.User, s.Pass, s.Host, s.Port)
	_conn, err := sql.Open("mysql", connString)
	if err == nil {
		_conn.Exec(fmt.Sprintf("USE %s", db))
	}
	return _conn, err
}

func (s SQL) GetMsSQLConn(db string) (*sql.DB, error) {
	if _conn != nil {
		return _conn, nil
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", s.Host, s.User, s.Pass, s.Port)
	_conn, err := sql.Open("mssql", connString)
	if err == nil {
		_conn.Exec(fmt.Sprintf("USE %s", db))
	}
	return _conn, err
}

func (s SQL) ClearTx(tx *sql.Tx) error {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		return err
	}
	return nil
}
