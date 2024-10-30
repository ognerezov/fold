package mem

import "fold/console"

type Batch interface {
	UpdateRow(data *[]Data) error
}

type StringMapper func(s string) (string, error)

type MapColumnBatch struct {
	InputColumn  int
	ResultColumn int
	Transform    StringMapper
}

func (p MapColumnBatch) UpdateRow(row *[]Data) error {
	r := *row

	hash, err := p.Transform(r[p.InputColumn].Str())

	if err != nil {
		console.RedPrintln(err.Error())
	} else {
		r[p.ResultColumn] = *FromString(hash)
		r[p.InputColumn] = *FromString("")
	}

	return err
}
