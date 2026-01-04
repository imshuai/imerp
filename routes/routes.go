package routes

import (
	"erp/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// API路由组
	api := r.Group("/api")
	{
		// 客户管理路由
		customers := api.Group("/customers")
		{
			customers.GET("", controllers.GetCustomers)
			customers.POST("", controllers.CreateCustomer)
			customers.GET("/:id", controllers.GetCustomer)
			customers.PUT("/:id", controllers.UpdateCustomer)
			customers.DELETE("/:id", controllers.DeleteCustomer)
			customers.GET("/:id/tasks", controllers.GetCustomerTasks)
			customers.GET("/:id/payments", controllers.GetCustomerPayments)
		}

		// 任务管理路由
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.GetTasks)
			tasks.POST("", controllers.CreateTask)
			tasks.GET("/:id", controllers.GetTask)
			tasks.PUT("/:id", controllers.UpdateTask)
			tasks.DELETE("/:id", controllers.DeleteTask)
		}

		// 协议管理路由
		agreements := api.Group("/agreements")
		{
			agreements.GET("", controllers.GetAgreements)
			agreements.POST("", controllers.CreateAgreement)
			agreements.GET("/:id", controllers.GetAgreement)
			agreements.PUT("/:id", controllers.UpdateAgreement)
			agreements.DELETE("/:id", controllers.DeleteAgreement)
		}

		// 收款管理路由
		payments := api.Group("/payments")
		{
			payments.GET("", controllers.GetPayments)
			payments.POST("", controllers.CreatePayment)
			payments.GET("/:id", controllers.GetPayment)
			payments.PUT("/:id", controllers.UpdatePayment)
			payments.DELETE("/:id", controllers.DeletePayment)
		}

		// 统计分析路由
		statistics := api.Group("/statistics")
		{
			statistics.GET("/overview", controllers.GetOverview)
			statistics.GET("/tasks", controllers.GetTaskStats)
			statistics.GET("/payments", controllers.GetPaymentStats)
		}
	}
}
