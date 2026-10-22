# auth — 基于角色的 API 权限（RBAC）

> 仓库：https://github.com/open4go/auth ｜ 模块：`github.com/open4go/auth` ｜ Go 1.24.0 ｜ 最近提交：2026-09-10（最新 tag v0.2.53）｜ 仓库提交了 vendor（2633 个文件）

## 项目简介
角色→API 路径→权限位（`req5rsp/cst.Permission`：Create=1、Delete=2、Update=4、Read=8、List=16、Show=32、Disable=64）的权限模型。角色和应用存在 Mongo 库 `sys_auth`；用户登录时把账号→角色、角色→路径权限缓存到 Redis（handler 名 `auth`）；`AccessMiddleware` 按请求方法和路径校验权限。WebSocket 接口只要求 GET 权限（List 或 Read）。

## 目录结构
```
auth/
├── middle.go        # AccessMiddleware（Gin）
├── role.go          # RoleManager：SignIn/SignOut/Verify/canAccess/FetchAllPaths/Menu/Reload/KeepAlive/IsOnline
├── perm.go          # WS 识别、GET 权限判断、权限格式化
├── app.go           # MyApp：应用信息缓存（LoadAppInfo 函数体已整体注释）
├── menu.go          # MenuTree / MyAccessData
├── serialization.go # LoginInfo base64(JSON) 序列化
├── conf.go          # jwt.expire（默认 30m）
├── db.go            # Redis handler "auth"
├── model/{account,app,role,tenant}  # Mongo 模型（嵌入 model.Model，bson:"_"）
├── model/db.go      # 旧的全局 MDB 连接（log.Fatal/panic）
└── perm_test.go middle_test.go
```

## 主要 API
```go
// 路由
r.Use(auth.AccessMiddleware())           // 依赖上游写入的 UserID 请求头
// 登录时
rm := auth.RoleManager{RedisPrefix: "system:auth"}
_ = rm.SignIn(ctx, accountID, roleIDs)    // 写 account:to:role 集合，TTL 2h
ok, _ := rm.Verify(ctx, "/v1/member", accountID, "GET", false)
paths2Attr, _, _ := rm.FetchAllPaths(ctx, accountID)
menu := rm.Menu(ctx, paths2Attr)
rm.Reload(ctx)                            // 全量重载所有角色权限
```
HTTP 方法到权限的映射（`translateHTTPMethodToPermission`）：单资源（`/:_id`）DELETE→Delete、PUT/PATCH→Update、GET→Read；集合 POST→Create、GET→List；**其他组合返回 0**。

## 依赖
`gin v1.9.1`、`go-redis/v9 v9.0.5`、`mongo-driver v1.17.3`、`open4go/{db v0.0.10, log v0.0.16, model v0.0.24, r3time, req5rsp v0.1.30}`、`r2day/db`

## 构建与测试
`go build`、`go vet` 通过；`go test ./...` 通过（`perm_test.go`、`middle_test.go`）。

## 与其他仓库的关系
依赖 `req5rsp/cst` 的权限常量、`db` 的连接池、`model` 基类；需要 `middle` 先把 `UserID` 写进请求头。被 r2day/auth-api 等服务使用（open4go 内只有已停更的 `take-out` 还引用 v0.0.5）。
