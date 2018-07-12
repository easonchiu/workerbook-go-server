package errgo

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "net/http"
  "time"
)

const (
  // 用户相关错误 100xxx
  ErrUserIdError             = "100001"
  ErrUserNotFound            = "100002"
  ErrUserReLogin             = "100003"
  ErrUsernameOrPasswordError = "100004"
  ErrUserList                = "100005"
  ErrUsernameEmpty           = "100006"
  ErrNicknameEmpty           = "100007"
  ErrPasswordEmpty           = "100008"
  ErrSameUsername            = "100009"
  ErrUsernameTooShort        = "100010"
  ErrUsernameTooLong         = "100011"
  ErrSameNickname            = "100013"
  ErrNicknameTooShort        = "100014"
  ErrNicknameTooLong         = "100015"
  ErrUserTitleIsEmpty        = "100016"
  ErrUserTitleTooLong        = "100017"
  ErrUserRoleError           = "100018"
  ErrCreateUserFailed        = "100019"
  ErrUpdateUserFailed        = "100020"
  ErrLoginFailed             = "100021"

  // 部门相关错误 101xxx
  ErrDepartmentIdError      = "101001"
  ErrDepartmentNotFound     = "101002"
  ErrDepartmentNameEmpty    = "101003"
  ErrSameDepartmentName     = "101004"
  ErrCreateDepartmentFailed = "101005"
  ErrUpdateDepartmentFailed = "101006"

  // 项目相关错误 102xxx
  ErrProjectIdError            = "102001"
  ErrProjectNotFound           = "102002"
  ErrProjectNameEmpty          = "102003"
  ErrProjectNameTooShort       = "102005"
  ErrProjectNameTooLong        = "102006"
  ErrProjectDeadlineEmpty      = "102007"
  ErrProjectDeadlineTooSoon    = "102008"
  ErrProjectDepartmentsEmpty   = "102009"
  ErrProjectDepartmentNotFound = "102010"
  ErrProjectWeightError        = "102011"
  ErrProjectDescriptionTooLong = "102012"
  ErrCreateProjectFailed       = "102013"
  ErrUpdateProjectFailed       = "102014"

  // 任务相关错误 103xxx
  ErrMissionIdError            = "103001"
  ErrMissionNotFound           = "103002"
  ErrMissionNameEmpty          = "103003"
  ErrMissionNameTooShort       = "103004"
  ErrMissionNameTooLong        = "103005"
  ErrMissionDeadlineEmpty      = "103006"
  ErrMissionDeadlineTooSoon    = "103007"
  ErrMissionDeadlineTooLate    = "103008"
  ErrMissionDescriptionTooLong = "103009"
  ErrCreateMissionFailed       = "103010"
  ErrUpdateMissionFailed       = "103011"

  // 日报相关错误 105xxx
  ErrDailyIdError  = "105001"
  ErrDailyNotFound = "105002"

  // 系统级错误 200xxx
  ErrSkipRange  = "200001"
  ErrLimitRange = "200002"

  // 默认错误
  ErrServerError = "999999"
)

type errType struct {
  Message string
  Status  int
  Code    string
}

// 默认错误
var DefaultError = errType{"系统错误", http.StatusInternalServerError, ErrServerError}

// 错误列表
var Error = map[string]errType{
  ErrUserIdError:             {"用户id错误", http.StatusOK, ""},
  ErrUserNotFound:            {"找不到该用户", http.StatusOK, ""},
  ErrUserReLogin:             {"请重新登录", http.StatusUnauthorized, ""},
  ErrUsernameOrPasswordError: {"帐号或密码错误", http.StatusOK, ""},
  ErrUsernameEmpty:           {"帐号不能为空", http.StatusOK, ""},
  ErrNicknameEmpty:           {"姓名不能为空", http.StatusOK, ""},
  ErrPasswordEmpty:           {"密码不能为空", http.StatusOK, ""},
  ErrSameUsername:            {"已有相同的姓名", http.StatusOK, ""},
  ErrLoginFailed:             {"登录失败，请重试", http.StatusOK, ""},
  ErrUserList:                {"获取用户列表失败", http.StatusOK, ""},
  ErrUsernameTooShort:        {"帐号不得少于6个字", http.StatusOK, ""},
  ErrUsernameTooLong:         {"帐号不得多于14个字", http.StatusOK, ""},
  ErrSameNickname:            {"已存在相同的姓名", http.StatusOK, ""},
  ErrNicknameTooShort:        {"姓名不得于少2个字", http.StatusOK, ""},
  ErrNicknameTooLong:         {"姓名不得多于14个字", http.StatusOK, ""},
  ErrUserTitleIsEmpty:        {"职称不能为空", http.StatusOK, ""},
  ErrUserTitleTooLong:        {"职称不得多于14个字", http.StatusOK, ""},
  ErrUserRoleError:           {"用户职位错误", http.StatusOK, ""},
  ErrCreateUserFailed:        {"创建用户失败", http.StatusOK, ""},
  ErrUpdateUserFailed:        {"更新用户失败", http.StatusOK, ""},

  ErrProjectIdError:            {"项目id错误", http.StatusOK, ""},
  ErrProjectNotFound:           {"找不到该项目", http.StatusOK, ""},
  ErrProjectNameEmpty:          {"项目名称不能为空", http.StatusOK, ""},
  ErrProjectNameTooShort:       {"项目名称不得少于4个字", http.StatusOK, ""},
  ErrProjectNameTooLong:        {"项目名称不得多于15个字", http.StatusOK, ""},
  ErrProjectDeadlineEmpty:      {"截至时间不能不空", http.StatusOK, ""},
  ErrProjectDeadlineTooSoon:    {"截至时间不能早于当前时间", http.StatusOK, ""},
  ErrProjectDepartmentsEmpty:   {"参与部门不能为空", http.StatusOK, ""},
  ErrProjectDepartmentNotFound: {"没有找到相关的参与部门", http.StatusOK, ""},
  ErrProjectWeightError:        {"权重号错误", http.StatusOK, ""},
  ErrProjectDescriptionTooLong: {"项目说明不能多于500字", http.StatusOK, ""},
  ErrCreateProjectFailed:       {"创建项目失败", http.StatusOK, ""},
  ErrUpdateProjectFailed:       {"更新项目失败", http.StatusOK, ""},

  ErrMissionIdError:            {"任务id错误", http.StatusOK, ""},
  ErrMissionNotFound:           {"找不到该任务", http.StatusOK, ""},
  ErrMissionNameEmpty:          {"任务名称不能为空", http.StatusOK, ""},
  ErrMissionNameTooShort:       {"任务名称不得少于4个字", http.StatusOK, ""},
  ErrMissionNameTooLong:        {"任务名称不得多于15个字", http.StatusOK, ""},
  ErrMissionDeadlineEmpty:      {"截至时间不能不空", http.StatusOK, ""},
  ErrMissionDeadlineTooSoon:    {"截至时间不能早于当前时间", http.StatusOK, ""},
  ErrMissionDeadlineTooLate:    {"截至时间不能晚于项目时间", http.StatusOK, ""},
  ErrMissionDescriptionTooLong: {"任务说明不能多于500字", http.StatusOK, ""},
  ErrCreateMissionFailed:       {"创建任务失败", http.StatusOK, ""},
  ErrUpdateMissionFailed:       {"更新任务失败", http.StatusOK, ""},

  ErrDailyIdError:  {"日报id错误", http.StatusOK, ""},
  ErrDailyNotFound: {"找不到相关日报", http.StatusOK, ""},

  ErrDepartmentIdError:      {"部门id错误", http.StatusOK, ""},
  ErrDepartmentNotFound:     {"找不到相关部门", http.StatusOK, ""},
  ErrDepartmentNameEmpty:    {"部门名称不能为空", http.StatusOK, ""},
  ErrSameDepartmentName:     {"已存在相同部门名称", http.StatusOK, ""},
  ErrCreateDepartmentFailed: {"创建部门失败", http.StatusOK, ""},
  ErrUpdateDepartmentFailed: {"更新部门失败", http.StatusOK, ""},

  ErrSkipRange:  {"skip取值范围错误", http.StatusInternalServerError, ""},
  ErrLimitRange: {"limit取值范围错误", http.StatusInternalServerError, ""},
}

// 根据错误码换取错误信息
func Get(no interface{}) errType {
  errStrNo := ""

  switch no.(type) {
  case string:
    errStrNo = no.(string)
  case error:
    errStrNo = no.(error).Error()
  }

  if errStrNo != "" && Error[errStrNo].Message != "" {
    err := Error[errStrNo]
    err.Code = errStrNo
    return err
  }

  return DefaultError
}

// 错误栈
var errStack []error

// 判断int是否小于一个值
func ErrorIfIntLessThen(val int, min int, errNo string) error {
  if val < min {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断int是否大于一个值
func ErrorIfIntMoreThen(val int, max int, errNo string) error {
  if val > max {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断一个值是否为objectId
func ErrorIfStringNotObjectId(id string, errNo string) error {
  if !bson.IsObjectIdHex(id) {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断字符串是否为空
func ErrorIfStringIsEmpty(str string, errNo string) error {
  if str == "" {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断int是否为0
func ErrorIfIntIsZero(val int, errNo string) error {
  if val == 0 {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断length是否小于
func ErrorIfLenLessThen(str string, length int, errNo string) error {
  if len([]rune(str)) < length {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断length是否大于
func ErrorIfLenMoreThen(str string, length int, errNo string) error {
  if len([]rune(str)) > length {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断时间是否早于
func ErrorIfTimeEarlierThen(t time.Time, t2 time.Time, errNo string) error {
  if t.Before(t2) == true {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 判断时间是否晚于
func ErrorIfTimeLaterThen(t time.Time, t2 time.Time, errNo string) error {
  if t.After(t2) == true {
    err := errors.New(errNo)
    errStack = append(errStack, err)
    return err
  }
  return nil
}

// 处理ErrorIf相关的错误
func HandleError(handle func(err interface{})) bool {
  if len(errStack) > 0 {
    first := errStack[0]
    errStack = errStack[1:]
    handle(first)
    return true
  }
  return false
}

// 清空ErrorStack
func ClearErrorStack() {
  errStack = nil
}

// 获取栈中的第一个错误
func PopError() error {
  if len(errStack) > 0 {
    first := errStack[0]
    errStack = errStack[1:]
    return first
  }
  return nil
}