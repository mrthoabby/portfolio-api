package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid UUID v4", "123e4567-e89b-12d3-a456-426614174000", true},
		{"valid UUID v4 uppercase", "123E4567-E89B-12D3-A456-426614174000", true},
		{"invalid - too short", "123e4567-e89b-12d3-a456", false},
		{"invalid - wrong format", "not-a-uuid", false},
		{"invalid - empty", "", false},
		{"invalid - special chars", "123e4567-e89b-12d3-a456-42661417400g", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUUID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal string", "Hello World", "Hello World"},
		{"leading/trailing spaces", "  Hello World  ", "Hello World"},
		{"HTML entities", "<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"quotes", `"Hello" 'World'`, "&#34;Hello&#34; &#39;World&#39;"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal email", "user@example.com", "user@example.com"},
		{"uppercase", "User@Example.COM", "user@example.com"},
		{"with spaces", "  user@example.com  ", "user@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStripHTMLTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no tags", "Hello World", "Hello World"},
		{"simple tag", "<b>Hello</b>", "Hello"},
		{"script tag", "<script>alert('xss')</script>Hello", "alert('xss')Hello"},
		{"multiple tags", "<div><p>Hello</p></div>", "Hello"},
		{"with attributes", `<a href="http://evil.com">Click</a>`, "Click"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripHTMLTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeContactInput(t *testing.T) {
	input := SanitizeContactInput(
		"  <b>John</b> Doe  ",
		"  JOHN@EXAMPLE.COM  ",
		"<script>alert('xss')</script>Hello there!",
	)

	assert.Equal(t, "John Doe", input.Name)
	assert.Equal(t, "john@example.com", input.Email)
	assert.Equal(t, "alert(&#39;xss&#39;)Hello there!", input.Message)
}

