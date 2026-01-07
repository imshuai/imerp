package controllers

import (
	"erp/models"
	"github.com/gin-gonic/gin"
)

// GetCustomerTypes 获取客户类型选项
func GetCustomerTypes(c *gin.Context) {
	types := models.GetCustomerTypeOptions()
	SuccessResponse(c, types)
}

// GetCreditRatings 获取信用等级选项
func GetCreditRatings(c *gin.Context) {
	ratings := models.GetCreditRatingOptions()
	SuccessResponse(c, ratings)
}

// GetAccountTypes 获取账户类型选项
func GetAccountTypes(c *gin.Context) {
	types := models.GetAccountTypeOptions()
	SuccessResponse(c, types)
}
