package testsupport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertContainsExactly[T any](t *testing.T, expected []T, actual []T) {
	assert.Len(t, actual, len(expected), "expected to contain exactly %d element(s)", len(expected))

	for _, expectedValue := range expected {
		assert.Contains(t, actual, expectedValue, "expected to contain %s", expectedValue)
	}
}
