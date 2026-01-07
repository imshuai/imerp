package controllers

import (
	"erp/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetPeople 测试获取人员列表
func TestGetPeople(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试人员
	createTestPerson(t, db, "张三", "13800138001", "500000199001011234", false)
	createTestPerson(t, db, "李四", "13800138002", "500000199002021234", true)

	w := makeRequest(router, "GET", "/api/people", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])

	items := data["items"].([]interface{})
	assert.Equal(t, 2, len(items))
}

// TestGetPersonByID 测试获取单个人员
func TestGetPersonByID(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试人员
	person := createTestPerson(t, db, "王五", "13800138003", "500000199003031234", false)

	w := makeRequest(router, "GET", fmt.Sprintf("/api/people/%d", person.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	personData := data["person"].(map[string]interface{})
	assert.Equal(t, "王五", personData["name"])
	assert.Equal(t, "13800138003", personData["phone"])
	assert.Equal(t, "500000199003031234", personData["id_card"])
}

// TestCreatePerson 测试创建人员
func TestCreatePerson(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	reqBody := map[string]interface{}{
		"name":            "测试人员",
		"phone":           "13900139000",
		"id_card":         "500000199001011235",
		"is_service_person": false,
	}

	w := makeRequest(router, "POST", "/api/people", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "测试人员", data["name"])
	assert.Equal(t, "13900139000", data["phone"])
	assert.Equal(t, "500000199001011235", data["id_card"])

	// 验证数据库
	var person models.Person
	db.First(&person, uint(data["id"].(float64)))
	assert.Equal(t, "测试人员", person.Name)
}

// TestCreateServicePerson 测试创建服务人员
func TestCreateServicePerson(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	reqBody := map[string]interface{}{
		"name":              "服务人员A",
		"phone":             "13900139001",
		"id_card":           "500000199001011236",
		"is_service_person": true,
	}

	w := makeRequest(router, "POST", "/api/people", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "服务人员A", data["name"])
	assert.Equal(t, true, data["is_service_person"])
	assert.Equal(t, "服务人员", data["type"])

	// 验证AdminUser已创建
	var adminUser models.AdminUser
	result := db.Where("person_id = ?", uint(data["id"].(float64))).First(&adminUser)
	assert.NoError(t, result.Error)
	assert.Equal(t, "服务人员A", adminUser.Username)
	assert.Equal(t, "service_person", adminUser.Role)
}

// TestUpdatePerson 测试更新人员
func TestUpdatePerson(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试人员
	person := createTestPerson(t, db, "原名", "13800138000", "500000199001011234", false)

	updateBody := map[string]interface{}{
		"name":              "新名",
		"phone":             "13900139000",
		"is_service_person": true,
	}

	w := makeRequest(router, "PUT", fmt.Sprintf("/api/people/%d", person.ID), updateBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "新名", data["name"])
	assert.Equal(t, "13900139000", data["phone"])
	assert.Equal(t, true, data["is_service_person"])

	// 验证数据库更新
	var updatedPerson models.Person
	db.First(&updatedPerson, person.ID)
	assert.Equal(t, "新名", updatedPerson.Name)
}

// TestDeletePerson 测试删除人员
func TestDeletePerson(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试人员
	person := createTestPerson(t, db, "待删除人员", "13800138000", "500000199001011234", false)

	w := makeRequest(router, "DELETE", fmt.Sprintf("/api/people/%d", person.ID), nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	// 验证数据库中已删除
	var deletedPerson models.Person
	result := db.First(&deletedPerson, person.ID)
	assert.Error(t, result.Error)
}

// TestSearchPeople 测试搜索人员
func TestSearchPeople(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建测试人员
	createTestPerson(t, db, "张三", "13800138001", "500000199001011234", false)
	createTestPerson(t, db, "张四", "13800138002", "500000199002021234", true)
	createTestPerson(t, db, "李三", "13800138003", "500000199003031234", false)

	// 测试关键词搜索
	w := makeRequest(router, "GET", "/api/people?keyword=张", nil)
	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])

	// 测试按服务人员筛选
	w2 := makeRequest(router, "GET", "/api/people?is_service_person=true", nil)
	assert.Equal(t, 200, w.Code)

	response2 := parseResponse(t, w2)
	data2 := response2["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data2["total"])
}

// TestPersonWithNewField 测试Person新增的TaxAgentCustomerIDs字段
func TestPersonWithNewField(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	reqBody := map[string]interface{}{
		"name":                     "办税人员",
		"phone":                    "13800138000",
		"id_card":                  "500000199001011234",
		"is_service_person":        false,
		"tax_agent_customer_ids":   "",
	}

	w := makeRequest(router, "POST", "/api/people", reqBody)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "", data["tax_agent_customer_ids"])

	// 验证数据库
	var person models.Person
	db.First(&person, uint(data["id"].(float64)))
	assert.Equal(t, "办税人员", person.Name)
}

// TestPasswordPolicy 测试密码策略 - 只有服务人员才有密码
func TestPasswordPolicy(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	// 创建普通人员
	normalPersonBody := map[string]interface{}{
		"name":              "普通人员",
		"phone":             "13800138001",
		"id_card":           "500000199001011234",
		"is_service_person": false,
	}

	w1 := makeRequest(router, "POST", "/api/people", normalPersonBody)
	assert.Equal(t, 200, w1.Code)

	response1 := parseResponse(t, w1)
	normalPersonID := uint(response1["data"].(map[string]interface{})["id"].(float64))

	// 验证普通人员没有AdminUser
	var adminUser models.AdminUser
	result := db.Where("person_id = ?", normalPersonID).First(&adminUser)
	assert.Error(t, result.Error) // 应该找不到

	// 创建服务人员
	servicePersonBody := map[string]interface{}{
		"name":              "服务人员",
		"phone":             "13800138002",
		"id_card":           "500000199002021234",
		"is_service_person": true,
	}

	w2 := makeRequest(router, "POST", "/api/people", servicePersonBody)
	assert.Equal(t, 200, w2.Code)

	response2 := parseResponse(t, w2)
	servicePersonID := uint(response2["data"].(map[string]interface{})["id"].(float64))

	// 验证服务人员有AdminUser和默认密码
	var serviceAdminUser models.AdminUser
	result2 := db.Where("person_id = ?", servicePersonID).First(&serviceAdminUser)
	assert.NoError(t, result2.Error)
	assert.Equal(t, "服务人员", serviceAdminUser.Username)
	assert.Equal(t, "service_person", serviceAdminUser.Role)
	assert.Equal(t, true, serviceAdminUser.MustChangePassword) // 默认需要修改密码
}
