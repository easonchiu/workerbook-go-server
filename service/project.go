package service

import (
  `errors`
  `gopkg.in/mgo.v2/bson`
  `time`
  `workerbook/db`
  `workerbook/model`
)

// Insert project info into database.
func CreateProject(data model.Project) error {
  db, close, err := db.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // check the data is error or not.
  if data.Name == "" {
    return errors.New("项目不能为空")
  }

  // supplement other data.
  data.CreateTime = time.Now()
  data.Status = 1

  // name must be the only.
  count, err := db.C(model.ProjectCollection).Find(bson.M{"name": data.Name}).Count()

  if err != nil {
    return err
  }

  if count > 0 {
    return errors.New("已存在相同的项目")
  }

  // create a new object id.
  if data.Id == "" {
    data.Id = bson.NewObjectId()
  }

  // insert it.
  err = db.C(model.ProjectCollection).Insert(data)

  if err != nil {
    return err
  }

  return nil
}

// Query group info by id.
func GetProjectInfoById(id bson.ObjectId) (model.Project, error) {
  db, close, err := db.CloneDB()

  data := model.Project{}

  if err != nil {
    return data, err
  } else {
    defer close()
  }

  err = db.C(model.ProjectCollection).FindId(id).One(&data)

  if err != nil {
    return data, err
  }

  return data, nil
}

// Query groups list with skip and limit.
func GetProjectsList(skip int, limit int, search bson.M) ([]model.Project, error) {
  db, close, err := db.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  data := []model.Project{}

  if limit < 0 {
    limit = 0
  }

  err = db.C(model.ProjectCollection).Find(search).Skip(skip).Limit(limit).All(&data)

  if err != nil {
    return nil, err
  }

  return data, nil
}
