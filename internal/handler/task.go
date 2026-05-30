package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type TaskHandler struct{ svc *service.TaskService }

func NewTaskHandler(svc *service.TaskService) *TaskHandler { return &TaskHandler{svc: svc} }

func (h *TaskHandler) List(c *gin.Context) {
	var q request.TaskQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := q.Normalize()
	tasks, total, err := h.svc.List(page, pageSize, q)
	if err != nil {
		response.InternalError(c, "failed to list tasks")
		return
	}
	response.SuccessWithPagination(c, tasks, page, pageSize, total)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	t, err := h.svc.GetByID(c.Param("id"))
	if err != nil {
		if err == service.ErrTaskNotFound {
			response.NotFound(c, "task not found")
			return
		}
		response.InternalError(c, "failed to get task")
		return
	}
	response.Success(c, t)
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req request.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	task := &model.Task{
		TaskType: req.TaskType, Title: req.Title, Description: req.Description,
		ProjectID: req.ProjectID, CustomerID: req.CustomerID,
		Priority: req.Priority, Status: model.TaskStatusInProgress,
	}
	if task.Priority == "" {
		task.Priority = model.TaskPriorityMedium
	}
	if req.DueDate != nil && *req.DueDate != "" {
		dt, _ := time.Parse("2006-01-02", *req.DueDate)
		task.DueDate = &dt
	}
	if err := h.svc.Create(task); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	var req request.UpdateTaskRequest
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
	if req.Status != nil {
		u["status"] = *req.Status
	}
	if req.Priority != nil {
		u["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		u["due_date"] = *req.DueDate
	}
	if req.EstimatedHours != nil {
		u["estimated_hours"] = *req.EstimatedHours
	}
	if len(u) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(c.Param("id"), u); err != nil {
		if err == service.ErrTaskNotFound {
			response.NotFound(c, "task not found")
			return
		}
		response.InternalError(c, "failed to update task")
		return
	}
	response.Success(c, nil)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	var req request.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(c.Param("id"), req.Status); err != nil {
		if err == service.ErrTaskNotFound {
			response.NotFound(c, "task not found")
			return
		}
		response.InternalError(c, "failed to update status")
		return
	}
	response.Success(c, nil)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		if err == service.ErrTaskNotFound {
			response.NotFound(c, "task not found")
			return
		}
		response.InternalError(c, "failed to delete task")
		return
	}
	response.Success(c, nil)
}

func (h *TaskHandler) ForceDelete(c *gin.Context) {
	if err := h.svc.ForceDelete(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
func (h *TaskHandler) Restore(c *gin.Context) {
	if err := h.svc.Restore(c.Param("id")); err != nil { response.InternalError(c, "failed"); return }
	response.Success(c, nil)
}
