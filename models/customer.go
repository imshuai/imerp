package models

import "time"

// Customer 客户信息
type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`           // 公司名称/个人姓名
	Contact   string    `json:"contact"`                         // 联系人
	Phone     string    `json:"phone"`                           // 联系电话
	Email     string    `json:"email"`                           // 邮箱
	Address   string    `json:"address"`                         // 地址
	TaxNumber string    `json:"tax_number"`                      // 税号
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Tasks      []Task      `json:"tasks,omitempty" gorm:"foreignKey:CustomerID"`
	Agreements []Agreement `json:"agreements,omitempty" gorm:"foreignKey:CustomerID"`
	Payments   []Payment   `json:"payments,omitempty" gorm:"foreignKey:CustomerID"`
}
