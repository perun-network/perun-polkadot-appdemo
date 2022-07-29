package cli

import "strings"

func buildHelp(commands []Command) string {
	var sb strings.Builder
	for _, c := range commands {
		help := c.Name + "\n" + c.Help
		sb.WriteString("\n" + help + "\n")
	}
	return sb.String()
}
