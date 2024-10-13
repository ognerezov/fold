package util

import (
	"net/http"
	"net/url"
)

func MapQuery(r *http.Request) (map[string][]string, error) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.RawQuery)

	if err != nil {
		q = make(url.Values)
	}

	return q, err
}

func OneValueQuery(col string, value string) map[string][]string {
	res := make(map[string][]string)
	values := make([]string, 1)
	values[0] = value
	res[col] = values
	return res
}
