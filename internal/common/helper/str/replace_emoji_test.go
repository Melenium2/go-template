package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceEmoji(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "should do nothing if string not contains emoji",
			str:      "simple string",
			expected: "simple string",
		},
		{
			name:     "should do nothing if provided string is empty",
			str:      "",
			expected: "",
		},
		{
			name:     "should replace only emojis to empty string",
			str:      "simple string with 🤩🤩🤩",
			expected: "simple string with ",
		},
		{
			name:     "should replace only emojis and do nothing with other symbols",
			str:      "simple string with 🤩🤩🤩 and ;:.,-+='\"[]{}<>%$#@!^&*()",
			expected: "simple string with  and ;:.,-+='\"[]{}<>%$#@!^&*()",
		},
		{
			name:     "should do nothing if string contains only simple symbols",
			str:      ";:.,-+='\"[]{}<>%$#@!^&*()",
			expected: ";:.,-+='\"[]{}<>%$#@!^&*()",
		},
	}

	t.Parallel()

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result := ReplaceEmoji(tc.str, "")
			assert.Equal(t, tc.expected, result)
		})
	}
}
