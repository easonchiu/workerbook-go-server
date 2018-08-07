package model

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
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

  // 公司
  CompanyName string `bson:"companyName"`

  // 公司id
  CompanyId string `bson:"companyId"`

  // 工号
  JobNumber string `bson:"jobNumber"`

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

  // 操作人
  Editor mgo.DBRef `bson:"editor,omitempty"`

  // 操作时间
  EditTime time.Time `bson:"editTime,omitempty"`
}

func (u *User) GetMap(forgets ... string) gin.H {
  data := gin.H{
    "id":       u.Id,
    "nickname": u.NickName,
    "email":    u.Email,
    "mobile":   u.Mobile,
    "role":     u.Role,
    "title":    u.Title,
    "department": gin.H{
      "id": u.Department.Id,
    },
    "createTime": u.CreateTime,
    "username":   u.UserName,
    "status":     u.Status,

    "exist": u.Exist,
    "editor": bson.M{
      "id": u.Editor.Id,
    },
    "editTime": u.EditTime,
  }

  if forgets != nil {
    if forgets[0] == REMEMBER {
      remember(data, forgets[1:]...)
    } else {
      forget(data, forgets...)
    }
  }

  return data
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

// 用户列表结构
type UserList struct {
  List  []*User
  Count int
  Limit int
  Skip  int
}

// 列表的迭代器
func (d *UserList) Each(fn func(*User) gin.H) gin.H {
  data := gin.H{}

  if d.Limit != 0 {
    data = gin.H{
      "count": d.Count,
      "limit": d.Limit,
      "skip":  d.Skip,
    }
  }

  var list []gin.H
  for _, item := range d.List {
    list = append(list, fn(item))
  }

  data["list"] = list

  return data
}

func (d *UserList) Find(id bson.ObjectId) *User {
  if d.List == nil {
    return nil
  }
  for _, item := range d.List {
    if item.Id == id {
      return item
    }
  }
  return nil
}

func (d *UserList) Ids() []bson.ObjectId {
  var list []bson.ObjectId
  for _, item := range d.List {
    list = append(list, item.Id)
  }
  return list
}
