package controller

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/tidwall/gjson"
  "net/http"
  "strconv"
  "strings"
  "time"
  "workerbook/errgo"
)

type Context struct {
  Ctx     *gin.Context
  RawData []byte
}

// create a new context
func CreateCtx(c *gin.Context) *Context {
  bytes, _ := c.GetRawData()
  return &Context{c, bytes}
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

  // 根据错误号获取错误内容（错误号是个string或error）
  err := errgo.Get(errNo)

  fmt.Println()
  fmt.Println(" >>> ERROR:", err.Message)
  fmt.Println(" >>> ERROR CODE:", err.Code)
  fmt.Println(" >>> REQUEST METHOD:", c.Ctx.Request.Method)
  fmt.Println(" >>> REQUEST URL:", c.Ctx.Request.URL.String())
  fmt.Println(" >>> USER AGENT:", c.Ctx.Request.UserAgent())
  fmt.Println(" >>> USER AUTH:", c.Ctx.Request.Header.Get("authorization"))
  fmt.Println()

  // 清除错误栈
  errgo.ClearErrorStack()

  c.Ctx.JSON(err.Status, gin.H{
    "msg":  err.Message,
    "code": err.Code,
    "data": nil,
  })
}

// get body by string
func (c *Context) getRaw(key string) (string, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return strings.TrimSpace(res.Str), res.Exists()
}

func (c *Context) getRawArray(key string) ([]gjson.Result, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Array(), res.Exists()
}

func (c *Context) getRawTime(key string) (time.Time, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Time(), res.Exists()
}

// get body by int
func (c *Context) getRawInt(key string) (int, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return int(res.Int()), res.Exists()
}

// get body by bool
func (c *Context) getRawBool(key string) (bool, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Bool(), res.Exists()
}

// get body by JSON
func (c *Context) getRawJSON(key string) (gjson.Result, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res, res.Exists()
}

// get params by string
func (c *Context) getParam(key string) (string, bool) {
  res, ok := c.Ctx.Params.Get(key)
  return res, ok
}

// get params by int
func (c *Context) getParamInt(key string) (int, bool) {
  res, ok := c.Ctx.Params.Get(key)
  intRes, _ := strconv.Atoi(res)
  return intRes, ok
}

// get params by bool
func (c *Context) getParamBool(key string) (bool, bool) {
  res, ok := c.Ctx.Params.Get(key)
  return res == "true", ok
}

// get query by string
func (c *Context) getQuery(key string) (string, bool) {
  res, ok := c.Ctx.GetQuery(key)
  return res, ok
}

func (c *Context) getQueryDefault(key string, def string) string {
  val, ok := c.getQuery(key)
  if !ok {
    return def
  }
  return val
}

// get query by int
func (c *Context) getQueryInt(key string) (int, bool) {
  res, ok := c.Ctx.GetQuery(key)
  if !ok {
    return 0, false
  }
  intRes, err := strconv.Atoi(res)
  if err != nil {
    return 0, false
  }
  return intRes, true
}

func (c *Context) getQueryIntDefault(key string, def int) int {
  val, ok := c.getQueryInt(key)
  if !ok {
    return def
  }
  return val
}

// get query by bool
func (c *Context) getQueryBool(key string) (bool, bool) {
  res, ok := c.Ctx.GetQuery(key)
  if !ok {
    return false, false
  }
  return res == "true", true
}

func (c *Context) getQueryBoolDefault(key string, def bool) bool {
  val, ok := c.getQueryBool(key)
  if !ok {
    return def
  }
  return val
}

// get query by JSON
func (c *Context) getQueryJSON(key string) (gjson.Result, bool) {
  res, ok := c.Ctx.GetQuery(key)
  if !ok {
    return gjson.Result{}, false
  }
  return gjson.Parse(res), true
}

// get value by string
func (c *Context) get(key string) (string, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return "", false
  }
  return res.(string), true
}

// get value by string
func (c *Context) getInt(key string) (int, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return 0, false
  }
  return res.(int), true
}

// get value by string
func (c *Context) getBool(key string) (bool, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return false, false
  }
  return res.(bool), true
}
