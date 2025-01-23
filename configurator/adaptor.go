package configurator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fold/controls"
	"fold/openapi"
	"fold/router"
	"fold/util"
	"net/http"
	"net/url"
	"slices"
)

type Adaptor struct {
	Port            int               `json:"port"`
	Path            string            `json:"path"`
	Method          string            `json:"method"`
	ParamsFromQuery []string          `json:"params_from_query"`
	QueryFromParams map[string]string `json:"query_from_params"`
	QueryFromQuery  map[string]string `json:"query_from_query"`
}

func (adaptor Adaptor) Do(data map[string]any, w http.ResponseWriter, r *http.Request) {
	service := App.services[adaptor.Port]

	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		router.ServerError(err, w)
		return
	}
	reader := bytes.NewReader(b)

	var method string
	if adaptor.Method != "" {
		method = adaptor.Method
	} else {
		method = r.Method
	}

	var path string
	if adaptor.Path != "" {
		path = adaptor.Path
	} else {
		path = r.URL.Path
	}

	rawQuery := adaptor.MapQuery(r)
	fmt.Println(path + rawQuery)

	request, err := http.NewRequest(method, path+rawQuery, reader)

	service.server.Handler.ServeHTTP(w, request)
}

func (adaptor Adaptor) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return []openapi.Parameter{}, openapi.JsonResponse(openapi.AnObject)
}

func (adaptor Adaptor) MapQuery(r *http.Request) string {
	q, _ := util.MapQuery(r)
	newQuery := make(url.Values)
	pathParams := make([]string, 0)
	paramNames := adaptor.ParamsFromQuery
	if paramNames == nil {
		paramNames = make([]string, 0)
	}
	for key, value := range q {
		if slices.Contains(paramNames, key) && len(value) > 0 {
			pathParams = append(pathParams, value[0])
		} else {
			newQuery[adaptor.MapQueryParam(key)] = value
		}
	}
	for key, value := range adaptor.QueryFromParams {
		var param string
		util.PathParamValue(r, key, &param)
		if param != "" {
			newQuery[value] = []string{param}
		}
	}
	if len(newQuery) == 0 {
		return ""
	}
	fmt.Println(util.ParamLiterals(pathParams, "/") + "?" + newQuery.Encode())
	return util.ParamLiterals(pathParams, "/") + "?" + newQuery.Encode()
}

func (adaptor Adaptor) MapQueryParam(key string) string {
	if adaptor.QueryFromQuery == nil || len(adaptor.QueryFromQuery) == 0 {
		return key
	}
	res, ok := adaptor.QueryFromQuery[key]
	if ok {
		return res
	}
	return key
}

func ConfigureAdaptor(_ string, config any) *controls.Control {
	var a Adaptor
	err := util.Restructure(config, &a)
	if err != nil {
		panic(err)
	}
	var res controls.Control
	res = a
	return &res
}
