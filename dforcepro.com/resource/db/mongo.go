package db

import (
	"errors"
	"reflect"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

var (
	_session *mgo.Session
	TxnC     = "txn"
)

type Mongo struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`

	MainDB   string `yaml:"main-db"`
	CommonDB string `yaml:"common-db"`
	_DB      string
	_C       string
}

func (m Mongo) DB(db string) *Mongo {
	m._DB = db
	return &m
}

func (m Mongo) C(c string) *Mongo {
	m._C = c
	return &m
}

func (m Mongo) GetTxnRunner() *txn.Runner {
	tcollection := m.connect().DB(m._DB).C(TxnC)
	return txn.NewRunner(tcollection)
}

func (m Mongo) connect() *mgo.Session {
	if _session != nil {
		return _session
	}
	dialInfo := mgo.DialInfo{}
	var address []string
	address = strings.SplitN(m.Host, ",", -1)
	if m.User != "" && m.Pass != "" {
		dialInfo.Addrs = address
		dialInfo.Mechanism = ""
		dialInfo.Username = m.User
		dialInfo.Password = m.Pass
	} else {
		dialInfo.Addrs = address
	}

	mySession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		panic(err)
	}

	_session = mySession
	_session.SetMode(mgo.Monotonic, true)

	return _session
}

func (m Mongo) Ping() error {
	session := m.connect()
	return session.Ping()
}

func (m Mongo) Refresh() {
	session := m.connect()
	session.Refresh()
}

func (m Mongo) Insert(docs ...interface{}) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.Insert(docs...)
}

// func (m Mongo) FindById(id string) *mgo.Query {
// 	// session := m.connect()
// 	// defer session.Close()
// 	// cl := session.DB(m._DB).C(m._C)
// 	// return cl.FindId(bson.M{"_id": id})
// }

func (m Mongo) Find(query interface{}) *mgo.Query {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.Find(query)
}
func (m Mongo) Update(field interface{}, update interface{}) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.Update(field, update)
}
func (m Mongo) UpdateAll(field interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.UpdateAll(field, update)
}
func (m Mongo) Upsert(field interface{}, update interface{}) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	_, err := cl.Upsert(field, update)
	return err
}
func (m Mongo) Remove(field interface{}) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.Remove(field)
}
func (m Mongo) RemoveAll(field interface{}) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	_, err := cl.RemoveAll(field)
	return err

}
func (m Mongo) FindId(id interface{}) *mgo.Query {
	myID, _ := GetObjectID(id)
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.FindId(myID)
}

func (m Mongo) UpdateId(id interface{}, update interface{}) error {
	myID, err := GetObjectID(id)
	if err != nil {
		return err
	}
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.UpdateId(myID, update)
}
func (m Mongo) Pipe(pipeline interface{}) *mgo.Pipe {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.Pipe(pipeline)
}

func (m Mongo) EnsureIndexKey(keys ...string) error {
	session := m.connect()
	cl := session.DB(m._DB).C(m._C)
	return cl.EnsureIndexKey(keys...)
}

func (m Mongo) Close() {
	if _session != nil {
		_session.Close()
		_session = nil
	}
}

func (m Mongo) CollectionNames() (names []string, err error) {
	session := m.connect()
	cl, err := session.DB(m._DB).CollectionNames()
	return cl, err
}

func GetObjectID(id interface{}) (*bson.ObjectId, error) {
	var myID bson.ObjectId
	switch dtype := reflect.TypeOf(id).String(); dtype {
	case "string":
		str := id.(string)
		if str == "" || !bson.IsObjectIdHex(str) {
			return nil, errors.New("id is error: " + str)
		}
		myID = bson.ObjectIdHex(str)
	case "bson.ObjectId":
		myID = id.(bson.ObjectId)
	default:
		return nil, errors.New("not support type: " + dtype)
	}
	return &myID, nil
}
