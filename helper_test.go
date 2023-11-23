package main

import "testing"

func TestIsisValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"with space", "google .ca", false},
		{"without space", "google.ca", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidInput(tt.input)
			if got != tt.want {
				t.Errorf("isValidInput(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
