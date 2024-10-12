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
