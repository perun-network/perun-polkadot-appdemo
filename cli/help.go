package cli

import "strings"

func buildHelp(commands []Command) string {
	var sb strings.Builder
	for _, c := range commands {
		sb.WriteString(c.Help)
		sb.WriteString("\n")
	}
	return sb.String()
}
