# auth 优化方案

| # | 优先级 | 问题 | 位置 | 建议 | 工作量 |
|---|---|---|---|---|---|
| 1 | 高 | **未映射的方法得到权限 0，等于放行**：`translateHTTPMethodToPermission` 对集合上的 PUT/PATCH/DELETE、单资源上的 POST、HEAD/OPTIONS 都返回 0；`canAccess` 判断 `roleAttr & 0 == 0` 恒为真，只要角色在该路径上有任意一条记录（哪怕只有 Show），就能执行这些操作（例如批量删除 `DELETE /v1/xxx`） | `role.go:482-508`、`canAccess` | 返回 0 时直接拒绝（`if expect == 0 { return false }`），或者显式映射：集合 DELETE→Delete、集合 PUT→Update；补表驱动测试 | S |
| 2 | 高 | **按路径判断 WS 会让写请求降级成 GET 校验**：`Verify` 中 `isWSEndpoint(path)` 只看路径里是否有 `ws`/`websocket` 段，不看是否真的是 WebSocket 升级请求。对这类路径发 POST/DELETE，只需要 List/Read 权限 | `role.go:171`、`middle.go` | 只用 `c.IsWebsocket()`（Upgrade 头）判断；路径规则仅作为额外条件 | S |
| 3 | 高 | `AccessMiddleware` 直接信任 `UserID` 请求头，客户端可以伪造身份 | `middle.go:43` | 只信任 `middle.JWTMiddleware` 解析后写入的值（例如改为从 ctx 读 `model.AccountKey`）；入口处先删除外部传入的身份头 | M |
| 4 | 中 | 权限不足返回 401，应该返回 403；校验出错（Redis 故障）也返回 401 | `middle.go` | 未登录返回 401，无权限返回 403，内部错误返回 503 | S |
| 5 | 中 | `account.Model` 的 `Password` 带 `json:"password"`，如果接口直接返回这个模型，会泄露密码哈希 | `model/account/model.go` | 改成 `json:"-"`，或者另外定义 DTO | S |
| 6 | 中 | `model/db.go` 是旧的全局 `MDB`（`log.Fatal`/`panic`），和 `db.DBPool` 重复；`app.go` 的 `InitApp` 出错会 panic，`LoadAppInfo` 函数体整段被注释 | `model/db.go`、`app.go` | 删掉旧连接和死代码；`InitApp` 改为返回 error | S |
| 7 | 中 | `SignOut` 只把 key 的过期时间改成 10 秒，期间仍然可以访问 | `role.go:159` | 直接 `Del`；同时把 JWT 的 jti 加入黑名单 | S |
| 8 | 中 | vendor 提交进了仓库（2633 个文件） | `vendor/` | 删除并加入 `.gitignore` | S |
| 9 | 中 | govulncheck 发现 8 个问题：x/net v0.29.0、x/text v0.18.0、mongo-driver v1.17.3（GO-2026-5327，修复版本 v1.17.7）、docker（来自旧版 log） | `go.mod` | 升级 log 到 v0.0.23+、mongo-driver 到 v1.17.7+、x/net、x/text | S |
| 10 | 低 | 依赖的 req5rsp v0.1.30 和 db v0.0.10 比最新版落后很多，只为了拿 `cst` 常量就引入了整个 req5rsp（以及 wechat SDK 等） | `go.mod` | 把 `cst` 拆成独立的无依赖模块 | M |
