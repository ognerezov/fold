package configurator

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"os"
	"regexp"
)

type RestConfig struct {
	Protocol  string
	Host      string
	Resources []RestResourceConfig
}

type RestResourceConfig struct {
	Id    string `json:"id"`
	Query string `json:"query"`
	Store string `json:"store"`
	Path  string `json:"path"`
}

var (
	WWWRegex = regexp.MustCompile("^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$|^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$")
)

func ConfigureResources(dataPath string) []*RestConfig {
	files, err := ReadDir(dataPath + "/www")
	if err != nil {
		return nil
	}
	resources := make([]*RestConfig, 0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		host := file.Name()
		if !WWWRegex.MatchString(host) {
			continue
		}
		var config *RestConfig
		config, err = ConfigureResource(fmt.Sprintf("%s/www/%s", dataPath, host), host)
		resources = append(resources, config)
	}

	if len(resources) == 0 {
		return nil
	}
	return resources
}

func ConfigureResource(dataPath string, host string) (*RestConfig, error) {
	f, err := os.OpenFile(dataPath+"/index.json", os.O_RDONLY, 0)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(f)

	decoder := json.NewDecoder(f)
	var config RestResourceConfig
	err = decoder.Decode(&config)
	if err != nil {
		console.RedPrintln("Invalid REST json config " + host)
		panic(err)
	}
	resource := &RestConfig{
		Protocol:  "https",
		Host:      host,
		Resources: make([]RestResourceConfig, 0),
	}
	if config.Path != "" {
		resource.Resources = append(resource.Resources, config)
	}
	fmt.Println(resource)
	return resource, nil
}
