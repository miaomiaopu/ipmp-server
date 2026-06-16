package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strconv"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var ErrWorkLogNotFound = errors.New("work log not found")

type WorkLogService struct{ repo *repository.WorkLogRepository }

func NewWorkLogService(repo *repository.WorkLogRepository) *WorkLogService {
	return &WorkLogService{repo: repo}
}

func (s *WorkLogService) List(page, pageSize int, userID, projectID *string, startDate, endDate string) ([]model.WorkLog, int64, error) {
	return s.repo.List(page, pageSize, userID, projectID, startDate, endDate)
}

func (s *WorkLogService) GetByID(id string) (*model.WorkLog, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkLogNotFound
		}
		return nil, err
	}
	return w, nil
}

func (s *WorkLogService) Create(w *model.WorkLog) error {
	if w.TaskID != nil && *w.TaskID != "" {
		task, err := s.repo.FindTaskForWorkLog(*w.TaskID)
		if err != nil {
			return err
		}
		w.ProjectID = task.ProjectID
		w.CustomerID = task.CustomerID
		if task.TaskType == model.TaskTypeProject && task.Project != nil {
			w.CustomerID = task.Project.CustomerID
		}
	}
	return s.repo.Create(w)
}

func (s *WorkLogService) Update(id string, u map[string]interface{}) error {
	w, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWorkLogNotFound
		}
		return err
	}
	for k, v := range u {
		switch k {
		case "log_date":
			if s, ok := v.(string); ok && s != "" {
				dt, _ := time.Parse("2006-01-02", s)
				w.LogDate = dt
			}
		case "hours":
			w.Hours = v.(float64)
		case "description":
			w.Description = v.(string)
		}
	}
	w.UpdatedAt = time.Now()
	return s.repo.Update(w)
}

func (s *WorkLogService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWorkLogNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}

// Stats 工时聚合统计
func (s *WorkLogService) Stats(userID *string, startDate, endDate, groupBy string) ([]model.WorkLog, error) {
	return s.repo.Stats(userID, startDate, endDate, groupBy)
}

func (s *WorkLogService) ExportCSV(userID, projectID *string, startDate, endDate string) ([]byte, error) {
	items, err := s.repo.FindForReport(userID, projectID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"日期", "用户ID", "任务", "项目", "客户", "工时", "描述"}); err != nil {
		return nil, err
	}
	for _, item := range items {
		taskName := ""
		if item.Task != nil {
			taskName = item.Task.Title
		}
		projectName := ""
		if item.Project != nil {
			projectName = item.Project.Name
		}
		customerName := ""
		if item.Customer != nil {
			customerName = item.Customer.Name
		}
		if err := writer.Write([]string{
			item.LogDate.Format("2006-01-02"),
			item.UserID,
			taskName,
			projectName,
			customerName,
			strconv.FormatFloat(item.Hours, 'f', -1, 64),
			item.Description,
		}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *WorkLogService) ForceDelete(id string) error { return s.repo.ForceDelete(id) }
func (s *WorkLogService) Restore(id string) error     { return s.repo.Restore(id) }
