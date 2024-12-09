package controls

import "errors"

type EchoControl string

func (id EchoControl) Do(data map[string]any) (any, error) {
	if data == nil {
		return nil, errors.New("request is empty")
	}

	return data, nil
}

func GetEcho(id string) *Control {
	var ctr Control
	ctr = EchoControl(id)
	return &ctr
}
