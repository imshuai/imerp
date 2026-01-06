package controllers

import (
	"erp/auth"
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

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"create",
		"person",
		nil, // resourceID is nil for create
		nil, // no old value
		person,
		func() error {
			if err := config.DB.Create(&person).Error; err != nil {
				return err
			}
			// 更新关联客户的ID字段
			updatePersonCustomerIDs(&person)
			// 如果是服务人员，创建对应的 AdminUser（默认密码 123456）
			if person.IsServicePerson {
				hashedPassword, _ := auth.HashPassword("123456")
				adminUser := models.AdminUser{
					Username:          person.Name,
					PasswordHash:      hashedPassword,
					Role:              "service_person",
					PersonID:          &person.ID,
					MustChangePassword: true,
				}
				if err := config.DB.Create(&adminUser).Error; err != nil {
					return err
				}
			}
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回创建的数据
	if !needsApproval {
		SuccessResponse(c, person)
	}
}

// GetPeople 获取人员列表
func GetPeople(c *gin.Context) {
	var people []models.Person
	var total int64

	// 获取查询参数
	isServicePerson := c.Query("is_service_person")
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.Person{})

	// 按是否服务人员筛选
	if isServicePerson != "" {
		query = query.Where("is_service_person = ?", isServicePerson == "true" || isServicePerson == "1")
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

	personID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"update",
		"person",
		&personID,
		person, // old value
		updateData, // new value
		func() error {
			// 更新字段
			config.DB.Model(&person).Updates(updateData)
			// 更新关联客户的ID字段
			updatePersonCustomerIDs(&updateData)
			// 重新获取更新后的数据
			config.DB.First(&person, id)
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回更新后的数据
	if !needsApproval {
		SuccessResponse(c, person)
	}
}

// DeletePerson 删除人员
func DeletePerson(c *gin.Context) {
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

	personID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"delete",
		"person",
		&personID,
		person, // old value
		nil, // no new value for delete
		func() error {
			return config.DB.Delete(&models.Person{}, id).Error
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回成功消息
	if !needsApproval {
		SuccessResponse(c, gin.H{"message": "Person deleted successfully"})
	}
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
