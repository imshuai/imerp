package models

import "time"

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "待处理"  // 待处理
	TaskStatusInProgress TaskStatus = "进行中"  // 进行中
	TaskStatusCompleted  TaskStatus = "已完成"  // 已完成
)

// Task 代办任务
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CustomerID  uint       `json:"customer_id" gorm:"not null"` // 关联客户
	Title       string     `json:"title" gorm:"not null"`       // 任务标题
	Description string     `json:"description"`                 // 任务描述
	Status      TaskStatus `json:"status"`                      // 待处理/进行中/已完成
	DueDate     *time.Time `json:"due_date"`                    // 截止日期
	CompletedAt *time.Time `json:"completed_at"`                // 完成日期
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}
