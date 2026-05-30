package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// List GET /customers
func (h *CustomerHandler) List(c *gin.Context) {
	var pagination request.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := pagination.Normalize()
	customers, total, err := h.svc.List(page, pageSize, pagination.Keyword, pagination.Status)
	if err != nil {
		response.InternalError(c, "failed to list customers")
		return
	}
	response.SuccessWithPagination(c, customers, page, pageSize, total)
}

// GetByID GET /customers/:id
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	customer, projectCount, err := h.svc.GetByID(id)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			response.NotFound(c, "customer not found")
			return
		}
		response.InternalError(c, "failed to get customer")
		return
	}
	response.Success(c, gin.H{
		"customer":       customer,
		"project_count":  projectCount,
	})
}

// Create POST /customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req request.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer := &model.Customer{
		CustomerCode:  req.CustomerCode,
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		ContactPhone:  model.EncryptedField(req.ContactPhone),
		ContactEmail:  model.EncryptedField(req.ContactEmail),
		Address:       model.EncryptedField(req.Address),
		Notes:         req.Notes,
		Status:        "active",
	}
	if err := h.svc.Create(customer); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, customer)
}

// Update POST /customers/:id/update
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 构建更新 map，仅更新非 nil 字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CustomerCode != nil {
		updates["customer_code"] = *req.CustomerCode
	}
	if req.ContactPerson != nil {
		updates["contact_person"] = *req.ContactPerson
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(id, updates); err != nil {
		if err == service.ErrCustomerNotFound {
			response.NotFound(c, "customer not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete POST /customers/:id/delete
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		if err == service.ErrCustomerNotFound {
			response.NotFound(c, "customer not found")
			return
		}
		response.InternalError(c, "failed to delete customer")
		return
	}
	response.Success(c, nil)
}

func (h *CustomerHandler) ForceDelete(c *gin.Context) {
	if err := h.svc.ForceDelete(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
func (h *CustomerHandler) Restore(c *gin.Context) {
	if err := h.svc.Restore(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
