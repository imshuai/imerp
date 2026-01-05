package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SaveUploadedFile 保存上传的文件到临时目录
func SaveUploadedFile(c *gin.Context, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("获取上传文件失败: %w", err)
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := []string{".xlsx", ".xls"}
	valid := false
	for _, ve := range validExts {
		if ext == ve {
			valid = true
			break
		}
	}
	if !valid {
		return "", fmt.Errorf("不支持的文件格式，请上传Excel文件（.xlsx或.xls）")
	}

	// 创建临时目录
	tempDir := os.TempDir()
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), time.Now().UnixNano(), ext)
	filePath := filepath.Join(tempDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return filePath, nil
}

// CleanupTempFile 清理临时文件
func CleanupTempFile(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
	}
}

// GetTimestamp 获取当前时间戳字符串
func GetTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// SanitizeFilename 清理文件名，移除不安全字符
func SanitizeFilename(filename string) string {
	// 替换不安全的字符
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := filename
	for _, char := range unsafe {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	return sanitized
}

// ValidateRequiredFields 验证必填字段
func ValidateRequiredFields(data map[string]interface{}, requiredFields []string) error {
	for _, field := range requiredFields {
		val, exists := data[field]
		if !exists || val == nil || val == "" {
			return fmt.Errorf("字段 %s 不能为空", field)
		}
	}
	return nil
}

// TruncateString 截断字符串到指定长度
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// JoinNonEmpty 连接非空字符串
func JoinNonEmpty(strs []string, sep string) string {
	var nonEmpty []string
	for _, s := range strs {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, sep)
}
