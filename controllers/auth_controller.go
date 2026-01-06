package controllers

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token             string `json:"token"`
	UserID            uint   `json:"user_id"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
	PersonID          *uint  `json:"person_id,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Login 用户登录（统一使用用户名+密码）
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "Invalid request: "+err.Error())
		return
	}

	// 查找用户（通过用户名）
	var adminUser models.AdminUser
	if err := config.DB.Where("username = ?", req.Username).First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, 401, "用户名或密码错误")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}

	// 验证密码
	if !auth.CheckPassword(req.Password, adminUser.PasswordHash) {
		ErrorResponse(c, 401, "用户名或密码错误")
		return
	}

	// 生成 token
	token, err := auth.GenerateToken(adminUser.ID, adminUser.Username, adminUser.Role, adminUser.PersonID)
	if err != nil {
		ErrorResponse(c, 500, "生成token失败")
		return
	}

	SuccessResponse(c, LoginResponse{
		Token:             token,
		UserID:            adminUser.ID,
		Username:          adminUser.Username,
		Role:              adminUser.Role,
		MustChangePassword: adminUser.MustChangePassword,
		PersonID:          adminUser.PersonID,
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
		ErrorResponse(c, 404, "用户不存在")
		return
	}

	// 验证旧密码
	if !auth.CheckPassword(req.OldPassword, adminUser.PasswordHash) {
		ErrorResponse(c, 400, "原密码错误")
		return
	}

	// 生成新密码哈希
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		ErrorResponse(c, 500, "密码加密失败")
		return
	}

	// 更新密码
	now := time.Now()
	adminUser.PasswordHash = hashedPassword
	adminUser.MustChangePassword = false
	adminUser.LastPasswordChange = &now

	if err := config.DB.Save(&adminUser).Error; err != nil {
		ErrorResponse(c, 500, "密码修改失败")
		return
	}

	SuccessResponse(c, gin.H{"message": "密码修改成功"})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var adminUser models.AdminUser
	if err := config.DB.Preload("Person").First(&adminUser, userID).Error; err != nil {
		ErrorResponse(c, 404, "用户不存在")
		return
	}

	// 构建响应
	response := gin.H{
		"id":                  adminUser.ID,
		"username":            adminUser.Username,
		"role":                adminUser.Role,
		"must_change_password": adminUser.MustChangePassword,
		"person":              adminUser.Person,
	}

	// 如果是服务人员，额外返回 IsManager 信息
	if role == "service_person" && adminUser.Person != nil {
		response["is_manager"] = adminUser.Person.IsManager
		response["is_service_person"] = adminUser.Person.IsServicePerson
	}

	SuccessResponse(c, response)
}
