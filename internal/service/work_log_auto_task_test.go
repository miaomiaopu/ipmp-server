package service

import (
	"testing"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
)

func TestBuildAutoTaskForWorkLogUsesProjectAssociation(t *testing.T) {
	projectID := "project-1"
	customerID := "customer-1"
	logDate := time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC)
	workLog := &model.WorkLog{
		ProjectID:   &projectID,
		CustomerID:  &customerID,
		LogDate:     logDate,
		Description: "客户沟通",
	}

	task := buildAutoTaskForWorkLog(workLog)

	if task.TaskType != model.TaskTypeProject {
		t.Fatalf("expected project task, got %q", task.TaskType)
	}
	if task.ProjectID == nil || *task.ProjectID != projectID {
		t.Fatalf("expected project_id %q, got %#v", projectID, task.ProjectID)
	}
	if task.CustomerID == nil || *task.CustomerID != customerID {
		t.Fatalf("expected customer_id %q, got %#v", customerID, task.CustomerID)
	}
	if task.Title != "[2026-06-16] 客户沟通" {
		t.Fatalf("expected title from description, got %q", task.Title)
	}
	if task.Status != model.TaskStatusDone {
		t.Fatalf("expected status %q, got %q", model.TaskStatusDone, task.Status)
	}
}

func TestBuildAutoTaskForWorkLogFallsBackToDailyTask(t *testing.T) {
	logDate := time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC)
	workLog := &model.WorkLog{LogDate: logDate}

	task := buildAutoTaskForWorkLog(workLog)

	if task.TaskType != model.TaskTypeDaily {
		t.Fatalf("expected daily task, got %q", task.TaskType)
	}
	if task.ProjectID != nil || task.CustomerID != nil {
		t.Fatalf("expected no relation for daily task, got project=%#v customer=%#v", task.ProjectID, task.CustomerID)
	}
	if task.Title != "[2026-06-16] 工时" {
		t.Fatalf("expected fallback title with log date, got %q", task.Title)
	}
	if task.Status != model.TaskStatusDone {
		t.Fatalf("expected status %q, got %q", model.TaskStatusDone, task.Status)
	}
}
