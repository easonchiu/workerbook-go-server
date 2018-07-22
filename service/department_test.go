package service

import (
  "gopkg.in/mgo.v2/bson"
  "testing"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/mongo"
)

func TestDepartment(t *testing.T) {
  db.ConnectDB()
  defer db.CloseDB()

  id := bson.NewObjectId()

  // 创建部门
  err := CreateDepartment(model.Department{
    Id:   id,
    Name: "测试部门",
  })

  if err != nil {
    t.Error(errgo.Get(err))
    return
  }

  // 修改部门
  err = UpdateDepartment(id.Hex(), bson.M{
    "name": "修改了部门",
  })

  if err != nil {
    t.Error(errgo.Get(err))
    return
  }

  // 删除部门
  err = DelDepartmentById(id.Hex())

  if err != nil {
    t.Error(errgo.Get(err))
    return
  }

}
