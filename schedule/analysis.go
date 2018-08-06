package schedule

import (
  "time"
  "workerbook/db"
  "workerbook/service"
)

// 保存昨天的用户统计数据
func SaveAnalysis() error {
  db.ConnectMgoDB()

  m, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  }

  defer closer()

  // 获取日期
  day := time.Now().Add(-time.Hour * 24)

  // 根据日期获取相关的日报数据
  err = service.SaveAnalysisByDay(m, day)

  if err != nil {
    return err
  }

  return nil
}
