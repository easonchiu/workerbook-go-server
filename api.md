



get  /users?skip=n&limit=m         获取用户列表
get  /users/:id                    获取单个用户
put  /users/:id                    修改用户信息
post /users                        添加用户

get  /projects?skip=n&limit=m    获取项目列表
get  /projects/:id               获取单个项目
post /projects                   添加项目
put  /projects/:id               修改项目

get  /departments?skip=n&limit=m     获取部门列表
get  /departments                    获取全部部门列表
get  /departments/:id                获取部门列表
put  /departments/:id                修改部门信息
post /departments                    添加部门

get  /profile               获取我的信息
post /login                 登录

