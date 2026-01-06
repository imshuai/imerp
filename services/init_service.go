package services

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"errors"
)

// InitService 初始化服务
type InitService struct{}

// NewInitService 创建初始化服务
func NewInitService() *InitService {
	return &InitService{}
}

// InitializeSuperAdmin 初始化超级管理员
func (s *InitService) InitializeSuperAdmin() error {
	// 检查是否已存在超级管理员
	var count int64
	config.DB.Model(&models.AdminUser{}).Where("role = ?", "super_admin").Count(&count)
	if count > 0 {
		return errors.New("super admin already exists")
	}

	// 创建超级管理员
	hashedPassword, err := auth.HashPassword("admin")
	if err != nil {
		return err
	}

	adminUser := models.AdminUser{
		Username:          "admin",
		PasswordHash:      hashedPassword,
		Role:              "super_admin",
		MustChangePassword: true,
	}

	if err := config.DB.Create(&adminUser).Error; err != nil {
		return err
	}

	return nil
}

// IsSuperAdminExists 检查超级管理员是否存在
func (s *InitService) IsSuperAdminExists() bool {
	var count int64
	config.DB.Model(&models.AdminUser{}).Where("role = ?", "super_admin").Count(&count)
	return count > 0
}
