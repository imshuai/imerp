package middleware

import (
	"erp/auth"
	"erp/config"
	"erp/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": 401, "message": "Authorization header required"})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息存入 context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("person_id", claims.PersonID)

		c.Next()
	}
}

// RequireRole 要求特定角色
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"code": 401, "message": "Unauthorized"})
			c.Abort()
			return
		}

		// 检查角色
		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(403, gin.H{"code": 403, "message": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin 要求超级管理员
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("super_admin")
}

// RequireManager 要求管理员（包括超级管理员和普通管理员）
func RequireManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"code": 401, "message": "Unauthorized"})
			c.Abort()
			return
		}

		if userRole != "super_admin" && userRole != "manager" {
			c.JSON(403, gin.H{"code": 403, "message": "Forbidden: manager permissions required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser 从 context 获取当前用户
func GetCurrentUser(c *gin.Context) (uint, string) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	return userID.(uint), role.(string)
}

// GetServicePerson 从 context 获取当前服务人员
func GetServicePerson(c *gin.Context) (*models.Person, error) {
	userID, role := GetCurrentUser(c)
	if role != "service_person" {
		return nil, nil
	}

	var person models.Person
	if err := config.DB.First(&person, userID).Error; err != nil {
		return nil, err
	}
	return &person, nil
}
