package models

import "time"

// PersonType 人员类型
type PersonType string

const (
	PersonTypeRepresentative PersonType = "representative"   // 法定代表人
	PersonTypeInvestor       PersonType = "investor"         // 投资人
	PersonTypeServicePerson  PersonType = "service_person"   // 服务人员
	PersonTypeMixed          PersonType = "mixed"            // 混合角色
)

// Person 人员信息
type Person struct {
	ID                        uint       `json:"id" gorm:"primaryKey"`
	Type                      PersonType `json:"type" gorm:"not null"`
	Name                      string     `json:"name" gorm:"not null"`
	Phone                     string     `json:"phone" gorm:"not null"`
	IDCard                    string     `json:"id_card" gorm:"unique"`
	Password                  string     `json:"-" gorm:""`
	RepresentativeCustomerIDs string     `json:"representative_customer_ids"` // 担任法人的企业ID，逗号分隔: "1,5,8"
	InvestorCustomerIDs       string     `json:"investor_customer_ids"`        // 持股的企业ID，逗号分隔: "1,2,3"
	ServiceCustomerIDs        string     `json:"service_customer_ids"`         // 服务的企业ID，逗号分隔: "1,4,7"
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}
