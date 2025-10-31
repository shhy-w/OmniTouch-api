# Cosy Example API

基于 [Cosy](https://github.com/uozi-tech/cosy) 框架的 Go API 开发示例项目，提供了完整的用户管理、权限控制、文件上传等核心功能，可作为 API 项目的开发起点。

## 📋 目录

- [技术栈](#技术栈)
- [功能特性](#功能特性)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [开发指南](#开发指南)
- [API 示例](#api-示例)
- [部署说明](#部署说明)

## 🛠 技术栈

- **框架核心**: [Cosy](https://github.com/uozi-tech/cosy) - 提供统一的 API 开发模式
- **HTTP 框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能的 HTTP Web 框架
- **ORM**: [GORM](https://gorm.io/) - 强大的 Go ORM 库
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **认证**: JWT (JSON Web Token)
- **权限控制**: ACL (Access Control List)
- **代码生成**: GORM Gen - 类型安全的查询构建器

## ✨ 功能特性

- ✅ **用户管理**: 用户注册、登录、CRUD 操作
- ✅ **权限控制**: 基于用户组的 ACL 权限系统
- ✅ **JWT 认证**: Token 机制，支持自动续期
- ✅ **限流保护**: 基于 Redis 的请求限流
- ✅ **文件上传**: 支持 OSS/MinIO 对象存储
- ✅ **邮件服务**: 邮件发送功能
- ✅ **审计日志**: 操作审计和日志收集（支持阿里云 SLS）
- ✅ **软删除**: 数据软删除机制
- ✅ **缓存机制**: Redis 缓存用户和权限信息
- ✅ **数据库视图**: 支持创建和管理数据库视图

## 📁 项目结构

```
cosy-example-api/
├── api/                    # API 路由和处理器
│   ├── admin/             # 管理员 API
│   ├── global/            # 全局 API (认证、用户信息等)
│   ├── settings/          # 设置相关 API
│   ├── upload/            # 文件上传 API
│   └── user/              # 用户相关 API
├── cmd/                    # 命令行工具
│   └── generate/          # 代码生成工具
├── internal/              # 内部包（不对外暴露）
│   ├── acl/               # 权限控制逻辑
│   ├── audit/             # 审计日志
│   ├── limiter/           # 限流器
│   ├── mail/              # 邮件服务
│   ├── upload/            # 文件上传处理
│   └── user/              # 用户相关业务逻辑
├── model/                 # 数据模型
│   └── view/              # 数据库视图定义
├── query/                 # GORM Gen 生成的类型安全查询
├── router/                # 路由配置和中间件
├── settings/              # 配置管理
├── main.go               # 应用入口
├── app.example.ini       # 配置文件示例
└── Dockerfile            # Docker 镜像构建文件
```

## 🚀 快速开始

### 前置要求

- Go 1.24+ 
- MySQL 8.0+
- Redis 6.0+
- 可选: 阿里云 OSS/MinIO (用于文件上传)
- 可选: 邮件服务配置 (用于邮件发送)

### 安装步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd cosy-example-api
```

2. **安装依赖**

```bash
go mod download
```

3. **配置数据库**

创建 MySQL 数据库（例如: `burn`）

4. **复制配置文件**

```bash
cp app.example.ini app.ini
```

5. **编辑配置文件**

编辑 `app.ini`，填写数据库、Redis 等配置信息：

```ini
[database]
Host        = 127.0.0.1
Port        = 3306
User        = root
Password    = your_password
Name        = burn
```

```ini
[redis]
Addr     = 127.0.0.1:6379
Password = 
DB       = 0
```

6. **生成查询代码（首次运行前）**

```bash
go run cmd/generate/generate.go
```

或者：

```bash
./gen.sh
```

7. **运行项目**

```bash
go run main.go
```

或者使用 Air 进行热重载开发：

```bash
air
```

项目将在 `http://127.0.0.1:9000` 启动。

### 初始管理员账户

首次启动后，系统会自动创建初始管理员账户：
- **邮箱**: `admin`
- **密码**: `admin`

**⚠️ 重要**: 首次登录后请立即修改管理员密码！

## ⚙️ 配置说明

### app.ini 配置详解

```ini
[app]
PageSize  = 20           # 默认分页大小
JwtSecret = your-secret  # JWT 签名密钥（建议使用随机字符串）

[server]
Host    = 127.0.0.1      # 服务监听地址
Port    = 9000          # 服务监听端口
RunMode = debug         # 运行模式: debug/release

[database]
Host        = 127.0.0.1  # 数据库地址
Port        = 3306       # 数据库端口
User        = root       # 数据库用户名
Password    =            # 数据库密码
Name        = burn       # 数据库名称
TablePrefix =            # 表前缀（可选）

[redis]
Addr     = 127.0.0.1:6379  # Redis 地址
Password =                 # Redis 密码（如果有）
DB       = 0               # Redis 数据库编号
Prefix   = burn            # Redis Key 前缀

[mail]
Host     = smtp.example.com  # SMTP 服务器地址
Port     = 587                # SMTP 端口
Email    = your@email.com     # 发件人邮箱
Password = your_password      # 邮箱密码

[oss]
AccessKeyId     =         # OSS AccessKey ID
AccessKeySecret =         # OSS AccessKey Secret
EndPoint        =         # OSS Endpoint
BucketName      =         # OSS Bucket 名称
BaseUrl         =         # OSS 访问基础 URL

[sls]
AccessKeyId     =         # 阿里云 SLS AccessKey ID
AccessKeySecret =         # 阿里云 SLS AccessKey Secret
EndPoint        =         # SLS Endpoint
ProjectName     =         # SLS 项目名称
LogStoreName    =         # SLS LogStore 名称
Source          =         # 日志来源标识

[crypto]
Chacha20Key   =           # ChaCha20 加密密钥
Chacha20Nonce =           # ChaCha20 Nonce
AESKey        =           # AES 加密密钥
AESIv         =           # AES 初始化向量
```

## 📖 开发指南

### 1. 创建数据模型

在 `model/` 目录下创建新的模型文件，例如 `article.go`:

```go
package model

type Article struct {
    Model
    Title   string `json:"title" cosy:"add:required;update:omitempty;list:fussy"`
    Content string `json:"content" cosy:"add:required;update:omitempty"`
    UserID  uint64 `json:"user_id" cosy:"list:eq" gorm:"index"`
    User    *User  `json:"user,omitempty" cosy:"item:preload;list:preload"`
}

// 在 model/model.go 的 GenerateAllModel() 中注册模型
func GenerateAllModel() []any {
    return []any{
        User{},
        UserGroup{},
        Article{},  // 新增
        // ... 其他模型
    }
}
```

**Cosy 标签说明：**
- `add:required` - 创建时必填
- `update:omitempty` - 更新时可选
- `list:fussy` - 列表查询支持模糊搜索
- `list:eq` - 列表查询支持等值查询
- `list:in` - 列表查询支持 IN 查询
- `item:preload` - 单个查询时预加载关联
- `list:preload` - 列表查询时预加载关联

### 2. 注册路由

在 `router/routers.go` 中注册路由：

```go
func InitRouter() {
    r := cosy.GetEngine()
    
    // 需要认证的路由组
    authRouter := r.Group("api", router.AuthRequired())
    {
        article.InitArticleRouter(authRouter)
    }
    
    // 管理员路由组
    adminRouter := r.Group("admin", router.AuthAdminRequired())
    {
        // ... 管理员路由
    }
}
```

### 3. 创建 API 处理器

在 `api/` 目录下创建 API 处理器，例如 `api/article/router.go`:

```go
package article

import (
    "git.uozi.org/uozi/cosy-example-api/model"
    "github.com/gin-gonic/gin"
    "github.com/uozi-tech/cosy"
)

func InitArticleRouter(g *gin.RouterGroup) {
    c := cosy.Api[model.Article]("articles")
    
    // 创建前钩子
    c.CreateHook(func(c *cosy.Ctx[model.Article]) {
        // 自动设置作者
        user := api.CurrentUser(c.GinContext)
        c.BeforeExecuteHook(func(ctx *cosy.Ctx[model.Article]) {
            ctx.Model.UserID = user.ID
        })
    })
    
    c.InitRouter(g)
}
```

### 4. 使用 Cosy Core API

Cosy 提供了便捷的 CRUD API：

```go
// 获取单个资源
func GetArticle(c *gin.Context) {
    cosy.Core[model.Article](c).Get()
}

// 获取列表（分页）
func GetArticleList(c *gin.Context) {
    cosy.Core[model.Article](c).PagingList()
}

// 创建资源
func CreateArticle(c *gin.Context) {
    cosy.Core[model.Article](c).Create()
}

// 更新资源
func UpdateArticle(c *gin.Context) {
    cosy.Core[model.Article](c).Modify()
}

// 删除资源（软删除）
func DeleteArticle(c *gin.Context) {
    cosy.Core[model.Article](c).Destroy()
}
```

### 5. 权限控制

使用 ACL 进行权限控制：

```go
import (
    "git.uozi.org/uozi/cosy-example-api/api"
    "git.uozi.org/uozi/cosy-example-api/internal/acl"
)

// 在路由中使用权限守卫
func InitArticleRouter(g *gin.RouterGroup) {
    c := cosy.Api[model.Article]("articles")
    
    // 应用权限控制
    api.CosyGuard[model.Article](c, acl.Article)
    
    c.InitRouter(g)
}

// 在处理器中检查权限
func SomeHandler(c *gin.Context) {
    if !api.Can(c, acl.Article, acl.Write) {
        c.JSON(403, gin.H{"message": "无权限"})
        return
    }
    // 处理逻辑
}
```

### 6. 使用查询构建器

GORM Gen 生成的类型安全查询：

```go
import "git.uozi.org/uozi/cosy-example-api/query"

// 查询用户
user, err := query.User.Where(query.User.Email.Eq("user@example.com")).First()

// 更新用户
err = query.User.Where(query.User.ID.Eq(1)).
    Update(query.User.Name, "新名称")

// 删除用户
err = query.User.Where(query.User.ID.Eq(1)).Delete()
```

### 7. 使用 Redis 缓存

```go
import "github.com/uozi-tech/cosy/redis"

// 设置缓存
err := redis.Set("key", "value", time.Hour)

// 获取缓存
value, err := redis.Get("key")

// 删除缓存
err := redis.Del("key")
```

### 8. 中间件开发

在 `router/middleware.go` 中添加自定义中间件：

```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}

// 使用中间件
adminRouter.Use(CustomMiddleware())
```

### 9. 限流中间件

```go
import (
    rate_limiter "git.uozi.com/uozi/rate-limiter-go"
)

func InitArticleRouter(g *gin.RouterGroup) {
    // 应用限流：每 IP 每分钟最多 10 次请求
    g.Use(router.LimiterMiddleware(&rate_limiter.LimitConf{
        Rate:  10,
        Burst: 10,
        Period: time.Minute,
    }))
    
    // ... 路由配置
}
```

## 📝 API 示例

### 用户登录

```bash
POST /admin/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

响应示例：
```json
{
  "message": "ok",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "用户名"
  }
}
```

**注意**: 
- Token 可以通过以下方式传递：
  - 请求头: `Token: <your-token>`
  - 查询参数: `?token=<base64-encoded-token>`（需要 base64 编码）
- Token 有效期：
  - 普通用户: 6 小时
  - 管理员: 15 天
- Token 自动刷新：当 Token 剩余时间少于 30 分钟时，系统会自动刷新并在响应头中返回新的 Token (`refresh-token`)

### 用户登出

```bash
DELETE /admin/logout
Token: <your-token>
```

### 获取当前用户信息

```bash
GET /admin/user
Token: <your-token>
```

### 更新当前用户信息

```bash
POST /admin/user
Token: <your-token>
Content-Type: application/json

{
  "name": "新用户名",
  "email": "newemail@example.com"
}
```

### 创建用户（管理员）

```bash
POST /admin/users
Token: <your-token>
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "用户名",
  "status": 1
}
```

### 使用 Cosy API 进行 CRUD

**创建资源**
```bash
POST /admin/articles
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "文章标题",
  "content": "文章内容"
}
```

**获取列表（支持分页和过滤）**
```bash
GET /admin/articles?page=1&page_size=20&title=搜索关键词
Authorization: Bearer <token>
```

**获取单个资源**
```bash
GET /admin/articles/1
Authorization: Bearer <token>
```

**更新资源**
```bash
PUT /admin/articles/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "新标题"
}
```

**删除资源**
```bash
DELETE /admin/articles/1
Authorization: Bearer <token>
```

## 🚢 部署说明

### Docker 部署

1. **构建镜像**

```bash
./build.sh
```

或手动构建：

```bash
GOOS=linux GOARCH=amd64 go build -o burn-api main.go
docker build -t cosy-api:latest .
```

2. **运行容器**

```bash
docker run -d \
  --name cosy-api \
  -p 9000:9000 \
  -v /path/to/app.ini:/config/app.ini \
  cosy-api:latest
```

### 使用 Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "9000:9000"
    volumes:
      - ./app.ini:/config/app.ini
    depends_on:
      - mysql
      - redis
  
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: burn
    ports:
      - "3306:3306"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

运行：

```bash
docker-compose up -d
```

## 🔧 开发工具

### 代码生成

生成 GORM Gen 查询代码：

```bash
go run cmd/generate/generate.go
```

或使用脚本：

```bash
./gen.sh
```

### 运行测试

```bash
go test ./...
```

或使用测试脚本：

```bash
./test.sh
```

## 📚 更多资源

- [Cosy 框架文档](https://github.com/uozi-tech/cosy)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [GORM Gen 文档](https://gorm.io/gen/)

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。

---

**提示**: 在开发过程中，建议使用 Air 进行热重载开发，修改代码后会自动重新编译运行，提升开发效率。

```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 运行 Air
air
```

如有问题或建议，欢迎提交 Issue！
