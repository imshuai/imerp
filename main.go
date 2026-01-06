package main

import (
	"erp/config"
	"erp/routes"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "config.yaml", "Path to config file")
	env := flag.String("env", "", "Runtime environment (development, production)")
	host := flag.String("host", "", "Server listen address")
	port := flag.Int("port", 0, "Server listen port")
	logLevel := flag.String("log-level", "", "Log level (debug, info, warn, error, fatal)")
	flag.Parse()

	// 加载应用配置
	if err := config.LoadAppConfig(*configPath); err != nil {
		os.Stderr.WriteString("Failed to load config: " + err.Error() + "\n")
		os.Exit(1)
	}

	// 命令行参数覆盖配置文件
	if *env != "" {
		config.AppConfig.Server.Env = *env
	}
	if *host != "" {
		config.AppConfig.Server.Host = *host
	}
	if *port > 0 {
		config.AppConfig.Server.Port = *port
	}
	if *logLevel != "" {
		config.AppConfig.Log.Level = *logLevel
	}

	// 重新计算 Gin 运行模式（如果 env 被覆盖）
	if *env != "" && config.AppConfig.Server.Mode == "" {
		if config.AppConfig.Server.Env == "production" {
			config.AppConfig.Server.Mode = "release"
		} else {
			config.AppConfig.Server.Mode = "debug"
		}
	}

	// 初始化日志系统
	if err := config.InitLogger(config.AppConfig.Log.Level); err != nil {
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer config.SyncLogger()

	// 输出启动日志
	if config.AppConfig.Server.Env == "production" {
		config.Info("Starting ERP server in production mode",
			zap.String("version", "v0.4.0"))
	} else {
		config.Info("Starting ERP server",
			zap.String("version", "v0.4.0"),
			zap.String("env", config.AppConfig.Server.Env),
			zap.String("gin_mode", config.AppConfig.Server.Mode))
	}

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		config.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 创建Gin实例
	r := routes.SetupGinEngine()

	// 优雅关闭处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		config.Info("Shutting down server gracefully...")
		config.SyncLogger()
		os.Exit(0)
	}()

	// 设置前端路由（放在最后，作为fallback）
	routes.SetupFrontendRoutes(r)

	// 启动服务
	addr := config.AppConfig.GetServerAddr()
	config.Info("Server listening",
		zap.String("addr", addr),
		zap.String("api", "http://"+addr+"/api"))
	if err := r.Run(addr); err != nil {
		config.Fatal("Failed to start server", zap.Error(err))
	}
}
