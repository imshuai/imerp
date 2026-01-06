package models

import "time"

// AuditLog 审计日志
type AuditLog struct {
	ID           uint                `json:"id" gorm:"primaryKey"`
	UserID       uint                `json:"user_id" gorm:"not null"`         // 操作人ID
	UserType     string              `json:"user_type" gorm:"not null"`       // 操作人类型: admin, manager, service_person
	ActionType   string              `json:"action_type" gorm:"not null"`     // 操作类型: create, update, delete, approve, reject
	ResourceType string              `json:"resource_type" gorm:"not null"`   // 资源类型: customer, task, agreement, payment, person
	ResourceID   *uint               `json:"resource_id"`                      // 资源ID
	OldValue     string              `json:"old_value"`                        // 修改前的值（JSON）
	NewValue     string              `json:"new_value"`                        // 修改后的值（JSON）
	Status       string              `json:"status" gorm:"default:'approved'"` // 状态: pending, approved, rejected
	ApprovedBy   *uint               `json:"approved_by"`                      // 审批人ID
	ApprovedAt   *time.Time          `json:"approved_at"`                      // 审批时间
	Reason       string              `json:"reason"`                           // 拒绝原因或备注
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`

	// 关联
	User *AdminUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
