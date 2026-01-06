package import_export

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CustomerImportService 客户导入服务
type CustomerImportService struct {
	db *gorm.DB
}

// NewCustomerImportService 创建客户导入服务
func NewCustomerImportService(db *gorm.DB) *CustomerImportService {
	return &CustomerImportService{db: db}
}

// ImportCustomersFromExcel 从Excel导入客户
func (s *CustomerImportService) ImportCustomersFromExcel(filePath string, strategy ImportStrategy) (*ImportResult, error) {
	// 打开Excel文件
	excelService, err := NewExcelServiceFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer excelService.Close()

	// 读取所有行
	rows, err := excelService.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("读取Excel数据失败: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel文件没有数据行")
	}

	// 解析表头
	headers := rows[0]
	colIndex := make(map[string]int)
	for i, header := range headers {
		colIndex[header] = i
	}

	// 验证必需的列
	requiredColumns := []string{"公司名称", "税号", "客户类型"}
	for _, col := range requiredColumns {
		if _, exists := colIndex[col]; !exists {
			return nil, fmt.Errorf("缺少必需的列: %s", col)
		}
	}

	result := &ImportResult{
		Total:  len(rows) - 1, // 减去表头
		Errors: []ImportError{},
	}

	// 处理数据行（从第2行开始）
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1

		// 解析行数据
		customer, parseErr := s.parseCustomerRow(row, colIndex, rowNum)
		if parseErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, *parseErr)
			continue
		}

		// 导入客户
		importErr := s.importCustomer(customer, strategy, rowNum, result)
		if importErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, *importErr)
		} else {
			result.Success++
		}
	}

	return result, nil
}

// CustomerRowData 客户行数据
type CustomerRowData struct {
	Name                  string
	Phone                 string
	Address               string
	TaxNumber             string
	CustomerType          string
	RegisteredCapital     float64
	LicenseRegistrationDate string // 执照登记日
	TaxRegistrationDate     string // 税务登记日
	TaxOffice               string // 税务所
	TaxAdministrator        string // 税务管理员
	TaxAdministratorPhone   string // 税务管理员联系电话
	TaxpayerType            string // 纳税人类型（一般纳税人/小规模纳税人）
	RepresentativeName   string
	RepresentativeIDCard string
	InvestorsInfo      string // 格式: 姓名:身份证号:持股比例;...
	ServicePeopleInfo  string // 格式: 姓名,姓名,姓名
	AgreementsInfo     string // 格式: 有效期起:有效期止:收费类型:收费金额|...
}

// InvestorInfo 投资人信息
type InvestorInfo struct {
	Name        string
	IDCard      string
	ShareRatio  float64
}

// AgreementInfo 协议信息
type AgreementInfo struct {
	StartDate  time.Time
	EndDate    time.Time
	FeeType    string
	FeeAmount  float64
}

// parseCustomerRow 解析客户行数据
func (s *CustomerImportService) parseCustomerRow(row []string, colIndex map[string]int, rowNum int) (*CustomerRowData, *ImportError) {
	getCell := func(colName string) string {
		idx, exists := colIndex[colName]
		if !exists || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	data := &CustomerRowData{
		Name:                 getCell("公司名称"),
		Phone:                getCell("联系电话"),
		Address:              getCell("地址"),
		TaxNumber:            getCell("税号"),
		CustomerType:         getCell("客户类型"),
		LicenseRegistrationDate: getCell("执照登记日"),
		TaxRegistrationDate:     getCell("税务登记日"),
		TaxOffice:               getCell("税务所"),
		TaxAdministrator:        getCell("税务管理员"),
		TaxAdministratorPhone:   getCell("税务管理员联系电话"),
		TaxpayerType:            getCell("纳税人类型"),
		RepresentativeName:  getCell("法定代表人姓名"),
		RepresentativeIDCard: getCell("法定代表人身份证"),
		InvestorsInfo:       getCell("投资人信息"),
		ServicePeopleInfo:   getCell("服务人员信息"),
		AgreementsInfo:      getCell("协议信息"),
	}

	// 解析注册资本
	capitalStr := getCell("注册资本")
	if capitalStr != "" {
		capital, err := ParseFloat(capitalStr)
		if err != nil {
			return nil, &ImportError{Row: rowNum, Column: "注册资本", Message: "注册资本必须是数字"}
		}
		data.RegisteredCapital = capital
	}

	// 验证必填字段
	if data.Name == "" {
		return nil, &ImportError{Row: rowNum, Column: "公司名称", Message: "公司名称不能为空"}
	}
	if data.TaxNumber == "" {
		return nil, &ImportError{Row: rowNum, Column: "税号", Message: "税号不能为空"}
	}
	if data.CustomerType == "" {
		return nil, &ImportError{Row: rowNum, Column: "客户类型", Message: "客户类型不能为空"}
	}

	// 验证客户类型枚举值
	validTypes := []string{"有限公司", "个人独资企业", "合伙企业", "个体工商户"}
	validType := false
	for _, vt := range validTypes {
		if data.CustomerType == vt {
			validType = true
			break
		}
	}
	if !validType {
		return nil, &ImportError{
			Row: rowNum, Column: "客户类型",
			Message: fmt.Sprintf("客户类型必须是以下之一: %s", strings.Join(validTypes, "、")),
		}
	}

	// 验证纳税人类型枚举值（如果填写了的话）
	if data.TaxpayerType != "" {
		validTaxpayerTypes := []string{"一般纳税人", "小规模纳税人"}
		validTaxpayerType := false
		for _, vt := range validTaxpayerTypes {
			if data.TaxpayerType == vt {
				validTaxpayerType = true
				break
			}
		}
		if !validTaxpayerType {
			return nil, &ImportError{
				Row: rowNum, Column: "纳税人类型",
				Message: fmt.Sprintf("纳税人类型必须是以下之一: %s", strings.Join(validTaxpayerTypes, "、")),
			}
		}
	}

	return data, nil
}

// importCustomer 导入单个客户
func (s *CustomerImportService) importCustomer(data *CustomerRowData, strategy ImportStrategy, rowNum int, result *ImportResult) *ImportError {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询是否已存在
	var existingCustomer map[string]interface{}
	err := tx.Raw("SELECT id FROM customers WHERE tax_number = ?", data.TaxNumber).Scan(&existingCustomer).Error
	isConflict := err == nil && existingCustomer != nil

	var customerID int64

	if isConflict {
		switch strategy {
		case StrategySkip:
			tx.Rollback()
			return nil
		case StrategyUpdate:
			// 更新已存在的记录
			err := tx.Table("customers").
				Where("tax_number = ?", data.TaxNumber).
				Updates(map[string]interface{}{
					"name":                     data.Name,
					"phone":                    data.Phone,
					"address":                  data.Address,
					"type":                     data.CustomerType,
					"registered_capital":       data.RegisteredCapital,
					"license_registration_date": data.LicenseRegistrationDate,
					"tax_registration_date":     data.TaxRegistrationDate,
					"tax_office":                data.TaxOffice,
					"tax_administrator":         data.TaxAdministrator,
					"tax_administrator_phone":   data.TaxAdministratorPhone,
					"taxpayer_type":             data.TaxpayerType,
				}).Error
			if err != nil {
				tx.Rollback()
				return &ImportError{Row: rowNum, Column: "", Message: fmt.Sprintf("更新客户失败: %v", err)}
			}
			// 获取更新后的ID
			tx.Raw("SELECT id FROM customers WHERE tax_number = ?", data.TaxNumber).Scan(&customerID)
		case StrategyCreateNew:
			// 修改税号后创建
			suffix := 1
			newTaxNumber := data.TaxNumber
			for {
				var count int64
				tx.Table("customers").Where("tax_number = ?", newTaxNumber).Count(&count)
				if count == 0 {
					break
				}
				suffix++
				newTaxNumber = fmt.Sprintf("%s_%d", data.TaxNumber, suffix)
			}
			data.TaxNumber = newTaxNumber
		}
	}

	// 创建新客户记录
	if customerID == 0 {
		err := tx.Table("customers").Create(map[string]interface{}{
			"name":                       data.Name,
			"phone":                      data.Phone,
			"address":                    data.Address,
			"tax_number":                 data.TaxNumber,
			"type":                       data.CustomerType,
			"registered_capital":         data.RegisteredCapital,
			"license_registration_date":  data.LicenseRegistrationDate,
			"tax_registration_date":      data.TaxRegistrationDate,
			"tax_office":                 data.TaxOffice,
			"tax_administrator":          data.TaxAdministrator,
			"tax_administrator_phone":    data.TaxAdministratorPhone,
			"taxpayer_type":              data.TaxpayerType,
			"representative_id":          nil,
			"investor_ids":               "",
		}).Error
		if err != nil {
			tx.Rollback()
			return &ImportError{Row: rowNum, Column: "", Message: fmt.Sprintf("创建客户失败: %v", err)}
		}

		// 获取新创建的ID
		tx.Raw("SELECT last_insert_rowid()").Scan(&customerID)
	}

	// 处理法定代表人
	if data.RepresentativeName != "" && data.RepresentativeIDCard != "" {
		repID, repErr := s.getOrCreateRepresentative(tx, data.RepresentativeName, data.RepresentativeIDCard, rowNum)
		if repErr != nil {
			tx.Rollback()
			return repErr
		}
		// 更新客户的法定代表人
		tx.Table("customers").Where("id = ?", customerID).Update("representative_id", repID)
	}

	// 处理投资人
	if data.InvestorsInfo != "" {
		investors, err := s.parseInvestorsInfo(data.InvestorsInfo, rowNum)
		if err != nil {
			tx.Rollback()
			return err
		}

		var investorIDs []string
		for _, investor := range investors {
			invID, invErr := s.getOrCreateInvestor(tx, investor.Name, investor.IDCard, rowNum)
			if invErr != nil {
				tx.Rollback()
				return invErr
			}
			investorIDs = append(investorIDs, strconv.FormatInt(invID, 10))
		}

		// 更新客户的所有投资人（JSON格式）
		tx.Table("customers").Where("id = ?", customerID).Update("investor_ids", strings.Join(investorIDs, ","))
	}

	// 处理服务人员（验证必须存在）
	if data.ServicePeopleInfo != "" {
		serviceNames := strings.Split(data.ServicePeopleInfo, ",")
		var serviceIDs []string

		for _, name := range serviceNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}

			var person map[string]interface{}
			err := tx.Raw("SELECT id FROM people WHERE name = ? AND type = ?", name, "服务人员").Scan(&person).Error
			if err != nil || person == nil {
				tx.Rollback()
				return &ImportError{Row: rowNum, Column: "服务人员信息", Message: fmt.Sprintf("服务人员 '%s' 不存在，请先创建", name)}
			}
			serviceIDs = append(serviceIDs, strconv.FormatInt(int64(person["id"].(int64)), 10))
		}

		// 更新人员的服务客户列表
		for _, sid := range serviceIDs {
			var existingPerson map[string]interface{}
			tx.Raw("SELECT service_customer_ids FROM people WHERE id = ?", sid).Scan(&existingPerson)
			if existingPerson != nil {
				oldIDs := existingPerson["service_customer_ids"].(string)
				if oldIDs == "" {
					oldIDs = strconv.FormatInt(customerID, 10)
				} else {
					oldIDs = oldIDs + "," + strconv.FormatInt(customerID, 10)
				}
				tx.Table("people").Where("id = ?", sid).Update("service_customer_ids", oldIDs)
			}
		}
	}

	// 处理协议信息
	if data.AgreementsInfo != "" {
		agreements, err := s.parseAgreementsInfo(data.AgreementsInfo, rowNum)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, agreement := range agreements {
			// 验证收费类型
			validFeeTypes := []string{"月度", "季度", "年度"}
			validFeeType := false
			for _, ft := range validFeeTypes {
				if agreement.FeeType == ft {
					validFeeType = true
					break
				}
			}
			if !validFeeType {
				tx.Rollback()
				return &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("收费类型 '%s' 无效，必须是: %s", agreement.FeeType, strings.Join(validFeeTypes, "、"))}
			}

			// 创建协议
			err := tx.Table("agreements").Create(map[string]interface{}{
				"customer_id":   customerID,
				"start_date":    agreement.StartDate,
				"end_date":      agreement.EndDate,
				"fee_type":      agreement.FeeType,
				"fee_amount":    agreement.FeeAmount,
				"status":        "有效",
			}).Error
			if err != nil {
				tx.Rollback()
				return &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("创建协议失败: %v", err)}
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return &ImportError{Row: rowNum, Column: "", Message: fmt.Sprintf("提交事务失败: %v", err)}
	}

	return nil
}

// getOrCreateRepresentative 获取或创建法定代表人
func (s *CustomerImportService) getOrCreateRepresentative(tx *gorm.DB, name, idCard string, rowNum int) (int64, *ImportError) {
	// 先查询是否存在
	var person map[string]interface{}
	err := tx.Raw("SELECT id FROM people WHERE id_card = ?", idCard).Scan(&person).Error
	if err == nil && person != nil {
		return int64(person["id"].(int64)), nil
	}

	// 不存在则创建
	err = tx.Table("people").Create(map[string]interface{}{
		"type":                         "法定代表人",
		"name":                         name,
		"id_card":                      idCard,
		"phone":                        "",
		"password":                     "123456",
		"representative_customer_ids":  "",
		"investor_customer_ids":        "",
		"service_customer_ids":         "",
	}).Error
	if err != nil {
		return 0, &ImportError{Row: rowNum, Column: "法定代表人", Message: fmt.Sprintf("创建法定代表人失败: %v", err)}
	}

	var newID int64
	tx.Raw("SELECT last_insert_rowid()").Scan(&newID)
	return newID, nil
}

// getOrCreateInvestor 获取或创建投资人
func (s *CustomerImportService) getOrCreateInvestor(tx *gorm.DB, name, idCard string, rowNum int) (int64, *ImportError) {
	// 先查询是否存在
	var person map[string]interface{}
	err := tx.Raw("SELECT id FROM people WHERE id_card = ?", idCard).Scan(&person).Error
	if err == nil && person != nil {
		return int64(person["id"].(int64)), nil
	}

	// 不存在则创建
	err = tx.Table("people").Create(map[string]interface{}{
		"type":                         "投资人",
		"name":                         name,
		"id_card":                      idCard,
		"phone":                        "",
		"password":                     "123456",
		"representative_customer_ids":  "",
		"investor_customer_ids":        "",
		"service_customer_ids":         "",
	}).Error
	if err != nil {
		return 0, &ImportError{Row: rowNum, Column: "投资人信息", Message: fmt.Sprintf("创建投资人失败: %v", err)}
	}

	var newID int64
	tx.Raw("SELECT last_insert_rowid()").Scan(&newID)
	return newID, nil
}

// parseInvestorsInfo 解析投资人信息
func (s *CustomerImportService) parseInvestorsInfo(info string, rowNum int) ([]InvestorInfo, *ImportError) {
	if info == "" {
		return nil, nil
	}

	var investors []InvestorInfo
	parts := strings.Split(info, ";")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 格式: 姓名:身份证号:持股比例
		fields := strings.Split(part, ":")
		if len(fields) != 3 {
			return nil, &ImportError{Row: rowNum, Column: "投资人信息", Message: fmt.Sprintf("投资人信息格式错误，应为: 姓名:身份证号:持股比例;..., 当前: %s", part)}
		}

		name := strings.TrimSpace(fields[0])
		idCard := strings.TrimSpace(fields[1])
		shareRatioStr := strings.TrimSpace(fields[2])

		if name == "" || idCard == "" || shareRatioStr == "" {
			return nil, &ImportError{Row: rowNum, Column: "投资人信息", Message: fmt.Sprintf("投资人信息不完整: %s", part)}
		}

		shareRatio, err := ParseFloat(shareRatioStr)
		if err != nil {
			return nil, &ImportError{Row: rowNum, Column: "投资人信息", Message: fmt.Sprintf("持股比例必须是数字: %s", shareRatioStr)}
		}

		investors = append(investors, InvestorInfo{
			Name:       name,
			IDCard:     idCard,
			ShareRatio: shareRatio,
		})
	}

	return investors, nil
}

// parseAgreementsInfo 解析协议信息
func (s *CustomerImportService) parseAgreementsInfo(info string, rowNum int) ([]AgreementInfo, *ImportError) {
	if info == "" {
		return nil, nil
	}

	var agreements []AgreementInfo
	parts := strings.Split(info, "|")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 格式: 有效期起:有效期止:收费类型:收费金额
		fields := strings.Split(part, ":")
		if len(fields) != 4 {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("协议信息格式错误，应为: 有效期起:有效期止:收费类型:收费金额|..., 当前: %s", part)}
		}

		startDateStr := strings.TrimSpace(fields[0])
		endDateStr := strings.TrimSpace(fields[1])
		feeType := strings.TrimSpace(fields[2])
		feeAmountStr := strings.TrimSpace(fields[3])

		startDate, err := ParseDate(startDateStr)
		if err != nil {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("有效期起格式错误: %s", startDateStr)}
		}

		endDate, err := ParseDate(endDateStr)
		if err != nil {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("有效期止格式错误: %s", endDateStr)}
		}

		if feeType == "" {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: "收费类型不能为空"}
		}

		feeAmount, err := ParseFloat(feeAmountStr)
		if err != nil {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: fmt.Sprintf("收费金额必须是数字: %s", feeAmountStr)}
		}

		if feeAmount < 0 {
			return nil, &ImportError{Row: rowNum, Column: "协议信息", Message: "收费金额不能为负数"}
		}

		agreements = append(agreements, AgreementInfo{
			StartDate: startDate,
			EndDate:   endDate,
			FeeType:   feeType,
			FeeAmount: feeAmount,
		})
	}

	return agreements, nil
}
