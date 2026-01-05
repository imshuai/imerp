# 代理记账ERP系统

基于Go语言和Vue 3开发的代理记账ERP系统，支持客户管理、人员管理、任务管理、协议管理和收款管理等功能。采用前后端分离架构，支持单文件部署。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://github.com/go-gorm/gorm)
- **数据库**: SQLite (支持迁移到MySQL)
- **Excel处理**: [excelize](https://github.com/xuri/excelize)

### 前端
- **框架**: [Vue 3](https://vuejs.org/)
- **构建工具**: [Vite](https://vitejs.dev/)
- **UI组件**: [Element Plus](https://element-plus.org/)
- **HTTP客户端**: [Axios](https://axios-http.com/)
- **路由**: [Vue Router](https://router.vuejs.org/)
- **状态管理**: [Pinia](https://pinia.vuejs.org/)
- **语言**: TypeScript

## 功能特性

### 核心模块
- **人员管理** - 统一管理法定代表人、投资人、服务人员
- **客户管理** - 企业信息，关联法定代表人、投资人、服务人员、协议
- **任务管理** - 客户代办任务，支持状态跟踪和截止日期
- **协议管理** - 代理记账协议，支持服务费和有效期管理
- **收款管理** - 收款记录，支持按时间范围筛选
- **统计分析** - 首页概览、任务统计、收款汇总
- **导入导出** - Excel批量导入/导出人员和客户数据

### 人员类型
- **法定代表人** - 企业法人代表
- **投资人** - 企业股东，支持持股比例和多次出资记录
- **服务人员** - 服务该客户的员工

### 客户类型
- 有限公司
- 个人独资企业
- 合伙企业
- 个体工商户

### 查询功能
- 关键词搜索（客户名称、税号等）
- 按人员搜索（法定代表人、投资人、服务人员）
- 状态筛选
- 时间范围查询
- 多条件组合筛选

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- Node.js 18+ 和 npm（开发前端时需要）

### 安装运行

#### 完整安装（包含前端）

```bash
# 克隆项目
git clone http://192.168.3.20/hashqq/erp.git
cd erp

# 安装后端依赖
go mod download

# 安装前端依赖并构建
cd frontend
npm install
npm run build
cd ..

# 运行服务
go run main.go
```

服务启动后监听在 `http://localhost:8080`

#### 仅运行后端（已包含预构建的前端）

```bash
# 克隆项目
git clone http://192.168.3.20/hashqq/erp.git
cd erp

# 安装后端依赖
go mod download

# 运行服务
go run main.go
```

#### 前端开发模式

```bash
# 在 frontend 目录下运行开发服务器
cd frontend
npm run dev

# 在另一个终端运行后端
cd ..
go run main.go
```

前端开发服务器运行在 `http://localhost:5173`，并自动代理API请求到后端。

### 编译

#### 生产构建（单文件部署）

```bash
# 1. 构建前端
cd frontend
npm run build
cd ..

# 2. 编译Go程序（前端资源已嵌入）
go build -o erp main.go

# 3. 运行
./erp
```

编译后的 `erp` 可执行文件已包含所有前端资源，可以直接部署到服务器运行。

#### 交叉编译

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o erp-linux-arm64 main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o erp.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o erp-macos main.go
```

## 项目结构

```
erp/
├── main.go                 # 程序入口
├── go.mod                  # 依赖管理
├── config/                 # 配置
│   └── database.go         # 数据库配置
├── models/                 # 数据模型
│   ├── person.go           # 人员信息
│   ├── customer.go         # 客户信息
│   ├── task.go             # 任务
│   ├── agreement.go        # 协议
│   └── payment.go          # 收款
├── controllers/            # 控制器
│   ├── common.go           # 通用响应
│   ├── person_controller.go    # 人员控制器
│   ├── customer_controller.go  # 客户控制器
│   ├── task_controller.go      # 任务控制器
│   ├── agreement_controller.go # 协议控制器
│   ├── payment_controller.go   # 收款控制器
│   ├── statistics_controller.go # 统计控制器
│   └── import_export_controller.go # 导入导出控制器
├── routes/                 # 路由
│   ├── routes.go           # API路由
│   └── frontend.go         # 前端路由
├── services/               # 服务层
│   └── import_export/      # 导入导出服务
│       ├── excel_service.go      # Excel基础服务
│       ├── template_service.go   # 模板生成服务
│       ├── people_import.go      # 人员导入服务
│       ├── customer_import.go    # 客户导入服务
│       └── export_service.go     # 导出服务
├── utils/                  # 工具函数
│   └── excel_utils.go      # Excel工具函数
├── embedded/               # 嵌入的静态资源
│   ├── static.go           # Go embed 文件
│   └── dist/               # 前端构建产物（git忽略）
├── frontend/               # 前端源码
│   ├── src/
│   │   ├── api/            # API调用封装
│   │   ├── components/     # Vue组件
│   │   ├── views/          # 页面视图
│   │   ├── router/         # 路由配置
│   │   ├── App.vue
│   │   └── main.ts
│   ├── public/
│   ├── index.html
│   ├── vite.config.ts
│   ├── package.json
│   └── tsconfig.json
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
| 人员 | `GET /api/people` | 获取人员列表 |
| 客户 | `GET /api/customers` | 获取客户列表 |
| 任务 | `GET /api/tasks` | 获取任务列表 |
| 协议 | `GET /api/agreements` | 获取协议列表 |
| 收款 | `GET /api/payments` | 获取收款记录 |
| 统计 | `GET /api/statistics/overview` | 首页统计 |
| 模板 | `GET /api/templates/:type` | 下载导入模板 |
| 导入 | `POST /api/import/people` | 导入人员 |
| 导入 | `POST /api/import/customers` | 导入客户 |
| 导出 | `GET /api/export/people` | 导出人员 |
| 导出 | `GET /api/export/customers` | 导出客户 |

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

### 数据模型

#### Person（人员）
- 存储法定代表人、投资人、服务人员信息
- 身份证号唯一约束
- 支持混合角色（一个人既是投资人又是服务人员等）

#### Customer（客户）
- 企业基础信息
- 关联法定代表人（一对一）
- 关联投资人（一对多，含持股比例和出资记录）
- 关联服务人员（多对多）
- 关联代理协议
- 注册资本

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
