package cli

import "bufio"

const Prefix = "> "

type IO struct {
	in      chan string
	out     chan string
	context map[string]interface{}
}

func NewIO() *IO {
	return &IO{
		in:      make(chan string),
		out:     make(chan string),
		context: make(map[string]interface{}),
	}
}

func (io IO) Run(reader *bufio.Reader, writer *bufio.Writer) error {
	errCh := make(chan error)
	defer close(errCh)

	go func() {
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				errCh <- err
				break
			}
			io.in <- input[:len(input)-1]
		}
	}()

	go func() {
		for text := range io.out {
			_, err := writer.WriteString(text)
			writer.Flush()
			if err != nil {
				errCh <- err
				break
			}
		}
	}()

	return <-errCh
}

func (io IO) Print(msg string) {
	io.out <- "\r" + msg + "\n"
}

func (io IO) PrintWithPrefix(msg string) {
	io.out <- "\r" + msg + "\n" + Prefix
}

func (io IO) PrintPrefix() {
	io.out <- "\r" + Prefix
}

func (io IO) SetContextValue(key string, value interface{}) {
	io.context[key] = value
}

func (io IO) ContextValue(key string) (interface{}, bool) {
	val, ok := io.context[key]
	return val, ok
}
