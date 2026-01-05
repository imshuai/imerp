package controllers

import (
	"erp/services/import_export"
	"erp/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImportExportController 导入导出控制器
type ImportExportController struct {
	db                *gorm.DB
	templateService   *import_export.TemplateService
	peopleImportSvc   *import_export.PeopleImportService
	customerImportSvc *import_export.CustomerImportService
	exportService     *import_export.ExportService
}

// NewImportExportController 创建导入导出控制器
func NewImportExportController(db *gorm.DB) *ImportExportController {
	return &ImportExportController{
		db:                db,
		templateService:   import_export.NewTemplateService(),
		peopleImportSvc:   import_export.NewPeopleImportService(db),
		customerImportSvc: import_export.NewCustomerImportService(db),
		exportService:     import_export.NewExportService(db),
	}
}

// DownloadTemplate 下载导入模板
// @Summary 下载导入模板
// @Description 下载人员或客户的Excel导入模板
// @Tags 导入导出
// @Param type path string true "模板类型 (people/customers)"
// @Success 200 {file} file "Excel文件"
// @Failure 400 {object} map[string]interface{} "错误的模板类型"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/templates/{type} [get]
func (ctrl *ImportExportController) DownloadTemplate(c *gin.Context) {
	templateType := c.Param("type")
	ctrl.templateService.DownloadTemplateResponse(c, templateType)
}

// ImportPeople 导入人员
// @Summary 导入人员
// @Description 从Excel文件导入人员数据
// @Tags 导入导出
// @Param file formData file true "Excel文件"
// @Param strategy formData string false "冲突策略: skip/update/create_new (默认: skip)"
// @Success 200 {object} map[string]interface{} "导入结果"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/import/people [post]
func (ctrl *ImportExportController) ImportPeople(c *gin.Context) {
	// 保存上传的文件
	filePath, err := utils.SaveUploadedFile(c, "file")
	if err != nil {
		c.JSON(400, gin.H{"code": 1, "message": err.Error()})
		return
	}
	defer utils.CleanupTempFile(filePath)

	// 获取冲突策略
	strategy := c.PostForm("strategy")
	if strategy == "" {
		strategy = "skip"
	}

	// 验证策略
	validStrategies := []import_export.ImportStrategy{
		import_export.StrategySkip,
		import_export.StrategyUpdate,
		import_export.StrategyCreateNew,
	}
	valid := false
	for _, vs := range validStrategies {
		if import_export.ImportStrategy(strategy) == vs {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(400, gin.H{"code": 1, "message": "无效的冲突策略，必须是: skip, update, create_new"})
		return
	}

	// 执行导入
	result, err := ctrl.peopleImportSvc.ImportPeopleFromExcel(filePath, import_export.ImportStrategy(strategy))
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "message": fmt.Sprintf("导入失败: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"message": "导入完成",
		"data": result,
	})
}

// ImportCustomers 导入客户
// @Summary 导入客户
// @Description 从Excel文件导入客户数据（包含关联人员和协议）
// @Tags 导入导出
// @Param file formData file true "Excel文件"
// @Param strategy formData string false "冲突策略: skip/update/create_new (默认: skip)"
// @Success 200 {object} map[string]interface{} "导入结果"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/import/customers [post]
func (ctrl *ImportExportController) ImportCustomers(c *gin.Context) {
	// 保存上传的文件
	filePath, err := utils.SaveUploadedFile(c, "file")
	if err != nil {
		c.JSON(400, gin.H{"code": 1, "message": err.Error()})
		return
	}
	defer utils.CleanupTempFile(filePath)

	// 获取冲突策略
	strategy := c.PostForm("strategy")
	if strategy == "" {
		strategy = "skip"
	}

	// 验证策略
	validStrategies := []import_export.ImportStrategy{
		import_export.StrategySkip,
		import_export.StrategyUpdate,
		import_export.StrategyCreateNew,
	}
	valid := false
	for _, vs := range validStrategies {
		if import_export.ImportStrategy(strategy) == vs {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(400, gin.H{"code": 1, "message": "无效的冲突策略，必须是: skip, update, create_new"})
		return
	}

	// 执行导入
	result, err := ctrl.customerImportSvc.ImportCustomersFromExcel(filePath, import_export.ImportStrategy(strategy))
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "message": fmt.Sprintf("导入失败: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"message": "导入完成",
		"data": result,
	})
}

// ExportPeople 导出人员
// @Summary 导出人员
// @Description 将所有人员数据导出为Excel文件
// @Tags 导入导出
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/export/people [get]
func (ctrl *ImportExportController) ExportPeople(c *gin.Context) {
	content, filename, err := ctrl.exportService.ExportPeopleToExcel()
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "message": fmt.Sprintf("导出失败: %v", err)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}

// ExportCustomers 导出客户
// @Summary 导出客户
// @Description 将所有客户数据（包含关联人员和协议）导出为Excel文件
// @Tags 导入导出
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/export/customers [get]
func (ctrl *ImportExportController) ExportCustomers(c *gin.Context) {
	content, filename, err := ctrl.exportService.ExportCustomersToExcel()
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "message": fmt.Sprintf("导出失败: %v", err)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}
