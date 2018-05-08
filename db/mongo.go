package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"errors"
	"net/http"
	"workerbook/conf"
)

var (
	Session *mgo.Session
	Mongo *mgo.DialInfo
	Connecting = false
)

func ConnectDB () {
	mongo, err := mgo.ParseURL(conf.DBUrl)

	s, err := mgo.Dial(conf.DBUrl)

	if err != nil {
		panic(err)
	}

	s.SetSafe(&mgo.Safe{})

	fmt.Println("Connect database successed.")

	Session = s
	Mongo = mongo
	Connecting = true
}

// get db with clone session
// must close the session after use !!!
//   e.g.  defer session.close()
func CloneDB() (*mgo.Database, func(), error) {
	if Connecting {
		session := Session.Clone()
		closeFn := func() {
			session.Close()
		}
		return session.DB(Mongo.Database), closeFn, nil
	}

	return nil, nil, errors.New(http.StatusText(http.StatusBadGateway))
}

// close db
func CloseDB() {
	if Connecting {
		Session.Close()
		Connecting = false
		fmt.Println("Database is closed.")
	} else {
		panic(errors.New("Database is not connected."))
	}
}