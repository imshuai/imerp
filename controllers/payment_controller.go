package controllers

import (
	"erp/config"
	"erp/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreatePayment 创建收款记录
func CreatePayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"create",
		"payment",
		nil, // resourceID is nil for create
		nil, // no old value
		payment,
		func() error {
			return config.DB.Create(&payment).Error
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回创建的数据
	if !needsApproval {
		SuccessResponse(c, payment)
	}
}

// GetPayments 获取收款记录列表
func GetPayments(c *gin.Context) {
	var payments []models.Payment
	var total int64

	// 获取查询参数
	customerID := c.Query("customer_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := config.DB.Model(&models.Payment{}).Preload("Customer").Preload("Agreement")

	// 按客户筛选
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// 按日期范围筛选
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("payment_date >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("payment_date <= ?", t)
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取列表，按日期倒序
	if err := query.Order("payment_date DESC").Find(&payments).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch payments: "+err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, payments)
}

// GetPayment 获取收款记录详情
func GetPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := config.DB.Preload("Customer").Preload("Agreement").First(&payment, id).Error; err != nil {
		ErrorResponse(c, 404, "Payment not found")
		return
	}

	SuccessResponse(c, payment)
}

// UpdatePayment 更新收款记录
func UpdatePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		ErrorResponse(c, 404, "Payment not found")
		return
	}

	var updateData models.Payment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	paymentID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"update",
		"payment",
		&paymentID,
		payment, // old value
		updateData, // new value
		func() error {
			// 更新字段
			config.DB.Model(&payment).Updates(updateData)
			// 重新获取更新后的数据
			config.DB.Preload("Customer").Preload("Agreement").First(&payment, id)
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回更新后的数据
	if !needsApproval {
		SuccessResponse(c, payment)
	}
}

// DeletePayment 删除收款记录
func DeletePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		ErrorResponse(c, 404, "Payment not found")
		return
	}

	paymentID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"delete",
		"payment",
		&paymentID,
		payment, // old value
		nil, // no new value for delete
		func() error {
			return config.DB.Delete(&models.Payment{}, id).Error
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回成功消息
	if !needsApproval {
		SuccessResponse(c, gin.H{"message": "Payment deleted successfully"})
	}
}
