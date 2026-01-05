package import_export

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// TemplateService 模板服务
type TemplateService struct {
	excelService *ExcelService
}

// NewTemplateService 创建模板服务
func NewTemplateService() *TemplateService {
	return &TemplateService{
		excelService: NewExcelService(),
	}
}

// GeneratePeopleTemplate 生成人员导入模板
func (s *TemplateService) GeneratePeopleTemplate() ([]byte, string, error) {
	defer s.excelService.Close()

	sheetName := "人员导入"
	s.excelService.CreateSheet(sheetName)
	s.excelService.SetActiveSheet(sheetName)

	// 设置表头
	headers := []string{
		"姓名", "类型", "电话", "身份证号", "登录密码",
	}
	if err := s.excelService.SetSheetHeader(sheetName, headers); err != nil {
		return nil, "", fmt.Errorf("设置表头失败: %w", err)
	}

	// 添加示例数据
	sampleData := [][]interface{}{
		{"张三", "法定代表人", "13800138000", "110101199001011234", "abc123"},
		{"李四", "投资人", "13900139000", "110101199002021234", "def456"},
		{"王五", "服务人员", "13700137000", "110101199003031234", "ghi789"},
		{"赵六", "混合角色", "13600136000", "110101199004041234", "jkl012"},
	}
	if err := s.excelService.WriteRows(sheetName, 2, sampleData); err != nil {
		return nil, "", fmt.Errorf("写入示例数据失败: %w", err)
	}

	// 为示例数据设置边框
	lastRow := 1 + len(sampleData)
	endCell1, _ := excelize.CoordinatesToCellName(1, lastRow)
	endCell2, _ := excelize.CoordinatesToCellName(5, lastRow)
	startCell2, _ := excelize.CoordinatesToCellName(1, 2)
	if err := s.excelService.SetBorderStyle(sheetName, startCell2, endCell2); err != nil {
		return nil, "", fmt.Errorf("设置边框失败: %w", err)
	}

	// 保存到临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "people_import_template.xlsx")
	if err := s.excelService.SaveAs(tempFile); err != nil {
		return nil, "", fmt.Errorf("保存模板失败: %w", err)
	}

	// 读取文件内容
	content, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 删除临时文件
	os.Remove(tempFile)

	return content, "人员导入模板.xlsx", nil
}

// GenerateCustomersTemplate 生成客户导入模板
func (s *TemplateService) GenerateCustomersTemplate() ([]byte, string, error) {
	defer s.excelService.Close()

	sheetName := "客户导入"
	s.excelService.CreateSheet(sheetName)
	s.excelService.SetActiveSheet(sheetName)

	// 设置表头
	headers := []string{
		"公司名称", "联系电话", "地址", "税号", "客户类型", "注册资本",
		"法定代表人姓名", "法定代表人身份证", "投资人信息", "服务人员信息", "协议信息",
	}
	if err := s.excelService.SetSheetHeader(sheetName, headers); err != nil {
		return nil, "", fmt.Errorf("设置表头失败: %w", err)
	}

	// 调整列宽（有些列需要更宽）
	s.excelService.SetColWidth(sheetName, "I", "I", 30)  // 投资人信息
	s.excelService.SetColWidth(sheetName, "J", "J", 20)  // 服务人员信息
	s.excelService.SetColWidth(sheetName, "K", "K", 40)  // 协议信息

	// 添加示例数据
	sampleData := [][]interface{}{
		{
			"某某科技有限公司", "13800138000", "北京市朝阳区某某街道123号",
			"91110000MA001234XX", "有限公司", 1000000,
			"张三", "110101199001011234",
			"李四:110101199002021234:51;王五:110101199003031234:49",
			"赵六,钱七",
			"2024-01-01:2024-12-31:月度:500",
		},
		{
			"某某商贸中心", "13900139000", "上海市浦东新区某某路456号",
			"91310000MA005678XX", "个体工商户", 50000,
			"", "", "",
			"孙八:110101199005051234:100",
			"赵六",
			"2024-01-01:2024-12-31:月度:300|2025-01-01:2025-12-31:月度:350",
		},
	}
	if err := s.excelService.WriteRows(sheetName, 2, sampleData); err != nil {
		return nil, "", fmt.Errorf("写入示例数据失败: %w", err)
	}

	// 为示例数据设置边框
	lastRow := 1 + len(sampleData)
	endCell1, _ := excelize.CoordinatesToCellName(1, lastRow)
	endCell2, _ := excelize.CoordinatesToCellName(11, lastRow)
	startCell2, _ := excelize.CoordinatesToCellName(1, 2)
	if err := s.excelService.SetBorderStyle(sheetName, startCell2, endCell2); err != nil {
		return nil, "", fmt.Errorf("设置边框失败: %w", err)
	}

	// 添加说明工作表
	s.excelService.CreateSheet("填写说明")
	s.excelService.SetActiveSheet("填写说明")

	instructions := [][]interface{}{
		{"字段", "说明", "示例", "是否必填"},
		{"公司名称", "企业的完整名称", "某某科技有限公司", "是"},
		{"联系电话", "企业联系电话", "13800138000", "否"},
		{"地址", "企业注册地址", "北京市朝阳区某某街道123号", "否"},
		{"税号", "纳税人识别号", "91110000MA001234XX", "是"},
		{"客户类型", "有限公司/个人独资企业/合伙企业/个体工商户", "有限公司", "是"},
		{"注册资本", "注册资本金额（数字）", "1000000", "否"},
		{"法定代表人姓名", "法定代表人姓名", "张三", "否"},
		{"法定代表人身份证", "法定代表人身份证号", "110101199001011234", "否"},
		{"投资人信息", "格式：姓名:身份证号:持股比例;姓名:身份证号:持股比例", "李四:110101199002021234:51;王五:110101199003031234:49", "否"},
		{"", "多个投资人用分号;分隔", "", ""},
		{"服务人员信息", "格式：姓名,姓名,姓名（逗号分隔）", "赵六,钱七", "否"},
		{"", "服务人员必须已存在于系统中", "", ""},
		{"协议信息", "格式：有效期起:有效期止:收费类型:收费金额", "2024-01-01:2024-12-31:月度:500", "否"},
		{"", "多个协议用竖线|分隔", "2024-01-01:2024-12-31:月度:500|2025-01-01:2025-12-31:月度:600", ""},
		{"", "收费类型：月度/季度/年度", "", ""},
	}

	for i, row := range instructions {
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			s.excelService.SetCellValue("填写说明", cell, val)
		}
	}

	// 设置说明表头样式
	s.excelService.SetHeaderStyleByRange("填写说明", "A1", "D1")
	s.excelService.SetBorderStyle("填写说明", "A1", fmt.Sprintf("D%d", len(instructions)))
	s.excelService.SetColWidth("填写说明", "A", "A", 20)
	s.excelService.SetColWidth("填写说明", "B", "B", 40)
	s.excelService.SetColWidth("填写说明", "C", "C", 30)
	s.excelService.SetColWidth("填写说明", "D", "D", 12)

	// 保存到临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "customers_import_template.xlsx")
	if err := s.excelService.SaveAs(tempFile); err != nil {
		return nil, "", fmt.Errorf("保存模板失败: %w", err)
	}

	// 读取文件内容
	content, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 删除临时文件
	os.Remove(tempFile)

	return content, "客户导入模板.xlsx", nil
}

// DownloadTemplateResponse 下载模板的响应处理
func (s *TemplateService) DownloadTemplateResponse(c *gin.Context, templateType string) {
	var content []byte
	var filename string
	var err error

	switch templateType {
	case "people":
		content, filename, err = s.GeneratePeopleTemplate()
	case "customers":
		content, filename, err = s.GenerateCustomersTemplate()
	default:
		c.JSON(400, gin.H{"code": 1, "message": "不支持的模板类型"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"code": 1, "message": fmt.Sprintf("生成模板失败: %v", err)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}
