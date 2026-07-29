# DCPM - 域控密码自助修改平台

Domain Controller Password Manager - 基于 Active Directory 的域账号密码自助管理 Web 平台。

## 功能特性

- **密码自助修改**：用户可在线修改域控密码，无需联系管理员
- **密码过期检测**：自动检测密码是否过期（通过 AD `data 773` 错误码）
- **密码永不过期识别**：自动识别 `userAccountControl` 中 `DONT_EXPIRE_PASSWORD` 标志
- **密码过期时间展示**：显示密码到期时间和剩余天数，临期预警
- **密码复杂度校验**：支持自定义密码策略（长度、大小写、数字、特殊字符）
- **JWT 会话管理**：Token 认证，安全可控
- **请求限流**：防止暴力破解（10次/5分钟）
- **审计日志**：记录密码修改操作

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + Gin + go-ldap + golang-jwt + Viper |
| 前端 | Vue 3 + TypeScript + Element Plus + Vite + Pinia |
| 协议 | LDAPS (636) / LDAP (389) |

## 项目结构

```
DCPM/
── server/                  # 后端服务
│   ├── config/              # 配置管理
│   │   ├── config.go        # 配置加载（支持环境变量覆盖）
│   │   └── config.example.yaml  # 配置模板
│   ├── handler/             # HTTP 处理器
│   ├── middleware/           # 中间件（JWT、限流、CORS、日志）
│   ├── model/               # 数据模型
│   ├── service/             # 业务逻辑（LDAP 操作、密码策略）
│   ├── utils/               # 工具函数（JWT、密码校验）
│   ├── certificate/         # LDAPS 证书目录
│   └── main.go              # 入口
├── web/                     # 前端应用
│   ├── src/
│   │   ├── api/             # API 请求封装
│   │   ├── router/          # 路由配置
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── views/           # 页面组件
│   │   ── types/           # TypeScript 类型定义
│   └── vite.config.ts       # Vite 配置
└── .gitignore
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- Active Directory 域控服务器

### 1. 克隆项目

```bash
git clone https://github.com/<your-username>/DCPM.git
cd DCPM
```

### 2. 配置后端

```bash
cd server/config
cp config.example.yaml config.yaml
```

编辑 `config.yaml`，填入你的域控信息：

```yaml
ldap:
  server: "ldaps://your-dc-server:636"
  base_dn: "dc=example,dc=com"
  domain_suffix: "example.com"
  username: "service-account@example.com"
  password: "your-service-account-password"
```

> **安全提示**：生产环境建议使用环境变量覆盖敏感配置：
> ```bash
> export LDAP_PASSWORD="your-actual-password"
> export JWT_SECRET="your-strong-random-secret"
> ```

### 3. 启动后端

```bash
cd server
go mod tidy
go run main.go
```

后端默认运行在 `http://localhost:8088`

### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

前端默认运行在 `http://localhost:8089`

## 配置说明

### 环境变量（优先级高于配置文件）

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `LDAP_SERVER` | LDAP 服务器地址 | `ldaps://dc.example.com:636` |
| `LDAP_USERNAME` | 服务账号 | `admin@example.com` |
| `LDAP_PASSWORD` | 服务账号密码 | `P@ssw0rd` |
| `JWT_SECRET` | JWT 签名密钥 | `random-secret-key` |
| `SERVER_PORT` | 后端服务端口 | `8088` |

### 密码策略配置

```yaml
password_policy:
  min_length: 14              # 最小长度
  require_uppercase: true     # 需要大写字母
  require_lowercase: true     # 需要小写字母
  require_digit: true         # 需要数字
  require_special: true       # 需要特殊字符
```

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/health` | 健康检查 | 否 |
| POST | `/api/v1/auth/login` | 用户登录 | 否（限流） |
| GET | `/api/v1/user/info` | 获取用户信息 | JWT |
| POST | `/api/v1/password/change` | 修改密码 | 否 |

## 部署建议

1. **生产环境**：将 `server.mode` 设为 `release`
2. **HTTPS**：前端通过 Nginx 反向代理，配置 SSL 证书
3. **证书验证**：将 `ldap.insecure` 设为 `false`，并配置 `cert_path`
4. **JWT 密钥**：使用强随机字符串，长度不少于 32 字符
5. **服务账号**：使用最小权限原则，仅授予密码修改权限

## License

MIT
