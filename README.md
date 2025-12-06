# 基层三基考试系统

基于Go语言开发的基层医院三基考试系统，支持用户管理、题库管理、试卷生成、在线考试和成绩统计等功能。

## 功能特性

- **用户管理**：手机号/身份证号注册登录，中文姓名验证，密码复杂度要求
- **角色权限**：站长、管理员、员工三种角色的权限控制
- **题库管理**：四个科目（基础知识、基本技能、专业知识、专业实践技能）的题库和题目管理
- **试卷生成**：根据题库随机生成试卷
- **考试功能**：考试开始、答题、提交和自动评分
- **成绩统计**：历史考试记录和统计数据

## 技术栈

- **后端**：Go 1.25, Gin框架
- **数据库**：MySQL 5.7
- **认证**：JWT
- **部署**：Docker

## 快速开始

### 环境要求

- Docker 19.03+ 
- Docker Compose 1.25+

### Docker部署

1. **克隆代码**

```bash
git clone https://github.com/hangbin2008/sanjicms.git
cd sanjicms
```

2. **配置环境变量**

复制 `.env.example` 文件为 `.env`，并根据需要修改配置：

```bash
cp .env.example .env
```

3. **启动服务**

使用 Docker Compose 启动应用程序和数据库服务（使用GitHub Container Registry上的预构建镜像）：

```bash
docker-compose up -d
```

系统使用的是预构建镜像：`ghcr.io/hangbin2008/sanjicms:latest`

4. **访问应用**

应用启动后，可通过以下地址访问：

- 首页：http://localhost:8080
- 登录页：http://localhost:8080/login
- 注册页：http://localhost:8080/register
- API文档：http://localhost:8080/api
- 健康检查：http://localhost:8080/health

### 手动部署

1. **安装依赖**

```bash
go mod download
```

2. **创建数据库**

使用MySQL客户端创建数据库，并执行初始化脚本：

```bash
mysql -u root -p < migrations/001_init_schema.sql
```

3. **配置环境变量**

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your-db-password
export DB_NAME=jiceng_sanji_exam
```

4. **运行应用**

```bash
go run cmd/server/main.go
```

## 项目结构

```
├── cmd/                  # 应用程序入口
│   └── server/           # 服务器入口
├── internal/             # 内部代码
│   ├── api/              # API控制器和路由
│   ├── db/               # 数据库连接
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   └── service/          # 业务逻辑
├── migrations/           # 数据库迁移脚本
├── pkg/                  # 公共包
│   └── config/           # 配置管理
├── static/               # 静态文件
├── templates/            # HTML模板
├── Dockerfile            # Docker构建文件
├── docker-compose.yml    # Docker Compose配置
├── .env.example          # 环境变量示例
├── go.mod                # Go模块依赖
└── go.sum                # 依赖校验
```

## API文档

### 公共API

- **POST /api/register** - 用户注册
- **POST /api/login** - 用户登录

### 受保护API

- **GET /api/user/me** - 获取当前用户信息
- **POST /api/banks** - 创建题库
- **GET /api/banks** - 获取题库列表
- **GET /api/banks/:id** - 获取题库详情
- **POST /api/questions** - 创建题目
- **GET /api/questions/bank/:bank_id** - 获取题库下的题目列表
- **GET /api/questions/:id** - 获取题目详情
- **POST /api/exams/generate** - 生成试卷
- **GET /api/exams/:id** - 获取试卷详情
- **POST /api/exams/:id/start** - 开始考试
- **POST /api/exams/submit** - 提交试卷
- **GET /api/records** - 获取考试记录列表
- **GET /api/records/:id** - 获取考试记录详情
- **GET /api/records/stats** - 获取考试统计数据

## 初始账号

系统初始化时会创建一个默认管理员账号：

- 用户名：admin
- 密码：Admin@123

## 开发指南

### 添加新功能

1. 在 `internal/models` 中定义数据模型
2. 在 `internal/service` 中实现业务逻辑
3. 在 `internal/api` 中添加API路由和控制器
4. 在 `templates` 中添加前端页面（如果需要）

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！
