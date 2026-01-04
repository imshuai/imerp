package controllers

import (
	"erp/config"
	"erp/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateAgreement 创建协议
func CreateAgreement(c *gin.Context) {
	var agreement models.Agreement
	if err := c.ShouldBindJSON(&agreement); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	if err := config.DB.Create(&agreement).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create agreement: "+err.Error())
		return
	}

	SuccessResponse(c, agreement)
}

// GetAgreements 获取协议列表
func GetAgreements(c *gin.Context) {
	var agreements []models.Agreement
	var total int64

	// 获取查询参数
	keyword := c.Query("keyword")
	status := c.Query("status")
	customerID := c.Query("customer_id")

	query := config.DB.Model(&models.Agreement{}).Preload("Customer")

	// 搜索功能
	if keyword != "" {
		query = query.Where("agreement_number LIKE ?", "%"+keyword+"%")
	}

	// 按状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 按客户筛选
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	if err := query.Find(&agreements).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch agreements: "+err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, agreements)
}

// GetAgreement 获取协议详情
func GetAgreement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid agreement ID")
		return
	}

	var agreement models.Agreement
	if err := config.DB.Preload("Customer").Preload("Payments").First(&agreement, id).Error; err != nil {
		ErrorResponse(c, 404, "Agreement not found")
		return
	}

	SuccessResponse(c, agreement)
}

// UpdateAgreement 更新协议
func UpdateAgreement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid agreement ID")
		return
	}

	var agreement models.Agreement
	if err := config.DB.First(&agreement, id).Error; err != nil {
		ErrorResponse(c, 404, "Agreement not found")
		return
	}

	var updateData models.Agreement
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 更新字段
	config.DB.Model(&agreement).Updates(updateData)

	// 重新获取更新后的数据
	config.DB.Preload("Customer").First(&agreement, id)

	SuccessResponse(c, agreement)
}

// DeleteAgreement 删除协议
func DeleteAgreement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid agreement ID")
		return
	}

	if err := config.DB.Delete(&models.Agreement{}, id).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete agreement: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Agreement deleted successfully"})
}
