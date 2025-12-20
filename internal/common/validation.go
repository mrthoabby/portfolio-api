package common

import (
	"html"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// SanitizeString removes potentially dangerous content from strings
// This includes HTML tags and trims whitespace
func SanitizeString(s string) string {
	// Trim whitespace
	s = strings.TrimSpace(s)

	// Escape HTML entities
	s = html.EscapeString(s)

	return s
}

// SanitizeEmail validates and sanitizes email addresses
func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}

// StripHTMLTags removes all HTML tags from a string
func StripHTMLTags(s string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")

	// Trim whitespace
	s = strings.TrimSpace(s)

	return s
}

// ValidateAndSanitizeContact sanitizes contact form inputs
type ContactInput struct {
	Name    string
	Email   string
	Message string
}

func SanitizeContactInput(name, email, message string) ContactInput {
	return ContactInput{
		Name:    SanitizeString(StripHTMLTags(name)),
		Email:   SanitizeEmail(email),
		Message: SanitizeString(StripHTMLTags(message)),
	}
}

