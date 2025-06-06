package recorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/migrations"
	"fold/util"
	"io"
	"net/http"
	"strings"
)

func (rad RecordApiDescription) Process(importer migrations.Importer) error {

	for _, invocation := range rad.Invocations {
		for name, scheme := range invocation.SecuritySchemes {
			secret, ok := rad.Credentials[name]
			if ok {
				console.YellowPrintln(fmt.Sprintf("Set secret for security scheem %s from json", name))
				scheme.Secret = secret
			} else {
				console.YellowPrintln(fmt.Sprintf("Input secret for security scheem %s:", name))
				scheme.Secret = console.ReadStr("Type secret %s: ")
			}
		}
	}

	for _, invocation := range rad.Invocations {
		var body io.Reader
		if invocation.Data != nil {
			b, e := json.MarshalIndent(invocation.Data, "", "    ")
			if e != nil {
				body = bytes.NewReader(b)
			}
		}
		url := invocation.Url
		hasQuery := strings.Contains(url, "?")
		for _, securityScheme := range invocation.SecuritySchemes {
			securityQuery := securityScheme.SecurityQuery()
			if securityQuery != nil {
				if hasQuery {
					url += "&" + *securityQuery
				} else {
					url += "?" + *securityQuery
				}
				break
			}
		}
		console.GreenPrintln("Proceed http request to " + url)
		request, err := http.NewRequest(invocation.Method, url, body)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Content-Type", "application/json")

		for _, securityScheme := range invocation.SecuritySchemes {
			securityHeaderName, securityHeaderValue := securityScheme.SecurityHeader()
			if securityHeaderName != nil {
				request.Header.Set(*securityHeaderName, *securityHeaderValue)
				break
			}
		}

		resp, err := util.SendRequest(request)
		if err != nil {
			return err
		}

		b, err := util.HandleBody(resp.Body)

		//should handle different types

		parts := strings.Split(request.URL.Path, "/")

		noSql, err := mem.FromBytes(b)

		if err != nil {
			panic(err)
		}

		if hasQuery {
			noSql = noSql.ToArray()
		}

		data := migrations.FileData{
			Binary:   b,
			Path:     strings.Join(parts[:len(parts)-1], "/"),
			Filename: parts[len(parts)-1] + ".json",
			MimeType: resp.Header.Get("Content-Type"),
			NoSql:    noSql,
		}

		err = importer.CreateFolder(data.Path)
		if err != nil {
			panic(err)
		}
		err = importer.SaveFile(data)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
