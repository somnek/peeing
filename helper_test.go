package main

import (
	"testing"
	"time"
)

func TestIsisValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"with space", "google .ca", false},
		{"without space", "google.ca", true},
		{"without dot", "googleca", false},
		{"empty input", "", false},
		{"just space", " ", false},
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

func TestConvertToBlockUnit(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  rune
	}{
		{"minimum", 0 * time.Millisecond, BARS[0]},
		{"mid", 200 * time.Millisecond, BARS[3]},
		{"max", 400 * time.Millisecond, BARS[7]},
		{"exceed max", 500 * time.Millisecond, BARS[7]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToBlockUnit(tt.input)
			if got != tt.want {
				t.Errorf("convertToBlockUnit(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})

	}
}
