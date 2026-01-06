package import_export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ImportStrategy 导入冲突策略
type ImportStrategy string

const (
	StrategySkip      ImportStrategy = "skip"       // 跳过已存在的记录
	StrategyUpdate    ImportStrategy = "update"     // 更新已存在的记录
	StrategyCreateNew ImportStrategy = "create_new" // 创建新记录（修改标识）
)

// ImportError 导入错误
type ImportError struct {
	Row     int    `json:"row"`     // 行号
	Column  string `json:"column"`  // 列名
	Message string `json:"message"` // 错误信息
}

// ImportResult 导入结果
type ImportResult struct {
	Total   int           `json:"total"`   // 总行数
	Success int           `json:"success"` // 成功数
	Failed  int           `json:"failed"`  // 失败数
	Errors  []ImportError `json:"errors"`  // 错误详情
}

// PeopleImportService 人员导入服务
type PeopleImportService struct {
	db *gorm.DB
}

// NewPeopleImportService 创建人员导入服务
func NewPeopleImportService(db *gorm.DB) *PeopleImportService {
	return &PeopleImportService{db: db}
}

// ImportPeopleFromExcel 从Excel导入人员
func (s *PeopleImportService) ImportPeopleFromExcel(filePath string, strategy ImportStrategy) (*ImportResult, error) {
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
	requiredColumns := []string{"姓名", "电话", "身份证号"}
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
		person, parseErr := s.parsePersonRow(row, colIndex, rowNum)
		if parseErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, *parseErr)
			continue
		}

		// 导入人员
		importErr := s.importPerson(person, strategy, rowNum, result)
		if importErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, *importErr)
		} else {
			result.Success++
		}
	}

	return result, nil
}

// PersonRowData 人员行数据
type PersonRowData struct {
	Name            string
	IsServicePerson bool
	Phone           string
	IDCard          string
	Password        string
}

// parsePersonRow 解析人员行数据
func (s *PeopleImportService) parsePersonRow(row []string, colIndex map[string]int, rowNum int) (*PersonRowData, *ImportError) {
	getCell := func(colName string) string {
		idx, exists := colIndex[colName]
		if !exists || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	data := &PersonRowData{
		Name:            getCell("姓名"),
		Phone:           getCell("电话"),
		IDCard:          getCell("身份证号"),
		Password:        getCell("登录密码"),
		IsServicePerson: false, // 默认值
	}

	// 解析"是否服务人员"列（如果存在）
	if isServicePersonCell := getCell("是否服务人员"); isServicePersonCell != "" {
		data.IsServicePerson = isServicePersonCell == "是" || isServicePersonCell == "true" || isServicePersonCell == "1" || isServicePersonCell == "TRUE"
	}

	// 验证必填字段
	if data.Name == "" {
		return nil, &ImportError{Row: rowNum, Column: "姓名", Message: "姓名不能为空"}
	}
	if data.Phone == "" {
		return nil, &ImportError{Row: rowNum, Column: "电话", Message: "电话不能为空"}
	}
	if data.IDCard == "" {
		return nil, &ImportError{Row: rowNum, Column: "身份证号", Message: "身份证号不能为空"}
	}

	// 如果没有密码，设置默认密码
	if data.Password == "" {
		data.Password = "123456"
	}

	return data, nil
}

// importPerson 导入单个人员
func (s *PeopleImportService) importPerson(data *PersonRowData, strategy ImportStrategy, rowNum int, result *ImportResult) *ImportError {
	// 查询是否已存在
	var existingPerson map[string]interface{}
	err := s.db.Raw("SELECT id FROM people WHERE id_card = ?", data.IDCard).Scan(&existingPerson).Error
	isConflict := err == nil && existingPerson != nil

	if isConflict {
		switch strategy {
		case StrategySkip:
			// 跳过已存在的记录
			return nil
		case StrategyUpdate:
			// 更新已存在的记录
			err := s.db.Table("people").
				Where("id_card = ?", data.IDCard).
				Updates(map[string]interface{}{
					"name":              data.Name,
					"is_service_person": data.IsServicePerson,
					"phone":             data.Phone,
					"password":          data.Password,
				}).Error
			if err != nil {
				return &ImportError{Row: rowNum, Column: "", Message: fmt.Sprintf("更新失败: %v", err)}
			}
			return nil
		case StrategyCreateNew:
			// 修改身份证号后创建
			suffix := 1
			newIDCard := data.IDCard
			for {
				var count int64
				s.db.Table("people").Where("id_card = ?", newIDCard).Count(&count)
				if count == 0 {
					break
				}
				suffix++
				newIDCard = fmt.Sprintf("%s_%d", data.IDCard, suffix)
			}
			data.IDCard = newIDCard
		}
	}

	// 插入新记录
	err = s.db.Table("people").Create(map[string]interface{}{
		"name":                        data.Name,
		"is_service_person":           data.IsServicePerson,
		"phone":                       data.Phone,
		"id_card":                     data.IDCard,
		"password":                    data.Password,
		"representative_customer_ids": "",
		"investor_customer_ids":       "",
		"service_customer_ids":        "",
	}).Error

	if err != nil {
		// 检查是否是唯一约束冲突
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate") {
			return &ImportError{
				Row: rowNum, Column: "身份证号",
				Message: fmt.Sprintf("身份证号 %s 已存在", data.IDCard),
			}
		}
		return &ImportError{Row: rowNum, Column: "", Message: fmt.Sprintf("创建失败: %v", err)}
	}

	return nil
}

// GetFileExt 获取文件扩展名
func GetFileExt(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return filename[idx:]
}

// ValidateExcelFile 验证是否为有效的Excel文件
func ValidateExcelFile(filename string, content []byte) error {
	ext := GetFileExt(filename)
	validExts := []string{".xlsx", ".xls"}
	valid := false
	for _, ve := range validExts {
		if ext == ve {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("不支持的文件格式，请上传Excel文件（.xlsx或.xls）")
	}

	// 验证文件内容
	_, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("无效的Excel文件: %w", err)
	}

	return nil
}
