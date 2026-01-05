package import_export

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelService Excel基础服务
type ExcelService struct {
	file *excelize.File
}

// NewExcelService 创建新的Excel服务
func NewExcelService() *ExcelService {
	return &ExcelService{
		file: excelize.NewFile(),
	}
}

// NewExcelServiceFromFile 从文件创建Excel服务
func NewExcelServiceFromFile(filePath string) (*ExcelService, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	return &ExcelService{
		file: file,
	}, nil
}

// Close 关闭Excel文件
func (s *ExcelService) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

// SaveAs 保存文件到指定路径
func (s *ExcelService) SaveAs(filePath string) error {
	return s.file.SaveAs(filePath)
}

// SetCellValue 设置单元格值
func (s *ExcelService) SetCellValue(sheet, cell string, value interface{}) error {
	return s.file.SetCellValue(sheet, cell, value)
}

// SetCellValueByRowCol 按行列设置单元格值
func (s *ExcelService) SetCellValueByRowCol(sheet string, row, col int, value interface{}) error {
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return fmt.Errorf("坐标转换失败: %w", err)
	}
	return s.file.SetCellValue(sheet, cell, value)
}

// GetCellValue 获取单元格值
func (s *ExcelService) GetCellValue(sheet, cell string) (string, error) {
	return s.file.GetCellValue(sheet, cell)
}

// GetCellValueByRowCol 按行列获取单元格值
func (s *ExcelService) GetCellValueByRowCol(sheet string, row, col int) (string, error) {
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return "", fmt.Errorf("坐标转换失败: %w", err)
	}
	return s.file.GetCellValue(sheet, cell)
}

// GetRow 获取指定行的所有列值
func (s *ExcelService) GetRow(sheet string, row int) ([]string, error) {
	rows, err := s.file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("读取行数据失败: %w", err)
	}
	if row < 1 || row > len(rows) {
		return nil, fmt.Errorf("行号 %d 超出范围", row)
	}
	return rows[row-1], nil
}

// GetRows 获取所有行数据
func (s *ExcelService) GetRows(sheet string) ([][]string, error) {
	rows, err := s.file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("读取行数据失败: %w", err)
	}
	return rows, nil
}

// GetRowCount 获取行数
func (s *ExcelService) GetRowCount(sheet string) (int, error) {
	rows, err := s.file.GetRows(sheet)
	if err != nil {
		return 0, fmt.Errorf("获取行数失败: %w", err)
	}
	return len(rows), nil
}

// CreateSheet 创建工作表
func (s *ExcelService) CreateSheet(sheetName string) (int, error) {
	index, err := s.file.NewSheet(sheetName)
	if err != nil {
		return 0, fmt.Errorf("创建工作表失败: %w", err)
	}
	return index, nil
}

// SetActiveSheet 设置活动工作表
func (s *ExcelService) SetActiveSheet(sheetName string) error {
	index, err := s.file.GetSheetIndex(sheetName)
	if err != nil {
		return fmt.Errorf("获取工作表索引失败: %w", err)
	}
	s.file.SetActiveSheet(index)
	return nil
}

// DeleteSheet 删除工作表
func (s *ExcelService) DeleteSheet(sheetName string) error {
	return s.file.DeleteSheet(sheetName)
}

// SetColWidth 设置列宽
func (s *ExcelService) SetColWidth(sheet, col string, width float64) error {
	return s.file.SetColWidth(sheet, col, col, width)
}

// SetRowHeight 设置行高
func (s *ExcelService) SetRowHeight(sheet string, row int, height float64) error {
	return s.file.SetRowHeight(sheet, row, height)
}

// SetHeaderStyle 设置表头样式
func (s *ExcelService) SetHeaderStyle(sheet, cell string) error {
	style, err := s.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("创建样式失败: %w", err)
	}
	return s.file.SetCellStyle(sheet, cell, cell, style)
}

// SetHeaderStyleByRange 按范围设置表头样式
func (s *ExcelService) SetHeaderStyleByRange(sheet, startCell, endCell string) error {
	style, err := s.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("创建样式失败: %w", err)
	}
	return s.file.SetCellStyle(sheet, startCell, endCell, style)
}

// SetBorderStyle 设置边框样式
func (s *ExcelService) SetBorderStyle(sheet, startCell, endCell string) error {
	style, err := s.file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("创建样式失败: %w", err)
	}
	return s.file.SetCellStyle(sheet, startCell, endCell, style)
}

// SetSheetHeader 设置表头
func (s *ExcelService) SetSheetHeader(sheet string, headers []string) error {
	if len(headers) == 0 {
		return fmt.Errorf("表头不能为空")
	}

	// 设置表头值
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return fmt.Errorf("坐标转换失败: %w", err)
		}
		if err := s.file.SetCellValue(sheet, cell, header); err != nil {
			return fmt.Errorf("设置表头值失败: %w", err)
		}
	}

	// 设置表头样式
	startCell, _ := excelize.CoordinatesToCellName(1, 1)
	endCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	if err := s.SetHeaderStyleByRange(sheet, startCell, endCell); err != nil {
		return err
	}

	// 设置边框
	if err := s.SetBorderStyle(sheet, startCell, endCell); err != nil {
		return err
	}

	// 设置列宽
	for i := range headers {
		col, _ := excelize.CoordinatesToCellName(i+1, 1)
		colName := col[:len(col)-1]
		if err := s.SetColWidth(sheet, colName, 15); err != nil {
			return err
		}
	}

	return nil
}

// WriteRow 写入一行数据
func (s *ExcelService) WriteRow(sheet string, row int, data []interface{}) error {
	for i, value := range data {
		cell, err := excelize.CoordinatesToCellName(i+1, row)
		if err != nil {
			return fmt.Errorf("坐标转换失败: %w", err)
		}
		if err := s.file.SetCellValue(sheet, cell, value); err != nil {
			return fmt.Errorf("设置单元格值失败: %w", err)
		}
	}
	return nil
}

// WriteRows 批量写入多行数据
func (s *ExcelService) WriteRows(sheet string, startRow int, data [][]interface{}) error {
	for i, row := range data {
		if err := s.WriteRow(sheet, startRow+i, row); err != nil {
			return err
		}
	}
	return nil
}

// ParseFloat 解析浮点数
func ParseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

// ParseInt 解析整数
func ParseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// ParseDate 解析日期字符串
func ParseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("日期字符串为空")
	}
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-1-2",
		"2006/1/2",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析日期格式: %s", s)
}

// FormatDate 格式化日期为字符串
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// GetFile 获取底层File对象（用于高级操作）
func (s *ExcelService) GetFile() *excelize.File {
	return s.file
}
