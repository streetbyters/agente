package errors

func New(text string, args ...interface{}) error {
	e := &PluggableError{s: text}

	if len(args) == 1 {
		e.Status = args[0].(int)
	} else if len(args) == 2 {
		e.Detail = args[1].(string)
	}

	return e
}

type PluggableError struct {
	Status int
	Detail string
	s      string
}

func (e *PluggableError) Error() string {
	return e.s
}
