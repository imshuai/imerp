package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(
		&Person{},
		&Customer{},
		&CustomerInvestor{},
		&BankAccount{},
		&Task{},
		&Agreement{},
		&Payment{},
		&AdminUser{},
		&AuditLog{},
	)
	assert.NoError(t, err)

	return db
}

// TestBankAccountModel 测试BankAccount模型
func TestBankAccountModel(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试客户
	customer := Customer{
		Name:     "测试公司",
		TaxNumber: "TAX001",
		Type:     CustomerTypeLimitedCompany,
	}
	db.Create(&customer)

	// 创建银行账户
	bankAccount := BankAccount{
		CustomerID:    customer.ID,
		BankName:      "中国银行",
		AccountNumber: "1234567890",
		BankCode:      "104100000004",
		ContactPhone:  "023-12345678",
		AccountType:   AccountTypeBasic,
	}

	err := db.Create(&bankAccount).Error
	assert.NoError(t, err)
	assert.NotZero(t, bankAccount.ID)

	// 查询验证
	var found BankAccount
	db.First(&found, bankAccount.ID)
	assert.Equal(t, "中国银行", found.BankName)
	assert.Equal(t, "1234567890", found.AccountNumber)
	assert.Equal(t, AccountTypeBasic, found.AccountType)
}

// TestCustomerInvestorModel 测试CustomerInvestor模型
func TestCustomerInvestorModel(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试客户
	customer := Customer{
		Name:     "测试公司",
		TaxNumber: "TAX001",
		Type:     CustomerTypeLimitedCompany,
	}
	db.Create(&customer)

	// 创建测试人员
	person := Person{
		Name:      "投资人",
		Phone:     "13800138000",
		IDCard:    "500000199001011234",
	}
	db.Create(&person)

	// 创建投资人关联
	customerInvestor := CustomerInvestor{
		CustomerID: customer.ID,
		PersonID:   person.ID,
		ShareRatio: 60.5,
	}

	err := db.Create(&customerInvestor).Error
	assert.NoError(t, err)
	assert.NotZero(t, customerInvestor.ID)

	// 查询验证
	var found CustomerInvestor
	db.First(&found, customerInvestor.ID)
	assert.Equal(t, customer.ID, found.CustomerID)
	assert.Equal(t, person.ID, found.PersonID)
	assert.Equal(t, 60.5, found.ShareRatio)
}

// TestCustomerInvestorInvestmentRecords 测试出资记录JSON处理
func TestCustomerInvestorInvestmentRecords(t *testing.T) {
	db := setupTestDB(t)

	customer := Customer{Name: "测试公司", TaxNumber: "TAX001", Type: CustomerTypeLimitedCompany}
	db.Create(&customer)

	person := Person{Name: "投资人", Phone: "13800138000", IDCard: "500000199001011234"}
	db.Create(&person)

	customerInvestor := CustomerInvestor{
		CustomerID: customer.ID,
		PersonID:   person.ID,
		ShareRatio: 50.0,
	}

	// 设置出资记录
	records := []InvestmentRecord{
		{Date: "2020-01-01", Amount: 500000},
		{Date: "2021-01-01", Amount: 300000},
	}
	err := customerInvestor.SetInvestmentRecords(records)
	assert.NoError(t, err)

	db.Create(&customerInvestor)

	// 查询并验证
	var found CustomerInvestor
	db.First(&found, customerInvestor.ID)

	retrievedRecords := found.GetInvestmentRecords()
	assert.Equal(t, 2, len(retrievedRecords))
	assert.Equal(t, "2020-01-01", retrievedRecords[0].Date)
	assert.Equal(t, 500000.0, retrievedRecords[0].Amount)
	assert.Equal(t, "2021-01-01", retrievedRecords[1].Date)
	assert.Equal(t, 300000.0, retrievedRecords[1].Amount)
}

// TestCustomerModelWithNewFields 测试Customer模型新增字段
func TestCustomerModelWithNewFields(t *testing.T) {
	db := setupTestDB(t)

	customer := Customer{
		Name:                 "完整信息公司",
		TaxNumber:            "TAX001",
		Type:                 CustomerTypeLimitedCompany,
		CreditRating:         CreditRatingA,
		SocialSecurityNumber: "500123456789",
		YukuaiBanPassword:    "test123",
		BusinessScope:        "软件开发、技术服务",
		TaxAgentIDs:          "1,2,3",
	}

	err := db.Create(&customer).Error
	assert.NoError(t, err)

	// 查询验证
	var found Customer
	db.First(&found, customer.ID)
	assert.Equal(t, CreditRatingA, found.CreditRating)
	assert.Equal(t, "500123456789", found.SocialSecurityNumber)
	assert.Equal(t, "test123", found.YukuaiBanPassword)
	assert.Equal(t, "软件开发、技术服务", found.BusinessScope)
	assert.Equal(t, "1,2,3", found.TaxAgentIDs)
}

// TestPersonModelWithNewField 测试Person模型新增字段
func TestPersonModelWithNewField(t *testing.T) {
	db := setupTestDB(t)

	person := Person{
		Name:                "测试人员",
		Phone:               "13800138000",
		IDCard:              "500000199001019999", // Unique ID card
		IsServicePerson:     false,
		TaxAgentCustomerIDs: "1,2",
	}

	err := db.Create(&person).Error
	assert.NoError(t, err)

	// 查询验证
	var found Person
	db.First(&found, person.ID)
	assert.Equal(t, "1,2", found.TaxAgentCustomerIDs)
}

// TestJSONDate 测试JSONDate类型
func TestJSONDate(t *testing.T) {
	// 测试空值
	var emptyDate JSONDate
	err := emptyDate.UnmarshalJSON([]byte(`""`))
	assert.NoError(t, err)
	assert.Equal(t, JSONDate(""), emptyDate)

	// 测试有效日期
	var validDate JSONDate
	err = validDate.UnmarshalJSON([]byte(`"2024-01-15"`))
	assert.NoError(t, err)
	assert.Equal(t, JSONDate("2024-01-15"), validDate)

	// 测试序列化
	jsonBytes, err := validDate.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2024-01-15"`), jsonBytes)

	// 测试无效日期
	var invalidDate JSONDate
	err = invalidDate.UnmarshalJSON([]byte(`"2024-13-45"`))
	assert.Error(t, err)

	// 测试ToTime
	timeValue := validDate.ToTime()
	assert.NotNil(t, timeValue)
	assert.Equal(t, 2024, timeValue.Year())
	assert.Equal(t, 1, int(timeValue.Month()))
	assert.Equal(t, 15, timeValue.Day())
}

// TestCustomerTypeValidation 测试客户类型验证
func TestCustomerTypeValidation(t *testing.T) {
	validTypes := []CustomerType{
		CustomerTypeLimitedCompany,
		CustomerTypeSoleProprietorship,
		CustomerTypePartnership,
		CustomerTypeIndividualBusiness,
	}

	for _, ct := range validTypes {
		assert.NotEqual(t, "", string(ct))
	}
}

// TestCreditRatingValidation 测试信用等级验证
func TestCreditRatingValidation(t *testing.T) {
	validRatings := []CreditRating{
		CreditRatingA,
		CreditRatingB,
		CreditRatingC,
		CreditRatingD,
		CreditRatingM,
	}

	for _, cr := range validRatings {
		assert.NotEqual(t, "", string(cr))
		assert.Equal(t, 1, len(string(cr)))
	}
}

// TestAccountTypeValidation 测试账户类型验证
func TestAccountTypeValidation(t *testing.T) {
	validTypes := []AccountType{
		AccountTypeBasic,
		AccountTypeGeneral,
		AccountTypeTemporary,
	}

	expectedNames := []string{"基本户", "一般户", "临时户"}
	for i, at := range validTypes {
		assert.Equal(t, expectedNames[i], string(at))
	}
}

// TestCustomerWithInvestorRelations 测试客户与投资人关联
func TestCustomerWithInvestorRelations(t *testing.T) {
	db := setupTestDB(t)

	// 创建客户
	customer := Customer{
		Name:      "测试公司",
		TaxNumber: "TAX001",
		Type:      CustomerTypeLimitedCompany,
	}
	db.Create(&customer)

	// 创建多个投资人
	person1 := Person{Name: "投资人A", Phone: "13800138001", IDCard: "500000199001011234"}
	person2 := Person{Name: "投资人B", Phone: "13800138002", IDCard: "500000199002021234"}
	db.Create(&person1)
	db.Create(&person2)

	// 创建投资人关联
	investor1 := CustomerInvestor{
		CustomerID: customer.ID,
		PersonID:   person1.ID,
		ShareRatio: 60.0,
	}
	investor2 := CustomerInvestor{
		CustomerID: customer.ID,
		PersonID:   person2.ID,
		ShareRatio: 40.0,
	}
	db.Create(&investor1)
	db.Create(&investor2)

	// 查询验证
	var investorRelations []CustomerInvestor
	db.Where("customer_id = ?", customer.ID).Find(&investorRelations)

	assert.Equal(t, 2, len(investorRelations))
	assert.Equal(t, 60.0, investorRelations[0].ShareRatio)
	assert.Equal(t, 40.0, investorRelations[1].ShareRatio)
}

// TestCustomerWithBankAccounts 测试客户与银行账户关联
func TestCustomerWithBankAccounts(t *testing.T) {
	db := setupTestDB(t)

	// 创建客户
	customer := Customer{
		Name:      "测试公司",
		TaxNumber: "TAX001",
		Type:      CustomerTypeLimitedCompany,
	}
	db.Create(&customer)

	// 创建多个银行账户
	account1 := BankAccount{
		CustomerID:    customer.ID,
		BankName:      "中国银行",
		AccountNumber: "1234567890",
		AccountType:   AccountTypeBasic,
	}
	account2 := BankAccount{
		CustomerID:    customer.ID,
		BankName:      "工商银行",
		AccountNumber: "0987654321",
		AccountType:   AccountTypeGeneral,
	}
	db.Create(&account1)
	db.Create(&account2)

	// 查询验证
	var bankAccounts []BankAccount
	db.Where("customer_id = ?", customer.ID).Find(&bankAccounts)

	assert.Equal(t, 2, len(bankAccounts))
	assert.Equal(t, "中国银行", bankAccounts[0].BankName)
	assert.Equal(t, "工商银行", bankAccounts[1].BankName)
}
