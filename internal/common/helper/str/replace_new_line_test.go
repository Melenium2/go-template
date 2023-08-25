package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceNewLine(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "should do nothing if string not contains any new line symbols",
			str:      "simple string",
			expected: "simple string",
		},
		{
			name:     "should replace '\n' symbol to whitespace",
			str:      "simple\nstring",
			expected: "simple string",
		},
		{
			name:     "should replace '\r\n' symbols to whitespace",
			str:      "simple\r\nstring",
			expected: "simple string",
		},
		{
			name:     "should replace '\r' symbols to whitespace",
			str:      "simple\rstring",
			expected: "simple string",
		},
		{
			name:     "should replace '\n' symbol at the end of the string",
			str:      "simple string\n",
			expected: "simple string",
		},
	}

	t.Parallel()

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result := ReplaceNewLine(tc.str, " ")
			assert.Equal(t, tc.expected, result)
		})
	}
}
