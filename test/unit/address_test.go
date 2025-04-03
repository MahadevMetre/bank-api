package unittest

import (
	"bankapi/models"
	"testing"
)

func TestTrimLine(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Normal case with multiple parts",
			input:    []string{" 123 Main St ", "  Apt 4B ", " New York  "},
			expected: "123 Main St Apt 4B New York",
		},
		{
			name:     "Case with extra spaces at the beginning and end of parts",
			input:    []string{"  John ", "  Doe "},
			expected: "John Doe",
		},
		{
			name:     "Single part address with extra spaces",
			input:    []string{" 1234 Elm St  "},
			expected: "1234 Elm St",
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: "",
		},
		{
			name:     "All empty strings",
			input:    []string{"", " ", "   "},
			expected: "",
		},
		{
			name:     "Empty and non-empty strings mixed",
			input:    []string{"", "  Address ", "  "},
			expected: "Address",
		},
		{
			name:     "One non-empty string among empty parts",
			input:    []string{"  ", "  OnlyPart "},
			expected: "OnlyPart",
		},
		{
			name:     "No spaces at all",
			input:    []string{"Some", "Street", "Name"},
			expected: "Some Street Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.TrimLine(tt.input...)
			if got != tt.expected {
				t.Errorf("trimLine() = %v, want %v", got, tt.expected)
			}
		})
	}
}
