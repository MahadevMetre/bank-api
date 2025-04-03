package unittest

import (
	"bankapi/utils"
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func TestHandleResponse(t *testing.T) {
	mockBody := "response body"
	mockError := errors.New("request error")

	tests := []struct {
		resp     *http.Response
		err      error
		expected string
		hasError bool
	}{
		{
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(mockBody)),
			},
			err:      nil,
			expected: mockBody,
			hasError: false,
		},
		{
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(mockBody)),
			},
			err:      nil,
			expected: mockBody,
			hasError: true,
		},
		{
			resp:     nil,
			err:      mockError,
			expected: "",
			hasError: true,
		},
	}

	for _, test := range tests {
		body, err := utils.HandleResponse(test.resp, test.err)
		if test.hasError && err == nil {
			t.Errorf("expected an error, got none")
		} else if !test.hasError && err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if string(body) != test.expected {
			t.Errorf("expected body %s, got %s", test.expected, string(body))
		}
	}
}

func TestStructToMap(t *testing.T) {
	p := Person{FirstName: "John", LastName: "Doe", Age: 30}
	expected := map[string]interface{}{
		"FirstName": "John",
		"LastName":  "Doe",
		"Age":       float64(30),
	}

	result, err := utils.StructToMap(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGenerateReferralCode(t *testing.T) {
	name := "John"
	referralCode, err := utils.GenerateReferralCode(name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(referralCode) != 6 {
		t.Errorf("expected referral code length 6, got %d", len(referralCode))
	}

	if referralCode[:1] != "J" {
		t.Errorf("expected referral code to start with 'J', got %s", referralCode[:1])
	}
}

func TestParseNameForDebitCard(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input:    "JAKIR ALI ASGAR HUSSAIN SAYYED ASGAR HUSSAIN Sayyed",
			expected: map[string]string{"first_name": "JAKIR", "middle_name": "ALI ASGAR HUSSAIN SAYYED ASGAR HUSSAIN", "last_name": "Sayyed"},
		},
		{
			input:    "SEEMA SHAILESH KHETLE khetle",
			expected: map[string]string{"first_name": "SEEMA", "middle_name": "SHAILESH KHETLE", "last_name": "khetle"},
		},
	}
	for _, test := range tests {
		result := utils.ParseName(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseName(%q) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseName(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input:    "John",
			expected: map[string]string{"first_name": "John"},
		},
		{
			input:    "John Doe",
			expected: map[string]string{"first_name": "John", "last_name": "Doe"},
		},
		{
			input:    "John Middle Doe",
			expected: map[string]string{"first_name": "John", "middle_name": "Middle", "last_name": "Doe"},
		},
		{
			input:    "John Multiple Middle Names Doe",
			expected: map[string]string{"first_name": "John", "middle_name": "Multiple Middle Names", "last_name": "Doe"},
		},
		{
			input:    "",
			expected: map[string]string{},
		},

		{
			input:    "John-Multiple-Middle-Names-Doe",
			expected: map[string]string{"first_name": "John-Multiple-Middle-Names-Doe"},
		},
	}

	for _, test := range tests {
		result := utils.ParseName(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("parseName(%q) = %v; expected %v", test.input, result, test.expected)
		}
	}
}
