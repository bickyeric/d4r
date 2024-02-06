package view

type ActionHandler func() error

type Action struct {
	name    string
	handler ActionHandler
}

func NewAction(name string, handler ActionHandler) Action {
	return Action{name, handler}
}
