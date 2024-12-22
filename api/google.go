package api

import (
	"encoding/json"
	"fold/console"
	"fold/db"
	"fold/util"
)

type GoogleSecret struct {
	ClientId                string   `json:"client_id"`
	ProjectId               string   `json:"project_id"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
	JavascriptOrigins       []string `json:"javascript_origins"`
}

type GoogleJson struct {
	Web                     *GoogleSecret `json:"web"`
	Installed               *GoogleSecret `json:"installed"`
	IOs                     *GoogleSecret `json:"ios"`
	Android                 *GoogleSecret `json:"android"`
	ClientId                string        `json:"client_id"`
	ProjectId               string        `json:"project_id"`
	AuthUri                 string        `json:"auth_uri"`
	TokenUri                string        `json:"token_uri"`
	AuthProviderX509CertUrl string        `json:"auth_provider_x509_cert_url"`
	ClientSecret            string        `json:"client_secret"`
	RedirectUris            []string      `json:"redirect_uris"`
	JavascriptOrigins       []string      `json:"javascript_origins"`
}

func (gs *GoogleSecret) WithoutSecret() *GoogleSecret {
	res := *gs
	res.ClientSecret = ""
	return &res
}

func (gj *GoogleJson) WithoutSecret() GoogleJson {
	res := *gj
	res.ClientSecret = ""
	if gj.Web != nil {
		res.Web = gj.Web.WithoutSecret()
	}
	if gj.Installed != nil {
		res.Installed = gj.Installed.WithoutSecret()
	}
	if gj.IOs != nil {
		res.IOs = gj.IOs.WithoutSecret()
	}
	if gj.Android != nil {
		res.Android = gj.Android.WithoutSecret()
	}
	return res
}

func (gj *GoogleJson) Attach(other *GoogleJson) *GoogleJson {
	if other == nil {
		return gj
	}
	if other.Web != nil {
		gj.Web = other.Web
	}
	if other.Installed != nil {
		gj.Installed = other.Installed
	}
	if other.IOs != nil {
		gj.IOs = other.IOs
	}
	if other.Android != nil {
		gj.Android = other.Android
	}

	if other.ClientId != "" {
		gj.ClientId = other.ClientId
	}
	if other.ProjectId != "" {
		gj.ProjectId = other.ProjectId
	}
	if other.AuthUri != "" {
		gj.AuthUri = other.AuthUri
	}
	if other.TokenUri != "" {
		gj.TokenUri = other.TokenUri
	}
	if other.AuthProviderX509CertUrl != "" {
		gj.AuthProviderX509CertUrl = other.AuthProviderX509CertUrl
	}
	if other.RedirectUris != nil && len(other.RedirectUris) > 0 {
		gj.RedirectUris = other.RedirectUris
	}

	if other.JavascriptOrigins != nil && len(other.JavascriptOrigins) > 0 {
		gj.JavascriptOrigins = other.JavascriptOrigins
	}

	return gj
}

func (gj *GoogleJson) Export(path string) error {
	data := gj.WithoutSecret()
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		console.RedPrintln(err.Error())
	} else {
		str := "export default\n" + string(bytes)
		err = util.Save(path+".js", []byte(str))
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
	return db.SaveJson(path+".json", gj.WithoutSecret())
}
