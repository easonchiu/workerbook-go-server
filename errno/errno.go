package errno

import "net/http"

// 注意：添加错误类型时不要颠倒顺序或在同一类别的中间插入
const (
  // 用户相关错误 100xxx
  ErrorUserIdError = 100001 + iota
  ErrorUserNotFound
  ErrorUserReLogin
  ErrorUserList
  ErrorUsernameEmpty
  ErrorPasswordEmpty
  ErrorSameUsername
  ErrorCreateUserFailed
  ErrorLoginFailed

  // 部门相关错误 101xxx
  ErrorDepartmentIdError = 101001 + iota
  ErrorDepartmentNotFound
  ErrorDepartmentNameEmpty
  ErrorCreateDepartmentFailed

  // 日报相关错误 102xxx
  ErrorDailyIdError = 102001 + iota
  ErrorDailyNotFound

  // 系统级错误 200xxx
  ErrorSkipRange = 200001 + iota
  ErrorLimitRange

  // 默认错误
  ErrorServerError = 999999
)

type ErrType struct {
  Message string
  Status  int
  Code    int
}

// 默认错误
var DefaultType = ErrType{"系统错误", http.StatusInternalServerError, ErrorServerError}

// 错误列表
var Error = map[int]ErrType{
  ErrorUserIdError: {"用户id错误", http.StatusOK, 0},
  ErrorUserNotFound: {"找不到该用户", http.StatusOK, 0},
  ErrorUserReLogin:  {"请重新登录", http.StatusUnauthorized, 0},
  ErrorUsernameEmpty: {"用户名不能为空", http.StatusOK, 0},
  ErrorPasswordEmpty: {"密码不能为空", http.StatusOK, 0},
  ErrorLoginFailed: {"登录失败，请重试", http.StatusOK, 0},
  ErrorUserList: {"获取用户列表失败", http.StatusOK, 0},
  ErrorSameUsername: {"已存在相同的用户名", http.StatusOK, 0},
  ErrorCreateUserFailed: {"创建用户失败", http.StatusOK, 0},

  ErrorDailyIdError: {"日报id错误", http.StatusOK, 0},
  ErrorDailyNotFound: {"找不到相关日报", http.StatusOK, 0},

  ErrorDepartmentIdError: {"部门id错误", http.StatusOK, 0},
  ErrorDepartmentNotFound: {"找不到相关部门", http.StatusOK, 0},
  ErrorDepartmentNameEmpty: {"部门名称不能为空", http.StatusOK, 0},
  ErrorCreateDepartmentFailed: {"创建部门失败", http.StatusOK, 0},

  ErrorSkipRange: {"skip取值范围错误", http.StatusInternalServerError, 0},
  ErrorLimitRange: {"limit取值范围错误", http.StatusInternalServerError, 0},
}
