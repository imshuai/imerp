package routes

import (
	"erp/controllers"
	"erp/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGinEngine 创建并配置Gin引擎
func SetupGinEngine(env string) *gin.Engine {
	// 设置Gin运行模式
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 自定义日志中间件 - 统一使用zap记录所有HTTP请求
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		config.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.String("error", param.ErrorMessage),
		)
		return ""
	}))

	// 恢复中间件 - 捕获panic并记录日志
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		config.Error("Panic recovered",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Any("error", recovered),
		)
		c.JSON(500, gin.H{"code": 1, "message": "Internal Server Error"})
	}))

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 设置路由
	SetupRoutes(r)

	return r
}

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
