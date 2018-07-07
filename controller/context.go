package controller

import (
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "net/http"
  "strconv"
  "workerbook/errno"
)

type Context struct {
  Ctx *gin.Context
  Err error // 这个err用于错误的合并处理，在下列代码中只会存一次err，并在处理后设为nil
  RawData []byte
}

// create a new context
func CreateCtx(c *gin.Context) *Context {
  bytes, _ := c.GetRawData()
  return &Context{c, nil, bytes}
}

// check value is less then some int value.
func (c *Context) ErrorIfIntLessThen(val int, min int, errNo string) error {
  if val < min {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check value is more then some int value.
func (c *Context) ErrorIfIntMoreThen(val int, max int, errNo string) error {
  if val > max {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check value is objectId or not
func (c *Context) ErrorIfStringNotObjectId(id string, errNo string) error {
  if !bson.IsObjectIdHex(id) {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check value is empty string
func (c *Context) ErrorIfStringIsEmpty(str string, errNo string) error {
  if str == "" {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check length of string is less then.
func (c *Context) ErrorIfLenLessThen(str string, length int, errNo string) error {
  if len(str) < length {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check length of string is more then.
func (c *Context) ErrorIfLenMoreThen(str string, length int, errNo string) error {
  if len(str) > length {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// check value is mgo not-found
func (c *Context) ErrorIfMgoNotFound(err error, errNo string) error {
  if err == mgo.ErrNotFound {
    err := errors.New(errNo)
    if c.Err == nil {
      c.Err = err
    }
    return err
  }
  return nil
}

// 处理ErrorIf相关的错误
func (c *Context) HandleErrorIf() bool {
  if c.Err != nil {
    c.Error(c.Err)
    return true
  }
  return false
}

// success handle
func (c *Context) Success(data gin.H) {
  respH := gin.H{
    "msg":  "ok",
    "code": 0,
  }

  if len(data) > 1 { // Almost the length is more than 1, so just check it first.
    respH["data"] = data
  } else if data["data"] != nil {
    respH["data"] = data["data"]
  } else if data != nil && len(data) > 0 {
    respH["data"] = data
  }

  status := http.StatusOK

  if data == nil {
    status = http.StatusNoContent
  }

  c.Ctx.JSON(status, respH)
}

// 处理错误
func (c *Context) Error(errNo interface{}) {

  errStrNo := ""

  switch errNo.(type) {
  case string:
    errStrNo = errNo.(string)
  case error:
    errStrNo = errNo.(error).Error()
  default:
    errStrNo = errno.DefaultType.Code
  }

  err := errno.Error[errStrNo]
  err.Code = errStrNo

  c.Err = nil

  if err.Code == "" {
    err = errno.DefaultType
  }

  fmt.Println()
  fmt.Println(" >>> ERROR:", err.Message)
  fmt.Println(" >>> ERROR CODE:", err.Code)
  fmt.Println(" >>> REQUEST METHOD:", c.Ctx.Request.Method)
  fmt.Println(" >>> REQUEST URL:", c.Ctx.Request.URL.String())
  fmt.Println(" >>> USER AGENT:", c.Ctx.Request.UserAgent())
  fmt.Println(" >>> USER AUTH:", c.Ctx.Request.Header.Get("authorization"))
  fmt.Println()

  c.Ctx.JSON(err.Status, gin.H{
    "msg":  err.Message,
    "code": err.Code,
  })
}

// get post data by string
func (c *Context) getRaw(key string) string {
  res := gjson.GetBytes(c.RawData, key)
  return res.Str
}

// get post data by int
func (c *Context) getRawInt(key string) int {
  res := gjson.GetBytes(c.RawData, key)
  return int(res.Int())
}

// get post data by bool
func (c *Context) getRawBool(key string) bool {
  res := gjson.GetBytes(c.RawData, key)
  return res.Str == "true"
}

// get post data by JSON
func (c *Context) getRawJSON(key string) gjson.Result {
  res := gjson.GetBytes(c.RawData, key)
  return res
}

// get params by string
func (c *Context) getParam(key string) string {
  res := c.Ctx.Param(key)
  return res
}

// get params by int
func (c *Context) getParamInt(key string) int {
  res := c.Ctx.Param(key)
  intRes, _ := strconv.Atoi(res)
  return intRes
}

// get params by bool
func (c *Context) getParamBool(key string) bool {
  res := c.Ctx.Param(key)
  return res == "true"
}

// get params by JSON
func (c *Context) getParamJSON(key string) gjson.Result {
  res := c.Ctx.Param(key)
  return gjson.Parse(res)
}

// get query by string
func (c *Context) getQuery(key string) string {
  res, _ := c.Ctx.GetQuery(key)
  return res
}

// get query by int
func (c *Context) getQueryInt(key string) int {
  res, _ := c.Ctx.GetQuery(key)
  intRes, _ := strconv.Atoi(res)
  return intRes
}

// get query by bool
func (c *Context) getQueryBool(key string) bool {
  res, exist := c.Ctx.GetQuery(key)
  if !exist {
    return false
  }
  return res == "true"
}

// get query by JSON
func (c *Context) getQueryJSON(key string) gjson.Result {
  res, exist := c.Ctx.GetQuery(key)
  if !exist {
    return gjson.Result{}
  }
  return gjson.Parse(res)
}

// get value by string
func (c *Context) get(key string) string {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return ""
  }
  return res.(string)
}

// get value by string
func (c *Context) getInt(key string) int {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return 0
  }
  intRes, _ := strconv.Atoi(res.(string))
  return intRes
}

// get value by string
func (c *Context) getBool(key string) bool {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return false
  }
  return res.(string) == "true"
}

// get value by JSON
func (c *Context) getJSON(key string) gjson.Result {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return gjson.Result{}
  }
  return gjson.Parse(res.(string))
}
