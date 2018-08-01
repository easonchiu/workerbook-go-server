package schedule

import (
  "fmt"
  "time"
  "workerbook/db"
  "workerbook/service"
)

func SaveUserAnalysis() error {
  db.ConnectMgoDB()

  m, closer, err := db.CloneMgoDB()

  if err != nil {
    fmt.Println(err)
    return err
  }

  defer closer()

  // 获取日期
  day := time.Now().Format("20060102")

  // 根据日期获取相关的日报数据
  err = service.SaveUserAnalysis(m, day)

  if err != nil {
    fmt.Println(err)
    return err
  }

  fmt.Println("ok")

  return nil
}
