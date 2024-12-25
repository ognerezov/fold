package controls

import "errors"

type GoogleAuthControl string

// TODO exchange google code for token. Use some json data for init or may be json itself
func (ga GoogleAuthControl) Do(data map[string]any) (any, error) {
	if data == nil {
		return nil, errors.New("request is empty")
	}

	return data, nil
}
