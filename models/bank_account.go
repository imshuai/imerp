package models

import "time"

// AccountType 账户类型
type AccountType string

const (
	AccountTypeBasic     AccountType = "基本户"
	AccountTypeGeneral   AccountType = "一般户"
	AccountTypeTemporary AccountType = "临时户"
)

// GetAccountTypeOptions 获取账户类型选项
func GetAccountTypeOptions() []AccountType {
	return []AccountType{
		AccountTypeBasic,
		AccountTypeGeneral,
		AccountTypeTemporary,
	}
}

// BankAccount 对公账户
type BankAccount struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	CustomerID    uint        `json:"customer_id" gorm:"not null"`
	BankName      string      `json:"bank_name" gorm:"not null"`       // 开户银行
	AccountNumber string      `json:"account_number" gorm:"not null"`  // 账号
	BankCode      string      `json:"bank_code"`                       // 开户行号
	ContactPhone  string      `json:"contact_phone"`                   // 联系电话
	AccountType   AccountType `json:"account_type" gorm:"not null"`    // 账户类型
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

	// 关联（通过查询加载，不使用外键）
	Customer *Customer `json:"customer,omitempty" gorm:"-"`
}
