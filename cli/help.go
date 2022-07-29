package cli

import "strings"

func buildHelp(commands []Command) string {
	var sb strings.Builder
	for _, c := range commands {
		help := buildHelpForCommand(c)
		sb.WriteString("\n" + help + "\n")
	}
	return sb.String()
}

func buildHelpForCommand(cmd Command) string {
	return cmd.Name + "\n" + cmd.Help
}
