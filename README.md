# GiveTrack（公益捐赠追踪平台）

> 项目类型：全栈 Web 应用

GiveTrack 支持公益组织注册发布公益项目，平台审核通过后按助学、助老等分类展示；用户注册后选择项目捐款并生成电子凭证，项目页实时展示筹款进度与捐赠记录，项目方定期上传执行进展，个人中心可查捐赠记录与累计金额，平台提供捐款金额与服务时长公益排行榜。

## 快速启动（Docker Compose 一键部署，首选）

```bash
# 1. 首次启动前复制环境变量（务必修改 JWT_SECRET）
cp .env.example .env

# 2. 启动全部服务
docker compose up -d

# 3. 查看健康状态
docker compose ps
```

访问地址：

- 前端：http://localhost:8235
- 后端 API：http://localhost:3235
- 存活检查：http://localhost:3235/healthz
- 就绪检查：http://localhost:3235/readyz
- Swagger：http://localhost:3235/docs/index.html

演示账号：

| 角色 | 账号 | 密码 |
| --- | --- | --- |
| 管理员 | admin | admin123 |
| 公益组织 | careorg | org123 |
| 普通用户 | donor1 | user123 |

## 本地开发

```bash
# 前端
cd frontend
npm install
npm run dev        # http://localhost:8235

# 后端
cd backend
go mod tidy
go run ./cmd/server
```

## 技术栈

| 层次 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Vite + Tailwind CSS |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT (github.com/golang-jwt/jwt/v5) + RBAC |
| 密码哈希 | golang.org/x/crypto/bcrypt |
| 参数校验 | github.com/go-playground/validator/v10 |
| 配置 | github.com/caarlos0/env/v11 |
| 日志 | Go 标准库 log/slog |
| API 文档 | github.com/swaggo/gin-swagger |

## 项目目录结构

```
.
├── docker-compose.yml
├── .env / .env.example
├── README.md
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf          # /api/ 反代到 backend:8080/api/v1/
│   └── src/                # api / components / context / pages
└── backend/
    ├── Dockerfile
    ├── database/init.sql
    ├── cmd/server/main.go
    └── internal/
        ├── config/  constants/  model/  repository/  service/  handler/
        ├── middleware/  router/  util/  database/
```

## 主要 API 列表

统一前缀 `/api/v1`，统一响应 `{ "code": 0, "message": "ok", "data": ... }`。

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | /auth/register | 注册（个人/组织） | - |
| POST | /auth/login | 登录 | - |
| GET | /auth/me | 当前用户 | JWT |
| PUT | /auth/profile | 更新资料 | JWT |
| GET | /projects | 项目列表（分类/状态/分页） | - |
| GET | /projects/:id | 项目详情+捐赠记录+进展 | - |
| POST | /projects | 发布项目 | org |
| GET | /projects/org/my | 我的项目 | org |
| GET/POST | /projects/:id/updates | 项目进展 | org |
| POST | /donations | 捐款并生成凭证 | JWT |
| GET | /donations/my | 我的捐赠 | JWT |
| GET | /donations/:id/certificate | 电子凭证 | JWT |
| GET | /ranking/donation | 捐款排行榜 | - |
| GET | /ranking/service | 服务时长排行榜 | - |
| GET | /ranking/stats | 平台统计 | - |
| GET | /admin/projects/pending | 待审核项目 | admin |
| POST | /admin/projects/:id/review | 项目审核 | admin |
| GET | /admin/organizations/pending | 待审核组织 | admin |
| POST | /admin/organizations/:id/review | 组织审核 | admin |
| GET | /healthz | 存活检查 | - |
| GET | /readyz | 就绪检查（DB ping） | - |

## 环境变量

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | Compose 项目名 | givetrack |
| APP_ENV | 运行环境 | development |
| SERVER_PORT | 后端端口（容器内） | 8080 |
| DB_HOST / DB_PORT | 数据库地址/端口 | db / 3306 |
| DB_NAME / DB_USER / DB_PASSWORD | 数据库配置 | givetrack |
| DB_ROOT_PASSWORD | MySQL root 密码 | givetrack_root_pwd |
| DB_MAX_OPEN_CONNS | 最大连接数 | 50 |
| DB_MAX_IDLE_CONNS | 最大空闲连接数 | 10 |
| DB_CONN_MAX_LIFETIME_MINUTES | 连接最大存活分钟数 | 60 |
| DB_CONNECT_RETRY_COUNT | 启动期连接重试次数 | 10 |
| DB_CONNECT_RETRY_INTERVAL_SECONDS | 重试间隔秒数 | 3 |
| JWT_SECRET | JWT 密钥 | 请修改（生产至少 32 位） |
| CORS_ALLOWED_ORIGINS | 允许的跨域来源 | *（开发）/ 显式来源（生产） |
| AUTH_RATE_LIMIT | 登录/注册窗口内每 IP 最大请求数 | 10 |
| AUTH_RATE_WINDOW_SECONDS | 限流窗口秒数 | 60 |
| FRONTEND_PORT | 前端宿主端口 | 8235 |
| BACKEND_PORT | 后端宿主端口 | 3235 |
| DB_PORT | MySQL 宿主端口 | 33335 |

## 生产安全提示

- 生产环境（`APP_ENV=production`）下，`JWT_SECRET` 必须为至少 32 位的随机值且不得使用默认值，`CORS_ALLOWED_ORIGINS` 禁止使用 `*`，否则服务会拒绝启动。
- 开发/测试环境使用 GORM AutoMigrate；生产环境建议使用 `backend/database/init.sql` 或独立 migration 工具管理表结构。

## Docker 部署说明

- 端口映射：前端 `8235:80`、后端 `3235:8080`、MySQL `33335:3306`
- 前端 Nginx 将 `/api/` 反向代理到 `http://backend:8080/api/v1/`，SPA 路由由 `try_files` 兜底
- 数据卷：`db_data` 持久化 MySQL 数据
- 数据库 healthcheck 通过后后端才启动；后端提供 `/healthz`（存活）与 `/readyz`（就绪）供编排检查
- 前端 Dockerfile 使用本地已构建的 `dist/` 产物（本地开发先执行 `npm run build`），避免容器内网络下载依赖

## License

MIT License
