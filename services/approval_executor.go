package services

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ApprovalExecutor 审批执行器
type ApprovalExecutor struct{}

// NewApprovalExecutor 创建审批执行器
func NewApprovalExecutor() *ApprovalExecutor {
	return &ApprovalExecutor{}
}

// ExecuteOperation 执行操作（用于审批通过时调用）
func (e *ApprovalExecutor) ExecuteOperation(log *models.AuditLog) error {
	switch log.ResourceType {
	case "customer":
		return e.executeCustomerOperation(log)
	case "task":
		return e.executeTaskOperation(log)
	case "agreement":
		return e.executeAgreementOperation(log)
	case "payment":
		return e.executePaymentOperation(log)
	case "person":
		return e.executePersonOperation(log)
	default:
		return fmt.Errorf("unknown resource type: %s", log.ResourceType)
	}
}

// executeCustomerOperation 执行客户操作
func (e *ApprovalExecutor) executeCustomerOperation(log *models.AuditLog) error {
	if log.ResourceID == nil {
		return errors.New("resource_id is required")
	}

	var newData map[string]interface{}
	if err := json.Unmarshal([]byte(log.NewValue), &newData); err != nil {
		return err
	}

	switch log.ActionType {
	case "create":
		// 创建客户
		customer := models.Customer{}
		if err := config.DB.Create(&customer).Error; err != nil {
			return err
		}
		// 更新 audit_log 的 resource_id
		config.DB.Model(log).Update("resource_id", customer.ID)

	case "update":
		// 更新客户
		if err := config.DB.Model(&models.Customer{}).Where("id = ?", *log.ResourceID).Updates(newData).Error; err != nil {
			return err
		}
	}

	return nil
}

// executeTaskOperation 执行任务操作
func (e *ApprovalExecutor) executeTaskOperation(log *models.AuditLog) error {
	if log.ResourceID == nil {
		return errors.New("resource_id is required")
	}

	var newData map[string]interface{}
	if err := json.Unmarshal([]byte(log.NewValue), &newData); err != nil {
		return err
	}

	switch log.ActionType {
	case "create":
		// 创建任务
		task := models.Task{}
		if err := config.DB.Create(&task).Error; err != nil {
			return err
		}
		config.DB.Model(log).Update("resource_id", task.ID)

	case "update":
		// 更新任务 - 排除关联字段
		delete(newData, "Customer")
		delete(newData, "customer")
		if err := config.DB.Model(&models.Task{}).Where("id = ?", *log.ResourceID).Updates(newData).Error; err != nil {
			return err
		}
	}

	return nil
}

// executeAgreementOperation 执行协议操作
func (e *ApprovalExecutor) executeAgreementOperation(log *models.AuditLog) error {
	if log.ResourceID == nil {
		return errors.New("resource_id is required")
	}

	var newData map[string]interface{}
	if err := json.Unmarshal([]byte(log.NewValue), &newData); err != nil {
		return err
	}

	switch log.ActionType {
	case "create":
		agreement := models.Agreement{}
		if err := config.DB.Create(&agreement).Error; err != nil {
			return err
		}
		config.DB.Model(log).Update("resource_id", agreement.ID)

	case "update":
		// 更新协议 - 排除关联字段
		delete(newData, "Customer")
		delete(newData, "customer")
		delete(newData, "Payments")
		delete(newData, "payments")
		if err := config.DB.Model(&models.Agreement{}).Where("id = ?", *log.ResourceID).Updates(newData).Error; err != nil {
			return err
		}
	}

	return nil
}

// executePaymentOperation 执行收款操作
func (e *ApprovalExecutor) executePaymentOperation(log *models.AuditLog) error {
	if log.ResourceID == nil {
		return errors.New("resource_id is required")
	}

	var newData map[string]interface{}
	if err := json.Unmarshal([]byte(log.NewValue), &newData); err != nil {
		return err
	}

	switch log.ActionType {
	case "create":
		payment := models.Payment{}
		if err := config.DB.Create(&payment).Error; err != nil {
			return err
		}
		config.DB.Model(log).Update("resource_id", payment.ID)

	case "update":
		// 更新收款 - 排除关联字段
		delete(newData, "Customer")
		delete(newData, "customer")
		delete(newData, "Agreement")
		delete(newData, "agreement")
		if err := config.DB.Model(&models.Payment{}).Where("id = ?", *log.ResourceID).Updates(newData).Error; err != nil {
			return err
		}
	}

	return nil
}

// executePersonOperation 执行人员操作
func (e *ApprovalExecutor) executePersonOperation(log *models.AuditLog) error {
	if log.ResourceID == nil {
		return errors.New("resource_id is required")
	}

	var newData map[string]interface{}
	if err := json.Unmarshal([]byte(log.NewValue), &newData); err != nil {
		return err
	}

	switch log.ActionType {
	case "create":
		// 创建人员
		person := models.Person{}
		if err := config.DB.Create(&person).Error; err != nil {
			return err
		}
		config.DB.Model(log).Update("resource_id", person.ID)
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

	case "update":
		// 更新人员
		if err := config.DB.Model(&models.Person{}).Where("id = ?", *log.ResourceID).Updates(newData).Error; err != nil {
			return err
		}
	}

	return nil
}
