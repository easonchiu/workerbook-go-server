package controller

import (
  `errors`
  `fmt`
  `github.com/gin-gonic/gin`
  `github.com/tidwall/gjson`
  `gopkg.in/mgo.v2/bson`
  `net/http`
  `strconv`
)

type Context struct {
  Ctx     *gin.Context
  RawData []byte
}

// create a new context
func CreateCtx(c *gin.Context) *Context {
  bytes, _ := c.GetRawData()

  ctx := &Context{
    Ctx:     c,
    RawData: bytes,
  }

  return ctx
}

// check value is less then some int value.
func (c *Context) CheckIntIsLessThen(val int, lessVal int, errString string) bool {
  if val < lessVal {
    c.Error(errString, 1)
    return false
  }
  return true
}

// check value is more then some int value.
func (c *Context) CheckIntIsMoreThen(val int, moreVal int, errString string) bool {
  if val > moreVal {
    c.Error(errString, 1)
    return false
  }
  return true
}

// check value is objectId or not
func (c *Context) CheckIsObjectIdHex(id string, errString string) bool {
  if !bson.IsObjectIdHex(id) {
    c.Error(errString, 9)
    return false
  }
  return true
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

// error handle
func (c *Context) Error(errString string, errCode int) {

  if errCode == 0 {
    errCode = 1
  }

  // create an error.
  err := errors.New(errString)

  fmt.Println()
  fmt.Println(" >>> ERROR:", err.Error())
  fmt.Println(" >>> ERROR CODE:", errCode)
  fmt.Println(" >>> REQUEST METHOD:", c.Ctx.Request.Method)
  fmt.Println(" >>> REQUEST URL:", c.Ctx.Request.URL.String())
  fmt.Println(" >>> USER AGENT:", c.Ctx.Request.UserAgent())
  fmt.Println(" >>> USER AUTH:", c.Ctx.Request.Header.Get("authorization"))
  fmt.Println()


  c.Ctx.JSON(http.StatusOK, gin.H{
    "msg":  err.Error(),
    "code": errCode,
  })

}

// forbidden handle
func (c *Context) Forbidden() {
  c.Ctx.JSON(http.StatusForbidden, gin.H{
    "msg":  "forbidden",
    "code": http.StatusForbidden,
    "data": nil,
  })
  c.Ctx.Abort()
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
