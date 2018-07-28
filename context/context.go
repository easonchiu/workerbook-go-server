package context

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
  "github.com/tidwall/gjson"
  "gopkg.in/mgo.v2"
  "net/http"
  "strconv"
  "strings"
  "time"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/errgo"
)

type New struct {
  Ctx         *gin.Context
  RawData     []byte
  MgoDB       *mgo.Database
  MgoDBCloser func()
  Redis       redis.Conn
}

func CreateCtx(fn func(*New)) func(*gin.Context) {
  return func(c *gin.Context) {
    if gin.IsDebugging() {
      fmt.Println()
      fmt.Println("------------------------------------------")
      fmt.Println()
    }

    // 创建上下文
    ctx, err := NewCtx(c)

    // 如果创建过程中有报错，返回错误
    if err != nil {
      ctx.Error(err)
      return
    }

    // defer
    defer ctx.Close()

    // 调用控制器
    fn(ctx)
  }
}

// 创建上下文，连接mgo与redis数据库
func NewCtx(c *gin.Context) (*New, error) {
  bytes, _ := c.GetRawData()
  mg, closer, err := db.CloneMgoDB()
  if err != nil {
    fmt.Println("[MGO] 😈 Error")
    return nil, err
  }
  if gin.IsDebugging() {
    fmt.Println("[MGO] 😄 OK")
  }
  rds := db.RedisPool.Get()
  if gin.IsDebugging() {
    fmt.Println("[RDS] 😄 OK")
  }
  return &New{
    c,
    bytes,
    mg,
    closer,
    rds,
  }, nil
}

// 创建不连接数据库的上下文
func NewBaseCtx(c *gin.Context) *New {
  bytes, _ := c.GetRawData()
  return &New{
    Ctx:     c,
    RawData: bytes,
  }
}

// 关闭数据库连接
func (c *New) Close() {
  c.MgoDBCloser()
  if gin.IsDebugging() {
    fmt.Println("[MGO] 👋 CLOSED")
  }
  c.Redis.Close()
  if gin.IsDebugging() {
    fmt.Println("[RDS] 👋 CLOSED")
  }
}

// success handle
func (c *New) Success(data gin.H) {
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
func (c *New) Error(errNo interface{}) {

  // 根据错误号获取错误内容（错误号是个string或error）
  err := errgo.Get(errNo)
  role := c.Ctx.GetInt(conf.OWN_ROLE)

  fmt.Println()
  fmt.Println(" >>> ERROR:", err.Message)
  fmt.Println(" >>> ERROR CODE:", err.Code)
  fmt.Println(" >>> REQUEST METHOD:", c.Ctx.Request.Method)
  fmt.Println(" >>> REQUEST URL:", c.Ctx.Request.URL.String())
  fmt.Println(" >>> USER AGENT:", c.Ctx.Request.UserAgent())
  fmt.Println(" >>> USER ROLE:", role)
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
func (c *New) GetRaw(key string) (string, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return strings.TrimSpace(res.Str), res.Exists()
}

func (c *New) GetRawArray(key string) ([]gjson.Result, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Array(), res.Exists()
}

func (c *New) GetRawTime(key string) (time.Time, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Time(), res.Exists()
}

// get body by int
func (c *New) GetRawInt(key string) (int, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return int(res.Int()), res.Exists()
}

// get body by bool
func (c *New) GetRawBool(key string) (bool, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res.Bool(), res.Exists()
}

// get body by JSON
func (c *New) GetRawJSON(key string) (gjson.Result, bool) {
  res := gjson.GetBytes(c.RawData, key)
  return res, res.Exists()
}

// get params by string
func (c *New) GetParam(key string) (string, bool) {
  res, ok := c.Ctx.Params.Get(key)
  return res, ok
}

// get params by int
func (c *New) GetParamInt(key string) (int, bool) {
  res, ok := c.Ctx.Params.Get(key)
  intRes, _ := strconv.Atoi(res)
  return intRes, ok
}

// get params by bool
func (c *New) GetParamBool(key string) (bool, bool) {
  res, ok := c.Ctx.Params.Get(key)
  return res == "true", ok
}

// get query by string
func (c *New) GetQuery(key string) (string, bool) {
  res, ok := c.Ctx.GetQuery(key)
  return res, ok
}

func (c *New) GetQueryDefault(key string, def string) string {
  val, ok := c.GetQuery(key)
  if !ok {
    return def
  }
  return val
}

// get query by int
func (c *New) GetQueryInt(key string) (int, bool) {
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

func (c *New) GetQueryIntDefault(key string, def int) int {
  val, ok := c.GetQueryInt(key)
  if !ok {
    return def
  }
  return val
}

// get query by bool
func (c *New) GetQueryBool(key string) (bool, bool) {
  res, ok := c.Ctx.GetQuery(key)
  if !ok {
    return false, false
  }
  return res == "true", true
}

func (c *New) GetQueryBoolDefault(key string, def bool) bool {
  val, ok := c.GetQueryBool(key)
  if !ok {
    return def
  }
  return val
}

// get query by JSON
func (c *New) GetQueryJSON(key string) (gjson.Result, bool) {
  res, ok := c.Ctx.GetQuery(key)
  if !ok {
    return gjson.Result{}, false
  }
  return gjson.Parse(res), true
}

// get value by string
func (c *New) Get(key string) (string, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return "", false
  }
  return res.(string), true
}

// get value by string
func (c *New) GetInt(key string) (int, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return 0, false
  }
  return res.(int), true
}

// get value by string
func (c *New) GetBool(key string) (bool, bool) {
  res, ok := c.Ctx.Get(key)
  if !ok {
    return false, false
  }
  return res.(bool), true
}
