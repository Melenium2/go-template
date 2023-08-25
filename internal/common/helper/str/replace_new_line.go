package str

import "strings"

func ReplaceNewLine2Whitespace(str string) string {
	return ReplaceNewLine(str, " ")
}

func ReplaceNewLine(str, replace string) string {
	if len(str) == 0 {
		return ""
	}

	res := strings.Fields(str)

	return strings.Join(res, replace)
}
