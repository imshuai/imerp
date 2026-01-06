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

// LogOperation 记录操作日志（所有操作直接记录，不需要审批）
func (s *AuditLogService) LogOperation(userID uint, userType, actionType, resourceType string, resourceID *uint, resourceName string, oldValue, newValue interface{}) (uint, error) {
	var oldValueJSON, newValueJSON string

	if oldValue != nil {
		// 清理关联字段后再序列化
		cleanedOld := cleanAssociations(oldValue)
		data, err := json.Marshal(cleanedOld)
		if err == nil {
			oldValueJSON = string(data)
		}
	}

	if newValue != nil {
		// 清理关联字段后再序列化
		cleanedNew := cleanAssociations(newValue)
		data, err := json.Marshal(cleanedNew)
		if err == nil {
			newValueJSON = string(data)
		}
	}

	auditLog := models.AuditLog{
		UserID:       userID,
		UserType:     userType,
		ActionType:   actionType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		OldValue:     oldValueJSON,
		NewValue:     newValueJSON,
	}

	if err := config.DB.Create(&auditLog).Error; err != nil {
		return 0, err
	}

	return auditLog.ID, nil
}

// cleanAssociations 清理对象中的关联字段，避免序列化整个关联对象
func cleanAssociations(v interface{}) interface{} {
	// 将对象转换为 map
	data, err := json.Marshal(v)
	if err != nil {
		return v
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return v
	}

	// 在删除关联字段之前，尝试将ID字段转换为名字字段
	// 处理 service_person_ids
	if servicePersons, ok := m["service_persons"].([]interface{}); ok && len(servicePersons) > 0 {
		var names []string
		for _, sp := range servicePersons {
			if spMap, ok := sp.(map[string]interface{}); ok {
				if name, ok := spMap["name"].(string); ok {
					names = append(names, name)
				}
			}
		}
		if len(names) > 0 {
			m["service_person_ids"] = names
		}
	}

	// 删除 GORM 关联对象字段（JSON 序列化后的 snake_case 格式）
	delete(m, "Customer")
	delete(m, "Agreement")
	delete(m, "Payments")
	delete(m, "Tasks")
	delete(m, "Representative")
	delete(m, "InvestorList")
	delete(m, "ServicePersons")
	delete(m, "Agreements")
	delete(m, "User")
	delete(m, "Person")

	// 删除由 loadCustomerRelations 等函数加载的关联列表字段
	delete(m, "representative")
	delete(m, "investor_list")
	delete(m, "service_persons")
	delete(m, "agreements_list")
	// 删除小写的关联字段（根据各模型的 json tag）
	delete(m, "customer")
	delete(m, "agreement")
	delete(m, "payments")
	delete(m, "tasks")
	delete(m, "user")
	delete(m, "person")

	return m
}
