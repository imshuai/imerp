package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// App 应用配置
type App struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Env  string `yaml:"env"`  // 运行环境: development, production
	Host string `yaml:"host"` // 监听地址
	Port int    `yaml:"port"` // 监听端口
	Mode string `yaml:"mode"` // Gin运行模式: debug, release, test
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"` // 日志级别: debug, info, warn, error, fatal
}

var AppConfig *App

// LoadAppConfig 加载应用配置
func LoadAppConfig(configPath string) error {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析 YAML
	appConfig := &App{}
	if err := yaml.Unmarshal(data, appConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	if appConfig.Server.Env == "" {
		appConfig.Server.Env = "development"
	}
	if appConfig.Server.Host == "" {
		appConfig.Server.Host = "0.0.0.0"
	}
	if appConfig.Server.Port == 0 {
		appConfig.Server.Port = 8080
	}
	// mode 的默认值在 main.go 中根据最终的 env 设置
	if appConfig.Log.Level == "" {
		appConfig.Log.Level = "info"
	}

	AppConfig = appConfig
	return nil
}

// GetServerAddr 获取服务器监听地址
func (c *App) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
