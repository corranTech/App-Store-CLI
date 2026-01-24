package asc

import "strings"

func compactWhitespace(input string) string {
	clean := sanitizeTerminal(input)
	return strings.Join(strings.Fields(clean), " ")
}

func escapeMarkdown(input string) string {
	clean := compactWhitespace(input)
	return strings.ReplaceAll(clean, "|", "\\|")
}
