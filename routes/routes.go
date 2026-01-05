package routes

import (
	"erp/controllers"
	"erp/config"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// 获取数据库连接
	db := config.DB

	// 创建导入导出控制器
	importExportCtrl := controllers.NewImportExportController(db)

	// API路由组
	api := r.Group("/api")
	{
		// 人员管理路由
		people := api.Group("/people")
		{
			people.GET("", controllers.GetPeople)
			people.POST("", controllers.CreatePerson)
			people.GET("/:id", controllers.GetPerson)
			people.PUT("/:id", controllers.UpdatePerson)
			people.DELETE("/:id", controllers.DeletePerson)
			people.GET("/:id/customers", controllers.GetPersonCustomers)
		}

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

		// 导入导出路由
		templates := api.Group("/templates")
		{
			templates.GET("/:type", importExportCtrl.DownloadTemplate)
		}

		importAPI := api.Group("/import")
		{
			importAPI.POST("/people", importExportCtrl.ImportPeople)
			importAPI.POST("/customers", importExportCtrl.ImportCustomers)
		}

		export := api.Group("/export")
		{
			export.GET("/people", importExportCtrl.ExportPeople)
			export.GET("/customers", importExportCtrl.ExportCustomers)
		}
	}
}
