package models

import "time"

// AuditLog 审计日志
type AuditLog struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	UserID        uint       `json:"user_id" gorm:"not null"`       // 操作人ID
	UserType      string     `json:"user_type" gorm:"not null"`     // 操作人类型: admin, service_person
	ActionType    string     `json:"action_type" gorm:"not null"`   // 操作类型: create, update, delete
	ResourceType  string     `json:"resource_type" gorm:"not null"` // 资源类型: customer, task, agreement, payment, person
	ResourceID    *uint      `json:"resource_id"`                    // 资源ID
	ResourceName  string     `json:"resource_name"`                  // 资源名称（客户名称、服务人员名称、任务标题等）
	OldValue      string     `json:"old_value"`                     // 修改前的值（JSON）
	NewValue      string     `json:"new_value"`                     // 修改后的值（JSON）
	Status        string     `json:"status" gorm:"default:'approved'"` // 状态（保留字段，暂不使用）
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联
	User *AdminUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
