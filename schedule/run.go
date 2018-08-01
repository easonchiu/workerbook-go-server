package schedule

import (
  "time"
)

func Start() {
  run(1, 0, time.Hour*24, func() {
    SaveUserAnalysis()
  })
}

func run(hour int, min int, duration time.Duration, cb func()) {
  tick := time.Tick(time.Second) // 每秒检查一次
  now := time.Now()
  next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)

  for {
    <-tick
    t := time.Now()
    if t.After(next) {
      next = next.Add(duration)
      cb()
    }
  }
}
