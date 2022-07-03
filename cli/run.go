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

		// Build command map.
		commandMap := make(map[string]int)
		for i, c := range commands {
			commandMap[c.Name] = i
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
			cmdIndex, ok := commandMap[name]
			if !ok {
				help := buildHelp(commands)
				msg := fmt.Sprintf("Invalid command: %v\n\nList of valid commands:\n\n%v", name, help)
				io.Print(msg)
				continue
			}

			// Execute command.
			commands[cmdIndex].Func(io, args)
		}
	}()

	return <-errCh
}
