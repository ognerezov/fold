package configurator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/path"
	"fold/security"
	"fold/util"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type RestConfig struct {
	Protocol            string
	Host                string
	Resources           []RestResourceConfig
	Credentials         security.Credentials
	Invocation          InvocationModel
	temporalCredentials *map[string]string
}

type RestResourceConfig struct {
	Id     string            `json:"id"`
	Query  map[string]string `json:"query"`
	Store  string            `json:"store"`
	Path   string            `json:"path"`
	Method string            `json:"method"`
	Body   any               `json:"body"`
	server *RestConfig
}

type RestIndexConfig struct {
	Id          string               `json:"id"`
	Query       map[string]string    `json:"query"`
	Store       string               `json:"store"`
	Path        string               `json:"path"`
	Method      string               `json:"method"`
	Body        any                  `json:"body"`
	Protocol    string               `json:"protocol"`
	Credentials security.Credentials `json:"credentials"`
	Invocation  InvocationModel      `json:"invocation"`
}

type InvocationModel struct {
	OnApplicationStart bool `json:"onApplicationStart"`
	OnRequest          bool `json:"onRequest"`
	Repeat             int  `json:"repeat"`
	OverrideIfExists   bool `json:"overrideIfExists"`
}

var (
	WWWRegex = regexp.MustCompile("^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$|^[a-zA-Z0-9[a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$")
)

func BuildConfig(dataPath string, host string) *RestConfig {
	rootIndex := dataPath + "/index.json"
	var config RestIndexConfig
	err := util.FromJson(rootIndex, &config)
	if err != nil {
		console.RedPrintln(err.Error())
		panic("Invalid index.json file: " + rootIndex)
		return nil
	}
	protocol := config.Protocol
	if protocol == "" {
		protocol = "https"
	}
	resource := &RestConfig{
		Protocol:    protocol,
		Host:        host,
		Resources:   make([]RestResourceConfig, 0),
		Invocation:  config.Invocation,
		Credentials: config.Credentials,
	}
	if config.Path == "" {
		return resource
	}
	method := config.Method
	if method == "" {
		method = "GET"
	}
	if config.Path != "" {
		resource.Resources = append(resource.Resources, RestResourceConfig{
			Id:     config.Id,
			Query:  config.Query,
			Store:  config.Store,
			Path:   config.Path,
			Method: method,
			Body:   config.Body,
			server: resource,
		})
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
		var c *RestResourceConfig
		c, err = FromJson(filename)
		c.server = resource
		if err == nil && c.Path != "" {
			resource.Resources = append(resource.Resources, *c)
		}
		return nil
	})
	return resource
}

func (c *RestConfig) Url() string {
	return fmt.Sprintf("%s://%s", c.Protocol, c.Host)
}

func (r *RestResourceConfig) Url() string {
	res := util.WithLeadingSlash(r.Path) + r.Q()

	return fmt.Sprintf("%s%s", r.server.Url(), res)
}

func (r *RestResourceConfig) Q() string {
	q := make(map[string]string)
	for k, v := range r.Query {
		q[k] = v
	}

	for k, v := range *r.server.temporalCredentials {
		q[k] = v
	}

	if len(q) == 0 {
		return ""
	}

	return "?" + util.EncodeQuery(&q)
}

func (r *RestResourceConfig) Request() (*http.Request, error) {
	var body io.Reader
	if r.Body == nil {
		b, e := json.MarshalIndent(r.Body, "", "    ")
		if e != nil {
			body = bytes.NewReader(b)
		}
	}
	url := r.Url()
	console.GreenPrintln("Proceed http request to " + url)
	res, err := http.NewRequest(r.Method, url, body)

	return res, err
}

func (r *RestResourceConfig) SaveJson() error {
	req, err := r.Request()
	resp, err := util.SendRequest(req)
	if err != nil {
		return err
	}
	defer util.HideBody(resp.Body)

	filePath := r.Filename()
	out, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer util.CloseFile(out)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *RestResourceConfig) Filename() string {
	return fmt.Sprintf("%s%s", path.Root, util.WithLeadingSlash(r.Store))
}

func (r *RestResourceConfig) SaveRequired() bool {
	if r.server.Invocation.OverrideIfExists {
		return true
	}

	return !util.DoesFileExist(r.Filename())
}

func (c *RestConfig) Start() {
	console.CyanPrintln("Starting rest handler " + c.Host)
	err := c.InitCredentials()
	if err != nil {
		console.RedPrintln(err.Error())
	}
	if c.Invocation.OnApplicationStart {
		console.CyanPrintln("Save data locally for " + c.Host)
		c.AllJsons()
	}
}

func (c *RestConfig) AllJsons() {
	for _, resource := range c.Resources {
		if !resource.SaveRequired() {
			console.YellowPrintln("File already exists " + resource.Filename())
			continue
		}
		err := resource.SaveJson()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
}

func (c *RestConfig) InitCredentials() error {
	if c.temporalCredentials != nil {
		return nil
	}
	return c.ResetCredentials()
}

func (c *RestConfig) ResetCredentials() error {
	console.CyanPrintln("Configure credentials for " + c.Host)
	if c.Credentials.Usage.ExchangeRequired {
		return errors.New("not implemented yet")
	}
	if c.Credentials.SecretInputRequired() {
		cred := map[string]string{}
		console.YellowPrintln("Input secrets required for " + c.Host)
		for k, v := range c.Credentials.Secrets {
			if v != "" {
				console.YellowPrintln("Secret exists " + k)
				cred[k] = v
			} else {
				console.YellowPrintln("Input secret " + k)
				cred[k] = console.ReadStr(fmt.Sprintf("Enter %s: ", k))
			}
		}
		c.Credentials.Secrets = cred
	}
	c.temporalCredentials = &c.Credentials.Secrets
	return nil
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
		config = BuildConfig(fmt.Sprintf("%s/www/%s", dataPath, host), host)
		if config != nil {
			resources = append(resources, config)
		}
	}

	if len(resources) == 0 {
		return nil
	}
	return resources
}

func FromJson(filename string) (*RestResourceConfig, error) {
	var config RestResourceConfig
	err := util.FromJson(filename, &config)
	if err != nil {
		return nil, err
	}
	if config.Method == "" {
		config.Method = "GET"
	} else {
		config.Method = strings.ToUpper(config.Method)
	}

	return &config, nil
}
