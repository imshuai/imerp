package controllers

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"erp/services"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateAdminUserRequest 创建管理员请求
type CreateAdminUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	PersonID uint   `json:"person_id" binding:"required"`
}

// SetManagerRequest 设置管理员请求
type SetManagerRequest struct {
	PersonID uint `json:"person_id" binding:"required"`
	IsManager bool `json:"is_manager"`
}

// ApprovalRequest 审批请求
type ApprovalRequest struct {
	LogID  uint   `json:"log_id" binding:"required"`
	Reason string `json:"reason"` // 拒绝原因
}

// GetServicePeople 获取服务人员列表（用于设置管理员）
func GetServicePeople(c *gin.Context) {
	var people []models.Person
	if err := config.DB.Where("is_service_person = ?", true).Find(&people).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, people)
}

// CreateAdminUser 创建管理员（仅超级管理员）
func CreateAdminUser(c *gin.Context) {
	var req CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var count int64
	config.DB.Model(&models.AdminUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		ErrorResponse(c, 400, "Username already exists")
		return
	}

	// 验证 Person 是否存在且是服务人员
	var person models.Person
	if err := config.DB.First(&person, req.PersonID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Person not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	if !person.IsServicePerson {
		ErrorResponse(c, 400, "Person is not a service person")
		return
	}

	// 检查该人员是否已经是管理员
	var existingAdmin models.AdminUser
	if err := config.DB.Where("person_id = ?", req.PersonID).First(&existingAdmin).Error; err == nil {
		ErrorResponse(c, 400, "This person is already a manager")
		return
	}

	// 生成密码哈希
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// 创建管理员
	adminUser := models.AdminUser{
		Username:          req.Username,
		PasswordHash:      hashedPassword,
		Role:              "manager",
		PersonID:          &req.PersonID,
		MustChangePassword: true,
	}

	if err := config.DB.Create(&adminUser).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create admin user")
		return
	}

	// 同时设置 Person 的 IsManager
	if err := config.DB.Model(&person).Update("is_manager", true).Error; err != nil {
		ErrorResponse(c, 500, "Failed to update person")
		return
	}

	SuccessResponse(c, adminUser)
}

// SetManager 设置/取消服务人员为管理员（仅超级管理员）
func SetManager(c *gin.Context) {
	var req SetManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	// 验证 Person 是否存在
	var person models.Person
	if err := config.DB.First(&person, req.PersonID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Person not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	// 如果设置为管理员
	if req.IsManager {
		if !person.IsServicePerson {
			ErrorResponse(c, 400, "Person is not a service person")
			return
		}

		// 检查是否已有 AdminUser 记录
		var existingAdmin models.AdminUser
		err := config.DB.Where("person_id = ?", req.PersonID).First(&existingAdmin).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的 AdminUser
			hashedPassword, _ := auth.HashPassword("manager")
			adminUser := models.AdminUser{
				Username:          person.Name,
				PasswordHash:      hashedPassword,
				Role:              "manager",
				PersonID:          &req.PersonID,
				MustChangePassword: true,
			}
			if err := config.DB.Create(&adminUser).Error; err != nil {
				ErrorResponse(c, 500, "Failed to create admin user")
				return
			}
		}
	} else {
		// 取消管理员权限，删除对应的 AdminUser
		config.DB.Where("person_id = ?", req.PersonID).Delete(&models.AdminUser{})
	}

	// 更新 Person 的 IsManager
	if err := config.DB.Model(&person).Update("is_manager", req.IsManager).Error; err != nil {
		ErrorResponse(c, 500, "Failed to update person")
		return
	}

	SuccessResponse(c, gin.H{"message": "Manager status updated successfully"})
}

// GetAdminUsers 获取管理员列表
func GetAdminUsers(c *gin.Context) {
	var users []models.AdminUser
	if err := config.DB.Preload("Person").Find(&users).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, users)
}

// DeleteAdminUser 删除管理员（仅超级管理员）
func DeleteAdminUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	var adminUser models.AdminUser
	if err := config.DB.First(&adminUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Admin user not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	// 不能删除超级管理员
	if adminUser.Role == "super_admin" {
		ErrorResponse(c, 400, "Cannot delete super admin")
		return
	}

	// 取消对应人员的 IsManager 状态
	if adminUser.PersonID != nil {
		config.DB.Model(&models.Person{}).Where("id = ?", *adminUser.PersonID).Update("is_manager", false)
	}

	// 删除 AdminUser
	if err := config.DB.Delete(&adminUser).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete admin user")
		return
	}

	SuccessResponse(c, gin.H{"message": "Admin user deleted successfully"})
}

// GetPendingApprovals 获取待审批列表
func GetPendingApprovals(c *gin.Context) {
	auditService := services.NewAuditLogService()

	logs, err := auditService.GetPendingLogs()
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, logs)
}

// ApproveOperation 审批通过
func ApproveOperation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	// 获取审批记录
	var auditLog models.AuditLog
	if err := config.DB.Preload("User").First(&auditLog, req.LogID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Approval log not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	// 检查状态
	if auditLog.Status != "pending" {
		ErrorResponse(c, 400, "Operation is not pending")
		return
	}

	// 更新状态为已审批
	auditService := services.NewAuditLogService()
	if err := auditService.ApproveLog(req.LogID, userID.(uint)); err != nil {
		ErrorResponse(c, 500, "Failed to approve: "+err.Error())
		return
	}

	// 执行实际操作
	executor := services.NewApprovalExecutor()
	if err := executor.ExecuteOperation(&auditLog); err != nil {
		// 执行失败，回滚审批状态
		config.DB.Model(&auditLog).Updates(map[string]interface{}{
			"status":       "pending",
			"approved_by":   nil,
			"approved_at":   nil,
		})
		ErrorResponse(c, 500, "Failed to execute operation: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Operation approved and executed successfully"})
}

// RejectOperation 审批拒绝
func RejectOperation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	auditService := services.NewAuditLogService()
	if err := auditService.RejectLog(req.LogID, userID.(uint), req.Reason); err != nil {
		ErrorResponse(c, 500, "Failed to reject: "+err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Operation rejected successfully"})
}

// GetAuditLogs 获取审计日志
func GetAuditLogs(c *gin.Context) {
	status := c.Query("status")
	offset := getIntQuery(c, "offset", 0)
	limit := getIntQuery(c, "limit", 20)

	auditService := services.NewAuditLogService()
	logs, total, err := auditService.GetAuditLogs(status, offset, limit)
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
