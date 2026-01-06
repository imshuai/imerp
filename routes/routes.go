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
		}

		// 公开路由 - 获取服务人员列表（用于登录页面）
		api.GET("/people", controllers.GetPeople)

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authenticated.GET("/user/me", controllers.GetCurrentUser)
			authenticated.POST("/user/change-password", controllers.ChangePassword)

			// 管理员路由（需要管理员权限）
			admin := authenticated.Group("/admin")
			admin.Use(middleware.RequireManager())
			{
				admin.GET("/users", controllers.GetAdminUsers)
				admin.GET("/service-people", controllers.GetServicePeople)

				// 仅超级管理员
				superAdmin := admin.Group("")
				superAdmin.Use(middleware.RequireSuperAdmin())
				{
					superAdmin.POST("/users", controllers.CreateAdminUser)
					superAdmin.DELETE("/users/:id", controllers.DeleteAdminUser)
					superAdmin.POST("/set-manager", controllers.SetManager)
				}

				// 审批管理
				admin.GET("/approvals/pending", controllers.GetPendingApprovals)
				admin.POST("/approvals/approve", controllers.ApproveOperation)
				admin.POST("/approvals/reject", controllers.RejectOperation)
				admin.GET("/audit-logs", controllers.GetAuditLogs)
			}

			// 人员管理路由
			people := authenticated.Group("/people")
			{
				people.GET("/:id", controllers.GetPerson)
				people.POST("", middleware.RequireManager(), controllers.CreatePerson)
				people.PUT("/:id", middleware.RequireManager(), controllers.UpdatePerson)
				people.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeletePerson)
				people.GET("/:id/customers", controllers.GetPersonCustomers)
			}

			// 客户管理路由
			customers := authenticated.Group("/customers")
			{
				customers.GET("", controllers.GetCustomers)
				customers.GET("/:id", controllers.GetCustomer)
				customers.POST("", middleware.RequireManager(), controllers.CreateCustomer)
				customers.PUT("/:id", middleware.RequireManager(), controllers.UpdateCustomer)
				customers.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteCustomer)
				customers.GET("/:id/tasks", controllers.GetCustomerTasks)
				customers.GET("/:id/payments", controllers.GetCustomerPayments)
			}

			// 任务管理路由
			tasks := authenticated.Group("/tasks")
			{
				tasks.GET("", controllers.GetTasks)
				tasks.GET("/:id", controllers.GetTask)
				tasks.POST("", controllers.CreateTask)  // 所有认证用户可创建
				tasks.PUT("/:id", middleware.RequireManager(), controllers.UpdateTask)
				tasks.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteTask)
			}

			// 协议管理路由
			agreements := authenticated.Group("/agreements")
			{
				agreements.GET("", controllers.GetAgreements)
				agreements.GET("/:id", controllers.GetAgreement)
				agreements.POST("", middleware.RequireManager(), controllers.CreateAgreement)
				agreements.PUT("/:id", middleware.RequireManager(), controllers.UpdateAgreement)
				agreements.DELETE("/:id", middleware.RequireSuperAdmin(), controllers.DeleteAgreement)
			}

			// 收款管理路由
			payments := authenticated.Group("/payments")
			{
				payments.GET("", controllers.GetPayments)
				payments.GET("/:id", controllers.GetPayment)
				payments.POST("", middleware.RequireManager(), controllers.CreatePayment)
				payments.PUT("/:id", middleware.RequireManager(), controllers.UpdatePayment)
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
