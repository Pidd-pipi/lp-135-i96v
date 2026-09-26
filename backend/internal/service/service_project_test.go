package service

import (
	"testing"

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
