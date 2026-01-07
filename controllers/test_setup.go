package controllers

import (
	"bytes"
	"encoding/json"
	"erp/config"
	"erp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestRouter 设置测试路由
func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&models.Person{},
		&models.Customer{},
		&models.CustomerInvestor{},
		&models.BankAccount{},
		&models.Task{},
		&models.Agreement{},
		&models.Payment{},
		&models.AdminUser{},
		&models.AuditLog{},
	)
	assert.NoError(t, err)

	config.DB = db

	router := gin.New()
	router.Use(gin.Recovery())

	// 添加测试用的认证中间件（设置必要的上下文）
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Next()
	})

	api := router.Group("/api")
	{
		api.GET("/customers/types", GetCustomerTypes)
		api.GET("/customers/credit-ratings", GetCreditRatings)
		api.GET("/bank-accounts/types", GetAccountTypes)
		api.GET("/people", GetPeople)

		api.GET("/customers", GetCustomers)
		api.GET("/customers/:id", GetCustomer)
		api.POST("/customers", CreateCustomer)
		api.PUT("/customers/:id", UpdateCustomer)
		api.DELETE("/customers/:id", DeleteCustomer)
		api.GET("/customers/:id/tasks", GetCustomerTasks)
		api.GET("/customers/:id/payments", GetCustomerPayments)

		api.GET("/people/:id", GetPerson)
		api.POST("/people", CreatePerson)
		api.PUT("/people/:id", UpdatePerson)
		api.DELETE("/people/:id", DeletePerson)
	}

	return router, db
}

// makeRequest 发起HTTP请求
func makeRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// parseResponse 解析响应
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

// createTestPerson 创建测试人员
func createTestPerson(t *testing.T, db *gorm.DB, name, phone, idCard string, isServicePerson bool) models.Person {
	person := models.Person{
		Name:            name,
		Phone:           phone,
		IDCard:          idCard,
		IsServicePerson: isServicePerson,
	}
	err := db.Create(&person).Error
	assert.NoError(t, err)
	return person
}

// createTestCustomer 创建测试客户
func createTestCustomer(t *testing.T, db *gorm.DB, name, taxNumber string, customerType models.CustomerType) models.Customer {
	customer := models.Customer{
		Name:      name,
		TaxNumber: taxNumber,
		Type:      customerType,
		Phone:     "13800138000",
		Address:   "测试地址",
	}
	err := db.Create(&customer).Error
	assert.NoError(t, err)
	return customer
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T, db *gorm.DB) {
	db.Exec("DELETE FROM customer_investors")
	db.Exec("DELETE FROM bank_accounts")
	db.Exec("DELETE FROM tasks")
	db.Exec("DELETE FROM payments")
	db.Exec("DELETE FROM agreements")
	db.Exec("DELETE FROM customers")
	db.Exec("DELETE FROM admin_users")
	db.Exec("DELETE FROM audit_logs")
	db.Exec("DELETE FROM people")
}
