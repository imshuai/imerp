package main

import (
	"erp/config"
	"erp/routes"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// 获取运行环境，默认为development
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 初始化日志系统
	if err := config.InitLogger(env); err != nil {
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer config.SyncLogger()

	// 设置Gin运行模式
	if env == "production" {
		config.Info("Starting ERP server in production mode",
			zap.String("version", "v0.4.0"))
	} else {
		config.Info("Starting ERP server in development mode",
			zap.String("version", "v0.4.0"),
			zap.String("env", env))
	}

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		config.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 创建Gin实例
	r := routes.SetupGinEngine(env)

	// 优雅关闭处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		config.Info("Shutting down server gracefully...")
		config.SyncLogger()
		os.Exit(0)
	}()

	// 启动服务
	addr := ":8080"
	config.Info("Server listening",
		zap.String("addr", addr),
		zap.String("frontend", "http://localhost:8080"),
		zap.String("api", "http://localhost:8080/api"))
	if err := r.Run(addr); err != nil {
		config.Fatal("Failed to start server", zap.Error(err))
	}
}
