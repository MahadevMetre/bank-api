package unittest

import (
	"bankapi/utils"
	"testing"
)

func TestCalculateTimeDifference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Test case 1: future date with full timestamp",
			input:    "2025-02-21 17:38:58.0",
			expected: "1 Minute", // Current time is 5:35 PM, so 5:38:58 PM is 3 minutes later.
		},
		{
			name:     "Test case 2: future date with full timestamp",
			input:    "2025-02-21 18:38:58.0",
			expected: "1 Hour", // Current time is 5:35 PM, so 6:38:58 PM is 1 hour and 3 minutes later.
		},
		{
			name:     "Test case 3: past date",
			input:    "2025-02-21 15:04:05",
			expected: "2025-02-21 15:04:05", // Past date should return the input as is.
		},
		{
			name:     "Test case 4: future date with AM/PM",
			input:    "23-FEB-2025 11.12:26.000000 PM",
			expected: "2 Hour", // Current time is 5:35 PM on 21-Feb-2025, so 11:12:26 PM on 23-Feb-2025 is 2 hours ahead.
		},
		{
			name:     "Test case 5: same day, 1 hour later",
			input:    "2025-02-21 18:38:00",
			expected: "57 Minute", // The input is exactly 1 hour later than 5:35 PM.
		},
		{
			name:     "Test case 6: valid time with month in uppercase",
			input:    "23-FEB-2025 10:00:00 AM",
			expected: "1 Day", // The current time is 5:35 PM on 21-Feb-2025, and the input is 10 AM on 23-Feb-2025, so 1 day later.
		},
		{
			name:     "Test case 7: invalid date format",
			input:    "invalid-date-format",
			expected: "invalid-date-format", // Invalid input should return the same string.
		},
		{
			name:     "Test case 8: future date with incorrect timezone format",
			input:    "2025-02-21T10:00:00+9999",
			expected: "2025-02-21T10:00:00+9999", // The input has an incorrect timezone format, so return it unchanged.
		},
		{
			name:     "Test case 9: very large time difference",
			input:    "2050-02-21T10:00:00Z",
			expected: "2050-02-21T10:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CalculateTimeDifference(tt.input)
			if got != tt.expected {
				t.Errorf("CalculateTimeDifference() = %v, want %v", got, tt.expected)
			}
		})
	}
}
