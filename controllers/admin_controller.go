package controllers

import (
	"erp/config"
	"erp/models"
	"erp/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LogOperation 记录操作到审计日志
func LogOperation(c *gin.Context, actionType, resourceType string, resourceID *uint, resourceName string, oldValue, newValue interface{}) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("role")

	auditService := services.NewAuditLogService()
	auditService.LogOperation(userID.(uint), userType.(string), actionType, resourceType, resourceID, resourceName, oldValue, newValue)
}

// DeleteAuditLog 删除审计日志（仅超级管理员）
func DeleteAuditLog(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		ErrorResponse(c, 400, "Invalid log ID")
		return
	}

	if err := config.DB.Delete(&models.AuditLog{}, id).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete audit log")
		return
	}

	SuccessResponse(c, gin.H{"message": "Audit log deleted successfully"})
}

// ClearAuditLogs 批量清理审计日志（仅超级管理员）
func ClearAuditLogs(c *gin.Context) {
	var req struct {
		StartDate string `json:"start_date"` // 开始日期，格式: YYYY-MM-DD
		EndDate   string `json:"end_date"`   // 结束日期，格式: YYYY-MM-DD
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request data")
		return
	}

	var result *gorm.DB
	// 如果指定了日期范围，按日期过滤
	if req.StartDate != "" && req.EndDate != "" {
		// 结束日期设置为当天的 23:59:59，以包含当天的所有记录
		endDateTime := req.EndDate + " 23:59:59"
		result = config.DB.Where("created_at >= ? AND created_at <= ?", req.StartDate, endDateTime).Delete(&models.AuditLog{})
	} else {
		// 清理全部日志
		result = config.DB.Where("1 = 1").Delete(&models.AuditLog{})
	}

	if result.Error != nil {
		ErrorResponse(c, 500, "Failed to clear audit logs")
		return
	}

	SuccessResponse(c, gin.H{
		"message": "Audit logs cleared successfully",
		"count":   result.RowsAffected,
	})
}

// GetAuditLogs 获取审计日志（所有登录用户可查看）
func GetAuditLogs(c *gin.Context) {
	offset := getIntQuery(c, "offset", 0)
	limit := getIntQuery(c, "limit", 20)

	var logs []models.AuditLog
	var total int64

	query := config.DB.Model(&models.AuditLog{}).Preload("User.Person")

	// 获取总数
	query.Count(&total)

	// 获取列表
	err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&logs).Error

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, logs)
}

// 辅助函数
func getIDParam(c *gin.Context) (uint, error) {
	var id struct {
		ID uint `uri:"id"`
	}
	if err := c.ShouldBindUri(&id); err != nil {
		return 0, err
	}
	return id.ID, nil
}

func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.DefaultQuery(key, "")
	if value == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}
