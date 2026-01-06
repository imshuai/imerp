package controllers

import (
	"erp/services"

	"github.com/gin-gonic/gin"
)

// CheckApprovalNeeded 检查是否需要审批（服务人员操作需要审批）
func CheckApprovalNeeded(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	return role == "service_person"
}

// LogForApproval 记录操作到审计日志，等待审批
func LogForApproval(c *gin.Context, actionType, resourceType string, resourceID *uint, oldValue, newValue interface{}) (uint, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	auditService := services.NewAuditLogService()
	return auditService.LogOperation(
		userID.(uint),
		role.(string),
		actionType,
		resourceType,
		resourceID,
		oldValue,
		newValue,
	)
}

// RespondWithApprovalResponse 返回需要审批的响应
func RespondWithApprovalResponse(c *gin.Context, message string) {
	SuccessResponse(c, gin.H{
		"requires_approval": true,
		"message":           message,
	})
}

// HandleOperationWithApproval 处理带审批的操作
// 如果用户是服务人员，记录到审计日志并返回需要审批（返回 true）
// 如果用户是管理员，直接执行操作（返回 false）
// 返回值: (needsApproval bool, err error)
func HandleOperationWithApproval(
	c *gin.Context,
	actionType, resourceType string,
	resourceID *uint,
	oldValue, newValue interface{},
	executeFunc func() error,
) (bool, error) {
	if CheckApprovalNeeded(c) {
		// 服务人员：记录到审计日志，等待审批
		_, err := LogForApproval(c, actionType, resourceType, resourceID, oldValue, newValue)
		if err != nil {
			return false, err
		}
		RespondWithApprovalResponse(c, "操作已提交，等待管理员审批")
		return true, nil
	}

	// 管理员：直接执行操作
	if err := executeFunc(); err != nil {
		return false, err
	}

	// 操作成功后记录审计日志（不需要审批，LogOperation内部会根据userType自动设置status）
	_, _ = LogForApproval(c, actionType, resourceType, resourceID, nil, nil)

	return false, nil
}
