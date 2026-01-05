package controllers

import (
	"erp/config"
	"erp/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreatePerson 创建人员
func CreatePerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	if err := config.DB.Create(&person).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create person: "+err.Error())
		return
	}

	// 更新关联客户的ID字段
	updatePersonCustomerIDs(&person)

	SuccessResponse(c, person)
}

// GetPeople 获取人员列表
func GetPeople(c *gin.Context) {
	var people []models.Person
	var total int64

	// 获取查询参数
	personType := c.Query("type")
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.Person{})

	// 按类型筛选
	if personType != "" {
		query = query.Where("type = ?", personType)
	}

	// 搜索功能（姓名、电话、身份证）
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ? OR id_card LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	if err := query.Find(&people).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch people: "+err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, people)
}

// GetServicePersonnel 获取服务人员列表（仅服务人员类型，包含客户数量）
func GetServicePersonnel(c *gin.Context) {
	var people []models.Person
	var total int64

	// 获取查询参数
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.Person{}).Where("type = ?", "服务人员")

	// 搜索功能（姓名、电话、身份证）
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ? OR id_card LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	if err := query.Find(&people).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch service personnel: "+err.Error())
		return
	}

	// 构建响应数据，添加客户数量
	type ServicePersonnelResponse struct {
		models.Person
		CustomerCount int `json:"customer_count"`
	}

	response := make([]ServicePersonnelResponse, len(people))
	for i, person := range people {
		customerCount := 0
		if person.ServiceCustomerIDs != "" {
			ids := StringToIDs(person.ServiceCustomerIDs)
			customerCount = len(ids)
		}
		response[i] = ServicePersonnelResponse{
			Person:        person,
			CustomerCount: customerCount,
		}
	}

	SuccessPaginatedResponse(c, total, response)
}

// CreateServicePersonnel 创建服务人员（默认类型为服务人员）
func CreateServicePersonnel(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 强制设置为服务人员类型
	person.Type = models.PersonTypeServicePerson

	if err := config.DB.Create(&person).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create service person: "+err.Error())
		return
	}

	// 更新关联客户的ID字段
	updatePersonCustomerIDs(&person)

	SuccessResponse(c, person)
}

// GetPerson 获取人员详情
func GetPerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid person ID")
		return
	}

	var person models.Person
	if err := config.DB.First(&person, id).Error; err != nil {
		ErrorResponse(c, 404, "Person not found")
		return
	}

	// 获取关联的企业信息
	customers := getPersonRelatedCustomers(&person)
	personData := map[string]interface{}{
		"person":         person,
		"customers":      customers,
		"representative": customers["representative"],
		"investor":       customers["investor"],
		"service":        customers["service"],
	}

	SuccessResponse(c, personData)
}

// UpdatePerson 更新人员
func UpdatePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid person ID")
		return
	}

	var person models.Person
	if err := config.DB.First(&person, id).Error; err != nil {
		ErrorResponse(c, 404, "Person not found")
		return
	}

	var updateData models.Person
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 更新字段
	config.DB.Model(&person).Updates(updateData)

	// 更新关联客户的ID字段
	updatePersonCustomerIDs(&updateData)

	// 重新获取更新后的数据
	config.DB.First(&person, id)

	SuccessResponse(c, person)
}

// DeletePerson 删除人员
func DeletePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid person ID")
		return
	}

	if err := config.DB.Delete(&models.Person{}, id).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete person: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Person deleted successfully"})
}

// GetPersonCustomers 获取人员关联的企业列表
func GetPersonCustomers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid person ID")
		return
	}

	var person models.Person
	if err := config.DB.First(&person, id).Error; err != nil {
		ErrorResponse(c, 404, "Person not found")
		return
	}

	customers := getPersonRelatedCustomers(&person)
	SuccessResponse(c, customers)
}

// ============ 辅助函数 ============

// getPersonRelatedCustomers 获取人员关联的所有企业
func getPersonRelatedCustomers(person *models.Person) map[string][]models.Customer {
	result := make(map[string][]models.Customer)

	// 作为法定代表人的企业
	if person.RepresentativeCustomerIDs != "" {
		var repCustomers []models.Customer
		ids := StringToIDs(person.RepresentativeCustomerIDs)
		config.DB.Where("id IN ?", ids).Find(&repCustomers)
		result["representative"] = repCustomers
	}

	// 作为投资人的企业
	if person.InvestorCustomerIDs != "" {
		var invCustomers []models.Customer
		ids := StringToIDs(person.InvestorCustomerIDs)
		config.DB.Where("id IN ?", ids).Find(&invCustomers)
		result["investor"] = invCustomers
	}

	// 作为服务人员的企业
	if person.ServiceCustomerIDs != "" {
		var svcCustomers []models.Customer
		ids := StringToIDs(person.ServiceCustomerIDs)
		config.DB.Where("id IN ?", ids).Find(&svcCustomers)
		result["service"] = svcCustomers
	}

	return result
}

// updatePersonCustomerIDs 更新人员关联的企业ID（反向同步）
func updatePersonCustomerIDs(person *models.Person) {
	// 根据person.type，更新对应的customer_id字段
	// 这里需要在customer创建/更新时同步调用
}

// IDsToString 将uint数组转为逗号分隔字符串
func IDsToString(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(strIDs, ",")
}

// StringToIDs 将逗号分隔字符串转为uint数组
func StringToIDs(s string) []uint {
	if s == "" {
		return []uint{}
	}
	parts := strings.Split(s, ",")
	ids := make([]uint, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if id, err := strconv.ParseUint(part, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}
