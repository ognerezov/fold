package controls

type Control interface {
	Do(map[string]any) (any, error)
}
