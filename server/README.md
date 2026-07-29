# DCPM Server - 域控密码自助修改平台后端

基于 Go + Gin 框架开发的 RESTful API 服务，通过 LDAPS 协议连接 Active Directory，为域用户提供安全的密码自助修改能力。

## 目录

- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [架构设计](#架构设计)
- [快速开始](#快速开始)
- [配置详解](#配置详解)
- [API 接口文档](#api-接口文档)
- [核心机制](#核心机制)
- [安全设计](#安全设计)
- [构建与部署](#构建与部署)

---

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.25+ | 高性能、编译型语言 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 | 轻量级 HTTP 框架 |
| LDAP 客户端 | [go-ldap](https://github.com/go-ldap/ldap) v3.4 | 支持 LDAPS/TLS |
| JWT 认证 | [golang-jwt](https://github.com/golang-jwt/jwt) v5 | Token 签发与验证 |
| 配置管理 | [Viper](https://github.com/spf13/viper) v1.21 | YAML 配置 + 环境变量覆盖 |

## 项目结构

```
server/
├── certificate/              # CA 证书目录
│   └── ca.crt                # LDAPS 连接用的 CA 根证书
├── config/
│   ├── config.example.yaml   # 配置模板（提交到版本控制）
│   ├── config.go             # 配置结构定义与加载逻辑
│   └── config.yaml           # 运行配置（不提交到版本控制）
├── handler/                  # HTTP 请求处理器层
│   ├── auth_handler.go       # 认证：登录、获取用户信息
│   └── password_handler.go   # 密码修改请求处理
├── middleware/               # Gin 中间件
│   ├── cors.go               # CORS 跨域处理
│   ├── jwt.go                # JWT Bearer Token 鉴权
│   ├── logger.go             # 审计日志记录
│   └── ratelimit.go          # 基于 IP 的请求限流
├── model/                    # 数据模型定义
│   └── request.go            # 请求/响应结构体
├── service/                  # 业务逻辑层
│   ├── ldap_service.go       # LDAP 连接管理、认证、密码操作
│   └── password_service.go   # 密码修改业务编排（校验→验证→执行）
├── utils/                    # 工具函数
│   ├── jwt.go                # JWT Token 生成与解析
│   └── password.go           # 密码复杂度校验与强度评估
├── go.mod                    # Go 模块依赖
├── go.sum                    # 依赖校验文件
└── main.go                   # 程序入口：配置加载、路由注册、服务启动
```

## 架构设计

系统采用经典的分层架构，职责清晰：

```
请求 → [CORS] → [限流] → [JWT鉴权] → [审计日志] → Handler → Service → LDAP
```

**分层说明：**

| 层级 | 职责 | 文件 |
|------|------|------|
| **Handler 层** | 接收 HTTP 请求、参数绑定与校验、调用 Service、返回响应 | `handler/*.go` |
| **Service 层** | 业务逻辑编排：密码复杂度校验 → 旧密码验证 → 执行修改 | `service/*.go` |
| **Model 层** | 定义请求/响应数据结构体 | `model/*.go` |
| **Middleware 层** | 横切关注点：跨域、限流、认证、日志 | `middleware/*.go` |
| **Utils 层** | 无状态工具函数：JWT 操作、密码规则校验 | `utils/*.go` |

**LDAP 连接流程：**

```
1. 读取配置（服务器地址、TLS 设置、CA 证书）
2. 建立 LDAPS 连接（端口 636）或 LDAP + StartTLS（端口 389）
3. 根据操作类型选择绑定方式：
   - 用户认证：使用 user@domain 格式 Bind
   - 信息查询/密码修改：使用服务账号 Bind
4. 执行 LDAP 操作后关闭连接
```

## 快速开始

### 环境要求

- Go 1.25+
- Active Directory 域控服务器（支持 LDAPS 636 或 LDAP 389）
- AD 服务账号（具有查询用户信息和修改密码的权限）

### 安装与启动

```bash
# 1. 进入后端目录
cd server

# 2. 下载依赖
go mod download

# 3. 复制配置文件
cp config/config.example.yaml config/config.yaml

# 4. 编辑 config/config.yaml，填入实际的 AD 连接信息
#    （详见下方"配置详解"章节）

# 5. 启动服务
go run main.go
```

服务默认运行在 `http://localhost:8088`。

## 配置详解

配置文件位于 `config/config.yaml`，基于 YAML 格式。所有配置项说明如下：

### 服务器配置

```yaml
server:
  port: 8088          # 监听端口
  mode: debug         # 运行模式：debug（开发）/ release（生产）
```

### JWT 配置

```yaml
jwt:
  secret: "your-jwt-secret-key"  # 签名密钥，生产环境必须使用强随机字符串
  expiry: 24h                    # Token 有效期，格式为 Go Duration（如 24h、30m）
```

### LDAP 配置

```yaml
ldap:
  server: "ldaps://your-dc-server:636"  # AD 服务器地址（ldaps:// 或 ldap://）
  base_dn: "dc=example,dc=com"          # 搜索基础 DN
  domain_suffix: "example.com"          # 域名后缀（用于 user@domain 格式登录）
  use_tls: true                         # 是否启用 TLS
  insecure: true                        # true=跳过证书验证（开发用），生产环境必须为 false
  user_filter: "(sAMAccountName=%s)"    # 用户搜索过滤器，%s 为用户名占位符
  username: "service@example.com"       # 服务账号（用于查询和密码修改操作）
  password: "service-password"          # 服务账号密码
  cert_path: "./certificate/ca.crt"     # CA 证书路径（insecure=false 时必需）
  max_pwd_age_days: 30                  # 密码最大使用天数（当无法从域控获取时的兜底值）
```

### 密码策略配置

```yaml
password_policy:
  min_length: 14          # 密码最小长度
  require_uppercase: true  # 必须包含大写字母
  require_lowercase: true  # 必须包含小写字母
  require_digit: true      # 必须包含数字
  require_special: true    # 必须包含特殊字符
```

> 注意：密码策略需与 AD 域控的组策略保持一致，否则用户在 AD 端修改密码时可能遇到不同的规则。

### 环境变量覆盖

以下敏感配置支持通过环境变量注入，**优先级高于配置文件**：

| 环境变量 | 对应配置项 | 说明 |
|---------|-----------|------|
| `LDAP_SERVER` | `ldap.server` | AD 服务器地址 |
| `LDAP_USERNAME` | `ldap.username` | 服务账号 |
| `LDAP_PASSWORD` | `ldap.password` | 服务账号密码 |
| `JWT_SECRET` | `jwt.secret` | JWT 签名密钥 |
| `SERVER_PORT` | `server.port` | 服务监听端口 |

生产环境推荐使用环境变量注入敏感信息，避免将密码写入配置文件。

## API 接口文档

所有接口统一前缀 `/api/v1`，响应格式：

```json
{
  "code": 200,
  "message": "描述信息",
  "data": { ... }
}
```

### 健康检查

```
GET /api/v1/health
```

**响应示例：**

```json
{ "code": 200, "message": "服务运行正常" }
```

---

### 用户登录

```
POST /api/v1/auth/login
Content-Type: application/json
```

**请求体：**

```json
{
  "username": "zhangsan",
  "password": "OldPassword123!"
}
```

**响应 — 登录成功（密码正常）：**

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "status": "ok",
    "message": "登录成功",
    "token": "eyJhbGciOiJIUzI1...",
    "username": "zhangsan",
    "password_expires_at": "2026-08-15 10:30:00",
    "days_remaining": 16,
    "password_never_expires": false
  }
}
```

**响应 — 密码已过期：**

```json
{
  "code": 200,
  "message": "密码已过期，请修改密码",
  "data": {
    "status": "expired",
    "message": "密码已过期，请修改密码"
  }
}
```

**响应 — 认证失败：**

```json
{
  "code": 401,
  "message": "账号或密码错误",
  "data": { "status": "error", "message": "账号或密码错误" }
}
```

**特殊场景：**

| 场景 | 处理方式 |
|------|---------|
| 密码过期 | AD 返回错误码 `data 773`，系统识别为过期状态，前端引导用户修改密码 |
| 账号锁定 | AD 返回错误码 `data 775`，提示用户联系管理员解锁 |
| 限流触发 | 同一 IP 在 5 分钟内超过 10 次请求，返回 429 状态码 |

---

### 获取用户信息

```
GET /api/v1/user/info
Authorization: Bearer <token>
```

**响应示例：**

```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "username": "zhangsan",
    "password_expires_at": "2026-08-15 10:30:00",
    "days_remaining": 16,
    "password_never_expires": false
  }
}
```

---

### 修改密码

```
POST /api/v1/password/change
Content-Type: application/json
```

**请求体：**

```json
{
  "username": "zhangsan",
  "old_password": "OldPassword123!",
  "new_password": "NewStrongPwd@456"
}
```

**响应 — 成功：**

```json
{ "code": 200, "message": "密码修改成功" }
```

**响应 — 失败（密码复杂度不满足）：**

```json
{ "code": 400, "message": "密码必须包含大写字母" }
```

**修改密码的业务流程：**

```
1. 参数校验（必填字段检查）
2. 新旧密码对比检查（不能相同）
3. 新密码复杂度校验（长度、字符类型）
4. 旧密码验证（通过 LDAP Bind 验证，密码过期时也视为通过）
5. 使用服务账号连接 AD，查找用户 DN
6. 将新密码编码为 UTF-16LE 格式（双引号包裹）
7. 通过 LDAP Modify 操作更新 unicodePwd 属性
8. 记录审计日志
```

## 核心机制

### 密码过期检测

系统通过以下方式检测密码过期状态：

1. **登录时检测**：用户 Bind 时，AD 返回 `data 773` 错误码表示密码已过期
2. **属性查询**：使用服务账号查询 `pwdLastSet`（Windows FILETIME 格式）和 `maxPwdAge`（域策略）
3. **时间计算**：将 FILETIME（100ns 间隔，起始于 1601-01-01）转换为标准时间，计算到期日与剩余天数
4. **永不过期识别**：检查 `userAccountControl` 属性中的 `DONT_EXPIRE_PASSWORD` 标志位（第 17 位，值 0x10000 = 65536）

### 密码复杂度校验

密码修改前，系统会在服务端进行双重校验：

| 校验项 | 规则 |
|--------|------|
| 最小长度 | 可配置，默认 14 位 |
| 用户名包含 | 密码中不能包含用户名 |
| 字符类别 | 大写、小写、数字、特殊字符至少满足 3 类 |
| 必须字符 | 根据 `password_policy` 配置逐项检查 |
| 新旧对比 | 新密码不能与旧密码相同 |

### 密码编码

AD 要求 `unicodePwd` 属性的值必须为 **UTF-16LE 编码并用双引号包裹**，系统自动完成此编码转换：

```
原始密码: MyP@ssw0rd
编码过程: "MyP@ssw0rd" → UTF-16LE 字节序列
```

## 安全设计

| 安全措施 | 实现方式 |
|---------|---------|
| **JWT 会话管理** | 登录成功签发 Token，后续请求通过 `Authorization: Bearer <token>` 鉴权 |
| **请求限流** | 基于 IP 的滑动窗口限流，登录接口限制为 10 次/5 分钟，自动清理过期记录 |
| **密码传输** | 全链路 LDAPS 加密传输，密码不以明文存储 |
| **服务账号隔离** | 用户认证使用用户自身凭据 Bind，查询/修改操作使用服务账号 Bind |
| **审计日志** | 所有密码修改操作记录用户名、IP、时间、结果 |
| **CORS 控制** | 中间件统一处理跨域请求 |
| **环境变量注入** | 敏感配置支持环境变量覆盖，避免密码写入文件 |

## 构建与部署

### 编译

```bash
# 编译为当前平台可执行文件
go build -o dcpm-server main.go

# 交叉编译为 Linux amd64
GOOS=linux GOARCH=amd64 go build -o dcpm-server main.go
```

### 生产部署清单

1. 将 `config/config.yaml` 中的 `server.mode` 设为 `release`
2. 将 `ldap.insecure` 设为 `false`，并确保 `cert_path` 指向有效的 CA 证书
3. 使用强随机字符串作为 `jwt.secret`
4. 通过环境变量注入 `LDAP_USERNAME`、`LDAP_PASSWORD`、`JWT_SECRET`
5. 将 `password_policy` 配置与域控组策略对齐
6. 使用 systemd / Docker / 进程管理工具托管进程

### 反向代理示例（Nginx）

```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8088;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```
