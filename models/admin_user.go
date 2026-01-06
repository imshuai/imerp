package models

import "time"

// AdminUser 管理员用户表（所有用户都需要通过此表登录）
type AdminUser struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Username          string     `json:"username" gorm:"uniqueIndex;not null"`        // 登录用户名
	PasswordHash      string     `json:"-" gorm:"not null"`                            // 密码哈希
	Role              string     `json:"role" gorm:"not null"`                          // 角色: super_admin, manager, service_person
	PersonID          *uint      `json:"person_id"`                                     // 关联的服务人员ID
	MustChangePassword bool      `json:"must_change_password" gorm:"default:true"`     // 是否必须修改密码
	LastPasswordChange *time.Time `json:"last_password_change"`                         // 最后修改密码时间
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// 关联
	Person *Person `json:"person,omitempty" gorm:"foreignKey:PersonID"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}
