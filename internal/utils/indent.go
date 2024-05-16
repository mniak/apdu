package utils

import (
	"strings"
)

func IndentString(text string, indentation string) string {
	return indentation + strings.ReplaceAll(text, "\n", "\n"+indentation)
}
