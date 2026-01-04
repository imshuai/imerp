package controllers

import (
	"erp/config"
	"erp/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCustomer 创建客户
func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	if err := config.DB.Create(&customer).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create customer: "+err.Error())
		return
	}

	SuccessResponse(c, customer)
}

// GetCustomers 获取客户列表
func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	var total int64

	// 获取查询参数
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.Customer{})

	// 搜索功能
	if keyword != "" {
		query = query.Where("name LIKE ? OR tax_number LIKE ? OR contact LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	if err := query.Find(&customers).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch customers: "+err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, customers)
}

// GetCustomer 获取客户详情
func GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := config.DB.Preload("Tasks").Preload("Agreements").Preload("Payments").First(&customer, id).Error; err != nil {
		ErrorResponse(c, 404, "Customer not found")
		return
	}

	SuccessResponse(c, customer)
}

// UpdateCustomer 更新客户
func UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		ErrorResponse(c, 404, "Customer not found")
		return
	}

	var updateData models.Customer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 更新字段
	config.DB.Model(&customer).Updates(updateData)

	// 重新获取更新后的数据
	config.DB.First(&customer, id)

	SuccessResponse(c, customer)
}

// DeleteCustomer 删除客户
func DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	if err := config.DB.Delete(&models.Customer{}, id).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete customer: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Customer deleted successfully"})
}

// GetCustomerTasks 获取客户的任务列表
func GetCustomerTasks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var tasks []models.Task
	if err := config.DB.Where("customer_id = ?", id).Find(&tasks).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch tasks: "+err.Error())
		return
	}

	SuccessResponse(c, tasks)
}

// GetCustomerPayments 获取客户的收款记录
func GetCustomerPayments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var payments []models.Payment
	if err := config.DB.Where("customer_id = ?", id).Find(&payments).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch payments: "+err.Error())
		return
	}

	SuccessResponse(c, payments)
}
