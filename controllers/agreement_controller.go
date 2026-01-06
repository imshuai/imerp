package controllers

import (
	"erp/config"
	"erp/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateAgreement 创建协议
func CreateAgreement(c *gin.Context) {
	var agreement models.Agreement
	if err := c.ShouldBindJSON(&agreement); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"create",
		"agreement",
		nil, // resourceID is nil for create
		nil, // no old value
		agreement,
		func() error {
			// 如果协议编号为空，自动生成
			if agreement.AgreementNumber == "" {
				number, err := GenerateAgreementNumber(config.DB)
				if err != nil {
					return err
				}
				agreement.AgreementNumber = number
			}
			return config.DB.Create(&agreement).Error
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回创建的数据
	if !needsApproval {
		SuccessResponse(c, agreement)
	}
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

	agreementID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"update",
		"agreement",
		&agreementID,
		agreement, // old value
		updateData, // new value
		func() error {
			// 更新字段
			config.DB.Model(&agreement).Updates(updateData)
			// 重新获取更新后的数据
			config.DB.Preload("Customer").First(&agreement, id)
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回更新后的数据
	if !needsApproval {
		SuccessResponse(c, agreement)
	}
}

// DeleteAgreement 删除协议
func DeleteAgreement(c *gin.Context) {
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

	agreementID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"delete",
		"agreement",
		&agreementID,
		agreement, // old value
		nil, // no new value for delete
		func() error {
			return config.DB.Delete(&models.Agreement{}, id).Error
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回成功消息
	if !needsApproval {
		SuccessResponse(c, gin.H{"message": "Agreement deleted successfully"})
	}
}

// ============ 辅助函数 ============

// GenerateAgreementNumber 生成协议编号
// 格式: AGT + YYYYMMDD + 4位序号
// 例如: AGT202601050001
func GenerateAgreementNumber(db *gorm.DB) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	// 查询当天最大的协议编号
	prefix := "AGT" + dateStr
	var lastAgreement models.Agreement
	err := db.Where("agreement_number LIKE ?", prefix+"%").Order("agreement_number DESC").First(&lastAgreement).Error

	serialNum := 1
	if err == nil {
		// 提取序号部分
		lastNumber := lastAgreement.AgreementNumber
		if len(lastNumber) >= len(prefix)+4 {
			serialStr := lastNumber[len(prefix):]
			if num, err := strconv.ParseInt(serialStr, 10, 32); err == nil {
				serialNum = int(num) + 1
			}
		}
	}

	// 生成新编号：AGT + YYYYMMDD + 4位序号（补零）
	newNumber := fmt.Sprintf("%s%04d", prefix, serialNum)
	return newNumber, nil
}
