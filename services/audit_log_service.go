package services

import (
	"erp/config"
	"erp/models"
	"encoding/json"
)

// AuditLogService 审计日志服务
type AuditLogService struct{}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService() *AuditLogService {
	return &AuditLogService{}
}

// LogOperation 记录操作日志
func (s *AuditLogService) LogOperation(userID uint, userType, actionType, resourceType string, resourceID *uint, oldValue, newValue interface{}) (uint, error) {
	var oldValueJSON, newValueJSON string

	if oldValue != nil {
		data, err := json.Marshal(oldValue)
		if err == nil {
			oldValueJSON = string(data)
		}
	}

	if newValue != nil {
		data, err := json.Marshal(newValue)
		if err == nil {
			newValueJSON = string(data)
		}
	}

	// 判断状态：如果是超级管理员或管理员操作，直接通过；普通服务人员操作需要审批
	status := "pending" // 默认为待审批
	if userType == "super_admin" || userType == "manager" {
		status = "approved"
	}

	auditLog := models.AuditLog{
		UserID:       userID,
		UserType:     userType,
		ActionType:   actionType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValue:     oldValueJSON,
		NewValue:     newValueJSON,
		Status:       status,
	}

	// 如果是自动审批，设置审批人和时间
	if status == "approved" {
		auditLog.ApprovedBy = &userID
		now := config.DB.NowFunc()
		auditLog.ApprovedAt = &now
	}

	if err := config.DB.Create(&auditLog).Error; err != nil {
		return 0, err
	}

	return auditLog.ID, nil
}

// LogAutoApproved 自动审批的日志（admin/manager 操作）
func (s *AuditLogService) LogAutoApproved(userID uint, userType, actionType, resourceType string, resourceID *uint, oldValue, newValue interface{}) (uint, error) {
	logID, err := s.LogOperation(userID, userType, actionType, resourceType, resourceID, oldValue, newValue)
	if err != nil {
		return 0, err
	}

	// 更新为已审批状态
	now := config.DB.NowFunc
	config.DB.Model(&models.AuditLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_by": userID,
		"approved_at": now,
	})

	return logID, nil
}

// GetPendingLogs 获取待审批的日志列表
func (s *AuditLogService) GetPendingLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := config.DB.Where("status = ?", "pending").
		Preload("User").
		Order("created_at desc").
		Find(&logs).Error
	return logs, err
}

// ApproveLog 审批通过
func (s *AuditLogService) ApproveLog(logID, approvedBy uint) error {
	now := config.DB.NowFunc()
	return config.DB.Model(&models.AuditLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_by": approvedBy,
		"approved_at": now,
	}).Error
}

// RejectLog 审批拒绝
func (s *AuditLogService) RejectLog(logID, approvedBy uint, reason string) error {
	now := config.DB.NowFunc()
	return config.DB.Model(&models.AuditLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"status":      "rejected",
		"approved_by": approvedBy,
		"approved_at": now,
		"reason":      reason,
	}).Error
}

// GetAuditLogs 获取审计日志列表（可筛选）
func (s *AuditLogService) GetAuditLogs(status string, offset, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := config.DB.Model(&models.AuditLog{}).Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&logs).Error

	return logs, total, err
}
