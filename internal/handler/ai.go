package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type AIHandler struct{ svc *service.AIService }

func NewAIHandler(svc *service.AIService) *AIHandler {
	return &AIHandler{svc: svc}
}

func (h *AIHandler) GenerateReport(c *gin.Context) {
	var req request.GenerateAIReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := h.svc.GenerateReport(c.GetString("user_id"), req.Content)
	if err != nil {
		response.InternalError(c, "failed to generate report")
		return
	}
	response.Success(c, data)
}

func (h *AIHandler) Summarize(c *gin.Context) {
	var req request.SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := h.svc.Summarize(c.GetString("user_id"), req.Content)
	if err != nil {
		response.InternalError(c, "failed to summarize")
		return
	}
	response.Success(c, data)
}
