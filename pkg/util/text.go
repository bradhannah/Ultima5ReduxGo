package util

import "strings"

func GetCenteredText(text string) string {
	const maxCharsPerLine = 16

	textLen := len(text)
	if textLen >= maxCharsPerLine {
		// If the text is longer than or equal to the max width, just truncate or return as is
		return text[:maxCharsPerLine]
	}

	// Calculate padding
	totalPadding := maxCharsPerLine - textLen
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	// Construct the centered text with spaces
	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
}
