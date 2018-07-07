package service

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/mongo"
)

// Insert department info into database.
func CreateDepartment(data model.Department) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // check the data is error or not.
  if data.Name == "" {
    return errors.New(errno.ErrDepartmentNameEmpty)
  }

  // supplement other data.
  data.CreateTime = time.Now()
  data.UserCount = 0

  // name must be the only.
  count, err := db.C(model.DepartmentCollection).Find(bson.M{"name": data.Name}).Count()

  if err != nil {
    return err
  }

  if count > 0 {
    return errors.New(errno.ErrSameDepartmentName)
  }

  // insert it.
  err = db.C(model.DepartmentCollection).Insert(data)

  if err != nil {
    return errors.New(errno.ErrCreateDepartmentFailed)
  }

  return nil
}

// Query group info by id.
// func GetGroupInfoById(id bson.ObjectId) (*model.Group, error) {
//   db, close, err := mongo.CloneDB()
//
//   data := new(model.Group)
//
//   if err != nil {
//     return data, err
//   } else {
//     defer close()
//   }
//
//   err = db.C(model.GroupCollection).FindId(id).One(data)
//
//   if err != nil {
//     return data, err
//   }
//
//   return data, nil
// }

// Query groups list with skip and limit.
// func GetGroupsList(skip int, limit int) (*[]model.Group, error) {
//   db, close, err := mongo.CloneDB()
//
//   if err != nil {
//     return nil, err
//   } else {
//     defer close()
//   }
//
//   data := new([]model.Group)
//
//   if limit < 0 {
//     limit = 0
//   }
//
//   err = db.C(model.GroupCollection).Find(bson.M{}).Skip(skip).Limit(limit).All(data)
//
//   if err != nil {
//     return nil, err
//   }
//
//   return data, nil
// }

// get count of group
// func GetCountOfGroup() (int, error) {
//   db, close, err := mongo.CloneDB()
//
//   if err != nil {
//     return 0, err
//   } else {
//     defer close()
//   }
//
//   return db.C(model.GroupCollection).Count()
// }

// update department count
func UpdateDepartmentCount(departmentId string) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  userCount, err := db.C(model.UserCollection).Find(bson.M{
    "departmentId": departmentId,
  }).Count()

  if err != nil {
    return err
  }

  db.C(model.DepartmentCollection).UpdateId(bson.ObjectIdHex(departmentId), bson.M{
    "$set": bson.M{
      "userCount": userCount,
    },
  })

  return nil
}
