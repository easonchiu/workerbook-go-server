package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"errors"
)

var (
	Session *mgo.Session
	Mongo *mgo.DialInfo
	Connecting bool = false
)

const dburl = "mongodb://localhost:27017/workerbook"

func ConnectDB () {
	mongo, err := mgo.ParseURL(dburl)

	s, err := mgo.Dial(dburl)

	if err != nil {
		panic(err)
	}

	s.SetSafe(&mgo.Safe{})

	fmt.Println("Connect database successed.")

	Session = s
	Mongo = mongo
	Connecting = true
}

// get database
func MuseDB() *mgo.Database {
	if Connecting {
		return Session.DB(Mongo.Database)
	}
	panic(errors.New("Database is not connected."))
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