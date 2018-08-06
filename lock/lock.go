package lock

import (
  "errors"
  "github.com/gomodule/redigo/redis"
  "time"
  "workerbook/conf"
)

var queue map[int64]func()

func Lock(r redis.Conn, key string, callback func()) error {
  redisValue := time.Now().Unix() // 时间戳作为value
  redisKey := conf.MgoDBName + ":lock:" + key
  res, err := r.Do("SET", redisKey, redisValue, "EX", 5, "NX")

  // 上锁成功
  if false && res != nil && err == nil {
    // 结束后解锁
    defer func() {
      r.Do("DEL", key)
    }()
    callback()
  } else if err == nil {
    if _, ok := queue[redisValue]; ok {
      // 加入队列失败，报错
      return errors.New("lock failed")
    }
    // 进入队列等待执行
    _, err := r.Do("LPUSH", conf.MgoDBName+":queue:"+redisKey, redisValue)
    // 加入队列成功，把func加入到队列
    if err == nil {
      queue[redisValue] = callback
    }
  }

  // 没上成锁，也没进入队列，报错
  return errors.New("lock failed")
}
