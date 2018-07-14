package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  `gopkg.in/mgo.v2/bson`
  `time`
)

// collection name
const UserCollection = "users"

// collection schema
type User struct {
  // id
  Id bson.ObjectId `bson:"_id,omitempty"`

  // 昵称
  NickName string `bson:"nickname,omitempty"`

  // 邮箱
  Email string `bson:"email,omitempty"`

  // 用户名
  UserName string `bson:"username,omitempty"`

  // 部门
  Department mgo.DBRef `bson:"department,omitempty"`

  // 手机号
  Mobile string `bson:"mobile,omitempty"`

  // 密码
  Password string `bson:"password,omitempty"`

  // 职位 1: 开发者， 2: 部门主管，3: 观察员, 99: 管理员
  Role int `bson:"role,omitempty"`

  // 职称
  Title string `bson:"title,omitempty"`

  // 创建时间
  CreateTime time.Time `bson:"createTime,omitempty"`

  // 用户状态 1: 正常 2: 停用
  Status int `bson:"status,omitempty"`

  // 是否存在
  Exist bool `bson:"exist"`
}

func (u User) GetMap(db *mgo.Database) gin.H {
  department := new(Department)
  db.FindRef(&u.Department).One(department)


  return gin.H{
    "id":             u.Id,
    "nickname":       u.NickName,
    "email":          u.Email,
    "role":           u.Role,
    "title":          u.Title,
    "createTime":     u.CreateTime,
    "username":       u.UserName,
    "departmentId":   u.Department.Id,
    "departmentName": department.Name,
    "status":         u.Status,
  }
}
