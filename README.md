# 代理记账ERP系统

基于Go语言开发的简易代理记账ERP系统后端API，支持客户管理、任务管理、协议管理和收款管理等功能。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://github.com/go-gorm/gorm)
- **数据库**: SQLite (支持迁移到MySQL)

## 功能特性

### 核心模块
- **客户管理** - 管理客户基础信息、联系人、税号等
- **任务管理** - 管理客户代办任务，支持状态跟踪和截止日期
- **协议管理** - 管理代理记账协议，支持服务费和有效期管理
- **收款管理** - 记录和查询收款情况，支持按时间范围筛选
- **统计分析** - 首页概览、任务统计、收款汇总

### 查询功能
- 关键词搜索（客户名称、税号等）
- 状态筛选
- 时间范围查询
- 多条件组合筛选

## 快速开始

### 环境要求

- Go 1.21 或更高版本

### 安装运行

```bash
# 克隆项目
git clone http://192.168.3.20/hashqq/erp.git
cd erp

# 安装依赖
go mod download

# 运行服务
go run main.go
```

服务启动后监听在 `http://localhost:8080`

### 编译

```bash
go build -o erp main.go
./erp
```

## 项目结构

```
erp/
├── main.go                 # 程序入口
├── go.mod                  # 依赖管理
├── config/                 # 配置
│   └── database.go         # 数据库配置
├── models/                 # 数据模型
│   ├── customer.go         # 客户
│   ├── task.go             # 任务
│   ├── agreement.go        # 协议
│   └── payment.go          # 收款
├── controllers/            # 控制器
│   ├── common.go           # 通用响应
│   ├── customer_controller.go
│   ├── task_controller.go
│   ├── agreement_controller.go
│   ├── payment_controller.go
│   └── statistics_controller.go
├── routes/                 # 路由
│   └── routes.go
├── docs/                   # 文档
│   └── api.md              # API文档
└── database/               # 数据库文件
    └── erp.db
```

## API文档

完整的API文档请查看 [docs/api.md](docs/api.md)

### 主要端点

| 模块 | 端点 | 说明 |
|------|------|------|
| 客户 | `GET /api/customers` | 获取客户列表 |
| 任务 | `GET /api/tasks` | 获取任务列表 |
| 协议 | `GET /api/agreements` | 获取协议列表 |
| 收款 | `GET /api/payments` | 获取收款记录 |
| 统计 | `GET /api/statistics/overview` | 首页统计 |

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 数据库

项目默认使用SQLite数据库，数据库文件位于 `database/erp.db`。

### 迁移到MySQL

修改 `config/database.go`，使用MySQL驱动：

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// 替换SQLite连接为MySQL
dsn := "user:password@tcp(127.0.0.1:3306)/erp?charset=utf8mb4&parseTime=True&loc=Local"
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

安装MySQL驱动：
```bash
go get gorm.io/driver/mysql
```

## 开发计划

查看 [TODO.md](TODO.md) 了解当前进度和待实现功能。

## License

MIT
