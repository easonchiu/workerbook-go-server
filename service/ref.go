package service

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "workerbook/cache"
  "workerbook/model"
)

func FindRef(mg *mgo.Database) func(mgo.DBRef) (gin.H, bool) {
  return func(ref mgo.DBRef) (gin.H, bool) {
    switch ref.Collection {
    case model.DepartmentCollection:
      department := new(model.Department)
      err := FindDepartmentRef(mg, ref, department)
      if err == nil {
        return department.GetMap(nil), true
      } else {
        return nil, false
      }
    case model.UserCollection:
      user := new(model.User)
      err := FindUserRef(mg, ref, user)
      if err == nil {
        return user.GetMap(nil), true
      } else {
        return nil, false
      }
    case model.ProjectCollection:
      project := new(model.Project)
      err := FindProjectRef(mg, ref, project)
      if err == nil {
        return project.GetMap(nil), true
      } else {
        return nil, false
      }
    case model.MissionCollection:
      mission := new(model.Mission)
      err := FindMissionRef(mg, ref, mission)
      if err == nil {
        return mission.GetMap(nil), true
      } else {
        return nil, false
      }
    default:
      return nil, false
    }

    return nil, false
  }

}

func FindDepartmentRef(mg *mgo.Database, ref mgo.DBRef, department *model.Department) error {
  if ok := cache.DepartmentGet(ref.Id.(bson.ObjectId).Hex(), department); !ok {
    return mg.FindRef(&ref).One(department)
  }
  return nil
}

func FindUserRef(mg *mgo.Database, ref mgo.DBRef, user *model.User) error {
  if ok := cache.UserGet(ref.Id.(bson.ObjectId).Hex(), user); !ok {
    return mg.FindRef(&ref).One(user)
  }
  return nil
}

func FindProjectRef(mg *mgo.Database, ref mgo.DBRef, project *model.Project) error {
  if ok := cache.ProjectGet(ref.Id.(bson.ObjectId).Hex(), project); !ok {
    return mg.FindRef(&ref).One(project)
  }
  return nil
}

func FindMissionRef(mg *mgo.Database, ref mgo.DBRef, mission *model.Mission) error {
  if ok := cache.MissionGet(ref.Id.(bson.ObjectId).Hex(), mission); !ok {
    return mg.FindRef(&ref).One(mission)
  }
  return nil
}
