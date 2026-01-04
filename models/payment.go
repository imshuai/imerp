package models

import "time"

// Payment 收款记录
type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CustomerID    uint      `json:"customer_id" gorm:"not null"` // 关联客户
	AgreementID   uint      `json:"agreement_id"`                // 关联协议 (可选)
	Amount        float64   `json:"amount" gorm:"not null"`      // 收款金额
	PaymentDate   time.Time `json:"payment_date"`                // 收款日期
	PaymentMethod string    `json:"payment_method"`              // 收款方式 (转账/现金/支票)
	Period        string    `json:"period"`                      // 费用所属期间 (如: 2024-01)
	Remark        string    `json:"remark"`                      // 备注
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	Customer  *Customer  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Agreement *Agreement `json:"agreement,omitempty" gorm:"foreignKey:AgreementID"`
}
