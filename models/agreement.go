package models

import "time"

// Agreement 代理记账协议
type Agreement struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CustomerID      uint      `json:"customer_id" gorm:"not null"` // 关联客户
	AgreementNumber string    `json:"agreement_number" gorm:"unique"` // 协议编号
	StartDate       time.Time `json:"start_date"`                    // 协议开始日期
	EndDate         time.Time `json:"end_date"`                      // 协议结束日期
	FeeType         string    `json:"fee_type"`                      // 收费类型 (月度/季度/年度)
	Amount          float64   `json:"amount"`                        // 服务费金额
	Status          string    `json:"status"`                        // active/expired/cancelled
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	Customer  *Customer  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Payments  []Payment  `json:"payments,omitempty" gorm:"foreignKey:AgreementID"`
}
