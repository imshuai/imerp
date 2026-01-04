package models

import "time"

// Task 代办任务
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CustomerID  uint       `json:"customer_id" gorm:"not null"` // 关联客户
	Title       string     `json:"title" gorm:"not null"`       // 任务标题
	Description string     `json:"description"`                 // 任务描述
	Status      string     `json:"status"`                      // pending/in_progress/completed
	DueDate     *time.Time `json:"due_date"`                    // 截止日期
	CompletedAt *time.Time `json:"completed_at"`                // 完成日期
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}
