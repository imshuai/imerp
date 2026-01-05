package import_export

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ExportService 导出服务
type ExportService struct {
	db *gorm.DB
}

// NewExportService 创建导出服务
func NewExportService(db *gorm.DB) *ExportService {
	return &ExportService{db: db}
}

// ExportPeopleToExcel 导出人员到Excel
func (s *ExportService) ExportPeopleToExcel() ([]byte, string, error) {
	excelService := NewExcelService()
	defer excelService.Close()

	sheetName := "人员列表"
	excelService.CreateSheet(sheetName)
	excelService.SetActiveSheet(sheetName)

	// 设置表头
	headers := []string{"姓名", "类型", "电话", "身份证号", "登录密码"}
	if err := excelService.SetSheetHeader(sheetName, headers); err != nil {
		return nil, "", fmt.Errorf("设置表头失败: %w", err)
	}

	// 查询所有人员
	var people []map[string]interface{}
	err := s.db.Table("people").
		Select("id, type, name, phone, id_card, password").
		Order("id ASC").
		Find(&people).Error
	if err != nil {
		return nil, "", fmt.Errorf("查询人员失败: %w", err)
	}

	// 写入数据
	data := make([][]interface{}, len(people))
	for i, person := range people {
		data[i] = []interface{}{
			person["name"],
			person["type"],
			person["phone"],
			person["id_card"],
			person["password"],
		}
	}

	if err := excelService.WriteRows(sheetName, 2, data); err != nil {
		return nil, "", fmt.Errorf("写入数据失败: %w", err)
	}

	// 设置数据边框
	if len(data) > 0 {
		startCell, _ := excelize.CoordinatesToCellName(1, 2)
		endCell, _ := excelize.CoordinatesToCellName(5, 2+len(data)-1)
		excelService.SetBorderStyle(sheetName, startCell, endCell)
	}

	// 保存到临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("人员导出_%s.xlsx", time.Now().Format("20060102_150405")))
	if err := excelService.SaveAs(tempFile); err != nil {
		return nil, "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 读取文件内容
	content, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 删除临时文件
	os.Remove(tempFile)

	filename := fmt.Sprintf("人员导出_%s.xlsx", time.Now().Format("20060102_150405"))
	return content, filename, nil
}

// ExportCustomersToExcel 导出客户到Excel
func (s *ExportService) ExportCustomersToExcel() ([]byte, string, error) {
	excelService := NewExcelService()
	defer excelService.Close()

	sheetName := "客户列表"
	excelService.CreateSheet(sheetName)
	excelService.SetActiveSheet(sheetName)

	// 设置表头
	headers := []string{
		"公司名称", "联系电话", "地址", "税号", "客户类型", "注册资本",
		"法定代表人", "投资人", "服务人员", "协议信息",
	}
	if err := excelService.SetSheetHeader(sheetName, headers); err != nil {
		return nil, "", fmt.Errorf("设置表头失败: %w", err)
	}

	// 查询所有客户
	var customers []map[string]interface{}
	err := s.db.Table("customers").
		Select("id, name, phone, address, tax_number, type, registered_capital, representative_id, investor_ids").
		Order("id ASC").
		Find(&customers).Error
	if err != nil {
		return nil, "", fmt.Errorf("查询客户失败: %w", err)
	}

	// 写入数据
	data := make([][]interface{}, len(customers))
	for i, customer := range customers {
		// 获取法定代表人信息
		repName := ""
		if customer["representative_id"] != nil {
			var rep map[string]interface{}
			s.db.Table("people").
				Select("name").
				Where("id = ?", customer["representative_id"]).
				Scan(&rep)
			if rep != nil {
				repName = rep["name"].(string)
			}
		}

		// 获取投资人信息
		investorsInfo := ""
		if customer["investor_ids"] != nil && customer["investor_ids"].(string) != "" {
			investorIDs := customer["investor_ids"].(string)
			var investors []map[string]interface{}
			s.db.Raw("SELECT id, name, id_card FROM people WHERE id IN ("+investorIDs+")").Scan(&investors)
			if len(investors) > 0 {
				var investorStrs []string
				for _, inv := range investors {
					investorStrs = append(investorStrs, fmt.Sprintf("%s:%s", inv["name"], inv["id_card"]))
				}
				investorsInfo = fmt.Sprintf("%s", investorStrs)
			}
		}

		// 获取服务人员信息
		serviceNames := ""
		customerID := customer["id"].(int64)
		var services []map[string]interface{}
		s.db.Raw(`
			SELECT p.name
			FROM people p
			WHERE p.type = '服务人员'
			AND p.id IN (
				SELECT CAST(value AS INTEGER)
				FROM json_each('["' || REPLACE(p.service_customer_ids, ',', '","') || '"]')
				WHERE p.service_customer_ids != ''
			)
			AND ',' || p.service_customer_ids || ',' LIKE '%,' || ? || ',%'
		`, customerID).Scan(&services)
		if len(services) > 0 {
			var names []string
			for _, svc := range services {
				names = append(names, svc["name"].(string))
			}
			serviceNames = fmt.Sprintf("%s", names)
		}

		// 获取协议信息
		agreementsInfo := ""
		var agreements []map[string]interface{}
		s.db.Table("agreements").
			Select("start_date, end_date, fee_type, fee_amount").
			Where("customer_id = ?", customerID).
			Order("start_date ASC").
			Scan(&agreements)
		if len(agreements) > 0 {
			var agrStrs []string
			for _, agr := range agreements {
				startDate := agr["start_date"].(time.Time).Format("2006-01-02")
				endDate := agr["end_date"].(time.Time).Format("2006-01-02")
				feeType := agr["fee_type"].(string)
				feeAmount := agr["fee_amount"]
				agrStrs = append(agrStrs, fmt.Sprintf("%s:%s:%s:%.0f", startDate, endDate, feeType, feeAmount))
			}
			agreementsInfo = fmt.Sprintf("%s", agrStrs)
		}

		data[i] = []interface{}{
			customer["name"],
			customer["phone"],
			customer["address"],
			customer["tax_number"],
			customer["type"],
			customer["registered_capital"],
			repName,
			investorsInfo,
			serviceNames,
			agreementsInfo,
		}
	}

	if err := excelService.WriteRows(sheetName, 2, data); err != nil {
		return nil, "", fmt.Errorf("写入数据失败: %w", err)
	}

	// 设置数据边框
	if len(data) > 0 {
		startCell, _ := excelize.CoordinatesToCellName(1, 2)
		endCell, _ := excelize.CoordinatesToCellName(11, 2+len(data)-1)
		excelService.SetBorderStyle(sheetName, startCell, endCell)
	}

	// 调整列宽
	excelService.SetColWidth(sheetName, "A", "A", 25)  // 公司名称
	excelService.SetColWidth(sheetName, "I", "I", 30)  // 投资人
	excelService.SetColWidth(sheetName, "J", "J", 20)  // 服务人员
	excelService.SetColWidth(sheetName, "K", "K", 40)  // 协议信息

	// 保存到临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("客户导出_%s.xlsx", time.Now().Format("20060102_150405")))
	if err := excelService.SaveAs(tempFile); err != nil {
		return nil, "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 读取文件内容
	content, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 删除临时文件
	os.Remove(tempFile)

	filename := fmt.Sprintf("客户导出_%s.xlsx", time.Now().Format("20060102_150405"))
	return content, filename, nil
}

// FormatInt64 格式化int64为字符串
func FormatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}

// FormatFloat 格式化float64为字符串
func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
