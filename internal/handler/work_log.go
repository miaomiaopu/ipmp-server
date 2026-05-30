package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type WorkLogHandler struct{ svc *service.WorkLogService }

func NewWorkLogHandler(svc *service.WorkLogService) *WorkLogHandler { return &WorkLogHandler{svc: svc} }

func (h *WorkLogHandler) List(c *gin.Context) {
	var q request.WorkLogQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := q.Normalize()
	items, total, err := h.svc.List(page, pageSize, q.UserID, q.ProjectID, q.StartDate, q.EndDate)
	if err != nil {
		response.InternalError(c, "failed to list work logs")
		return
	}
	response.SuccessWithPagination(c, items, page, pageSize, total)
}

func (h *WorkLogHandler) GetByID(c *gin.Context) {
	w, err := h.svc.GetByID(c.Param("id"))
	if err != nil {
		if err == service.ErrWorkLogNotFound {
			response.NotFound(c, "work log not found")
			return
		}
		response.InternalError(c, "failed to get work log")
		return
	}
	response.Success(c, w)
}

func (h *WorkLogHandler) Create(c *gin.Context) {
	var req request.CreateWorkLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dt, _ := time.Parse("2006-01-02", req.LogDate)
	userID, _ := c.Get("user_id")

	w := &model.WorkLog{
		TaskID:      req.TaskID,
		UserID:      userID.(string),
		LogDate:     dt,
		Hours:       req.Hours,
		Description: req.Description,
	}
	// 冗余填充 project_id/customer_id（从关联 task 推导，加速报表查询）
	// TODO: 从 task 查询关联的 project/customer 并填充
	if err := h.svc.Create(w); err != nil {
		response.InternalError(c, "failed to create work log")
		return
	}
	response.Success(c, w)
}

func (h *WorkLogHandler) Update(c *gin.Context) {
	var req request.UpdateWorkLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	u := map[string]interface{}{}
	if req.LogDate != nil {
		u["log_date"] = *req.LogDate
	}
	if req.Hours != nil {
		u["hours"] = *req.Hours
	}
	if req.Description != nil {
		u["description"] = *req.Description
	}
	if len(u) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(c.Param("id"), u); err != nil {
		if err == service.ErrWorkLogNotFound {
			response.NotFound(c, "work log not found")
			return
		}
		response.InternalError(c, "failed to update work log")
		return
	}
	response.Success(c, nil)
}

func (h *WorkLogHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		if err == service.ErrWorkLogNotFound {
			response.NotFound(c, "work log not found")
			return
		}
		response.InternalError(c, "failed to delete work log")
		return
	}
	response.Success(c, nil)
}

func (h *WorkLogHandler) Stats(c *gin.Context) {
	userID := c.Query("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	groupBy := c.Query("group_by")
	var uid *string
	if userID != "" {
		uid = &userID
	}
	items, err := h.svc.Stats(uid, startDate, endDate, groupBy)
	if err != nil {
		response.InternalError(c, "failed to get stats")
		return
	}
	response.Success(c, items)
}

func (h *WorkLogHandler) ForceDelete(c *gin.Context) {
	if err := h.svc.ForceDelete(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
func (h *WorkLogHandler) Restore(c *gin.Context) {
	if err := h.svc.Restore(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
