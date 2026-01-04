# 代理记账ERP API文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **数据格式**: JSON
- **字符编码**: UTF-8

## 统一响应格式

```json
{
  "code": 0,           // 0表示成功，非0表示错误
  "message": "success",
  "data": {}           // 响应数据
}
```

## 客户管理 API

### 1. 获取客户列表

**请求**
```
GET /api/customers
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词（匹配名称、税号、联系人、电话） |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "items": [
      {
        "id": 1,
        "name": "某某科技有限公司",
        "contact": "张三",
        "phone": "13800138000",
        "email": "contact@example.com",
        "address": "北京市朝阳区xxx",
        "tax_number": "91110000xxxxxxxx",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 2. 创建客户

**请求**
```
POST /api/customers
Content-Type: application/json
```

**请求体**
```json
{
  "name": "某某科技有限公司",
  "contact": "张三",
  "phone": "13800138000",
  "email": "contact@example.com",
  "address": "北京市朝阳区xxx",
  "tax_number": "91110000xxxxxxxx"
}
```

### 3. 获取客户详情

**请求**
```
GET /api/customers/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "某某科技有限公司",
    "contact": "张三",
    "phone": "13800138000",
    "email": "contact@example.com",
    "address": "北京市朝阳区xxx",
    "tax_number": "91110000xxxxxxxx",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "tasks": [],
    "agreements": [],
    "payments": []
  }
}
```

### 4. 更新客户

**请求**
```
PUT /api/customers/:id
Content-Type: application/json
```

### 5. 删除客户

**请求**
```
DELETE /api/customers/:id
```

### 6. 获取客户的任务列表

**请求**
```
GET /api/customers/:id/tasks
```

### 7. 获取客户的收款记录

**请求**
```
GET /api/customers/:id/payments
```

---

## 任务管理 API

### 1. 获取任务列表

**请求**
```
GET /api/tasks
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词（匹配标题、描述） |
| status | string | 否 | 状态筛选 (pending/in_progress/completed) |
| customer_id | int | 否 | 按客户筛选 |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "id": 1,
        "customer_id": 1,
        "title": "完成1月纳税申报",
        "description": "申报增值税、企业所得税",
        "status": "pending",
        "due_date": "2024-01-15T00:00:00Z",
        "completed_at": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司"
        }
      }
    ]
  }
}
```

### 2. 创建任务

**请求**
```
POST /api/tasks
Content-Type: application/json
```

**请求体**
```json
{
  "customer_id": 1,
  "title": "完成1月纳税申报",
  "description": "申报增值税、企业所得税",
  "status": "pending",
  "due_date": "2024-01-15T00:00:00Z"
}
```

### 3. 获取任务详情

**请求**
```
GET /api/tasks/:id
```

### 4. 更新任务

**请求**
```
PUT /api/tasks/:id
Content-Type: application/json
```

**注意**: 当状态更新为 `completed` 时，系统会自动设置 `completed_at` 时间。

### 5. 删除任务

**请求**
```
DELETE /api/tasks/:id
```

---

## 协议管理 API

### 1. 获取协议列表

**请求**
```
GET /api/agreements
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索协议编号 |
| status | string | 否 | 状态筛选 (active/expired/cancelled) |
| customer_id | int | 否 | 按客户筛选 |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 3,
    "items": [
      {
        "id": 1,
        "customer_id": 1,
        "agreement_number": "AGR2024001",
        "start_date": "2024-01-01T00:00:00Z",
        "end_date": "2024-12-31T00:00:00Z",
        "fee_type": "月度",
        "amount": 500.00,
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司"
        }
      }
    ]
  }
}
```

### 2. 创建协议

**请求**
```
POST /api/agreements
Content-Type: application/json
```

**请求体**
```json
{
  "customer_id": 1,
  "agreement_number": "AGR2024001",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z",
  "fee_type": "月度",
  "amount": 500.00,
  "status": "active"
}
```

### 3. 获取协议详情

**请求**
```
GET /api/agreements/:id
```

### 4. 更新协议

**请求**
```
PUT /api/agreements/:id
Content-Type: application/json
```

### 5. 删除协议

**请求**
```
DELETE /api/agreements/:id
```

---

## 收款管理 API

### 1. 获取收款记录列表

**请求**
```
GET /api/payments
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_id | int | 否 | 按客户筛选 |
| start_date | string | 否 | 开始日期 (格式: 2024-01-01) |
| end_date | string | 否 | 结束日期 (格式: 2024-12-31) |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 12,
    "items": [
      {
        "id": 1,
        "customer_id": 1,
        "agreement_id": 1,
        "amount": 500.00,
        "payment_date": "2024-01-05T00:00:00Z",
        "payment_method": "转账",
        "period": "2024-01",
        "remark": "1月代理记账服务费",
        "created_at": "2024-01-05T00:00:00Z",
        "updated_at": "2024-01-05T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司"
        },
        "agreement": {
          "id": 1,
          "agreement_number": "AGR2024001"
        }
      }
    ]
  }
}
```

### 2. 创建收款记录

**请求**
```
POST /api/payments
Content-Type: application/json
```

**请求体**
```json
{
  "customer_id": 1,
  "agreement_id": 1,
  "amount": 500.00,
  "payment_date": "2024-01-05T00:00:00Z",
  "payment_method": "转账",
  "period": "2024-01",
  "remark": "1月代理记账服务费"
}
```

### 3. 获取收款记录详情

**请求**
```
GET /api/payments/:id
```

### 4. 更新收款记录

**请求**
```
PUT /api/payments/:id
Content-Type: application/json
```

### 5. 删除收款记录

**请求**
```
DELETE /api/payments/:id
```

---

## 统计分析 API

### 1. 获取首页概览统计

**请求**
```
GET /api/statistics/overview
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "customer_count": 25,
    "pending_task_count": 8,
    "active_agreement_count": 20,
    "monthly_payment": 12500.00,
    "yearly_payment": 150000.00
  }
}
```

### 2. 获取任务统计

**请求**
```
GET /api/statistics/tasks
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "pending": 5,
    "in_progress": 3,
    "completed": 42
  }
}
```

### 3. 获取收款统计

**请求**
```
GET /api/statistics/payments
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 开始日期 (格式: 2024-01-01) |
| end_date | string | 否 | 结束日期 (格式: 2024-12-31) |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_amount": 12500.00,
    "count": 25
  }
}
```

---

## 运行项目

```bash
# 启动服务
go run main.go

# 服务运行在 http://localhost:8080
```

## 数据库

数据库文件位于 `database/erp.db`，可以使用 SQLite 客户端工具查看数据。

## 后期迁移到MySQL

1. 修改 `config/database.go` 中的驱动和DSN
2. 更改 `go.mod` 中的依赖（使用 `gorm.io/driver/mysql`）
3. 重新运行自动迁移即可
