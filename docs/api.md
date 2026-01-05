# 代理记账ERP API文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **数据格式**: JSON
- **字符编码**: UTF-8

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 人员管理 API

### 1. 获取人员列表

**请求**
```
GET /api/people
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 人员类型 (representative/investor/service_person/mixed) |
| keyword | string | 否 | 搜索关键词（匹配姓名、电话、身份证） |

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
        "type": "representative",
        "name": "张三",
        "phone": "13800138000",
        "id_card": "110101199001011234",
        "representative_customer_ids": "1,5",
        "investor_customer_ids": "",
        "service_customer_ids": "",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 2. 创建人员

**请求**
```
POST /api/people
Content-Type: application/json
```

**请求体**
```json
{
  "type": "representative",
  "name": "张三",
  "phone": "13800138000",
  "id_card": "110101199001011234",
  "password": "abc123"
}
```

### 3. 获取人员详情

**请求**
```
GET /api/people/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "person": { ... },
    "customers": {
      "representative": [],
      "investor": [],
      "service": []
    }
  }
}
```

### 4. 更新人员

**请求**
```
PUT /api/people/:id
Content-Type: application/json
```

### 5. 删除人员

**请求**
```
DELETE /api/people/:id
```

### 6. 获取人员关联的企业

**请求**
```
GET /api/people/:id/customers
```

---

## 客户管理 API

### 1. 获取客户列表

**请求**
```
GET /api/customers
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词（匹配名称、税号、电话） |
| representative | string | 否 | 按法定代表人搜索 |
| investor | string | 否 | 按投资人搜索 |
| service_person | string | 否 | 按服务人员搜索 |

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
        "phone": "13800138000",
        "address": "北京市朝阳区xxx",
        "tax_number": "91110000xxxxxxxx",
        "type": "limited_company",
        "representative_id": 1,
        "investors": [{"person_id": 2, "share_ratio": 51.0}],
        "service_person_ids": "5,6",
        "agreement_ids": "1,3",
        "registered_capital": 1000000.00,
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
  "phone": "13800138000",
  "address": "北京市朝阳区xxx",
  "tax_number": "91110000xxxxxxxx",
  "type": "limited_company",
  "representative_id": 1,
  "investors": [
    {"person_id": 2, "share_ratio": 51.0},
    {"person_id": 3, "share_ratio": 49.0}
  ],
  "service_person_ids": "5,6",
  "agreement_ids": "1,3",
  "registered_capital": 1000000.00
}
```

**客户类型 (type)**
- `limited_company` - 有限公司
- `sole_proprietorship` - 个人独资企业
- `partnership` - 合伙企业
- `individual_business` - 个体工商户

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
    "phone": "13800138000",
    "address": "北京市朝阳区xxx",
    "tax_number": "91110000xxxxxxxx",
    "type": "limited_company",
    "representative_id": 1,
    "representative": {
      "id": 1,
      "name": "张三",
      "phone": "13800138000"
    },
    "investor_list": [
      {"id": 2, "name": "李四", "phone": "13900139000"}
    ],
    "service_persons": [
      {"id": 5, "name": "王五", "phone": "13700137000"}
    ],
    "agreements_list": [...],
    "tasks": [],
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

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tasks | 获取任务列表 |
| POST | /api/tasks | 创建任务 |
| GET | /api/tasks/:id | 获取任务详情 |
| PUT | /api/tasks/:id | 更新任务 |
| DELETE | /api/tasks/:id | 删除任务 |

---

## 协议管理 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/agreements | 获取协议列表 |
| POST | /api/agreements | 创建协议 |
| GET | /api/agreements/:id | 获取协议详情 |
| PUT | /api/agreements/:id | 更新协议 |
| DELETE | /api/agreements/:id | 删除协议 |

---

## 收款管理 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/payments | 获取收款记录 |
| POST | /api/payments | 创建收款记录 |
| GET | /api/payments/:id | 获取收款详情 |
| PUT | /api/payments/:id | 更新收款记录 |
| DELETE | /api/payments/:id | 删除收款记录 |

---

## 统计分析 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/statistics/overview | 首页概览统计 |
| GET | /api/statistics/tasks | 任务统计 |
| GET | /api/statistics/payments | 收款统计 |

---

## 导入导出 API

### 1. 下载导入模板

**请求**
```
GET /api/templates/:type
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 模板类型 (people/customers) |

**响应**
- 返回Excel文件下载

### 2. 导入人员

**请求**
```
POST /api/import/people
Content-Type: multipart/form-data
```

**表单参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | Excel文件 |
| strategy | string | 否 | 冲突策略 (skip/update/create_new) |

**strategy 说明**
- `skip` - 跳过已存在的记录（默认）
- `update` - 更新已存在的记录
- `create_new` - 修改标识后创建新记录

**响应示例**
```json
{
  "code": 0,
  "message": "导入完成",
  "data": {
    "total": 10,
    "success": 8,
    "failed": 2,
    "errors": [
      {"row": 3, "column": "身份证号", "message": "身份证号已存在"},
      {"row": 7, "column": "类型", "message": "类型无效"}
    ]
  }
}
```

### 3. 导入客户

**请求**
```
POST /api/import/customers
Content-Type: multipart/form-data
```

**表单参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | Excel文件 |
| strategy | string | 否 | 冲突策略 (skip/update/create_new) |

**客户导入说明**
- 法定代表人：不存在则自动创建
- 投资人：不存在则自动创建
- 服务人员：必须已存在，否则报错
- 协议：随客户一起创建

**响应示例**
```json
{
  "code": 0,
  "message": "导入完成",
  "data": {
    "total": 5,
    "success": 5,
    "failed": 0,
    "errors": []
  }
}
```

### 4. 导出人员

**请求**
```
GET /api/export/people
```

**响应**
- 返回Excel文件下载

### 5. 导出客户

**请求**
```
GET /api/export/customers
```

**响应**
- 返回Excel文件下载（包含关联人员和协议信息）

---

## 数据模型

### Person (人员)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| type | string | 人员类型 |
| name | string | 姓名 |
| phone | string | 电话 |
| id_card | string | 身份证号（唯一） |
| password | string | 登录密码 |
| representative_customer_ids | string | 担任法人的企业ID（逗号分隔） |
| investor_customer_ids | string | 持股的企业ID（逗号分隔） |
| service_customer_ids | string | 服务的企业ID（逗号分隔） |

### Customer (客户)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 公司名称 |
| phone | string | 联系电话 |
| address | string | 地址 |
| tax_number | string | 税号 |
| type | string | 客户类型 |
| representative_id | uint | 法定代表人ID |
| investors | json | 投资人JSON数组 |
| service_person_ids | string | 服务人员ID（逗号分隔） |
| agreement_ids | string | 代理协议ID（逗号分隔） |
| registered_capital | float64 | 注册资本 |

**investors JSON格式**
```json
[
  {
    "person_id": 1,
    "share_ratio": 25.5,
    "investment_records": [
      {"date": "2024-01-01", "amount": 500000}
    ]
  }
]
```
