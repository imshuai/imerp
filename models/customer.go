package models

import (
	"encoding/json"
	"time"
	"gorm.io/datatypes"
)

// JSONDate 日期类型，数据库存储为字符串，JSON序列化为YYYY-MM-DD格式
type JSONDate string

// UnmarshalJSON 实现json.Unmarshaler接口，解析YYYY-MM-DD格式
func (d *JSONDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*d = ""
		return nil
	}
	// 验证日期格式
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = JSONDate(s)
	return nil
}

// MarshalJSON 实现json.Marshaler接口，输出YYYY-MM-DD格式
func (d JSONDate) MarshalJSON() ([]byte, error) {
	if d == "" {
		return json.Marshal("")
	}
	return json.Marshal(string(d))
}

// ToTime 转换为time.Time
func (d JSONDate) ToTime() *time.Time {
	if d == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", string(d))
	if err != nil {
		return nil
	}
	return &t
}

// JSONDateFromTime 从time.Time创建JSONDate
func JSONDateFromTime(t *time.Time) JSONDate {
	if t == nil {
		return ""
	}
	return JSONDate(t.Format("2006-01-02"))
}

// CustomerType 客户类型
type CustomerType string

const (
	CustomerTypeLimitedCompany     CustomerType = "有限公司"      // 有限公司
	CustomerTypeSoleProprietorship CustomerType = "个人独资企业"  // 个人独资企业
	CustomerTypePartnership        CustomerType = "合伙企业"          // 合伙企业
	CustomerTypeIndividualBusiness CustomerType = "个体工商户"  // 个体工商户
)

// InvestorInfo 投资人信息（JSON结构）
type InvestorInfo struct {
	PersonID          uint                `json:"person_id"`
	ShareRatio        float64             `json:"share_ratio"`         // 持股比例
	InvestmentRecords []InvestmentRecord  `json:"investment_records,omitempty"` // 出资记录（可选）
}

// InvestmentRecord 出资记录
type InvestmentRecord struct {
	Date   string  `json:"date"`   // 出资日期
	Amount float64 `json:"amount"` // 出资金额
}

// Customer 客户信息
type Customer struct {
	ID                      uint          `json:"id" gorm:"primaryKey"`
	Name                    string        `json:"name" gorm:"not null"` // 公司名称/个人姓名
	Phone                   string        `json:"phone"`                // 联系电话
	Address                 string        `json:"address"`              // 地址
	TaxNumber               string        `json:"tax_number"`           // 税号
	Type                    CustomerType  `json:"type" gorm:"not null"` // 客户类型
	RepresentativeID        *uint         `json:"representative_id"`    // 法定代表人ID
	Investors               datatypes.JSON `json:"investors"`           // 投资人JSON数组
	ServicePersonIDs        string        `json:"service_person_ids"`  // 服务人员ID，逗号分隔: "5,6"
	AgreementIDs            string        `json:"agreement_ids"`       // 代理协议ID，逗号分隔: "1,3,5"
	RegisteredCapital       float64       `json:"registered_capital"`   // 注册资本
	LicenseRegistrationDate  JSONDate     `json:"license_registration_date" gorm:"type:varchar(10)"`   // 执照登记日
	TaxRegistrationDate      JSONDate     `json:"tax_registration_date" gorm:"type:varchar(10)"`       // 税务登记日
	TaxOffice                string        `json:"tax_office"`                  // 税务所
	TaxAdministrator         string        `json:"tax_administrator"`           // 税务管理员
	TaxAdministratorPhone    string        `json:"tax_administrator_phone"`     // 税务管理员联系电话
	TaxpayerType             string        `json:"taxpayer_type"`               // 纳税人类型（一般纳税人/小规模纳税人）
	CreatedAt               time.Time     `json:"created_at"`
	UpdatedAt               time.Time     `json:"updated_at"`

	// 关联（通过查询加载，不存储在数据库）
	Representative *Person     `json:"representative,omitempty" gorm:"-"`
	InvestorList   []Person    `json:"investor_list,omitempty" gorm:"-"`
	ServicePersons []Person    `json:"service_persons,omitempty" gorm:"-"`
	Agreements     []Agreement `json:"agreements_list,omitempty" gorm:"-"`

	// 原有关联
	Tasks    []Task    `json:"tasks,omitempty" gorm:"foreignKey:CustomerID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:CustomerID"`
}
