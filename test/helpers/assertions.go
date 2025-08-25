package helpers

import (
	"strings"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// AssertPositionEquals compares two Position structs for equality
func AssertPositionEquals(t *testing.T, expected, actual references.Position) {
	t.Helper()
	if expected.X != actual.X || expected.Y != actual.Y {
		t.Errorf("Position mismatch: expected (%d,%d), got (%d,%d)",
			expected.X, expected.Y, actual.X, actual.Y)
	}
}

// AssertStringContains checks if actual string contains expected substring
func AssertStringContains(t *testing.T, expected, actual string) {
	t.Helper()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected string to contain '%s', got: %s", expected, actual)
	}
}

// AssertInventoryQuantity checks inventory quantity matches expected value
func AssertInventoryQuantity(t *testing.T, expected uint16, actual uint16, itemName string) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %s quantity %d, got %d", itemName, expected, actual)
	}
}
