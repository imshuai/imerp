package controllers

import (
	"erp/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateCustomer 测试创建客户
func TestCreateCustomer(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试数据
	reqBody := map[string]interface{}{
		"name":              "测试公司",
		"tax_number":        "91500000MA123456X",
		"type":              "有限公司",
		"phone":             "13800138000",
		"address":           "重庆市渝中区",
		"registered_capital": 1000000.00,
		"taxpayer_type":     "一般纳税人",
		"credit_rating":     "A",
		"social_security_number": "123456",
		"business_scope":    "软件开发",
	}

	w := makeRequest(router, "POST", "/api/customers", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "测试公司", data["name"])
	assert.Equal(t, "91500000MA123456X", data["tax_number"])
	assert.Equal(t, "有限公司", data["type"])
	assert.Equal(t, "A", data["credit_rating"])

	// 验证数据库中的记录
	var customer models.Customer
	result := db.First(&customer, uint(data["id"].(float64)))
	assert.NoError(t, result.Error)
	assert.Equal(t, "测试公司", customer.Name)
}

// TestGetCustomers 测试获取客户列表
func TestGetCustomers(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试客户
	createTestCustomer(t, db, "公司A", "TAX001", models.CustomerTypeLimitedCompany)
	createTestCustomer(t, db, "公司B", "TAX002", models.CustomerTypeIndividualBusiness)

	w := makeRequest(router, "GET", "/api/customers", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])

	items := data["items"].([]interface{})
	assert.Equal(t, 2, len(items))
}

// TestGetCustomerByID 测试获取单个客户
func TestGetCustomerByID(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试客户
	customer := createTestCustomer(t, db, "测试公司", "TAX001", models.CustomerTypeLimitedCompany)

	w := makeRequest(router, "GET", fmt.Sprintf("/api/customers/%d", customer.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "测试公司", data["name"])
	assert.Equal(t, "TAX001", data["tax_number"])
}

// TestGetCustomerNotFound 测试获取不存在的客户
func TestGetCustomerNotFound(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	w := makeRequest(router, "GET", "/api/customers/999", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.NotEqual(t, float64(0), response["code"])
}

// TestUpdateCustomer 测试更新客户
func TestUpdateCustomer(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试客户
	customer := createTestCustomer(t, db, "原公司名", "TAX001", models.CustomerTypeLimitedCompany)

	updateBody := map[string]interface{}{
		"name":          "新公司名",
		"phone":         "13900139000",
		"credit_rating": "B",
	}

	w := makeRequest(router, "PUT", fmt.Sprintf("/api/customers/%d", customer.ID), updateBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "新公司名", data["name"])
	assert.Equal(t, "13900139000", data["phone"])
	assert.Equal(t, "B", data["credit_rating"])

	// 验证数据库更新
	var updatedCustomer models.Customer
	db.First(&updatedCustomer, customer.ID)
	assert.Equal(t, "新公司名", updatedCustomer.Name)
}

// TestDeleteCustomer 测试删除客户
func TestDeleteCustomer(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试客户
	customer := createTestCustomer(t, db, "待删除公司", "TAX001", models.CustomerTypeLimitedCompany)

	w := makeRequest(router, "DELETE", fmt.Sprintf("/api/customers/%d", customer.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	// 验证数据库中已删除
	var deletedCustomer models.Customer
	result := db.First(&deletedCustomer, customer.ID)
	assert.Error(t, result.Error)
}

// TestCustomerWithNewFields 测试新增字段
func TestCustomerWithNewFields(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	reqBody := map[string]interface{}{
		"name":                   "完整信息公司",
		"tax_number":             "TAX001",
		"type":                   "有限公司",
		"credit_rating":          "A",
		"social_security_number": "500123456789",
		"yukuai_ban_password":    "test123",
		"business_scope":         "软件开发、技术服务",
		"tax_agent_ids":          "",
	}

	w := makeRequest(router, "POST", "/api/customers", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "A", data["credit_rating"])
	assert.Equal(t, "500123456789", data["social_security_number"])
	assert.Equal(t, "test123", data["yukuai_ban_password"])
	assert.Equal(t, "软件开发、技术服务", data["business_scope"])
}

// TestCustomerWithBankAccount 测试客户关联银行账户
func TestCustomerWithBankAccount(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建客户
	customer := createTestCustomer(t, db, "银行账户测试公司", "TAX001", models.CustomerTypeLimitedCompany)

	// 创建银行账户
	bankAccount := models.BankAccount{
		CustomerID:    customer.ID,
		BankName:      "中国银行",
		AccountNumber: "1234567890",
		BankCode:      "104100000004",
		ContactPhone:  "023-12345678",
		AccountType:   models.AccountTypeBasic,
	}
	db.Create(&bankAccount)

	// 获取客户详情
	w := makeRequest(router, "GET", fmt.Sprintf("/api/customers/%d", customer.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	bankAccounts, ok := data["bank_accounts"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(bankAccounts))

	firstAccount := bankAccounts[0].(map[string]interface{})
	assert.Equal(t, "中国银行", firstAccount["bank_name"])
	assert.Equal(t, "1234567890", firstAccount["account_number"])
}

// TestCustomerWithInvestor 测试客户关联投资人
func TestCustomerWithInvestor(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建客户
	customer := createTestCustomer(t, db, "投资人测试公司", "TAX001", models.CustomerTypeLimitedCompany)

	// 创建投资人（人员）
	investor := createTestPerson(t, db, "投资人A", "13800000001", "500000199001011234", false)

	// 创建投资人关联记录
	investorRelation := models.CustomerInvestor{
		CustomerID: customer.ID,
		PersonID:   investor.ID,
		ShareRatio: 60.5,
	}
	db.Create(&investorRelation)

	// 获取客户详情
	w := makeRequest(router, "GET", fmt.Sprintf("/api/customers/%d", customer.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	investorRelations, ok := data["investor_relations"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(investorRelations))

	firstRelation := investorRelations[0].(map[string]interface{})
	assert.Equal(t, float64(60.5), firstRelation["share_ratio"])

	// 验证关联的Person信息
	person, ok := firstRelation["person"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "投资人A", person["name"])
}

// TestCustomerWithTaxAgent 测试客户关联办税人
func TestCustomerWithTaxAgent(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建办税人
	taxAgent1 := createTestPerson(t, db, "办税人A", "13800000001", "500000199001011234", false)
	taxAgent2 := createTestPerson(t, db, "办税人B", "13800000002", "500000199002021234", false)

	reqBody := map[string]interface{}{
		"name":         "办税人测试公司",
		"tax_number":   "TAX001",
		"type":         "有限公司",
		"tax_agent_ids": fmt.Sprintf("%d,%d", taxAgent1.ID, taxAgent2.ID),
	}

	w := makeRequest(router, "POST", "/api/customers", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	taxAgents, ok := data["tax_agents"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(taxAgents))

	// 验证办税人信息
	firstAgent := taxAgents[0].(map[string]interface{})
	assert.True(t, firstAgent["name"] == "办税人A" || firstAgent["name"] == "办税人B")
}

// TestSearchCustomers 测试搜索客户
func TestSearchCustomers(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试客户
	createTestCustomer(t, db, "重庆科技公司", "TAX001", models.CustomerTypeLimitedCompany)
	createTestCustomer(t, db, "成都贸易公司", "TAX002", models.CustomerTypeLimitedCompany)
	createTestCustomer(t, db, "北京服务公司", "TAX003", models.CustomerTypeLimitedCompany)

	// 测试关键词搜索
	w := makeRequest(router, "GET", "/api/customers?keyword=重庆", nil)
	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])

	items := data["items"].([]interface{})
	customer := items[0].(map[string]interface{})
	assert.Contains(t, customer["name"], "重庆")
}
