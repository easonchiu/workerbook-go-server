package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// collection name
const DailyCollection = "dailies"

// daily schema
type DailyItem struct {
	Id				bson.ObjectId	`json:"id" bson:"_id"`
	// 项目归属
	Record			string			`json:"record"`
	// 进度
	Progress		int				`json:"progress"`
	// 项目归属名称
	Pname			string			`json:"pname"`
	// 项目归属id
	Pid				string			`json:"pid"`
}

// collection schema
type Daily struct {
	Id				bson.ObjectId	`json:"id" bson:"_id"`
	// 用户id
	Uid				string			`json:"uid"`
	// 用户的分组id
	Gid				string			`json:"gid"`
	// 日期
	Day				string			`json:"day"`
	// 日报数据
	DailyList		[]DailyItem		`json:"dailyList" bson:"dailyList"`
	// 发布时间
	CreateTime		time.Time		`json:"createTime"`
	// 更新时间
	UpdateTime		time.Time		`json:"updateTime"`
}