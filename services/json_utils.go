package services

import (
	"regexp"
	"strings"
)

// extractJSON extracts JSON from markdown code blocks or returns the original string
func extractJSON(text string) string {
	// Try to extract JSON from markdown code blocks (```json ... ```)
	jsonBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*({.*?})\\s*```")
	matches := jsonBlockRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find JSON object directly
	jsonRegex := regexp.MustCompile("(?s)({\\s*\".*?})")
	matches = jsonRegex.FindStringSubmatch(text)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[0])
	}

	// Return original if no JSON found
	return strings.TrimSpace(text)
}