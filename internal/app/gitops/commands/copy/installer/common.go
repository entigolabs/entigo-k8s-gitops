package installer

import (
	"fmt"
	"strings"
)

func formatCommaSeparatedString(str string) string {
	splits := strings.Split(str, ",")
	formattedString := strings.Join(splits, ", ")
	return fmt.Sprintf("[%s]", formattedString)
}
