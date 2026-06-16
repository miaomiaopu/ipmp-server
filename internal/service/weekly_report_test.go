package service

import (
	"strings"
	"testing"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
)

func TestBuildWeeklyReportContentSummarizesWorkLogs(t *testing.T) {
	start := time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC)
	logs := []model.WorkLog{
		{
			LogDate:     start,
			Hours:       2,
			Description: "完成客户资料整理",
			Task:        &model.Task{Title: "客户管理"},
			Project:     &model.Project{Name: "IPMP"},
		},
		{
			LogDate:     start.AddDate(0, 0, 1),
			Hours:       1.5,
			Description: "修复工时统计",
			Task:        &model.Task{Title: "工时概览"},
			Project:     &model.Project{Name: "IPMP"},
		},
	}

	content := BuildWeeklyReportContent(model.ReportTypePersonal, logs, start, end)
	for _, want := range []string{
		"2026-06-08 至 2026-06-14",
		"总工时: 3.5h",
		"客户管理",
		"完成客户资料整理",
		"工时概览",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("weekly report content should contain %q, got:\n%s", want, content)
		}
	}
}
