package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type WeeklyReportHandler struct{ svc *service.WeeklyReportService }

func NewWeeklyReportHandler(svc *service.WeeklyReportService) *WeeklyReportHandler {
	return &WeeklyReportHandler{svc: svc}
}

func (h *WeeklyReportHandler) List(c *gin.Context) {
	var q request.WeeklyReportQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := q.Normalize()
	userID := c.GetString("user_id")
	items, total, err := h.svc.List(page, pageSize, &userID)
	if err != nil {
		response.InternalError(c, "failed to list weekly reports")
		return
	}
	response.SuccessWithPagination(c, items, page, pageSize, total)
}

func (h *WeeklyReportHandler) GetByID(c *gin.Context) {
	item, err := h.svc.GetByID(c.Param("id"))
	if err != nil {
		if err == service.ErrWeeklyReportNotFound {
			response.NotFound(c, "weekly report not found")
			return
		}
		response.InternalError(c, "failed to get weekly report")
		return
	}
	response.Success(c, item)
}

func (h *WeeklyReportHandler) Generate(c *gin.Context) {
	var req request.GenerateWeeklyReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	item, err := h.svc.Generate(c.GetString("user_id"), service.GenerateWeeklyReportInput{
		ReportType: req.ReportType,
		WeekStart:  req.WeekStart,
		WeekEnd:    req.WeekEnd,
		ProjectID:  req.ProjectID,
		UseAI:      req.UseAI,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *WeeklyReportHandler) Update(c *gin.Context) {
	var req request.UpdateWeeklyReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	status := ""
	if req.Status != nil {
		status = *req.Status
	}
	if status == model.ReportStatusFinal {
		response.BadRequest(c, "use finalize endpoint for final status")
		return
	}
	if err := h.svc.UpdateContent(c.Param("id"), req.Content, status); err != nil {
		if err == service.ErrWeeklyReportNotFound {
			response.NotFound(c, "weekly report not found")
			return
		}
		response.InternalError(c, "failed to update weekly report")
		return
	}
	response.Success(c, nil)
}

func (h *WeeklyReportHandler) Review(c *gin.Context) {
	if err := h.svc.Review(c.Param("id")); err != nil {
		if err == service.ErrWeeklyReportNotFound {
			response.NotFound(c, "weekly report not found")
			return
		}
		response.InternalError(c, "failed to review weekly report")
		return
	}
	response.Success(c, nil)
}

func (h *WeeklyReportHandler) Finalize(c *gin.Context) {
	if err := h.svc.Finalize(c.Param("id")); err != nil {
		if err == service.ErrWeeklyReportNotFound {
			response.NotFound(c, "weekly report not found")
			return
		}
		response.InternalError(c, "failed to finalize weekly report")
		return
	}
	response.Success(c, nil)
}
