package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run(init func(io IO) error, commands []Command) error {
	errCh := make(chan error)
	defer close(errCh)

	// Setup IO.
	io := *NewIO()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	go func() {
		errCh <- io.Run(reader, writer)
	}()

	go func() {
		// Initialize.
		err := init(io)
		if err != nil {
			errCh <- fmt.Errorf("intialize: %w", err)
			return
		}

		// Add help command.
		commandMap := make(map[string]Command)
		helpName := "help"
		helpCommand := Command{
			Name: helpName,
			Func: func(io IO, args []string) {
				if len(args) == 1 {
					name := args[0]
					cmd, ok := commandMap[name]
					if ok {
						help := buildHelpForCommand(cmd)
						io.Print(help)
					} else {
						io.Print(fmt.Sprintf("Unknown command '%v'.\nEnter '%v' to see a list of valid commands.", name, helpName))
					}
				} else {
					help := buildHelp(commands)
					io.Print(help)
				}
			},
			Help: "Show list of available commands.",
		}
		commands = append(commands, helpCommand)

		// Build command map.
		for _, c := range commands {
			commandMap[c.Name] = c
		}

		// Run command loop.
		io.out <- Prefix
		for input := range io.in {
			if len(input) == 0 {
				continue
			}

			// Parse arguments.
			tokens := strings.Split(input, " ")
			name := tokens[0]
			args := tokens[1:]

			// Get command.
			cmd, ok := commandMap[name]
			if !ok {
				msg := fmt.Sprintf("Invalid command: %v\nEnter '%v' to show a list of valid commands.", name, helpCommand.Name)
				io.Print(msg)
				continue
			}

			// Execute command.
			cmd.Func(io, args)
		}
	}()

	return <-errCh
}
