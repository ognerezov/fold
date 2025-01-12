package interfaces

type Action interface {
	Act()
}

type Phase struct {
	Actions []Action
}

func NewPhase(actions ...Action) *Phase {
	phase := Phase{}

	phase.Actions = actions
	if phase.Actions == nil {
		phase.Actions = []Action{}
	}
	return &phase
}

func (p *Phase) Append(actions ...Action) {
	if actions == nil || len(actions) == 0 {
		return
	}

	p.Actions = append(p.Actions, actions...)
}

func (p *Phase) Act() {
	for _, a := range p.Actions {
		a.Act()
	}
}
