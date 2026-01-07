package routes

import (
	"erp/controllers"
	"erp/config"
	"erp/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupGinEngine 创建并配置Gin引擎
func SetupGinEngine() *gin.Engine {
	// 设置Gin运行模式
	if config.AppConfig != nil {
		mode := config.AppConfig.Server.Mode
		if mode == "" {
			if config.AppConfig.Server.Env == "production" {
				gin.SetMode(gin.ReleaseMode)
			}
		} else {
			gin.SetMode(mode)
		}
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
		// 认证路由（无需token）
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.GET("/users", controllers.GetLoginUsers) // 获取可登录用户列表（公开）
		}

		// 公开路由 - 获取人员列表（用于登录页面）
		api.GET("/people", controllers.GetPeople)

		// 枚举值获取路由（无需认证，供前端使用）
		api.GET("/customers/types", controllers.GetCustomerTypes)
		api.GET("/customers/credit-ratings", controllers.GetCreditRatings)
		api.GET("/bank-accounts/types", controllers.GetAccountTypes)

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authenticated.GET("/user/me", controllers.GetCurrentUser)
			authenticated.POST("/user/change-password", controllers.ChangePassword)

			// 管理员路由（仅超级管理员）
			admin := authenticated.Group("/admin")
			admin.Use(middleware.RequireSuperAdmin())
			{
				admin.DELETE("/audit-logs/:id", controllers.DeleteAuditLog)
				admin.POST("/audit-logs/clear", controllers.ClearAuditLogs)
			}

			// 审计日志（所有登录用户可查看）
			authenticated.GET("/audit-logs", controllers.GetAuditLogs)

			// 人员管理路由（所有认证用户可操作）
			people := authenticated.Group("/people")
			{
				people.GET("/:id", controllers.GetPerson)
				people.POST("", controllers.CreatePerson)
				people.PUT("/:id", controllers.UpdatePerson)
				people.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeletePerson)
				people.GET("/:id/customers", controllers.GetPersonCustomers)
			}

			// 客户管理路由（所有认证用户可操作）
			customers := authenticated.Group("/customers")
			{
				customers.GET("", controllers.GetCustomers)
				customers.GET("/:id", controllers.GetCustomer)
				customers.POST("", controllers.CreateCustomer)
				customers.PUT("/:id", controllers.UpdateCustomer)
				customers.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteCustomer)
				customers.GET("/:id/tasks", controllers.GetCustomerTasks)
				customers.GET("/:id/payments", controllers.GetCustomerPayments)
			}

			// 任务管理路由（所有认证用户可操作）
			tasks := authenticated.Group("/tasks")
			{
				tasks.GET("", controllers.GetTasks)
				tasks.GET("/:id", controllers.GetTask)
				tasks.POST("", controllers.CreateTask)
				tasks.PUT("/:id", controllers.UpdateTask)
				tasks.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteTask)
			}

			// 协议管理路由（所有认证用户可操作）
			agreements := authenticated.Group("/agreements")
			{
				agreements.GET("", controllers.GetAgreements)
				agreements.GET("/:id", controllers.GetAgreement)
				agreements.POST("", controllers.CreateAgreement)
				agreements.PUT("/:id", controllers.UpdateAgreement)
				agreements.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteAgreement)
			}

			// 收款管理路由（所有认证用户可操作）
			payments := authenticated.Group("/payments")
			{
				payments.GET("", controllers.GetPayments)
				payments.GET("/:id", controllers.GetPayment)
				payments.POST("", controllers.CreatePayment)
				payments.PUT("/:id", controllers.UpdatePayment)
				payments.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeletePayment)
			}

			// 统计分析路由
			statistics := authenticated.Group("/statistics")
			{
				statistics.GET("/overview", controllers.GetOverview)
				statistics.GET("/tasks", controllers.GetTaskStats)
				statistics.GET("/payments", controllers.GetPaymentStats)
			}

			// 导入导出路由
			templates := authenticated.Group("/templates")
			{
				templates.GET("/:type", importExportCtrl.DownloadTemplate)
			}

			importAPI := authenticated.Group("/import")
			{
				importAPI.POST("/people", importExportCtrl.ImportPeople)
				importAPI.POST("/customers", importExportCtrl.ImportCustomers)
			}

			export := authenticated.Group("/export")
			{
				export.GET("/people", importExportCtrl.ExportPeople)
				export.GET("/customers", importExportCtrl.ExportCustomers)
			}
		}
	}
}
