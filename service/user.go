package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "github.com/jwt-go"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/errgo"
  "workerbook/model"
)

// 创建用户
func CreateUser(data model.User, departmentId string) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

  // 用户状态
  data.Status = 1

  // check
  errgo.ErrorIfStringIsEmpty(data.NickName, errgo.ErrNicknameEmpty)
  errgo.ErrorIfLenLessThen(data.NickName, 2, errgo.ErrNicknameTooShort)
  errgo.ErrorIfLenMoreThen(data.NickName, 14, errgo.ErrNicknameTooLong)
  errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  errgo.ErrorIfStringIsEmpty(data.Title, errgo.ErrUserTitleIsEmpty)
  errgo.ErrorIfLenMoreThen(data.Title, 14, errgo.ErrUserTitleTooLong)
  if data.Role != 1 && data.Role != 2 && data.Role != 3 {
    return errors.New(errgo.ErrUserRoleError)
  }
  errgo.ErrorIfStringIsEmpty(data.UserName, errgo.ErrUsernameEmpty)
  errgo.ErrorIfLenLessThen(data.UserName, 6, errgo.ErrUsernameTooShort)
  errgo.ErrorIfLenMoreThen(data.UserName, 14, errgo.ErrUsernameTooLong)
  errgo.ErrorIfStringIsEmpty(data.Password, errgo.ErrPasswordEmpty)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // save department id
  data.Department = mgo.DBRef{
    Id:         bson.ObjectIdHex(departmentId),
    Collection: model.DepartmentCollection,
    Database:   conf.MgoDBName,
  }

  // username must be the only.
  count, err := mg.C(model.UserCollection).Find(model.User{
    UserName: data.UserName,
    Exist:    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameUsername)
  }

  // nickname must be the only.
  count, err = mg.C(model.UserCollection).Find(model.User{
    NickName: data.NickName,
    Exist:    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameNickname)
  }

  // department must be exist.
  department := new(model.Department)
  mg.FindRef(&data.Department).One(department)

  if department.Name == "" {
    return errors.New(errgo.ErrDepartmentNotFound)
  }

  // insert it.
  err = mg.C(model.UserCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount()

  return nil
}

// 更新用户
func UpdateUser(id string, data bson.M) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if nickname, ok := data["nickname"]; ok {
    errgo.ErrorIfLenLessThen(nickname.(string), 2, errgo.ErrNicknameTooShort)
    errgo.ErrorIfLenMoreThen(nickname.(string), 14, errgo.ErrNicknameTooLong)
  }

  if departmentId, ok := data["department.$id"]; ok {
    errgo.ErrorIfStringNotObjectId(departmentId.(string), errgo.ErrDepartmentIdError)
  }

  if title, ok := data["title"]; ok {
    errgo.ErrorIfLenMoreThen(title.(string), 14, errgo.ErrUserTitleTooLong)
  }

  if role, ok := data["role"]; ok {
    role := role.(int)
    if role != 1 && role != 2 && role != 3 {
      return errors.New(errgo.ErrUserRoleError)
    }
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 姓名唯一
  if nickname, ok := data["nickname"]; ok {
    count, err := mg.C(model.UserCollection).Find(bson.M{
      "nickname": nickname,
      "exist":    true,
      "_id": bson.M{
        "$ne": bson.ObjectIdHex(id),
      },
    }).Count()

    if err != nil {
      return errors.New(errgo.ErrUpdateUserFailed)
    }

    if count > 0 {
      return errors.New(errgo.ErrSameNickname)
    }
  }

  // 部门必须存在
  if departmentId, ok := data["department.$id"]; ok {
    ref := mgo.DBRef{
      Id:         bson.ObjectIdHex(departmentId.(string)),
      Collection: model.DepartmentCollection,
      Database:   conf.MgoDBName,
    }

    count, err := mg.FindRef(&ref).Select(bson.M{"exist": true}).Count()

    if err != nil {
      return errors.New(errgo.ErrUpdateUserFailed)
    }

    if count == 0 {
      return errors.New(errgo.ErrDepartmentNotFound)
    }

    data["department.$id"] = bson.ObjectIdHex(departmentId.(string))
  }

  // 先要清缓存，清成功后才可以更新数据
  err = cache.UserDel(id)

  if err != nil {
    return errors.New(errgo.ErrUpdateUserFailed)
  }

  // 更新数据
  data["exist"] = true
  err = mg.C(model.UserCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount()

  return nil
}

// 根据id删除用户
func DelUserById(id string) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 清除缓存，缓存清成功才可以清数据，不然会有脏数据
  err = cache.UserDel(id)

  if err != nil {
    return errors.New(errgo.ErrDeleteUserFailed)
  }

  // 删除数据
  err = mg.C(model.UserCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": bson.M{
      "exist": false,
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrUserNotFound)
    }
    return err
  }

  // 更新部门人数
  err = UpdateDepartmentsUserCount()

  if err != nil {
    return errors.New(errgo.ErrDeleteUserFailed)
  }

  return nil
}

// 根据id查找用户信息
func GetUserInfoById(id string, refs ... string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.User)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.UserGet(id, data)
  if !rok {
    err = mg.C(model.UserCollection).Find(bson.M{
      "_id":   bson.ObjectIdHex(id),
      "exist": true,
    }).One(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUserNotFound)
    }
    return nil, err
  }

  // 存到缓存
  if !rok {
    cache.UserSet(id, data)
  }

  return data.GetMap(FindRef(mg), refs...), nil
}

// 用户登录并返回用户id
func UserLogin(username string, password string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringIsEmpty(username, errgo.ErrUsernameEmpty)
  errgo.ErrorIfStringIsEmpty(password, errgo.ErrPasswordEmpty)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.User)

  err = mg.C(model.UserCollection).Find(bson.M{
    "username": username,
    "password": password,
    "exist":    true,
  }).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUsernameOrPasswordError)
    }
    return nil, err
  } else {

    // create jwt
    departmentId := data.Department.Id.(bson.ObjectId)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "iss":          "workerbook",
      "iat":          time.Now().Unix(),
      "departmentId": departmentId.Hex(),
      "id":           data.Id.Hex(),
      "role":         data.Role,
    })

    tokenStr, _ := token.SignedString(conf.JwtSecret)

    return gin.H{
      "data": tokenStr,
    }, nil
  }
}

// 获取用户列表
func GetUsersList(skip int, limit int, query bson.M, refs ... string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  if skip != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.User)
  query["exist"] = true

  // find it
  if skip == 0 && limit == 0 {
    err = mg.C(model.UserCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = mg.C(model.UserCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUserNotFound)
    }
    return nil, err
  }

  // result
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(FindRef(mg), refs...))
  }

  if skip == 0 && limit == 0 {
    return gin.H{
      "list": list,
    }, nil
  }

  // get count
  count, err := mg.C(model.UserCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrUserNotFound)
  }

  return gin.H{
    "list":  list,
    "count": count,
    "skip":  skip,
    "limit": limit,
  }, nil
}
