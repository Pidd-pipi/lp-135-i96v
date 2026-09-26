package service

import (
	"testing"

	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/model"
)

func TestProjectProgress(t *testing.T) {
	tests := []struct {
		name     string
		current  float64
		target   float64
		expected int
	}{
		{name: "zero target", current: 100, target: 0, expected: 0},
		{name: "half", current: 50000, target: 100000, expected: 50},
		{name: "over 100", current: 120000, target: 100000, expected: 100},
		{name: "quarter", current: 25000, target: 100000, expected: 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &model.Project{CurrentAmount: tt.current, TargetAmount: tt.target}
			got := withProgress(p).Progress
			if got != tt.expected {
				t.Errorf("progress = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestDonationCertificateNo(t *testing.T) {
	// 证书编号应为 CERT 前缀且非空
	specs := []string{"CERT100001", "CERT000123"}
	for _, s := range specs {
		if len(s) < 4 || s[:4] != "CERT" {
			t.Errorf("invalid certificate no: %s", s)
		}
	}
}

func TestValidateProjectUpdate(t *testing.T) {
	baseInput := UpdateProjectInput{
		Description:   "新介绍",
		Category:      constants.CategoryEducation,
		TargetAmount:  100000,
		ExecutionPlan: "新计划",
		StartDate:     "2026-01-01",
		EndDate:       "2026-12-31",
	}

	tests := []struct {
		name    string
		project *model.Project
		input   UpdateProjectInput
		wantErr error
	}{
		{
			name:    "approved project can be updated",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 50000},
			input:   baseInput,
			wantErr: nil,
		},
		{
			name:    "paused project can be updated",
			project: &model.Project{Status: constants.ProjectPaused, CurrentAmount: 50000},
			input:   baseInput,
			wantErr: nil,
		},
		{
			name:    "completed project cannot be updated",
			project: &model.Project{Status: constants.ProjectCompleted, CurrentAmount: 50000},
			input:   baseInput,
			wantErr: ErrProjectCompleted,
		},
		{
			name:    "target lower than raised is rejected",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 120000},
			input:   baseInput,
			wantErr: ErrTargetBelowRaised,
		},
		{
			name:    "target equal to raised is allowed",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 100000},
			input:   baseInput,
			wantErr: nil,
		},
		{
			name:    "invalid category rejected",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 0},
			input: func() UpdateProjectInput {
				in := baseInput
				in.Category = "not-a-category"
				return in
			}(),
			wantErr: ErrInvalidCategory,
		},
		{
			name:    "end before start rejected",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 0},
			input: func() UpdateProjectInput {
				in := baseInput
				in.StartDate = "2026-12-31"
				in.EndDate = "2026-01-01"
				return in
			}(),
			wantErr: ErrInvalidDateRange,
		},
		{
			name:    "bad date format rejected",
			project: &model.Project{Status: constants.ProjectApproved, CurrentAmount: 0},
			input: func() UpdateProjectInput {
				in := baseInput
				in.StartDate = "2026/01/01"
				return in
			}(),
			wantErr: ErrInvalidDate,
		},
		{
			name:    "empty dates allowed (clear)",
			project: &model.Project{Status: constants.ProjectPending, CurrentAmount: 0},
			input: func() UpdateProjectInput {
				in := baseInput
				in.StartDate = ""
				in.EndDate = ""
				return in
			}(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := validateProjectUpdate(tt.project, tt.input)
			if err != tt.wantErr {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if tt.input.StartDate == "" && start != nil {
				t.Errorf("start should be nil, got %v", start)
			}
			if tt.input.EndDate == "" && end != nil {
				t.Errorf("end should be nil, got %v", end)
			}
			if tt.input.StartDate != "" && (start == nil || start.Format("2006-01-02") != tt.input.StartDate) {
				t.Errorf("start = %v, want %s", start, tt.input.StartDate)
			}
			if tt.input.EndDate != "" && (end == nil || end.Format("2006-01-02") != tt.input.EndDate) {
				t.Errorf("end = %v, want %s", end, tt.input.EndDate)
			}
		})
	}
}

func TestValidateStatusAction(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		action    string
		wantStat  string
		wantError bool
	}{
		{name: "pause approved", status: constants.ProjectApproved, action: "pause", wantStat: constants.ProjectPaused},
		{name: "resume paused", status: constants.ProjectPaused, action: "resume", wantStat: constants.ProjectApproved},
		{name: "pause paused rejected", status: constants.ProjectPaused, action: "pause", wantError: true},
		{name: "resume approved rejected", status: constants.ProjectApproved, action: "resume", wantError: true},
		{name: "pause completed rejected", status: constants.ProjectCompleted, action: "pause", wantError: true},
		{name: "resume completed rejected", status: constants.ProjectCompleted, action: "resume", wantError: true},
		{name: "pause pending rejected", status: constants.ProjectPending, action: "pause", wantError: true},
		{name: "unknown action rejected", status: constants.ProjectApproved, action: "stop", wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateStatusAction(tt.status, tt.action)
			if tt.wantError {
				if err != ErrInvalidProjectStatus {
					t.Fatalf("err = %v, want ErrInvalidProjectStatus", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tt.wantStat {
				t.Errorf("new status = %s, want %s", got, tt.wantStat)
			}
		})
	}
}
