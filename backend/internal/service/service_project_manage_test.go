package service

import (
	"testing"
	"time"

	"github.com/givetrack/givetrack/internal/constants"
)

func TestCents(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want int64
	}{
		{name: "integer", in: 1000, want: 100000},
		{name: "two decimals", in: 99.99, want: 9999},
		{name: "float noise rounds", in: 0.29, want: 29},
		{name: "zero", in: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cents(tt.in); got != tt.want {
				t.Errorf("cents(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidCategory(t *testing.T) {
	for _, c := range []string{
		constants.CategoryEducation, constants.CategoryElderly, constants.CategoryMedical,
		constants.CategoryDisaster, constants.CategoryEnvironment, constants.CategoryOther,
	} {
		if !validCategory(c) {
			t.Errorf("category %q should be valid", c)
		}
	}
	for _, c := range []string{"", "hacker", "EDUCATION"} {
		if validCategory(c) {
			t.Errorf("category %q should be invalid", c)
		}
	}
}

func TestParseProjectDates(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		s, e, err := parseProjectDates("", "")
		if err != nil || s != nil || e != nil {
			t.Fatalf("got (%v,%v,%v), want nil,nil,nil", s, e, err)
		}
	})
	t.Run("valid range", func(t *testing.T) {
		s, e, err := parseProjectDates("2026-01-01", "2026-12-31")
		if err != nil || s == nil || e == nil {
			t.Fatalf("unexpected: %v %v %v", s, e, err)
		}
	})
	t.Run("same day allowed", func(t *testing.T) {
		if _, _, err := parseProjectDates("2026-01-01", "2026-01-01"); err != nil {
			t.Fatalf("same day should be allowed, got %v", err)
		}
	})
	t.Run("end before start rejected", func(t *testing.T) {
		if _, _, err := parseProjectDates("2026-12-31", "2026-01-01"); err == nil {
			t.Fatal("expected error when endDate is earlier than startDate")
		}
	})
	t.Run("bad format rejected", func(t *testing.T) {
		if _, _, err := parseProjectDates("2026/01/01", ""); err == nil {
			t.Fatal("expected error for malformed startDate")
		}
		if _, _, err := parseProjectDates("", "31-01-2026"); err == nil {
			t.Fatal("expected error for malformed endDate")
		}
	})
}

// TestStatusTransitionRules 校验停募/重开的状态机规则。
func TestStatusTransitionRules(t *testing.T) {
	canPause := func(status string) bool { return status == constants.ProjectApproved }
	canReopen := func(status string) bool { return status == constants.ProjectPaused }
	canEdit := func(status string) bool { return editableStatuses[status] }

	ok := map[string]bool{
		constants.ProjectApproved:  true,
		constants.ProjectPending:   false,
		constants.ProjectRejected:  false,
		constants.ProjectCompleted: false,
		constants.ProjectPaused:    false,
	}
	for status, want := range ok {
		if canPause(status) != want {
			t.Errorf("pause allowed from %s = %v, want %v", status, !want, want)
		}
	}
	reopen := map[string]bool{
		constants.ProjectPaused:    true,
		constants.ProjectApproved:  false,
		constants.ProjectCompleted: false,
		constants.ProjectPending:   false,
	}
	for status, want := range reopen {
		if canReopen(status) != want {
			t.Errorf("reopen allowed from %s = %v, want %v", status, !want, want)
		}
	}
	edit := map[string]bool{
		constants.ProjectPending:   true,
		constants.ProjectApproved:  true,
		constants.ProjectPaused:    true,
		constants.ProjectCompleted: false,
		constants.ProjectRejected:  false,
	}
	for status, want := range edit {
		if canEdit(status) != want {
			t.Errorf("edit allowed in %s = %v, want %v", status, !want, want)
		}
	}
}

// TestTargetAmountGuard 目标金额低于已筹金额必须拒绝（以分为单位精确比较）。
func TestTargetAmountGuard(t *testing.T) {
	current := 65000.0
	cases := []struct {
		target float64
		allow  bool
	}{
		{target: 65000, allow: true},  // 等于已筹：允许
		{target: 100000, allow: true}, // 高于已筹：允许
		{target: 64999.99, allow: false},
	}
	for _, c := range cases {
		got := cents(c.target) >= cents(current)
		if got != c.allow {
			t.Errorf("target %.2f vs raised %.2f allowed=%v, want %v", c.target, current, got, c.allow)
		}
	}
}

func TestParseProjectDateValues(t *testing.T) {
	s, e, err := parseProjectDates("2026-03-01", "2026-04-01")
	if err != nil {
		t.Fatal(err)
	}
	wantS := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	wantE := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	if !s.Equal(wantS) || !e.Equal(wantE) {
		t.Errorf("parsed dates = %v, %v, want %v, %v", s, e, wantS, wantE)
	}
}
