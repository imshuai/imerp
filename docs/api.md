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
| is_service_person | boolean | 否 | 是否为服务人员 |
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
        "is_service_person": false,
        "name": "张三",
        "phone": "13800138000",
        "id_card": "110101199001011234",
        "password": "abc123",
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
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| is_service_person | boolean | 否 | 是否为服务人员（默认 false） |
| name | string | 是 | 姓名 |
| phone | string | 是 | 电话 |
| id_card | string | 是 | 身份证号（唯一） |
| password | string | 否 | 登录密码 |

**请求体示例**
```json
{
  "is_service_person": false,
  "name": "张三",
  "phone": "13800138000",
  "id_card": "110101199001011234",
  "password": "abc123"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "is_service_person": false,
    "name": "张三",
    "phone": "13800138000",
    "id_card": "110101199001011234",
    "password": "abc123",
    "created_at": "2024-01-01T00:00:00Z"
  }
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
    "person": {
      "id": 1,
      "is_service_person": false,
      "name": "张三",
      "phone": "13800138000",
      "id_card": "110101199001011234",
      "representative_customer_ids": "1,5",
      "investor_customer_ids": "",
      "service_customer_ids": ""
    },
    "customers": {
      "representative": [
        {"id": 1, "name": "某某科技有限公司", "tax_number": "91110000MA001234XX"},
        {"id": 5, "name": "某某商贸中心", "tax_number": "91310000MA005678XX"}
      ],
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

**请求体示例**
```json
{
  "name": "张三丰",
  "phone": "13800138001",
  "password": "newpass123"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "is_service_person": false,
    "name": "张三丰",
    "phone": "13800138001",
    "id_card": "110101199001011234",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

### 5. 删除人员

**请求**
```
DELETE /api/people/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "人员删除成功"
  }
}
```

### 6. 获取人员关联的企业

**请求**
```
GET /api/people/:id/customers
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "person_id": 1,
    "person_name": "张三",
    "customers": {
      "as_representative": [
        {
          "id": 1,
          "name": "某某科技有限公司",
          "tax_number": "91110000MA001234XX",
          "type": "有限公司"
        }
      ],
      "as_investor": [],
      "as_service": []
    }
  }
}
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
        "type": "有限公司",
        "representative_id": 1,
        "investor_ids": "2,3",
        "service_person_ids": "5,6",
        "agreement_ids": "1,3",
        "registered_capital": 1000000,
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
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 公司名称 |
| phone | string | 否 | 联系电话 |
| address | string | 否 | 地址 |
| tax_number | string | 是 | 税号（唯一） |
| type | string | 是 | 客户类型 |
| representative_id | uint | 否 | 法定代表人ID |
| investors | array | 否 | 投资人数组 |
| service_person_ids | string | 否 | 服务人员ID（逗号分隔） |
| agreement_ids | string | 否 | 代理协议ID（逗号分隔） |
| registered_capital | float64 | 否 | 注册资本 |

**请求体示例**
```json
{
  "name": "某某科技有限公司",
  "phone": "13800138000",
  "address": "北京市朝阳区xxx",
  "tax_number": "91110000xxxxxxxx",
  "type": "有限公司",
  "representative_id": 1,
  "investors": [
    {"person_id": 2, "share_ratio": 51},
    {"person_id": 3, "share_ratio": 49}
  ],
  "service_person_ids": "5,6",
  "agreement_ids": "1,3",
  "registered_capital": 1000000
}
```

**客户类型 (type)**
- `有限公司` - 有限责任公司
- `个人独资企业` - 个人独资企业
- `合伙企业` - 合伙企业
- `个体工商户` - 个体工商户

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "某某科技有限公司",
    "tax_number": "91110000xxxxxxxx",
    "type": "有限公司",
    "created_at": "2024-01-01T00:00:00Z"
  }
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
    "phone": "13800138000",
    "address": "北京市朝阳区xxx",
    "tax_number": "91110000xxxxxxxx",
    "type": "有限公司",
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
    "agreements_list": [
      {
        "id": 1,
        "agreement_number": "AG2024001",
        "start_date": "2024-01-01",
        "end_date": "2024-12-31",
        "fee_type": "月度",
        "amount": 500
      }
    ],
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

**请求体示例**
```json
{
  "name": "某某科技集团有限公司",
  "phone": "13800138001",
  "address": "北京市朝阳区新地址",
  "registered_capital": 2000000
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "某某科技集团有限公司",
    "phone": "13800138001",
    "address": "北京市朝阳区新地址",
    "tax_number": "91110000xxxxxxxx",
    "type": "有限公司",
    "registered_capital": 2000000,
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

### 5. 删除客户

**请求**
```
DELETE /api/customers/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "客户删除成功"
  }
}
```

### 6. 获取客户的任务列表

**请求**
```
GET /api/customers/:id/tasks
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "customer_id": 1,
    "customer_name": "某某科技有限公司",
    "tasks": [
      {
        "id": 1,
        "title": "月度报税",
        "description": "完成1月份税务申报",
        "status": "待处理",
        "due_date": "2024-02-15T00:00:00Z",
        "created_at": "2024-01-15T00:00:00Z"
      },
      {
        "id": 2,
        "title": "年检申报",
        "description": "完成年度企业年报",
        "status": "进行中",
        "due_date": "2024-06-30T00:00:00Z",
        "created_at": "2024-01-10T00:00:00Z"
      }
    ]
  }
}
```

### 7. 获取客户的收款记录

**请求**
```
GET /api/customers/:id/payments
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "customer_id": 1,
    "customer_name": "某某科技有限公司",
    "payments": [
      {
        "id": 1,
        "amount": 500,
        "payment_date": "2024-01-15T00:00:00Z",
        "payment_method": "转账",
        "period": "2024-01",
        "remark": "1月服务费"
      },
      {
        "id": 2,
        "amount": 500,
        "payment_date": "2024-02-15T00:00:00Z",
        "payment_method": "转账",
        "period": "2024-02",
        "remark": "2月服务费"
      }
    ]
  }
}
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
| status | string | 否 | 任务状态 (待处理/进行中/已完成) |
| customer_id | int | 否 | 按客户ID筛选 |

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
        "customer_id": 1,
        "title": "月度报税",
        "description": "完成1月份税务申报",
        "status": "待处理",
        "due_date": "2024-02-15T00:00:00Z",
        "completed_at": null,
        "created_at": "2024-01-15T00:00:00Z",
        "updated_at": "2024-01-15T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司",
          "tax_number": "91110000MA001234XX"
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
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_id | uint | 是 | 关联客户ID |
| title | string | 是 | 任务标题 |
| description | string | 否 | 任务描述 |
| status | string | 否 | 任务状态 |
| due_date | string | 否 | 截止日期 (ISO 8601格式) |

**请求体示例**
```json
{
  "customer_id": 1,
  "title": "月度报税",
  "description": "完成1月份税务申报工作",
  "status": "待处理",
  "due_date": "2024-02-15T00:00:00Z"
}
```

**任务状态 (status)**
- `待处理` - 待处理的任务
- `进行中` - 正在进行的任务
- `已完成` - 已完成的任务

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "title": "月度报税",
    "description": "完成1月份税务申报工作",
    "status": "待处理",
    "due_date": "2024-02-15T00:00:00Z",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

### 3. 获取任务详情

**请求**
```
GET /api/tasks/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "title": "月度报税",
    "description": "完成1月份税务申报工作",
    "status": "待处理",
    "due_date": "2024-02-15T00:00:00Z",
    "completed_at": null,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司",
      "tax_number": "91110000MA001234XX",
      "phone": "13800138000"
    }
  }
}
```

### 4. 更新任务

**请求**
```
PUT /api/tasks/:id
Content-Type: application/json
```

**请求体示例**
```json
{
  "title": "月度报税（加急）",
  "description": "完成1月份税务申报工作，需要加急处理",
  "status": "进行中",
  "due_date": "2024-02-10T00:00:00Z"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "title": "月度报税（加急）",
    "description": "完成1月份税务申报工作，需要加急处理",
    "status": "进行中",
    "due_date": "2024-02-10T00:00:00Z",
    "updated_at": "2024-01-16T14:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司"
    }
  }
}
```

> **注意**: 当任务状态从非"已完成"变为"已完成"时，系统会自动设置 `completed_at` 时间戳。

### 5. 删除任务

**请求**
```
DELETE /api/tasks/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "任务删除成功"
  }
}
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
| status | string | 否 | 协议状态 (有效/已过期/已取消) |
| customer_id | int | 否 | 按客户ID筛选 |

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
        "agreement_number": "AG2024001",
        "start_date": "2024-01-01T00:00:00Z",
        "end_date": "2024-12-31T00:00:00Z",
        "fee_type": "月度",
        "amount": 500,
        "status": "有效",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司",
          "tax_number": "91110000MA001234XX"
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
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_id | uint | 是 | 关联客户ID |
| agreement_number | string | 是 | 协议编号（唯一） |
| start_date | string | 是 | 协议开始日期 (ISO 8601格式) |
| end_date | string | 是 | 协议结束日期 (ISO 8601格式) |
| fee_type | string | 是 | 收费类型 |
| amount | float64 | 是 | 服务费金额 |
| status | string | 否 | 协议状态 |

**请求体示例**
```json
{
  "customer_id": 1,
  "agreement_number": "AG2024001",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z",
  "fee_type": "月度",
  "amount": 500,
  "status": "有效"
}
```

**收费类型 (fee_type)**
- `月度` - 按月收费
- `季度` - 按季度收费
- `年度` - 按年收费

**协议状态 (status)**
- `有效` - 协议有效
- `已过期` - 协议已过期
- `已取消` - 协议已取消

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_number": "AG2024001",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T00:00:00Z",
    "fee_type": "月度",
    "amount": 500,
    "status": "有效",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 3. 获取协议详情

**请求**
```
GET /api/agreements/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_number": "AG2024001",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T00:00:00Z",
    "fee_type": "月度",
    "amount": 500,
    "status": "有效",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司",
      "tax_number": "91110000MA001234XX"
    },
    "payments": [
      {
        "id": 1,
        "amount": 500,
        "payment_date": "2024-01-15T00:00:00Z",
        "payment_method": "转账"
      }
    ]
  }
}
```

### 4. 更新协议

**请求**
```
PUT /api/agreements/:id
Content-Type: application/json
```

**请求体示例**
```json
{
  "end_date": "2025-12-31T00:00:00Z",
  "amount": 600,
  "status": "有效"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_number": "AG2024001",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2025-12-31T00:00:00Z",
    "fee_type": "月度",
    "amount": 600,
    "status": "有效",
    "updated_at": "2024-06-01T14:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司"
    }
  }
}
```

### 5. 删除协议

**请求**
```
DELETE /api/agreements/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "协议删除成功"
  }
}
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
| customer_id | int | 否 | 按客户ID筛选 |
| start_date | string | 否 | 起始日期 (格式: 2006-01-02) |
| end_date | string | 否 | 结束日期 (格式: 2006-01-02) |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 20,
    "items": [
      {
        "id": 1,
        "customer_id": 1,
        "agreement_id": 1,
        "amount": 500,
        "payment_date": "2024-01-15T00:00:00Z",
        "payment_method": "转账",
        "period": "2024-01",
        "remark": "1月服务费",
        "created_at": "2024-01-15T00:00:00Z",
        "updated_at": "2024-01-15T00:00:00Z",
        "customer": {
          "id": 1,
          "name": "某某科技有限公司",
          "tax_number": "91110000MA001234XX"
        },
        "agreement": {
          "id": 1,
          "agreement_number": "AG2024001",
          "fee_type": "月度"
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
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_id | uint | 是 | 关联客户ID |
| agreement_id | uint | 否 | 关联协议ID |
| amount | float64 | 是 | 收款金额 |
| payment_date | string | 是 | 收款日期 (ISO 8601格式) |
| payment_method | string | 否 | 收款方式 |
| period | string | 否 | 费用所属期间 (如: 2024-01) |
| remark | string | 否 | 备注 |

**请求体示例**
```json
{
  "customer_id": 1,
  "agreement_id": 1,
  "amount": 500,
  "payment_date": "2024-01-15T00:00:00Z",
  "payment_method": "转账",
  "period": "2024-01",
  "remark": "1月服务费"
}
```

**收款方式 (payment_method)**
- `转账` - 银行转账
- `现金` - 现金支付
- `支票` - 支票支付
- `其他` - 其他方式

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_id": 1,
    "amount": 500,
    "payment_date": "2024-01-15T00:00:00Z",
    "payment_method": "转账",
    "period": "2024-01",
    "remark": "1月服务费",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

### 3. 获取收款记录详情

**请求**
```
GET /api/payments/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_id": 1,
    "amount": 500,
    "payment_date": "2024-01-15T00:00:00Z",
    "payment_method": "转账",
    "period": "2024-01",
    "remark": "1月服务费",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司",
      "tax_number": "91110000MA001234XX",
      "phone": "13800138000"
    },
    "agreement": {
      "id": 1,
      "agreement_number": "AG2024001",
      "fee_type": "月度",
      "amount": 500
    }
  }
}
```

### 4. 更新收款记录

**请求**
```
PUT /api/payments/:id
Content-Type: application/json
```

**请求体示例**
```json
{
  "amount": 600,
  "payment_date": "2024-01-16T00:00:00Z",
  "remark": "1月服务费（调整后）"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "agreement_id": 1,
    "amount": 600,
    "payment_date": "2024-01-16T00:00:00Z",
    "payment_method": "转账",
    "period": "2024-01",
    "remark": "1月服务费（调整后）",
    "updated_at": "2024-01-16T14:00:00Z",
    "customer": {
      "id": 1,
      "name": "某某科技有限公司"
    },
    "agreement": {
      "id": 1,
      "agreement_number": "AG2024001"
    }
  }
}
```

### 5. 删除收款记录

**请求**
```
DELETE /api/payments/:id
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "收款记录删除成功"
  }
}
```

---

## 统计分析 API

### 1. 首页概览统计

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
    "customer_count": 50,
    "pending_task_count": 15,
    "active_agreement_count": 45,
    "monthly_payment": 25000,
    "yearly_payment": 180000
  }
}
```

**响应字段说明**
| 字段 | 类型 | 说明 |
|------|------|------|
| customer_count | int64 | 客户总数 |
| pending_task_count | int64 | 待办任务数（未完成的任务） |
| active_agreement_count | int64 | 有效协议数 |
| monthly_payment | float64 | 本月收款总额 |
| yearly_payment | float64 | 本年收款总额 |

### 2. 任务统计

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
    "pending": 10,
    "in_progress": 5,
    "completed": 100
  }
}
```

**响应字段说明**
| 字段 | 类型 | 说明 |
|------|------|------|
| pending | int64 | 待处理任务数 |
| in_progress | int64 | 进行中任务数 |
| completed | int64 | 已完成任务数 |

### 3. 收款统计

**请求**
```
GET /api/statistics/payments
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 起始日期 (格式: 2006-01-02) |
| end_date | string | 否 | 结束日期 (格式: 2006-01-02) |

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_amount": 50000,
    "count": 50
  }
}
```

**响应字段说明**
| 字段 | 类型 | 说明 |
|------|------|------|
| total_amount | float64 | 收款总金额 |
| count | int64 | 收款记录数 |

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
| is_service_person | boolean | 是否为服务人员 |
| name | string | 姓名 |
| phone | string | 电话 |
| id_card | string | 身份证号（唯一） |
| password | string | 登录密码 |
| representative_customer_ids | string | 担任法人的企业ID（逗号分隔） |
| investor_customer_ids | string | 持股的企业ID（逗号分隔） |
| service_customer_ids | string | 服务的企业ID（逗号分隔） |

**人员角色说明：**
- **服务人员**: `is_service_person = true` 的人员
- **法定代表人**: `representative_customer_ids` 不为空的人员
- **投资人**: `investor_customer_ids` 不为空的人员
- 一个人可以同时担任多个角色（既是投资人又是服务人员等）

### Customer (客户)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 公司名称 |
| phone | string | 联系电话 |
| address | string | 地址 |
| tax_number | string | 税号 |
| type | string | 客户类型（有限公司/个人独资企业/合伙企业/个体工商户） |
| representative_id | uint | 法定代表人ID |
| investor_ids | string | 投资人ID（逗号分隔） |
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

### Task (任务)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| customer_id | uint | 关联客户ID |
| title | string | 任务标题 |
| description | string | 任务描述 |
| status | string | 任务状态（待处理/进行中/已完成） |
| due_date | timestamp | 截止日期 |
| completed_at | timestamp | 完成日期 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| customer | Customer | 关联客户信息 |

### Agreement (协议)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| customer_id | uint | 关联客户ID |
| agreement_number | string | 协议编号（唯一） |
| start_date | date | 协议开始日期 |
| end_date | date | 协议结束日期 |
| fee_type | string | 收费类型（月度/季度/年度） |
| amount | float64 | 服务费金额 |
| status | string | 协议状态（有效/已过期/已取消） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| customer | Customer | 关联客户信息 |
| payments | Payment[] | 关联收款记录 |

### Payment (收款记录)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| customer_id | uint | 关联客户ID |
| agreement_id | uint | 关联协议ID（可选） |
| amount | float64 | 收款金额 |
| payment_date | date | 收款日期 |
| payment_method | string | 收款方式（转账/现金/支票/其他） |
| period | string | 费用所属期间（如: 2024-01） |
| remark | string | 备注 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| customer | Customer | 关联客户信息 |
| agreement | Agreement | 关联协议信息 |
