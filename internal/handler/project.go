package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// List GET /projects
func (h *ProjectHandler) List(c *gin.Context) {
	var query request.ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := query.Normalize()
	projects, total, err := h.svc.List(page, pageSize, query.Keyword, query.Status, query.CustomerID)
	if err != nil {
		response.InternalError(c, "failed to list projects")
		return
	}
	response.SuccessWithPagination(c, projects, page, pageSize, total)
}

// GetByID GET /projects/:id
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	detail, err := h.svc.GetByID(id)
	if err != nil {
		if err == service.ErrProjectNotFound {
			response.NotFound(c, "project not found")
			return
		}
		response.InternalError(c, "failed to get project")
		return
	}
	response.Success(c, detail)
}

// Create POST /projects
func (h *ProjectHandler) Create(c *gin.Context) {
	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	project := &model.Project{
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
		CustomerID:  req.CustomerID,
		ManagerID:   req.ManagerID,
		Status:      req.Status,
		Description: req.Description,
	}
	if req.Status == "" {
		project.Status = model.ProjectStatusPlanning
	}
	if req.StartDate != nil && *req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		project.StartDate = &t
	}
	if req.GoLiveDate != nil && *req.GoLiveDate != "" {
		t, _ := time.Parse("2006-01-02", *req.GoLiveDate)
		project.GoLiveDate = &t
	}
	if req.CompletionDate != nil && *req.CompletionDate != "" {
		t, _ := time.Parse("2006-01-02", *req.CompletionDate)
		project.CompletionDate = &t
	}
	if err := h.svc.Create(project); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, project)
}

// Update POST /projects/:id/update
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ProjectCode != nil {
		updates["project_code"] = *req.ProjectCode
	}
	if req.CustomerID != nil {
		updates["customer_id"] = *req.CustomerID
	}
	if req.ManagerID != nil {
		updates["manager_id"] = *req.ManagerID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.GoLiveDate != nil {
		updates["go_live_date"] = *req.GoLiveDate
	}
	if req.CompletionDate != nil {
		updates["completion_date"] = *req.CompletionDate
	}
	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(id, updates); err != nil {
		if err == service.ErrProjectNotFound {
			response.NotFound(c, "project not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete POST /projects/:id/delete
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		if err == service.ErrProjectNotFound {
			response.NotFound(c, "project not found")
			return
		}
		response.InternalError(c, "failed to delete project")
		return
	}
	response.Success(c, nil)
}

func (h *ProjectHandler) ForceDelete(c *gin.Context) {
	if err := h.svc.ForceDelete(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
func (h *ProjectHandler) Restore(c *gin.Context) {
	if err := h.svc.Restore(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
