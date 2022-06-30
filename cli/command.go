package cli

type Command struct {
	Name string
	Func func(IO, []string)
	Help string
}
