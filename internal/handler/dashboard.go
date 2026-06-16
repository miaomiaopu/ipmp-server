package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type DashboardHandler struct{ svc *service.DashboardService }

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	data, err := h.svc.Stats(c.GetString("user_id"), c.GetString("role"))
	if err != nil {
		response.InternalError(c, "failed to get dashboard stats")
		return
	}
	response.Success(c, data)
}

func (h *DashboardHandler) ThisWeek(c *gin.Context) {
	data, err := h.svc.ThisWeek(c.GetString("user_id"), c.GetString("role"))
	if err != nil {
		response.InternalError(c, "failed to get this week dashboard")
		return
	}
	response.Success(c, data)
}
