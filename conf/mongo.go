package conf

import (
	"gopkg.in/mgo.v2"
	"fmt"
)


var (
	Session *mgo.Session
	Mongo *mgo.DialInfo
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
}

func DB() {

}