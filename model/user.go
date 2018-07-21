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
  NickName string `bson:"nickname"`

  // 邮箱
  Email string `bson:"email"`

  // 用户名
  UserName string `bson:"username"`

  // 部门
  Department mgo.DBRef `bson:"department"`

  // 手机号
  Mobile string `bson:"mobile"`

  // 密码
  Password string `bson:"password"`

  // 职位 1: 开发者， 2: 部门管理者，3: 观察员, 4: 项目管理者 99: 管理员
  Role int `bson:"role"`

  // 职称
  Title string `bson:"title"`

  // 创建时间
  CreateTime time.Time `bson:"createTime"`

  // 用户状态 1: 正常 2: 停用
  Status int `bson:"status"`

  // 是否存在
  Exist bool `bson:"exist"`
}

func (u User) GetMap(db *mgo.Database) gin.H {
  department := new(Department)
  db.FindRef(&u.Department).One(department)

  return gin.H{
    "id":         u.Id,
    "nickname":   u.NickName,
    "email":      u.Email,
    "role":       u.Role,
    "title":      u.Title,
    "createTime": u.CreateTime,
    "username":   u.UserName,
    "department": gin.H{
      "id":   department.Id,
      "name": department.Name,
    },
    "status": u.Status,
  }
}

/*
                          1: 开发者 2: 部门管理者 3: 观察员 4: 项目管理者 99: 管理员
创建/修改/删除用户           x        x           x        x           √
创建/修改/删除部门           x        x           x        x           √
创建/修改/删除项目           x        x           x        √           √
暂停/归档项目               x        x           x        √           √
创建/修改/删除任务           x        √ own       x        √           √
查看任务/部门数据            √        √           √        √           √
添加/删除日报               √        √           x        x           x
查看项目                   √ own    √ own       √        √           √
重置密码                   x        x           x        x           √

*/
