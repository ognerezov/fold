package recorder

import (
	"errors"
	"fmt"
	"fold/migrations"
	"fold/openapi"
)

type Invocation struct {
	Url             string                            `json:"url"`
	Method          string                            `json:"method"`
	Data            any                               `json:"data"`
	SecuritySchemes map[string]openapi.SecuritySchema `json:"securitySchemes,omitempty"`
	ChangePath      string                            `json:"changePath,omitempty"`
	SaveAs          string                            `json:"saveAs,omitempty"`
	Headers         map[string]string                 `json:"headers,omitempty"`
	Sanitize        map[string]Sanitizer              `json:"sanitize,omitempty"`
}

type RecordApiDescription struct {
	Port        int               `json:"port"`
	Invocations []Invocation      `json:"invocations"`
	Credentials map[string]string `json:"credentials,omitempty"`
}

func Record(dataPath string, config *RecordApiDescription) error {
	if config == nil {
		return errors.New("config not found")
	}
	importer, _ := migrations.FsImporter(dataPath)
	err := importer.CreateFolder(fmt.Sprintf("%v", config.Port))
	if err != nil {
		panic(err)
	}
	importer, _ = migrations.FsImporter(fmt.Sprintf("%s/%v", dataPath, config.Port))
	return config.Process(importer)
}
