package str

import (
	"strings"
)

// http://apps.timwhitlock.info/emoji/tables/unicode
var emojiUnicodeRanges = [][]uint{
	{0x1F600, 0x1F92F},
	{0x1F300, 0x1F5FF},
	{0x1F680, 0x1F6FF},
	{0x1F190, 0x1F1FF},
	{0x2702, 0x27B0},
	{0x1F926, 0x1FA9F},
	{0x200d, 0x200d},
	{0x2640, 0x2642},
	{0x2600, 0x2B55},
	{0x23cf, 0x23cf},
	{0x23e9, 0x23e9},
	{0x231a, 0x231a},
	{0xfe0f, 0xfe0f},
}

var replaceExceptions = map[rune]struct{}{
	'$':  {},
	'&':  {},
	'+':  {},
	':':  {},
	';':  {},
	'=':  {},
	'?':  {},
	'@':  {},
	'#':  {},
	'|':  {},
	'\'': {},
	'"':  {},
	'<':  {},
	'>':  {},
	'.':  {},
	'-':  {},
	'^':  {},
	'*':  {},
	'(':  {},
	'}':  {},
	'{':  {},
	',':  {},
	')':  {},
	'%':  {},
	'!':  {},
}

func ReplaceEmoji2Empty(str string) string {
	return ReplaceEmoji(str, "")
}

// ReplaceEmoji replaces all emojis from the string to provided `replace` string.
// This function can not replace some emojis, but most of all can be replaced.
func ReplaceEmoji(str, replace string) string {
	if len(str) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.Grow(len(str))

	for _, r := range str {
		currStr := string(r)
		code := uint(r)

		_, ok := replaceExceptions[r]

		for j := 0; j < len(emojiUnicodeRanges) && !ok; j++ {
			leftBoundary, rightBoundary := emojiUnicodeRanges[j][0], emojiUnicodeRanges[j][1]

			if leftBoundary <= code && code <= rightBoundary {
				builder.Grow(len(replace))

				currStr = replace

				break
			}
		}

		builder.WriteString(currStr)
	}

	return builder.String()
}
