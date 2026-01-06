package controllers

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PersonID *uint  `json:"person_id"` // 普通服务人员登录时使用
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token             string `json:"token"`
	UserID            uint   `json:"user_id"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	// 普通服务人员登录（无密码）
	if req.PersonID != nil {
		loginServicePerson(c, *req.PersonID)
		return
	}

	// admin 或管理员登录（需要密码）
	loginAdmin(c, req.Username, req.Password)
}

// loginServicePerson 普通服务人员登录
func loginServicePerson(c *gin.Context, personID uint) {
	var person models.Person
	if err := config.DB.First(&person, personID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Service person not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	if !person.IsServicePerson {
		ErrorResponse(c, 400, "Not a service person")
		return
	}

	// 生成 token（role 为 service_person）
	token, err := auth.GenerateToken(person.ID, person.Name, "service_person", &person.ID)
	if err != nil {
		ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	SuccessResponse(c, LoginResponse{
		Token:             token,
		UserID:            person.ID,
		Username:          person.Name,
		Role:              "service_person",
		MustChangePassword: false,
	})
}

// loginAdmin admin 或管理员登录
func loginAdmin(c *gin.Context, username, password string) {
	var adminUser models.AdminUser
	if err := config.DB.Where("username = ?", username).First(&adminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 401, "Invalid username or password")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	// 验证密码
	if !auth.CheckPassword(password, adminUser.PasswordHash) {
		ErrorResponse(c, 401, "Invalid username or password")
		return
	}

	// 生成 token
	token, err := auth.GenerateToken(adminUser.ID, adminUser.Username, adminUser.Role, adminUser.PersonID)
	if err != nil {
		ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	SuccessResponse(c, LoginResponse{
		Token:             token,
		UserID:            adminUser.ID,
		Username:          adminUser.Username,
		Role:              adminUser.Role,
		MustChangePassword: adminUser.MustChangePassword,
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	// 从 context 获取用户信息（由中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, 401, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	var adminUser models.AdminUser
	if err := config.DB.First(&adminUser, userID).Error; err != nil {
		ErrorResponse(c, 404, "User not found")
		return
	}

	// 验证旧密码
	if !auth.CheckPassword(req.OldPassword, adminUser.PasswordHash) {
		ErrorResponse(c, 400, "Incorrect old password")
		return
	}

	// 生成新密码哈希
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// 更新密码
	now := time.Now()
	adminUser.PasswordHash = hashedPassword
	adminUser.MustChangePassword = false
	adminUser.LastPasswordChange = &now

	if err := config.DB.Save(&adminUser).Error; err != nil {
		ErrorResponse(c, 500, "Failed to update password")
		return
	}

	SuccessResponse(c, gin.H{"message": "Password changed successfully"})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role == "service_person" {
		// 服务人员
		var person models.Person
		if err := config.DB.First(&person, userID).Error; err != nil {
			ErrorResponse(c, 404, "User not found")
			return
		}
		SuccessResponse(c, gin.H{
			"id":               person.ID,
			"name":             person.Name,
			"role":             "service_person",
			"is_manager":       person.IsManager,
			"is_service_person": person.IsServicePerson,
		})
	} else {
		// admin 或管理员
		var adminUser models.AdminUser
		if err := config.DB.Preload("Person").First(&adminUser, userID).Error; err != nil {
			ErrorResponse(c, 404, "User not found")
			return
		}
		SuccessResponse(c, gin.H{
			"id":                adminUser.ID,
			"username":          adminUser.Username,
			"role":              adminUser.Role,
			"must_change_password": adminUser.MustChangePassword,
			"person":            adminUser.Person,
		})
	}
}
