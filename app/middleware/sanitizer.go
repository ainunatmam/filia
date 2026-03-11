package middleware

import (
	"html"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SanitizeHTML sanitizes HTML content to prevent XSS
func SanitizeHTML(content string) string {
	// First escape HTML
	content = html.EscapeString(content)

	// Remove potentially dangerous patterns
	dangerousPatterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`onerror=`,
		`onload=`,
		`onclick=`,
		`onmouseover=`,
		`<iframe[^>]*>.*?</iframe>`,
		`<object[^>]*>.*?</object>`,
		`<embed[^>]*>`,
	}

	for _, pattern := range dangerousPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		content = re.ReplaceAllString(content, "")
	}

	return content
}

// sanitize performs general input sanitization
func sanitize(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Escape HTML
	input = html.EscapeString(input)

	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\t' || r >= 32 {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SanitizeInput middleware provides input sanitization utilities
// Note: For full sanitization, apply sanitize() to each input field in handlers
func SanitizeInput() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add sanitization helpers to context
		c.Locals("sanitize", func(input string) string {
			return sanitize(input)
		})

		c.Locals("sanitizeHTML", func(content string) string {
			return SanitizeHTML(content)
		})

		return c.Next()
	}
}
