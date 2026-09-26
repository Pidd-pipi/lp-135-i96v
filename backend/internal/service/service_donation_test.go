package service

import (
	"strings"
	"testing"
)

func TestDonationCreateInputValidation(t *testing.T) {
	svc := &DonationService{}
	tests := []struct {
		name    string
		in      CreateInput
		wantErr string
	}{
		{name: "missing project", in: CreateInput{ProjectID: 0, Amount: 10}, wantErr: "projectId"},
		{name: "zero amount", in: CreateInput{ProjectID: 1, Amount: 0}, wantErr: "amount"},
		{name: "negative amount", in: CreateInput{ProjectID: 1, Amount: -5}, wantErr: "amount"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(1, tt.in)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewReferenceNo(t *testing.T) {
	tests := []string{"CERT", "TXN", "ORDER"}
	for _, prefix := range tests {
		t.Run(prefix, func(t *testing.T) {
			no, err := newReferenceNo(prefix)
			if err != nil {
				t.Fatalf("newReferenceNo(%q): %v", prefix, err)
			}
			if !strings.HasPrefix(no, prefix) {
				t.Fatalf("reference %q missing prefix %q", no, prefix)
			}
			if len(no) != len(prefix)+12 {
				t.Fatalf("reference %q has unexpected length %d", no, len(no))
			}
		})
	}
}
