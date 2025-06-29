package recorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fold/configurator"
	"fold/console"
	"fold/mem"
	"fold/migrations"
	"fold/util"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strings"
)

func (rad RecordApiDescription) Process(importer migrations.Importer) error {

	for _, invocation := range rad.Invocations {
		for name := range invocation.SecuritySchemes {
			_, ok := rad.Credentials[name]
			if ok {
				console.YellowPrintln(fmt.Sprintf("Set secret for security scheem %s from json", name))
			} else {
				console.YellowPrintln(fmt.Sprintf("Input secret for security scheem %s:", name))
				rad.Credentials[name] = strings.TrimSpace(console.ReadStr(fmt.Sprintf("Type secret %s: ", name)))
				fmt.Println(rad.Credentials)
			}
		}
	}

	for _, invocation := range rad.Invocations {
		var body io.Reader
		if invocation.Data != nil {
			b, e := json.MarshalIndent(invocation.Data, "", "    ")
			if e != nil {
				panic(e)
			}
			body = bytes.NewReader(b)
		}
		url := invocation.Url
		hasQuery := strings.Contains(url, "?")
		for name, securityScheme := range invocation.SecuritySchemes {
			securityQuery := securityScheme.SecurityQuery(rad.Credentials[name])
			if securityQuery != nil {
				if hasQuery {
					url += "&" + *securityQuery
				} else {
					url += "?" + *securityQuery
				}
				break
			}
		}
		method := strings.ToUpper(invocation.Method)
		console.GreenPrintln(fmt.Sprintf("Proceed http request to %s: %s", method, url))
		request, err := http.NewRequest(method, url, body)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Content-Type", "application/json")
		if invocation.Headers != nil {
			for key, value := range invocation.Headers {
				request.Header.Set(key, value)
			}
		}

		for name, securityScheme := range invocation.SecuritySchemes {
			securityHeaderName, securityHeaderValue := securityScheme.SecurityHeader(rad.Credentials[name])
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

		if invocation.Sanitize != nil {
			for k, sanitizer := range invocation.Sanitize {
				if sanitizer.Method != Randomize && sanitizer.Method != "" && sanitizer.Method != Erase {
					panic("Unknown sanitize method: " + sanitizer.Method)
				}
				valueFunction := func(_ any, _ *[]string) any {
					return nil
				}
				switch sanitizer.Method {
				case Randomize:
					combine := sanitizer.Combine
					if combine == 0 {
						combine = 1
					}
					values := util.RestructureArrays(sanitizer.Values, combine)
					if values != nil {
						valueFunction = func(v any, path *[]string) any {

							res := ""
							// Should add more aggregate functions here
							for _, value := range values {
								separator := " "
								if res == "" {
									separator = ""
								}
								res = fmt.Sprintf("%v%s%v", res, separator, value[rand.Intn(len(value))])
							}
							if sanitizer.Parents != nil {
								for _, parent := range sanitizer.Parents {
									// Not modifying value if a path doesn't match
									if !slices.Contains(*path, parent) {
										return v
									}
								}
							}
							return res
						}
					}

				}
				noSql.Replace(k, valueFunction)
			}
			b, err = json.MarshalIndent(noSql.Val(), "", "  ")
			if err != nil {
				panic(err)
			}
		}

		if hasQuery {
			noSql = noSql.ToArray()
		}

		path := strings.Join(parts[:len(parts)-1], "/")
		filename := parts[len(parts)-1] + ".json"
		if method != "GET" {
			path = fmt.Sprintf("%s/%s", configurator.RawRoutesFolder, path)
			filename = fmt.Sprintf("%s%s%s", invocation.Method, configurator.RawSeparator, filename)
		}

		data := migrations.FileData{
			Binary:   b,
			Path:     path,
			Filename: filename,
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
