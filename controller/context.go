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

  c.Ctx.Done()
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

  // 清楚错误栈
  errgo.ClearErrorStack()

  c.Ctx.JSON(err.Status, gin.H{
    "msg":  err.Message,
    "code": err.Code,
    "data": nil,
  })

  c.Ctx.Done()
}

// get body by string
func (c *Context) getRaw(key string) string {
  res := gjson.GetBytes(c.RawData, key)
  return strings.TrimSpace(res.Str)
}

func (c *Context) getRawArray(key string) []gjson.Result {
  res := gjson.GetBytes(c.RawData, key)
  return res.Array()
}

func (c *Context) getRawTime(key string) time.Time {
  res := gjson.GetBytes(c.RawData, key)
  return res.Time()
}

// get body by int
func (c *Context) getRawInt(key string) int {
  res := gjson.GetBytes(c.RawData, key)
  return int(res.Int())
}

// get body by bool
func (c *Context) getRawBool(key string) bool {
  res := gjson.GetBytes(c.RawData, key)
  return res.Str == "true"
}

// get body by JSON
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

func (c *Context) getQueryDefault(key string, def string) string {
  val := c.getQuery(key)
  if val == "" {
    return def
  }
  return val
}

// get query by int
func (c *Context) getQueryInt(key string) int {
  res, _ := c.Ctx.GetQuery(key)
  intRes, _ := strconv.Atoi(res)
  return intRes
}

func (c *Context) getQueryIntDefault(key string, def int) int {
  val := c.getQueryInt(key)
  if val == 0 {
    return def
  }
  return val
}

// get query by bool
func (c *Context) getQueryBool(key string) bool {
  res, exist := c.Ctx.GetQuery(key)
  if !exist {
    return false
  }
  return res == "true"
}

func (c *Context) getQueryBoolDefault(key string, def bool) bool {
  val := c.getQueryBool(key)
  if val == false {
    return def
  }
  return val
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
  return res.(int)
}

// get value by string
func (c *Context) getBool(key string) bool {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return false
  }
  return res.(bool)
}

// get value by JSON
func (c *Context) getJSON(key string) gjson.Result {
  res, exist := c.Ctx.Get(key)
  if !exist {
    return gjson.Result{}
  }
  return gjson.Parse(res.(string))
}
