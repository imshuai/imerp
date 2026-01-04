package controllers

import (
	"erp/config"
	"erp/models"
	"time"

	"github.com/gin-gonic/gin"
)

// OverviewStats 首页概览统计
type OverviewStats struct {
	CustomerCount      int64   `json:"customer_count"`
	PendingTaskCount   int64   `json:"pending_task_count"`
	ActiveAgreementCount int64 `json:"active_agreement_count"`
	MonthlyPayment     float64 `json:"monthly_payment"`
	YearlyPayment      float64 `json:"yearly_payment"`
}

// TaskStats 任务统计
type TaskStats struct {
	Pending      int64 `json:"pending"`
	InProgress   int64 `json:"in_progress"`
	Completed    int64 `json:"completed"`
}

// PaymentStats 收款统计
type PaymentStats struct {
	TotalAmount   float64 `json:"total_amount"`
	Count         int64   `json:"count"`
}

// GetOverview 获取首页概览统计
func GetOverview(c *gin.Context) {
	var stats OverviewStats

	// 客户总数
	config.DB.Model(&models.Customer{}).Count(&stats.CustomerCount)

	// 待办任务数
	config.DB.Model(&models.Task{}).Where("status != ?", "completed").Count(&stats.PendingTaskCount)

	// 有效协议数
	config.DB.Model(&models.Agreement{}).Where("status = ?", "active").Count(&stats.ActiveAgreementCount)

	// 本月收款
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	config.DB.Model(&models.Payment{}).
		Where("payment_date >= ?", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.MonthlyPayment)

	// 本年收款
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	config.DB.Model(&models.Payment{}).
		Where("payment_date >= ?", startOfYear).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.YearlyPayment)

	SuccessResponse(c, stats)
}

// GetTaskStats 获取任务统计
func GetTaskStats(c *gin.Context) {
	var stats TaskStats

	config.DB.Model(&models.Task{}).Where("status = ?", "pending").Count(&stats.Pending)
	config.DB.Model(&models.Task{}).Where("status = ?", "in_progress").Count(&stats.InProgress)
	config.DB.Model(&models.Task{}).Where("status = ?", "completed").Count(&stats.Completed)

	SuccessResponse(c, stats)
}

// GetPaymentStats 获取收款统计
func GetPaymentStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var stats PaymentStats

	query := config.DB.Model(&models.Payment{})

	// 按时间范围筛选
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("payment_date >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("payment_date <= ?", t)
		}
	}

	// 统计
	query.Count(&stats.Count)
	query.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount)

	SuccessResponse(c, stats)
}
