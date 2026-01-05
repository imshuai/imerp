package models

import "time"

// FeeType 收费类型
type FeeType string

const (
	FeeTypeMonthly   FeeType = "月度"   // 月度
	FeeTypeQuarterly  FeeType = "季度"   // 季度
	FeeTypeYearly    FeeType = "年度"   // 年度
)

// AgreementStatus 协议状态
type AgreementStatus string

const (
	AgreementStatusActive    AgreementStatus = "有效"     // 有效
	AgreementStatusExpired   AgreementStatus = "已过期"   // 已过期
	AgreementStatusCancelled AgreementStatus = "已取消"   // 已取消
)

// Agreement 代理记账协议
type Agreement struct {
	ID              uint             `json:"id" gorm:"primaryKey"`
	CustomerID      uint             `json:"customer_id" gorm:"not null"` // 关联客户
	AgreementNumber string           `json:"agreement_number" gorm:"unique"` // 协议编号
	StartDate       time.Time        `json:"start_date"`                    // 协议开始日期
	EndDate         time.Time        `json:"end_date"`                      // 协议结束日期
	FeeType         FeeType          `json:"fee_type"`                      // 收费类型
	Amount          float64          `json:"amount"`                        // 服务费金额
	Status          AgreementStatus  `json:"status"`                        // 协议状态
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`

	// 关联
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:AgreementID"`
}
