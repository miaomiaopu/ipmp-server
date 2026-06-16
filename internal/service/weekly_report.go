package service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var ErrWeeklyReportNotFound = errors.New("weekly report not found")

type WeeklyReportService struct {
	repo        *repository.WeeklyReportRepository
	workLogRepo *repository.WorkLogRepository
}

type GenerateWeeklyReportInput struct {
	ReportType string
	WeekStart  string
	WeekEnd    string
	ProjectID  *string
	UseAI      bool
}

func NewWeeklyReportService(repo *repository.WeeklyReportRepository, workLogRepo *repository.WorkLogRepository) *WeeklyReportService {
	return &WeeklyReportService{repo: repo, workLogRepo: workLogRepo}
}

func (s *WeeklyReportService) List(page, pageSize int, userID *string) ([]model.WeeklyReport, int64, error) {
	return s.repo.List(page, pageSize, userID)
}

func (s *WeeklyReportService) GetByID(id string) (*model.WeeklyReport, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWeeklyReportNotFound
		}
		return nil, err
	}
	return item, nil
}

func (s *WeeklyReportService) Generate(userID string, input GenerateWeeklyReportInput) (*model.WeeklyReport, error) {
	start, err := time.Parse("2006-01-02", input.WeekStart)
	if err != nil {
		return nil, fmt.Errorf("invalid week_start")
	}
	end, err := time.Parse("2006-01-02", input.WeekEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid week_end")
	}
	reportType := input.ReportType
	if reportType == "" {
		reportType = model.ReportTypePersonal
	}
	var filterUserID *string
	var filterProjectID *string
	switch reportType {
	case model.ReportTypePersonal:
		filterUserID = &userID
	case model.ReportTypeProject:
		if input.ProjectID == nil || *input.ProjectID == "" {
			return nil, fmt.Errorf("project_id required for project report")
		}
		filterProjectID = input.ProjectID
	default:
		return nil, fmt.Errorf("invalid report_type")
	}

	logs, err := s.workLogRepo.FindForReport(filterUserID, filterProjectID, input.WeekStart, input.WeekEnd)
	if err != nil {
		return nil, err
	}
	content := BuildWeeklyReportContent(reportType, logs, start, end)
	var raw *string
	if input.UseAI {
		aiRaw := MockAIReport(content)
		raw = &aiRaw
		content = aiRaw
	}
	report := &model.WeeklyReport{
		UserID:       userID,
		WeekStart:    start,
		WeekEnd:      end,
		ReportType:   reportType,
		ProjectID:    input.ProjectID,
		Content:      content,
		AIRawContent: raw,
		Status:       model.ReportStatusDraft,
	}
	if err := s.repo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *WeeklyReportService) UpdateContent(id, content, status string) error {
	item, err := s.GetByID(id)
	if err != nil {
		return err
	}
	item.Content = content
	if status != "" {
		item.Status = status
	}
	return s.repo.Update(item)
}

func (s *WeeklyReportService) Review(id string) error {
	item, err := s.GetByID(id)
	if err != nil {
		return err
	}
	item.Status = model.ReportStatusReviewed
	return s.repo.Update(item)
}

func (s *WeeklyReportService) Finalize(id string) error {
	item, err := s.GetByID(id)
	if err != nil {
		return err
	}
	item.Status = model.ReportStatusFinal
	return s.repo.Update(item)
}

func BuildWeeklyReportContent(reportType string, logs []model.WorkLog, weekStart, weekEnd time.Time) string {
	sort.SliceStable(logs, func(i, j int) bool {
		return logs[i].LogDate.Before(logs[j].LogDate)
	})
	total := 0.0
	var b strings.Builder
	title := "个人周报"
	if reportType == model.ReportTypeProject {
		title = "项目周报"
	}
	b.WriteString(title)
	b.WriteString("\n周期: ")
	b.WriteString(weekStart.Format("2006-01-02"))
	b.WriteString(" 至 ")
	b.WriteString(weekEnd.Format("2006-01-02"))
	b.WriteString("\n\n")
	for _, item := range logs {
		total += item.Hours
	}
	b.WriteString("总工时: ")
	b.WriteString(formatHours(total))
	b.WriteString("h\n\n")
	if len(logs) == 0 {
		b.WriteString("本周期暂无工时记录。\n")
		return b.String()
	}
	b.WriteString("工作明细:\n")
	for _, item := range logs {
		taskName := "-"
		if item.Task != nil && item.Task.Title != "" {
			taskName = item.Task.Title
		}
		projectName := "-"
		if item.Project != nil && item.Project.Name != "" {
			projectName = item.Project.Name
		}
		b.WriteString("- ")
		b.WriteString(item.LogDate.Format("2006-01-02"))
		b.WriteString(" [")
		b.WriteString(projectName)
		b.WriteString(" / ")
		b.WriteString(taskName)
		b.WriteString("] ")
		b.WriteString(formatHours(item.Hours))
		b.WriteString("h")
		if item.Description != "" {
			b.WriteString(": ")
			b.WriteString(item.Description)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func MockAIReport(content string) string {
	return "AI mock 周报草稿\n\n" + content
}

func formatHours(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
