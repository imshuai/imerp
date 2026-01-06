package controllers

import (
	"erp/config"
	"erp/models"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCustomer 创建客户
func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"create",
		"customer",
		nil, // resourceID is nil for create
		nil, // no old value
		customer,
		func() error {
			// 执行创建操作
			if err := config.DB.Create(&customer).Error; err != nil {
				return err
			}
			// 同步更新Person表的关联字段
			syncPersonRelations(&customer)
			// 加载关联数据
			loadCustomerRelations(&customer)
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回创建的数据
	if !needsApproval {
		SuccessResponse(c, customer)
	}
}

// GetCustomers 获取客户列表
func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	var total int64

	// 获取查询参数
	keyword := c.Query("keyword")
	representative := c.Query("representative")
	investor := c.Query("investor")
	servicePerson := c.Query("service_person")

	query := config.DB.Model(&models.Customer{})

	// 按名称/税号/电话搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR tax_number LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按法定代表人搜索
	if representative != "" {
		// 先查找符合条件的人员ID
		var personIDs []uint
		config.DB.Model(&models.Person{}).
			Where("name LIKE ? OR phone LIKE ? OR id_card LIKE ?",
				"%"+representative+"%", "%"+representative+"%", "%"+representative+"%").
			Pluck("id", &personIDs)

		if len(personIDs) > 0 {
			query = query.Where("representative_id IN ?", personIDs)
		} else {
			// 没有找到匹配的人员，返回空结果
			SuccessPaginatedResponse(c, 0, []models.Customer{})
			return
		}
	}

	// 按投资人搜索
	if investor != "" {
		// 通过JSON字段搜索
		var customerIDs []uint
		config.DB.Raw(`
			SELECT id FROM customers
			WHERE investors IS NOT NULL
			AND EXISTS (
				SELECT 1 FROM json_each(investors)
				WHERE json_valid(investors)
				AND CAST(json_extract(value, '$.person_id') AS INTEGER) IN (
					SELECT id FROM people
					WHERE name LIKE ? OR phone LIKE ? OR id_card LIKE ?
				)
			)
		`, "%"+investor+"%", "%"+investor+"%", "%"+investor+"%").Scan(&customerIDs)

		if len(customerIDs) > 0 {
			query = query.Where("id IN ?", customerIDs)
		} else {
			SuccessPaginatedResponse(c, 0, []models.Customer{})
			return
		}
	}

	// 按服务人员搜索
	if servicePerson != "" {
		// 先查找符合条件的人员ID
		var personIDs []uint
		config.DB.Model(&models.Person{}).
			Where("is_service_person = ? AND (name LIKE ? OR phone LIKE ?)",
				true, "%"+servicePerson+"%", "%"+servicePerson+"%").
			Pluck("id", &personIDs)

		if len(personIDs) > 0 {
			// 构建逗号分隔ID的搜索条件
			var conditions []string
			var args []interface{}
			for _, pid := range personIDs {
				conditions = append(conditions, ", " + strconv.Itoa(int(pid)) + ",")
				conditions = append(conditions, strconv.Itoa(int(pid)) + ",")
				conditions = append(conditions, "," + strconv.Itoa(int(pid)))
				conditions = append(conditions, strconv.Itoa(int(pid)))
				args = append(args, pid)
			}
			query = query.Where("service_person_ids LIKE ?", "%"+strconv.Itoa(int(personIDs[0]))+"%")
		} else {
			SuccessPaginatedResponse(c, 0, []models.Customer{})
			return
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取列表（关联数据通过loadCustomerRelations手动加载）
	if err := query.Find(&customers).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch customers: "+err.Error())
		return
	}

	// 为每个客户加载关联的服务人员信息
	for i := range customers {
		loadCustomerRelations(&customers[i])
	}

	SuccessPaginatedResponse(c, total, customers)
}

// GetCustomer 获取客户详情
func GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var customer models.Customer
	// 使用Preload一次性加载所有关联数据
	if err := config.DB.Preload("Tasks").Preload("Payments").First(&customer, id).Error; err != nil {
		ErrorResponse(c, 404, "Customer not found")
		return
	}

	// 加载关联的人员信息
	loadCustomerRelations(&customer)

	SuccessResponse(c, customer)
}

// UpdateCustomer 更新客户
func UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		ErrorResponse(c, 404, "Customer not found")
		return
	}

	var updateData models.Customer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	customerID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"update",
		"customer",
		&customerID,
		customer, // old value
		updateData, // new value
		func() error {
			// 更新字段
			if err := config.DB.Model(&customer).Updates(updateData).Error; err != nil {
				return err
			}
			// 同步更新Person表的关联字段
			syncPersonRelations(&updateData)
			// 重新获取更新后的数据
			config.DB.First(&customer, id)
			loadCustomerRelations(&customer)
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回更新后的数据
	if !needsApproval {
		SuccessResponse(c, customer)
	}
}

// DeleteCustomer 删除客户
func DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		ErrorResponse(c, 404, "Customer not found")
		return
	}

	customerID := uint(id)

	// 使用审批流程处理
	needsApproval, err := HandleOperationWithApproval(
		c,
		"delete",
		"customer",
		&customerID,
		customer, // old value
		nil, // no new value for delete
		func() error {
			if err := config.DB.Delete(&models.Customer{}, id).Error; err != nil {
				return err
			}
			// 清理Person表中的关联ID
			config.DB.Model(&models.Person{}).
				Where("representative_customer_ids LIKE ?", "%,"+strconv.Itoa(int(customerID))+",%").
				Update("representative_customer_ids", gorm.Expr("REPLACE(representative_customer_ids, ?, '')", ","+strconv.Itoa(int(customerID))+","))

			config.DB.Model(&models.Person{}).
				Where("investor_customer_ids LIKE ?", "%,"+strconv.Itoa(int(customerID))+",%").
				Update("investor_customer_ids", gorm.Expr("REPLACE(investor_customer_ids, ?, '')", ","+strconv.Itoa(int(customerID))+","))

			config.DB.Model(&models.Person{}).
				Where("service_customer_ids LIKE ?", "%,"+strconv.Itoa(int(customerID))+",%").
				Update("service_customer_ids", gorm.Expr("REPLACE(service_customer_ids, ?, '')", ","+strconv.Itoa(int(customerID))+","))
			return nil
		},
	)

	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	// 如果不需要审批（管理员操作），返回成功消息
	if !needsApproval {
		SuccessResponse(c, gin.H{"message": "Customer deleted successfully"})
	}
}

// GetCustomerTasks 获取客户的任务列表
func GetCustomerTasks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var tasks []models.Task
	if err := config.DB.Where("customer_id = ?", id).Find(&tasks).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch tasks: "+err.Error())
		return
	}

	SuccessResponse(c, tasks)
}

// GetCustomerPayments 获取客户的收款记录
func GetCustomerPayments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, 400, "Invalid customer ID")
		return
	}

	var payments []models.Payment
	if err := config.DB.Where("customer_id = ?", id).Find(&payments).Error; err != nil {
		ErrorResponse(c, 500, "Failed to fetch payments: "+err.Error())
		return
	}

	SuccessResponse(c, payments)
}

// ============ 辅助函数 ============

// loadCustomerRelations 加载客户关联的人员和协议信息
func loadCustomerRelations(customer *models.Customer) {
	// 加载法定代表人
	if customer.RepresentativeID != nil {
		var rep models.Person
		if config.DB.First(&rep, *customer.RepresentativeID).Error == nil {
			customer.Representative = &rep
		}
	}

	// 加载投资人列表
	// investors字段可能是JSON字符串（前端JSON.stringify的结果），需要先解析为字符串再解析为数组
	if customer.Investors != nil && len(customer.Investors) > 0 {
		var investorInfos []models.InvestorInfo
		var investorsJSON []byte

		// 尝试直接解析（如果investors已经是JSON数组）
		if err := json.Unmarshal(customer.Investors, &investorInfos); err != nil {
			// 如果失败，尝试先解析为字符串（如果investors是JSON字符串）
			var investorsStr string
			if err2 := json.Unmarshal(customer.Investors, &investorsStr); err2 == nil {
				investorsJSON = []byte(investorsStr)
			} else {
				// 如果还是失败，直接使用原始字节数组
				investorsJSON = customer.Investors
			}
		} else {
			// 直接解析成功，使用原始字节数组
			investorsJSON = customer.Investors
		}

		// 解析投资人信息
		if len(investorsJSON) > 0 {
			if err := json.Unmarshal(investorsJSON, &investorInfos); err == nil {
				var investorIDs []uint
				for _, info := range investorInfos {
					investorIDs = append(investorIDs, info.PersonID)
				}
				if len(investorIDs) > 0 {
					var investors []models.Person
					config.DB.Where("id IN ?", investorIDs).Find(&investors)
					customer.InvestorList = investors
				}
			}
		}
	}

	// 加载服务人员列表
	if customer.ServicePersonIDs != "" {
		ids := StringToIDs(customer.ServicePersonIDs)
		if len(ids) > 0 {
			var servicePersons []models.Person
			config.DB.Where("id IN ?", ids).Find(&servicePersons)
			customer.ServicePersons = servicePersons
		}
	}

	// 加载协议列表
	if customer.AgreementIDs != "" {
		ids := StringToIDs(customer.AgreementIDs)
		if len(ids) > 0 {
			var agreements []models.Agreement
			config.DB.Where("id IN ?", ids).Find(&agreements)
			customer.Agreements = agreements
		}
	}
}

// syncPersonRelations 同步更新Person表的客户关联字段
func syncPersonRelations(customer *models.Customer) {
	customerID := customer.ID

	// 更新法定代表人关联
	if customer.RepresentativeID != nil {
		var rep models.Person
		if config.DB.First(&rep, *customer.RepresentativeID).Error == nil {
			ids := StringToIDs(rep.RepresentativeCustomerIDs)
			ids = appendUniqueID(ids, customerID)
			rep.RepresentativeCustomerIDs = IDsToString(ids)
			config.DB.Save(&rep)
		}
	}

	// 更新投资人关联
	if customer.Investors != nil {
		var investorInfos []models.InvestorInfo
		if err := json.Unmarshal(customer.Investors, &investorInfos); err == nil {
			for _, info := range investorInfos {
				var inv models.Person
				if config.DB.First(&inv, info.PersonID).Error == nil {
					ids := StringToIDs(inv.InvestorCustomerIDs)
					ids = appendUniqueID(ids, customerID)
					inv.InvestorCustomerIDs = IDsToString(ids)
					config.DB.Save(&inv)
				}
			}
		}
	}

	// 更新服务人员关联
	if customer.ServicePersonIDs != "" {
		ids := StringToIDs(customer.ServicePersonIDs)
		for _, personID := range ids {
			var sp models.Person
			if config.DB.First(&sp, personID).Error == nil {
				customerIDs := StringToIDs(sp.ServiceCustomerIDs)
				customerIDs = appendUniqueID(customerIDs, customerID)
				sp.ServiceCustomerIDs = IDsToString(customerIDs)
				config.DB.Save(&sp)
			}
		}
	}
}

// appendUniqueID 追加ID到数组，避免重复
func appendUniqueID(ids []uint, newID uint) []uint {
	for _, id := range ids {
		if id == newID {
			return ids
		}
	}
	return append(ids, newID)
}
