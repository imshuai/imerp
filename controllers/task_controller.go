package controllers

import (
	"erp/config"
	"erp/models"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	if err := config.DB.Create(&task).Error; err != nil {
		ErrorResponse(c, 500, "Failed to create task: "+err.Error())
		return
	}

	// 记录操作日志
	LogOperation(c, "create", "task", &task.ID, task.Title, nil, task)

	SuccessResponse(c, task)
}

// GetTasks 获取任务列表
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	var total int64

	// 获取查询参数
	keyword := c.Query("keyword")
	status := c.Query("status")
	customerID := c.Query("customer_id")

	query := config.DB.Model(&models.Task{}).Preload("Customer")

	// 搜索功能
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 按客户筛选
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// 获取总数
	query.Count(&total)

	// 获取列表
	if err := query.Find(&tasks).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch tasks: "+err.Error())
		return
	}

	SuccessPaginatedResponse(c, total, tasks)
}

// GetTask 获取任务详情
func GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	var task models.Task
	if err := config.DB.Preload("Customer").First(&task, id).Error; err != nil {
		ErrorResponse(c, 404, "Task not found")
		return
	}

	SuccessResponse(c, task)
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		ErrorResponse(c, 404, "Task not found")
		return
	}

	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	taskID := uint(id)

	// 使用 JSON 深拷贝保存旧值
	oldValueJSON, _ := json.Marshal(task)
	var oldValueMap map[string]interface{}
	json.Unmarshal(oldValueJSON, &oldValueMap)

	// 如果状态变为已完成，设置完成时间
	if updateData.Status == models.TaskStatusCompleted && task.Status != models.TaskStatusCompleted {
		now := time.Now()
		updateData.CompletedAt = &now
	}

	// 更新字段
	config.DB.Model(&task).Updates(updateData)

	// 重新获取更新后的数据
	config.DB.Preload("Customer").First(&task, id)

	// 记录操作日志
	LogOperation(c, "update", "task", &taskID, task.Title, oldValueMap, task)

	SuccessResponse(c, task)
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		ErrorResponse(c, 404, "Task not found")
		return
	}

	taskID := uint(id)

	if err := config.DB.Delete(&models.Task{}, id).Error; err != nil {
		ErrorResponse(c, 500, "Failed to delete task: "+err.Error())
		return
	}

	// 记录操作日志
	LogOperation(c, "delete", "task", &taskID, task.Title, task, nil)

	SuccessResponse(c, gin.H{"message": "Task deleted successfully"})
}
