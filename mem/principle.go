package mem

type Principle struct {
	roles []string
	id    string
}

var Guest = Principle{
	roles: []string{"guest"},
}
