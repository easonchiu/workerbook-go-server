package errgo

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// 错误栈
var errStack []error

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

// // 处理ErrorIf相关的错误
// func HandleError(handle func(err interface{})) bool {
//   if len(errStack) > 0 {
//     first := errStack[0]
//     errStack = errStack[1:]
//     handle(first)
//     return true
//   }
//   return false
// }

// 清空错误栈
func ClearErrorStack() {
  errStack = nil
}

// 弹出栈中的第一个错误
func PopError() error {
  if len(errStack) > 0 {
    first := errStack[0]
    errStack = errStack[1:]
    return first
  }
  return nil
}