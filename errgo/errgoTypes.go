package errgo

import "net/http"

// 错误数据的结构
type errType struct {
  Message string
  Status  int
  Code    string
}

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
  ErrDeleteProjectFailed       = "102015"

  // 任务相关错误 103xxx
  ErrMissionIdError            = "103001"
  ErrMissionNotFound           = "103002"
  ErrMissionNameEmpty          = "103003"
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
  ErrDeleteProjectFailed:       {"删除项目失败", http.StatusOK, ""},

  ErrMissionIdError:            {"任务id错误", http.StatusOK, ""},
  ErrMissionNotFound:           {"找不到该任务", http.StatusOK, ""},
  ErrMissionNameEmpty:          {"任务名称不能为空", http.StatusOK, ""},
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
