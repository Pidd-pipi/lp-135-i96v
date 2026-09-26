package util

import "testing"

func TestNormalizePage(t *testing.T) {
	tests := []struct {
		name string
		p    int
		ps   int
		wP   int
		wPS  int
	}{
		{name: "defaults", p: 0, ps: 0, wP: 1, wPS: 10},
		{name: "negative", p: -1, ps: -5, wP: 1, wPS: 10},
		{name: "valid", p: 2, ps: 20, wP: 2, wPS: 20},
		{name: "oversize", p: 3, ps: 500, wP: 3, wPS: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, gotPS := NormalizePage(tt.p, tt.ps)
			if gotP != tt.wP || gotPS != tt.wPS {
				t.Fatalf("NormalizePage(%d,%d)=%d,%d want %d,%d", tt.p, tt.ps, gotP, gotPS, tt.wP, tt.wPS)
			}
		})
	}
}
