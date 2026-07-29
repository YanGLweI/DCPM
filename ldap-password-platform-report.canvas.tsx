import { Divider, Grid, H1, H2, H3, Stack, Stat, Table, Text } from 'qoder/canvas';

export default function LDAPPasswordPlatformReport() {
  return (
    <Stack gap={20}>
      <H1>LDAP 域控密码自助修改平台 - 项目完成报告</H1>
      <Text tone="secondary">
        基于 Go+Gin 后端和 Vue3+Element Plus 前端的域控密码自助修改平台，支持密码过期检测、JWT 会话管理、密码过期时间查询与展示、以及密码修改功能。
      </Text>

      <Divider />

      <H2>项目概览</H2>
      <Grid columns={4} gap={16}>
        <Stat value="20" label="开发任务" />
        <Stat value="15" label="后端文件" tone="info" />
        <Stat value="12" label="前端文件" tone="info" />
        <Stat value="4" label="API 接口" tone="success" />
      </Grid>

      <Divider />

      <H2>后端实现 (Go + Gin)</H2>
      <Grid columns={3} gap={12}>
        <Stat value="8088" label="服务端口" />
        <Stat value="14" label="密码最小长度" />
        <Stat value="10/5min" label="登录限流" tone="warning" />
      </Grid>

      <H3>核心功能</H3>
      <Table
        headers={['模块', '文件', '功能描述']}
        rows={[
          ['配置管理', 'config/config.go + config.yaml', 'Viper 配置加载，LDAP/JWT/密码策略全部可配置'],
          ['LDAP 服务', 'service/ldap_service.go', 'LDAPS 连接、认证、密码过期检测(data 773)、密码修改(UTF-16LE)'],
          ['密码服务', 'service/password_service.go', '密码复杂度校验、新旧密码比对、业务逻辑编排'],
          ['认证接口', 'handler/auth_handler.go', '登录验证、JWT 签发、密码过期时间查询'],
          ['密码接口', 'handler/password_handler.go', '密码修改、审计日志记录'],
          ['JWT 工具', 'utils/jwt.go', 'Token 生成与解析 (HS256)'],
          ['密码校验', 'utils/password.go', '复杂度校验(>=14位, 3/4类字符)'],
          ['限流中间件', 'middleware/ratelimit.go', 'IP 限流 10次/5分钟，自动清理过期记录'],
          ['JWT 中间件', 'middleware/jwt.go', 'Bearer Token 验证，Context 注入用户名'],
          ['CORS 中间件', 'middleware/cors.go', '跨域请求支持'],
          ['日志中间件', 'middleware/logger.go', '请求日志(状态码/耗时/IP/路径)'],
        ]}
      />

      <H3>API 接口</H3>
      <Table
        headers={['方法', '路径', '认证', '功能']}
        rows={[
          ['POST', '/api/v1/auth/login', '无 (限流)', '登录验证，判断密码是否过期，签发 JWT'],
          ['POST', '/api/v1/password/change', '无', '修改密码 (过期/未过期均可调用)'],
          ['GET', '/api/v1/user/info', 'JWT', '获取用户信息 + 密码过期时间'],
          ['GET', '/api/v1/health', '无', '健康检查'],
        ]}
        rowTone={[undefined, undefined, 'info', undefined]}
      />

      <Divider />

      <H2>前端实现 (Vue3 + Element Plus + TypeScript)</H2>
      <Grid columns={3} gap={12}>
        <Stat value="8089" label="开发端口" />
        <Stat value="3" label="页面数量" />
        <Stat value="Proxy" label="API 代理" tone="info" />
      </Grid>

      <H3>页面组件</H3>
      <Table
        headers={['页面', '文件', '功能']}
        rows={[
          ['登录页', 'views/LoginView.vue', '居中卡片布局，域账号/密码输入，自动判断密码过期并跳转'],
          ['密码修改页', 'views/ChangePassword.vue', '过期警告提示，密码强度指示器，实时复杂度校验，修改成功跳转登录'],
          ['账号管理页', 'views/AccountManage.vue', '密码过期时间展示，剩余天数警告(<=7天橙色/过期红色)，修改密码，退出登录'],
        ]}
      />

      <H3>前端架构</H3>
      <Table
        headers={['模块', '文件', '功能']}
        rows={[
          ['API 封装', 'api/index.ts', 'Axios 实例 + JWT 拦截器 + 401 自动跳转'],
          ['状态管理', 'stores/user.ts', 'Pinia Store: 登录/登出/用户信息/密码过期状态'],
          ['路由配置', 'router/index.ts', 'Vue Router + 路由守卫(JWT 检查)'],
          ['类型定义', 'types/index.ts', 'TypeScript 接口定义'],
        ]}
      />

      <Divider />

      <H2>安全特性</H2>
      <Grid columns={3} gap={12}>
        <Stat value="LDAPS" label="加密通信" tone="success" />
        <Stat value="JWT" label="会话管理" tone="success" />
        <Stat value="审计日志" label="操作追踪" tone="success" />
      </Grid>
      <Table
        headers={['安全措施', '实现方式']}
        rows={[
          ['传输安全', 'LDAPS (636端口) + CA 证书验证'],
          ['会话管理', 'JWT Token (HS256, 24h 有效期)'],
          ['请求限流', '登录接口 IP 限流 10次/5分钟'],
          ['审计日志', '密码修改操作记录用户名 + IP 地址'],
          ['CORS 控制', '中间件配置跨域访问'],
          ['密码策略', '最少14位，4类字符至少3类，不含用户名，新旧不同'],
          ['服务账号', '配置文件存储，不硬编码在代码中'],
        ]}
      />

      <Divider />

      <H2>项目状态</H2>
      <Grid columns={2} gap={16}>
        <Stat value="通过" label="后端编译" tone="success" />
        <Stat value="通过" label="前端编译" tone="success" />
      </Grid>

      <Text tone="secondary" size="small">
        项目已完整实现 Spec 中的全部需求，前后端均编译通过，可直接部署使用。
      </Text>
    </Stack>
  );
}
