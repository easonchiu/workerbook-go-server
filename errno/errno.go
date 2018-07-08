package errno

import (
  "net/http"
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
  ErrSameNickname            = "100010"
  ErrCreateUserFailed        = "100011"
  ErrLoginFailed             = "100012"

  // 部门相关错误 101xxx
  ErrDepartmentIdError      = "101001"
  ErrDepartmentNotFound     = "101002"
  ErrDepartmentNameEmpty    = "101003"
  ErrSameDepartmentName     = "101004"
  ErrCreateDepartmentFailed = "101005"

  // 日报相关错误 102xxx
  ErrDailyIdError  = "102001"
  ErrDailyNotFound = "102002"

  // 系统级错误 200xxx
  ErrSkipRange  = "200001"
  ErrLimitRange = "200002"

  // 默认错误
  ErrServerError = "999999"
)

type ErrType struct {
  Message string
  Status  int
  Code    string
}

// 默认错误
var DefaultType = ErrType{"系统错误", http.StatusInternalServerError, ErrServerError}

// 错误列表
var Error = map[string]ErrType{
  ErrUserIdError:             {"用户id错误", http.StatusOK, ""},
  ErrUserNotFound:            {"找不到该用户", http.StatusOK, ""},
  ErrUserReLogin:             {"请重新登录", http.StatusUnauthorized, ""},
  ErrUsernameOrPasswordError: {"帐号或密码错误", http.StatusOK, ""},
  ErrUsernameEmpty:           {"帐号不能为空", http.StatusOK, ""},
  ErrNicknameEmpty:           {"姓名不能为空", http.StatusOK, ""},
  ErrPasswordEmpty:           {"密码不能为空", http.StatusOK, ""},
  ErrLoginFailed:             {"登录失败，请重试", http.StatusOK, ""},
  ErrUserList:                {"获取用户列表失败", http.StatusOK, ""},
  ErrSameUsername:            {"该帐号已存在", http.StatusOK, ""},
  ErrSameNickname:            {"已存在相同的姓名", http.StatusOK, ""},
  ErrCreateUserFailed:        {"创建用户失败", http.StatusOK, ""},

  ErrDailyIdError:  {"日报id错误", http.StatusOK, ""},
  ErrDailyNotFound: {"找不到相关日报", http.StatusOK, ""},

  ErrDepartmentIdError:      {"部门id错误", http.StatusOK, ""},
  ErrDepartmentNotFound:     {"找不到相关部门", http.StatusOK, ""},
  ErrDepartmentNameEmpty:    {"部门名称不能为空", http.StatusOK, ""},
  ErrSameDepartmentName:     {"已存在相同部门名称", http.StatusOK, ""},
  ErrCreateDepartmentFailed: {"创建部门失败", http.StatusOK, ""},

  ErrSkipRange:  {"skip取值范围错误", http.StatusInternalServerError, ""},
  ErrLimitRange: {"limit取值范围错误", http.StatusInternalServerError, ""},
}
