package configurator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/path"
	"fold/util"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type RestConfig struct {
	Protocol  string
	Host      string
	Resources []RestResourceConfig
}

type RestResourceConfig struct {
	Id     string `json:"id"`
	Query  string `json:"query"`
	Store  string `json:"store"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   any    `json:"body"`
	server *RestConfig
}

var (
	WWWRegex   = regexp.MustCompile("^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$|^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$")
	httpClient = &http.Client{}
)

func (c *RestConfig) Url() string {
	return fmt.Sprintf("%s://%s", c.Protocol, c.Host)
}

func (r *RestResourceConfig) Url() string {
	res := util.WithLeadingSlash(r.Path)

	if r.Query != "" {
		res += "?" + r.Query
	}
	return fmt.Sprintf("%s%s", r.server.Url(), res)
}

func (r *RestResourceConfig) Request() (*http.Request, error) {
	var body io.Reader
	if r.Body == nil {
		b, e := json.MarshalIndent(r.Body, "", "    ")
		if e != nil {
			body = bytes.NewReader(b)
		}
	}
	res, err := http.NewRequest(r.Method, r.Url(), body)

	return res, err
}

func (r *RestResourceConfig) SaveJson() error {
	req, err := r.Request()
	resp, err := httpClient.Do(req)
	fmt.Println(resp)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(resp.Body)

	filePath := fmt.Sprintf("%s%s", path.Root, util.WithLeadingSlash(r.Store))
	out, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		console.RedPrintln(err.Error())
		return err
	}
	defer util.CloseFie(out)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *RestConfig) AllJsons() {
	for _, resource := range c.Resources {
		err := resource.SaveJson()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
}

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

func FromJson(filename string) (*RestResourceConfig, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
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
	if config.Method == "" {
		config.Method = "GET"
	} else {
		config.Method = strings.ToUpper(config.Method)
	}

	return &config, nil
}

func ConfigureResource(dataPath string, host string) (*RestConfig, error) {
	rootIndex := dataPath + "/index.json"
	c, err := FromJson(rootIndex)
	if err != nil {
		console.RedPrintln("Invalid REST json config " + host)
		panic(err)
		return nil, err
	}
	resource := &RestConfig{
		Protocol:  "https",
		Host:      host,
		Resources: make([]RestResourceConfig, 0),
	}
	config := *c
	config.server = resource
	if config.Path != "" {
		resource.Resources = append(resource.Resources, config)
	}

	clean := path.CreateRootCleaner(dataPath)
	err = path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, filename, extension := path.Structure(dataPath, p, info, clean)
		if extension != ".json" || filename == rootIndex {
			return nil
		}
		c, err = FromJson(filename)
		c.server = resource
		if err == nil && c.Path != "" {
			resource.Resources = append(resource.Resources, *c)
		}
		return nil
	})
	return resource, nil
}
