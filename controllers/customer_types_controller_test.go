package controllers

import (
	"erp/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetCustomerTypes 测试获取客户类型选项
func TestGetCustomerTypes(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	w := makeRequest(router, "GET", "/api/customers/types", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(data), 4)

	// 验证包含所有客户类型
	types := make([]string, len(data))
	for i, v := range data {
		types[i] = v.(string)
	}

	expectedTypes := []string{"有限公司", "个人独资企业", "合伙企业", "个体工商户"}
	for _, expected := range expectedTypes {
		assert.Contains(t, types, expected)
	}
}

// TestGetCreditRatings 测试获取信用等级选项
func TestGetCreditRatings(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	w := makeRequest(router, "GET", "/api/customers/credit-ratings", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 5, len(data))

	// 验证包含所有信用等级
	ratings := make([]string, len(data))
	for i, v := range data {
		ratings[i] = v.(string)
	}

	expectedRatings := []string{"A", "B", "C", "D", "M"}
	for _, expected := range expectedRatings {
		assert.Contains(t, ratings, expected)
	}
}

// TestGetAccountTypes 测试获取账户类型选项
func TestGetAccountTypes(t *testing.T) {
	router, db := setupTestRouter(t)
	defer cleanupTestData(t, db)

	w := makeRequest(router, "GET", "/api/bank-accounts/types", nil)

	assert.Equal(t, 200, w.Code)

	response := parseResponse(t, w)
	assert.Equal(t, float64(0), response["code"])

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(data))

	// 验证包含所有账户类型
	types := make([]string, len(data))
	for i, v := range data {
		types[i] = v.(string)
	}

	expectedTypes := []string{"基本户", "一般户", "临时户"}
	for _, expected := range expectedTypes {
		assert.Contains(t, types, expected)
	}
}

// TestCustomerTypeModel 测试CustomerType模型函数
func TestCustomerTypeModel(t *testing.T) {
	// 测试GetCustomerTypeOptions
	types := models.GetCustomerTypeOptions()
	assert.Equal(t, 4, len(types))

	expectedTypes := []models.CustomerType{
		models.CustomerTypeLimitedCompany,
		models.CustomerTypeSoleProprietorship,
		models.CustomerTypePartnership,
		models.CustomerTypeIndividualBusiness,
	}

	for i, expected := range expectedTypes {
		assert.Equal(t, expected, types[i])
	}

	// 验证常量值
	assert.Equal(t, models.CustomerType("有限公司"), models.CustomerTypeLimitedCompany)
	assert.Equal(t, models.CustomerType("个人独资企业"), models.CustomerTypeSoleProprietorship)
	assert.Equal(t, models.CustomerType("合伙企业"), models.CustomerTypePartnership)
	assert.Equal(t, models.CustomerType("个体工商户"), models.CustomerTypeIndividualBusiness)
}

// TestCreditRatingModel 测试CreditRating模型函数
func TestCreditRatingModel(t *testing.T) {
	// 测试GetCreditRatingOptions
	ratings := models.GetCreditRatingOptions()
	assert.Equal(t, 5, len(ratings))

	expectedRatings := []models.CreditRating{
		models.CreditRatingA,
		models.CreditRatingB,
		models.CreditRatingC,
		models.CreditRatingD,
		models.CreditRatingM,
	}

	for i, expected := range expectedRatings {
		assert.Equal(t, expected, ratings[i])
	}

	// 验证常量值
	assert.Equal(t, models.CreditRating("A"), models.CreditRatingA)
	assert.Equal(t, models.CreditRating("B"), models.CreditRatingB)
	assert.Equal(t, models.CreditRating("C"), models.CreditRatingC)
	assert.Equal(t, models.CreditRating("D"), models.CreditRatingD)
	assert.Equal(t, models.CreditRating("M"), models.CreditRatingM)
}

// TestAccountTypeModel 测试AccountType模型函数
func TestAccountTypeModel(t *testing.T) {
	// 测试GetAccountTypeOptions
	types := models.GetAccountTypeOptions()
	assert.Equal(t, 3, len(types))

	expectedTypes := []models.AccountType{
		models.AccountTypeBasic,
		models.AccountTypeGeneral,
		models.AccountTypeTemporary,
	}

	for i, expected := range expectedTypes {
		assert.Equal(t, expected, types[i])
	}

	// 验证常量值
	assert.Equal(t, models.AccountType("基本户"), models.AccountTypeBasic)
	assert.Equal(t, models.AccountType("一般户"), models.AccountTypeGeneral)
	assert.Equal(t, models.AccountType("临时户"), models.AccountTypeTemporary)
}
