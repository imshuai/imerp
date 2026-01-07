package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// InvestmentRecord 出资记录
type InvestmentRecord struct {
	Date   string  `json:"date"`   // 出资日期
	Amount float64 `json:"amount"` // 出资金额
}

// CustomerInvestor 客户投资人关联表
type CustomerInvestor struct {
	ID                uint          `json:"id" gorm:"primaryKey"`
	CustomerID        uint          `json:"customer_id" gorm:"not null"`
	PersonID          uint          `json:"person_id" gorm:"not null"`
	ShareRatio        float64       `json:"share_ratio"`                    // 持股比例
	InvestmentRecords datatypes.JSON `json:"investment_records,omitempty"` // 出资记录JSON
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`

	// 关联（通过查询加载，不使用外键）
	Customer *Customer `json:"customer,omitempty" gorm:"-"`
	Person   *Person   `json:"person,omitempty" gorm:"-"`
}

// GetInvestmentRecords 解析出资记录
func (ci *CustomerInvestor) GetInvestmentRecords() []InvestmentRecord {
	if ci.InvestmentRecords == nil {
		return []InvestmentRecord{}
	}
	var records []InvestmentRecord
	json.Unmarshal(ci.InvestmentRecords, &records)
	return records
}

// SetInvestmentRecords 设置出资记录
func (ci *CustomerInvestor) SetInvestmentRecords(records []InvestmentRecord) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	ci.InvestmentRecords = data
	return nil
}
