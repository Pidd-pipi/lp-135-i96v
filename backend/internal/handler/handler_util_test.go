package handler

import "testing"

func TestAtoi(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{name: "positive", in: "12", want: 12},
		{name: "zero", in: "0", want: 0},
		{name: "negative", in: "-3", want: -3},
		{name: "invalid", in: "abc", want: 0},
		{name: "empty", in: "", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atoi(tt.in); got != tt.want {
				t.Fatalf("atoi(%q)=%d want %d", tt.in, got, tt.want)
			}
		})
	}
}
