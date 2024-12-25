package util

import (
	"fold/console"
	"goji.io/pat"
	"io"
	"net/http"
	"net/url"
)

var (
	httpClient = &http.Client{}
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

func EncodeQuery(m *map[string]string) string {
	if m == nil || len(*m) == 0 {
		return ""
	}

	values := url.Values{}

	for k, v := range *m {
		values.Add(k, v)
	}
	return values.Encode()
}

func HideBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		console.RedPrintln(err.Error())
	}
}

func SendRequest(req *http.Request) (*http.Response, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}

	return resp, nil
}

func PathParamValue(req *http.Request, name string, out *string) {
	defer func() {
		if recover() != nil {
			*out = ""
		}
	}()
	*out = pat.Param(req, name)
}
