package main

import (
	"fmt"
	"erp/config"
	"erp/models"
)

func main() {
	config.InitLogger("development")
	defer config.SyncLogger()

	if err := config.InitDatabase(); err != nil {
		panic(err)
	}

	var count int64

	// 测试查询
	config.DB.Model(&models.Task{}).Where("status = ?", models.TaskStatusPending).Count(&count)
	fmt.Println("使用常量查询待处理任务:", count)

	config.DB.Model(&models.Task{}).Where("status = ?", "待处理").Count(&count)
	fmt.Println("直接使用字符串查询:", count)

	// 列出所有任务
	var tasks []models.Task
	config.DB.Find(&tasks)
	for _, t := range tasks {
		fmt.Printf("ID: %d, Title: %s, Status: %s\n", t.ID, t.Title, t.Status)
	}
}
