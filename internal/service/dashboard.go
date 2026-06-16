package service

import (
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
)

type DashboardService struct {
	repo        *repository.DashboardRepository
	workLogRepo *repository.WorkLogRepository
}

func NewDashboardService(repo *repository.DashboardRepository, workLogRepo *repository.WorkLogRepository) *DashboardService {
	return &DashboardService{repo: repo, workLogRepo: workLogRepo}
}

func (s *DashboardService) Stats(userID, role string) (map[string]interface{}, error) {
	start, end := currentWeekRange(time.Now())
	var uid *string
	if role != model.RoleAdmin {
		uid = &userID
	}
	thisWeekHours, err := s.repo.SumWorkHours(uid, start, end)
	if err != nil {
		return nil, err
	}
	customers, err := s.repo.CountCustomers("")
	if err != nil {
		return nil, err
	}
	projects, err := s.repo.CountProjects("")
	if err != nil {
		return nil, err
	}
	activeProjects, err := s.repo.CountProjects(model.ProjectStatusInProgress)
	if err != nil {
		return nil, err
	}
	openTasks, err := s.repo.CountTasks(model.TaskStatusDone)
	if err != nil {
		return nil, err
	}
	openRequirements, err := s.repo.CountRequirements(model.ReqStatusDone)
	if err != nil {
		return nil, err
	}
	stats := map[string]interface{}{
		"customers":         customers,
		"projects":          projects,
		"active_projects":   activeProjects,
		"open_tasks":        openTasks,
		"open_requirements": openRequirements,
		"this_week_hours":   thisWeekHours,
	}
	if role == model.RoleAdmin {
		users, err := s.repo.CountUsers("")
		if err != nil {
			return nil, err
		}
		activeUsers, err := s.repo.CountUsers(model.UserStatusActive)
		if err != nil {
			return nil, err
		}
		stats["users"] = users
		stats["active_users"] = activeUsers
	}
	return stats, nil
}

func (s *DashboardService) ThisWeek(userID, role string) (map[string]interface{}, error) {
	start, end := currentWeekRange(time.Now())
	var uid *string
	if role != model.RoleAdmin {
		uid = &userID
	}
	logs, err := s.workLogRepo.FindForReport(uid, nil, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	total := 0.0
	for _, item := range logs {
		total += item.Hours
	}
	return map[string]interface{}{
		"week_start":  start.Format("2006-01-02"),
		"week_end":    end.Format("2006-01-02"),
		"total_hours": total,
		"work_logs":   logs,
	}, nil
}

func currentWeekRange(now time.Time) (time.Time, time.Time) {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-weekday)
	return start, start.AddDate(0, 0, 6)
}
