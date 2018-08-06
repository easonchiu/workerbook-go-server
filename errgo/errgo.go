package errgo

import (
  "errors"
  "gopkg.in/mgo.v2/bson"
  "time"
)

// 错误栈
type Stack []error

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

// 创建
func Create() *Stack {
  return new(Stack)
}

// 判断int是否小于一个值
func (s *Stack) ErrorIfIntLessThen(val int, min int, errNo string) error {
  if val < min {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断int是否大于一个值
func (s *Stack) ErrorIfIntMoreThen(val int, max int, errNo string) error {
  if val > max {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断一个值是否为objectId
func (s *Stack) ErrorIfStringNotObjectId(id string, errNo string) error {
  if !bson.IsObjectIdHex(id) {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断字符串是否为空
func (s *Stack) ErrorIfStringIsEmpty(str string, errNo string) error {
  if str == "" {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断int是否为0
func (s *Stack) ErrorIfIntIsZero(val int, errNo string) error {
  if val == 0 {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断length是否小于
func (s *Stack) ErrorIfLenLessThen(str string, length int, errNo string) error {
  if len([]rune(str)) < length {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断length是否大于
func (s *Stack) ErrorIfLenMoreThen(str string, length int, errNo string) error {
  if len([]rune(str)) > length {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断时间是否早于
func (s *Stack) ErrorIfTimeEarlierThen(t time.Time, t2 time.Time, errNo string) error {
  if t.Before(t2) == true {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 判断时间是否晚于
func (s *Stack) ErrorIfTimeLaterThen(t time.Time, t2 time.Time, errNo string) error {
  if t.After(t2) == true {
    err := errors.New(errNo)
    *s = append(*s, err)
    return err
  }
  return nil
}

// 清空错误栈
func (s *Stack) ClearErrorStack() {
  *s = nil
}

// 弹出栈中的第一个错误(默认情况下弹出后就清空栈了)
func (s *Stack) PopError(clear ... bool) error {
  if len(*s) > 0 {
    first := (*s)[0]
    if clear != nil && clear[0] == false {
      *s = (*s)[1:]
    } else {
      s.ClearErrorStack()
    }
    return first
  }
  return nil
}
