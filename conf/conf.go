package conf

const (
  DBName = "workerbook"
  DBUrl  = "mongodb://localhost:27017/workerbook"
)

var JwtSecret = []byte("wb2018_qweasdzxc!@#")

// 1: 开发者 2: 部门管理者 3: 观察员 4: 项目管理者 99: 管理员
const (
  RoleAdmin  = 99 // 管理员
  RolePM     = 4  // 项目管理者
  RoleOB     = 3  // 观察员
  RoleLeader = 2  // 部门管理者
  RoleDev    = 1  // 开发者
)