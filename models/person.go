package models

import (
	"time"

	"gorm.io/gorm"
)

// Person 人员信息
type Person struct {
	ID                        uint       `json:"id" gorm:"primaryKey"`
	Type                      string     `json:"type" gorm:"default:'普通人员'"` // 兼容旧数据，根据IsServicePerson自动设置
	IsServicePerson           bool       `json:"is_service_person" gorm:"default:false"`
	Name                      string     `json:"name" gorm:"not null"`
	Phone                     string     `json:"phone" gorm:"not null"`
	IDCard                    string     `json:"id_card" gorm:"unique"`
	Password                  string     `json:"password"`                       // 电子税务局登录密码
	RepresentativeCustomerIDs string     `json:"representative_customer_ids"`    // 担任法人的企业ID，逗号分隔: "1,5,8"
	InvestorCustomerIDs       string     `json:"investor_customer_ids"`          // 持股的企业ID，逗号分隔: "1,2,3"
	ServiceCustomerIDs        string     `json:"service_customer_ids"`           // 服务的企业ID，逗号分隔: "1,4,7"
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

// BeforeCreate GORM hook - 创建前自动设置Type字段
func (p *Person) BeforeCreate(tx *gorm.DB) error {
	p.setType()
	return nil
}

// BeforeUpdate GORM hook - 更新前自动设置Type字段
func (p *Person) BeforeUpdate(tx *gorm.DB) error {
	p.setType()
	return nil
}

// setType 根据IsServicePerson自动设置Type字段
func (p *Person) setType() {
	if p.IsServicePerson {
		p.Type = "服务人员"
	} else {
		p.Type = "普通人员"
	}
}
