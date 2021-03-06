package service

import (
  "errors"
  "github.com/jwt-go"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/util"
)

// 创建用户
func CreateUser(ctx *context.New, data model.User, departmentId string) error {

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data.Editor = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.UserCollection,
    Id:         bson.ObjectIdHex(ownUserId),
  }
  data.EditTime = time.Now()

  // 用户状态
  data.Status = 1

  // check
  ctx.Errgo.ErrorIfStringIsEmpty(data.NickName, errgo.ErrNicknameEmpty)
  ctx.Errgo.ErrorIfLenLessThen(data.NickName, 2, errgo.ErrNicknameTooShort)
  ctx.Errgo.ErrorIfLenMoreThen(data.NickName, 14, errgo.ErrNicknameTooLong)
  ctx.Errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  ctx.Errgo.ErrorIfStringIsEmpty(data.Title, errgo.ErrUserTitleIsEmpty)
  ctx.Errgo.ErrorIfLenMoreThen(data.Title, 14, errgo.ErrUserTitleTooLong)
  if data.Role != 1 && data.Role != 2 && data.Role != 3 {
    return errors.New(errgo.ErrUserRoleError)
  }
  ctx.Errgo.ErrorIfStringIsEmpty(data.UserName, errgo.ErrUsernameEmpty)
  ctx.Errgo.ErrorIfLenLessThen(data.UserName, 6, errgo.ErrUsernameTooShort)
  ctx.Errgo.ErrorIfLenMoreThen(data.UserName, 14, errgo.ErrUsernameTooLong)
  ctx.Errgo.ErrorIfStringIsEmpty(data.Password, errgo.ErrPasswordEmpty)

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // username must be the only.
  count, err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
    "username": data.UserName,
    "exist":    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameUsername)
  }

  // nickname must be the only.
  count, err = ctx.MgoDB.C(model.UserCollection).Find(bson.M{
    "nickname": data.NickName,
    "exist":    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameNickname)
  }

  // department must be exist.
  if _, err := GetDepartmentInfoById(ctx, departmentId); err != nil {
    return err
  }

  // save department id
  data.Department = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.DepartmentCollection,
    Id:         bson.ObjectIdHex(departmentId),
  }

  // insert it.
  err = ctx.MgoDB.C(model.UserCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount(ctx)

  return nil
}

// 更新用户
func UpdateUser(ctx *context.New, id string, data bson.M) error {

  if data == nil {
    return errors.New(errgo.ErrServerError)
  }

  // 限制更新字段
  util.Only(
    data,
    util.Keys{
      "nickname":       util.TypeString,
      "email":          util.TypeString,
      "jobNumber":      util.TypeString,
      "department.$id": util.TypeString,
      "mobile":         util.TypeString,
      "role":           util.TypeInt,
      "title":          util.TypeString,
      "status":         util.TypeInt,
      "exist":          util.TypeBool,
      "editor":         util.TypeString,
      "editTime":       util.TypeTime,
    },
  )

  // check
  ctx.Errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if nickname, ok := util.GetString(data, "nickname"); ok {
    ctx.Errgo.ErrorIfLenLessThen(nickname, 2, errgo.ErrNicknameTooShort)
    ctx.Errgo.ErrorIfLenMoreThen(nickname, 14, errgo.ErrNicknameTooLong)
  }

  if id, ok := util.GetString(data, "department.$id"); ok {
    ctx.Errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)
  }

  if title, ok := util.GetString(data, "title"); ok {
    ctx.Errgo.ErrorIfLenMoreThen(title, 14, errgo.ErrUserTitleTooLong)
  }

  if role, ok := util.GetInt(data, "role"); ok {
    if role != 1 && role != 2 && role != 3 {
      return errors.New(errgo.ErrUserRoleError)
    }
  }

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // 姓名唯一
  if nickname, ok := util.GetString(data, "nickname"); ok {
    count, err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
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
  if departmentId, ok := util.GetString(data, "department.$id"); ok {
    _, err := GetDepartmentInfoById(ctx, departmentId)

    if err != nil {
      return err
    }

    data["department.$id"] = bson.ObjectIdHex(departmentId)
  }

  // 先要清缓存，清成功后才可以更新数据
  err := cache.UserDel(ctx.Redis, id)

  if err != nil {
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteUserFailed)
    }
    return errors.New(errgo.ErrUpdateUserFailed)
  }

  // 更新数据
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data["editor.$id"] = bson.ObjectIdHex(ownUserId)
  data["editTime"] = time.Now()

  err = ctx.MgoDB.C(model.UserCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrUserNotFound)
    }
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteUserFailed)
    }
    return errors.New(errgo.ErrUpdateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount(ctx)

  return nil
}

// 根据id查找用户信息
func GetUserInfoById(ctx *context.New, id string) (*model.User, error) {

  // check
  ctx.Errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  data := new(model.User)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.UserGet(ctx.Redis, id, data)
  if !rok {
    err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
      "_id":   bson.ObjectIdHex(id),
      "exist": true,
    }).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrUserNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.UserSet(ctx.Redis, data)
  }

  return data, nil
}

// 用户登录并返回用户token
func UserLogin(ctx *context.New, username string, password string) (string, error) {

  // check
  ctx.Errgo.ErrorIfStringIsEmpty(username, errgo.ErrUsernameEmpty)
  ctx.Errgo.ErrorIfStringIsEmpty(password, errgo.ErrPasswordEmpty)

  if err := ctx.Errgo.PopError(); err != nil {
    return "", err
  }

  data := new(model.User)

  err := ctx.MgoDB.C(model.UserCollection).Find(bson.M{
    "username": username,
    "password": password,
    "exist":    true,
  }).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return "", errors.New(errgo.ErrUsernameOrPasswordError)
    }
    return "", err
  } else {

    // create jwt
    departmentId := bson.NewObjectId()
    if data.Department.Id != nil {
      departmentId = data.Department.Id.(bson.ObjectId)
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "iss":                  "WorkerBook",
      "iat":                  time.Now().Unix(),
      conf.OWN_DEPARTMENT_ID: departmentId.Hex(),
      conf.OWN_USER_ID:       data.Id.Hex(),
      conf.OWN_ROLE:          data.Role,
    })

    tokenStr, _ := token.SignedString(conf.JwtSecret)

    return tokenStr, nil
  }
}

// 获取用户列表(当limit都为0时，查找全部)
func GetUsersList(ctx *context.New, skip int, limit int, query bson.M) (*model.UserList, error) {

  // check
  if skip != 0 {
    ctx.Errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    ctx.Errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    ctx.Errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  data := new([]*model.User)
  query["exist"] = true

  // find it
  var err error
  if skip == 0 && limit == 0 {
    err = ctx.MgoDB.C(model.UserCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = ctx.MgoDB.C(model.UserCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUserNotFound)
    }
    return nil, err
  }

  // result
  if skip == 0 && limit == 0 {
    return &model.UserList{
      Count: len(*data),
      List:  *data,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.UserCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrUserNotFound)
  }

  return &model.UserList{
    List:  *data,
    Count: count,
    Skip:  skip,
    Limit: limit,
  }, nil
}
