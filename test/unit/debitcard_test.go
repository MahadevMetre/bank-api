package unittest

import (
	debitcard "bankapi/stores/debit_card"
	"testing"
)

// TestGetName tests the getName function with various cases.
func TestGetName(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Valid full name",
			input:    []string{"John", "Doe", ""},
			expected: "John Doe",
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    []string{"   "},
			expected: "",
		},
		{
			name:     "Name with spaces around",
			input:    []string{"   John   ", "  Doe  ", ""},
			expected: "John Doe",
		},
		{
			name:     "Long name within limit",
			input:    []string{"Alexander", "Hamilton", ""},
			expected: "Alexander Hamilton",
		},
		{
			name:     "Exceeds character limit",
			input:    []string{"A very long name that will exceed the fifty characters limit", "with more text", "even more text"},
			expected: "",
		},
		{
			name:     "Name exactly at 50 characters",
			input:    []string{"ThisNameHasExactlyFiftyCharacters1234567890", "", ""},
			expected: "ThisNameHasExactlyFiftyCharacters1234567890",
		},
		{
			name:     "Mixed empty and valid names",
			input:    []string{"M", "  ", "John"},
			expected: "M John",
		},
		{
			name:     "Multiple names trimmed correctly",
			input:    []string{"  Alice   ", "   Bob  ", " Charlie "},
			expected: "Alice Bob Charlie",
		},
		{
			name:     "Single long word exceeding 50 characters",
			input:    []string{"SupercalifragilisticexpialidociousSuperLong", " hzjxfgbcfa .asfigtycxvb ASL:ejdd ncugasb", "cZ:clxjMIODGHD"},
			expected: "SupercalifragilisticexpialidociousSuperLong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := debitcard.NameOnDebitCard(tt.input...)
			if result != tt.expected {
				t.Errorf("Expected: %q, Got: %q", tt.expected, result)
			}
		})
	}
}
