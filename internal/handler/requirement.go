package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type RequirementHandler struct{ svc *service.RequirementService }

func NewRequirementHandler(svc *service.RequirementService) *RequirementHandler {
	return &RequirementHandler{svc: svc}
}

func (h *RequirementHandler) List(c *gin.Context) {
	var q request.RequirementQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := q.Normalize()
	items, total, err := h.svc.List(page, pageSize, q.ReqType, q.ProjectID, q.CustomerID, q.ScheduledDate, q.Status, q.Keyword)
	if err != nil {
		response.InternalError(c, "failed to list requirements")
		return
	}
	response.SuccessWithPagination(c, items, page, pageSize, total)
}

func (h *RequirementHandler) GetByID(c *gin.Context) {
	m, err := h.svc.GetByID(c.Param("id"))
	if err != nil {
		if err == service.ErrRequirementNotFound {
			response.NotFound(c, "requirement not found")
			return
		}
		response.InternalError(c, "failed to get requirement")
		return
	}
	response.Success(c, m)
}

func (h *RequirementHandler) Create(c *gin.Context) {
	var req request.CreateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m := &model.Requirement{
		ReqType: req.ReqType, Title: req.Title, Description: req.Description,
		ProjectID: req.ProjectID, CustomerID: req.CustomerID,
		RequirementCode: req.RequirementCode,
		Priority:        req.Priority, Status: model.ReqStatusPending,
	}
	if m.Priority == "" {
		m.Priority = model.TaskPriorityMedium
	}
	if req.ScheduledDate != nil && *req.ScheduledDate != "" {
		dt, _ := time.Parse("2006-01-02", *req.ScheduledDate)
		m.ScheduledDate = &dt
	}
	if err := h.svc.Create(m); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, m)
}

func (h *RequirementHandler) Update(c *gin.Context) {
	var req request.UpdateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	u := map[string]interface{}{}
	if req.Title != nil {
		u["title"] = *req.Title
	}
	if req.Description != nil {
		u["description"] = *req.Description
	}
	if req.ProjectID != nil {
		u["project_id"] = *req.ProjectID
	}
	if req.CustomerID != nil {
		u["customer_id"] = *req.CustomerID
	}
	if req.RequirementCode != nil {
		u["requirement_code"] = *req.RequirementCode
	}
	if req.Priority != nil {
		u["priority"] = *req.Priority
	}
	if req.Status != nil {
		u["status"] = *req.Status
	}
	if req.ScheduledDate != nil {
		u["scheduled_date"] = *req.ScheduledDate
	}
	if len(u) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(c.Param("id"), u); err != nil {
		if err == service.ErrRequirementNotFound {
			response.NotFound(c, "requirement not found")
			return
		}
		response.InternalError(c, "failed to update requirement")
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		if err == service.ErrRequirementNotFound {
			response.NotFound(c, "requirement not found")
			return
		}
		response.InternalError(c, "failed to delete requirement")
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) ForceDelete(c *gin.Context) {
	if err := h.svc.ForceDelete(c.Param("id")); err != nil {
		response.InternalError(c, "failed")
		return
	}
	response.Success(c, nil)
}
func (h *RequirementHandler) Restore(c *gin.Context) {
	if err := h.svc.Restore(c.Param("id")); err != nil {
		response.InternalError(c, "failed")
		return
	}
	response.Success(c, nil)
}
