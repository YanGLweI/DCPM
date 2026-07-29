# DCPM Web - 域控密码自助修改平台前端

基于 Vue 3 + TypeScript + Element Plus 构建的单页应用（SPA），为域用户提供密码自助修改、密码到期查询、账户信息查看等交互功能。

## 目录

- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [架构设计](#架构设计)
- [快速开始](#快速开始)
- [页面功能详解](#页面功能详解)
- [状态管理](#状态管理)
- [API 对接](#api-对接)
- [路由与导航](#路由与导航)
- [构建与部署](#构建与部署)

---

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | [Vue 3](https://vuejs.org/) | Composition API + `<script setup>` 语法 |
| 语言 | [TypeScript](https://www.typescriptlang.org/) v6.0 | 类型安全的 JavaScript 超集 |
| UI 组件库 | [Element Plus](https://element-plus.org/) v2.14 | Vue 3 企业级 UI 组件库 |
| 构建工具 | [Vite](https://vite.dev/) v8.1 | 极速开发服务器 + 构建 |
| 状态管理 | [Pinia](https://pinia.vuejs.org/) v4 | Vue 3 官方推荐状态管理 |
| 路由 | [Vue Router](https://router.vuejs.org/) v4 | 路由守卫 + 懒加载 |
| HTTP 客户端 | [Axios](https://axios-http.com/) v1.18 | 请求/响应拦截器 |
| 图标 | [@element-plus/icons-vue](https://element-plus.org/zh-CN/component/icon.html) | Element Plus 图标库 |

## 项目结构

```
web/
├── public/                   # 静态资源（不经过构建处理）
│   ├── favicon.svg           # 站点图标
│   └── icons.svg             # 图标集
├── src/
│   ├── api/                  # API 请求封装层
│   │   └── index.ts          # Axios 实例、拦截器、各接口请求函数
│   ├── assets/               # 静态资源（经过构建处理）
│   │   ├── hero.png          # 页面背景图
│   │   ├── vite.svg          # Vite logo
│   │   └── vue.svg           # Vue logo
│   ├── components/           # 公共组件
│   │   └── HelloWorld.vue    # 示例组件
│   ├── router/               # 路由配置
│   │   └── index.ts          # 路由定义、路由守卫
│   ├── stores/               # Pinia 状态管理
│   │   └── user.ts           # 用户状态：Token、登录/登出、密码到期信息
│   ├── types/                # TypeScript 类型定义
│   │   └── index.ts          # API 响应类型、请求类型、用户信息类型
│   ├── views/                # 页面视图组件
│   │   ├── LoginView.vue     # 登录页
│   │   ├── ChangePassword.vue # 密码修改页（支持过期强制修改）
│   │   └── AccountManage.vue # 账户管理页（信息展示 + 密码修改）
│   ├── App.vue               # 根组件（路由出口）
│   ├── main.ts               # 应用入口：挂载插件、创建实例
│   └── style.css             # 全局样式
├── index.html                # HTML 入口模板
├── package.json              # 依赖与脚本定义
├── tsconfig.json             # TypeScript 配置（项目级）
├── tsconfig.app.json         # TypeScript 配置（应用级）
├── tsconfig.node.json        # TypeScript 配置（Node 级，Vite 配置用）
└── vite.config.ts            # Vite 配置：端口、代理
```

## 架构设计

### 分层结构

```
用户操作 → View 视图组件 → Store 状态管理 → API 请求层 → 后端服务
                                    ↓
                              localStorage（持久化 Token/用户名）
```

**各层职责：**

| 层级 | 职责 | 文件 |
|------|------|------|
| **Views 层** | 页面 UI 渲染、表单校验、用户交互反馈 | `views/*.vue` |
| **Stores 层** | 管理用户状态（Token、密码到期信息）、封装登录/登出逻辑 | `stores/user.ts` |
| **API 层** | 封装 Axios 实例、请求/响应拦截器、各 API 调用函数 | `api/index.ts` |
| **Types 层** | 定义 TypeScript 接口，确保前后端数据结构一致 | `types/index.ts` |
| **Router 层** | 路由定义、路由守卫（登录态检查、过期跳转） | `router/index.ts` |

### 请求处理流程

```
1. View 调用 Store 方法（如 userStore.login）
2. Store 调用 API 函数（如 loginApi）
3. Axios 请求拦截器自动添加 Authorization Header
4. 发送请求到 /api/v1/* → Vite 代理转发到后端 localhost:8088
5. Axios 响应拦截器处理 401（自动清除 Token 并跳转登录页）
6. Store 更新状态并持久化到 localStorage
7. View 根据状态变化更新 UI
```

## 快速开始

### 环境要求

- Node.js 18+
- npm 9+
- 后端服务已运行（默认 `http://localhost:8088`）

### 安装与启动

```bash
# 1. 进入前端目录
cd web

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev
```

开发服务器运行在 `http://localhost:8089`，已配置代理自动将 `/api` 请求转发到后端。

### 可用脚本

| 命令 | 说明 |
|------|------|
| `npm run dev` | 启动开发服务器（端口 8089，支持热更新） |
| `npm run build` | TypeScript 类型检查 + 生产构建，输出到 `dist/` |
| `npm run preview` | 预览生产构建结果 |

## 页面功能详解

### 登录页 (`/login`)

**文件**: `src/views/LoginView.vue`

**功能：**
- 域账号 + 密码表单输入，带必填校验
- 登录成功后跳转到账户管理页
- 密码过期时自动跳转到密码修改页（携带 `?expired=true` 参数）
- 密码过期场景下，用户名通过 `sessionStorage` 临时传递到修改密码页
- 已登录用户访问登录页会自动重定向到账户管理页

**UI 特性：**
- 渐变背景 + 卡片式登录框
- 输入框带图标前缀（用户/锁）
- 密码输入支持显示/隐藏切换
- 登录按钮加载状态

### 密码修改页 (`/change-password`)

**文件**: `src/views/ChangePassword.vue`

**功能：**
- 支持两种场景：主动修改密码 和 密码过期强制修改
- 过期场景显示醒目的红色提示标签和警告横幅
- 从 `sessionStorage` 自动读取临时用户名（过期场景）
- 实时密码强度检测（弱/中/强/很强），带进度条可视化
- 表单校验：旧密码必填、新密码最少 14 位、确认密码一致性检查
- 修改成功后清除临时数据，跳转登录页

**密码要求提示：**
> 至少 14 位，包含大写字母、小写字母、数字、特殊字符中的至少 3 类

### 账户管理页 (`/account`)

**文件**: `src/views/AccountManage.vue`

**功能：**
- 展示用户账号信息（用户名、密码状态、过期时间、剩余天数）
- 进入页面自动拉取最新的密码到期信息
- 密码状态智能标签：
  - **正常**（绿色）— 密码未过期且剩余天数 > 7 天
  - **即将过期**（橙色）— 剩余天数 ≤ 7 天
  - **已过期**（红色）— 剩余天数 = 0
  - **永不过期**（灰色）— 账号设置了密码永不过期
- 状态告警横幅：根据密码状态显示不同级别的通知
- 内嵌密码修改表单（含强度检测、确认密码校验）
- 修改成功后弹窗确认是否退出重新登录
- 退出登录功能（二次确认）

**UI 特性：**
- 双卡片布局：上方信息展示 + 下方密码修改
- 描述列表（Descriptions）展示账号详情
- 剩余天数临近过期时高亮显示

## 状态管理

**文件**: `src/stores/user.ts`

使用 Pinia 管理用户全局状态：

| 状态 | 类型 | 说明 |
|------|------|------|
| `token` | `string` | JWT Token，持久化到 localStorage |
| `username` | `string` | 当前用户名，持久化到 localStorage |
| `passwordExpiresAt` | `string` | 密码到期时间 |
| `daysRemaining` | `number` | 密码剩余天数 |
| `passwordNeverExpires` | `boolean` | 密码是否永不过期 |

**计算属性：**

| 属性 | 说明 |
|------|------|
| `isLoggedIn` | 是否已登录（Token 非空） |
| `isExpired` | 密码是否已过期（非永不过期 且 剩余天数 ≤ 0） |
| `isExpiringSoon` | 密码是否即将过期（非永不过期 且 剩余天数 ≤ 7） |

**核心方法：**

| 方法 | 说明 |
|------|------|
| `login(loginData)` | 执行登录，处理成功/过期/失败三种情况 |
| `fetchUserInfo()` | 从后端拉取最新用户信息并更新状态 |
| `logout()` | 清除状态和 localStorage，跳转到登录页 |

## API 对接

**文件**: `src/api/index.ts`

### Axios 配置

```typescript
baseURL: '/api/v1'    // 所有请求统一前缀
timeout: 30000        // 超时时间 30 秒
```

### 请求拦截器

自动从 `localStorage` 读取 Token，添加到请求头：

```
Authorization: Bearer <token>
```

### 响应拦截器

捕获 401 响应，自动清除本地 Token 和用户名并跳转到登录页。

### API 函数

| 函数 | 方法 | 路径 | 说明 |
|------|------|------|------|
| `login(data)` | POST | `/auth/login` | 用户登录 |
| `changePassword(data)` | POST | `/password/change` | 修改密码 |
| `getUserInfo()` | GET | `/user/info` | 获取用户信息（需 JWT） |
| `healthCheck()` | GET | `/health` | 健康检查 |

### TypeScript 类型定义

**文件**: `src/types/index.ts`

```typescript
// 统一 API 响应结构
interface ApiResponse<T> { code: number; message: string; data?: T }

// 登录请求
interface LoginRequest { username: string; password: string }

// 登录响应（含状态、Token、密码到期信息）
interface LoginResponseData {
  status: 'ok' | 'expired' | 'error'
  message: string
  token?: string
  username?: string
  password_expires_at?: string
  days_remaining?: number
  password_never_expires?: boolean
}

// 修改密码请求
interface ChangePasswordRequest { username: string; old_password: string; new_password: string }

// 用户信息
interface UserInfo { username: string; password_expires_at: string; days_remaining: number; password_never_expires: boolean }
```

## 路由与导航

**文件**: `src/router/index.ts`

### 路由表

| 路径 | 组件 | 需要认证 | 说明 |
|------|------|---------|------|
| `/` | — | — | 重定向到 `/login` |
| `/login` | `LoginView.vue` | 否 | 登录页（懒加载） |
| `/change-password` | `ChangePassword.vue` | 否* | 密码修改页（懒加载） |
| `/account` | `AccountManage.vue` | 是 | 账户管理页（懒加载） |

> *密码修改页不要求 JWT 认证，因为密码过期场景下用户尚未获取 Token。

### 路由守卫逻辑

```
访问需要认证的路由（/account）且无 Token → 重定向到 /login
已登录用户访问 /login → 重定向到 /account
其他情况 → 正常放行
```

## 构建与部署

### 开发代理配置

**文件**: `vite.config.ts`

开发环境下，Vite 自动将 `/api` 请求代理到后端：

```typescript
server: {
  port: 8089,
  proxy: {
    '/api': {
      target: 'http://localhost:8088',
      changeOrigin: true
    }
  }
}
```

### 生产构建

```bash
npm run build
```

构建产物输出到 `dist/` 目录，包含压缩后的 HTML、CSS、JS 文件。

### 生产部署（Nginx 示例）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态资源
    location / {
        root /path/to/web/dist;
        try_files $uri $uri/ /index.html;  # SPA History 模式
    }

    # API 反向代理到后端
    location /api/ {
        proxy_pass http://127.0.0.1:8088;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

> 注意：由于使用 `createWebHistory()` 模式，Nginx 必须配置 `try_files` 回退到 `index.html`，否则刷新页面会 404。

### 开发注意事项

- 项目使用 Vue 3 `<script setup>` 语法，需要熟悉 Composition API
- 所有组件使用 TypeScript，新增功能需保持类型定义完整
- 密码强度检测在前端仅作为用户体验提示，实际校验以后端为准
- `localStorage` 存储 Token 和用户名，`sessionStorage` 临时存储密码过期场景的用户名
