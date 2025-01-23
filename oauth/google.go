package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/controls"
	"fold/db"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	"fold/util"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"net/url"
	"strings"
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
	RegistrationAllowed     bool
	Drive                   *drive.Service
	Sheets                  *sheets.Service

	Type              string `json:"type"`
	PrivateKeyId      string `json:"private_key_id"`
	PrivateKey        string `json:"private_key"`
	ClientEmail       string `json:"client_email"`
	ClientX509CertUrl string `json:"client_x509_cert_url"`
	UniverseDomain    string `json:"universe_domain"`
}

type CodeTokenRequest struct {
	Code        string `json:"code" validate:"required"`
	RedirectUri string `json:"redirect_uri"`
}

type CodeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
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

func (gj *GoogleJson) AnyClient() *GoogleSecret {
	if gj.Web != nil {
		return gj.Web
	}
	if gj.Installed != nil {
		return gj.Installed
	}
	if gj.IOs != nil {
		return gj.IOs
	}
	if gj.Android != nil {
		return gj.Android
	}
	return nil
}

func (gj *GoogleJson) ServerSecret() string {
	if gj.ClientSecret != "" {
		return gj.ClientSecret
	}
	client := gj.AnyClient()
	if client == nil {
		return ""
	}
	return client.ClientSecret
}

func (gj *GoogleJson) TokenUrl() string {
	if gj.TokenUri != "" {
		return gj.TokenUri
	}
	client := gj.AnyClient()
	if client == nil {
		return ""
	}
	return client.TokenUri
}

func (gj *GoogleJson) Id() string {
	if gj.ClientId != "" {
		return gj.ClientId
	}
	client := gj.AnyClient()
	if client == nil {
		return ""
	}
	return client.ClientId
}

func (gj *GoogleJson) Do(data map[string]any, w http.ResponseWriter, _ *http.Request) {
	codeErrMsg := "code is missing in the request"
	if data == nil {
		router.BadRequest(errors.New(codeErrMsg), w)
		return
	}
	var tokenReq CodeTokenRequest
	err := util.Restructure(data, &tokenReq)
	if err != nil {
		router.BadRequest(errors.New(fmt.Sprintf("error decoding request: %v", err.Error())), w)
		return
	}
	req, err := gj.TokenRequest(tokenReq.Code, tokenReq.RedirectUri)
	if err != nil {
		router.BadRequest(err, w)
		return
	}
	fmt.Println(tokenReq.Code)
	resp, err := util.SendRequest(req)
	if err != nil {
		router.ServerError(err, w)
		return
	}
	defer util.HideBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		var errResp map[string]any
		err = decoder.Decode(&errResp)
		err = errors.New(resp.Status)
		router.ServerError(err, w)
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var tokenResp CodeTokenResponse
	err = decoder.Decode(&tokenResp)

	iss := "fold"
	if err != nil {
		router.ServerError(err, w)
		return
	}
	provider := "google"
	userInfo, err := tokenResp.ExchangeForToken("https://www.googleapis.com/oauth2/v3/userinfo", iss)
	if err != nil {
		router.Unauthorized(err, w)
		return
	}
	RemovePreviousRecords(userInfo.Email, provider)
	StoreOauth(tokenResp.GetRow(userInfo.Token, userInfo.Email, provider))
	if gj.RegistrationAllowed {
		table, ok := mem.TheStore.GetTable(util.UserPath)
		if !ok {
			console.YellowPrintln("User table not found")
		} else {
			if !table.Has(userInfo.Email) {
				console.YellowPrintln("Registering user " + userInfo.Email)
				_, err = table.Insert(userInfo.Map())
				if err != nil {
					console.RedPrintln(err.Error())
				}
			}
		}

	}
	router.WriteResponse(api.LoginResponse{Token: userInfo.Token, Iss: iss}, w)
}

func (gj *GoogleJson) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return []openapi.Parameter{
			{
				Name: "redirect_uri",
				Schema: openapi.Schema{
					Type: "string",
				},
			},
			{
				Name: "code",
				Schema: openapi.Schema{
					Type: "string",
				},
				Required: true,
			},
		}, openapi.JsonResponse(openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"token": {
					Type: "string",
				},
				"iss": {
					Type: "string",
				},
			},
		})
}

func (gj *GoogleJson) TokenRequest(code string, redirectUri string) (*http.Request, error) {
	uri := gj.TokenUrl()
	if uri == "" {
		return nil, errors.New("token uri is empty")
	}
	console.GreenPrintln("Proceed http request to " + uri)
	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)
	data.Set("client_id", gj.Id())
	data.Set("client_secret", gj.ServerSecret())
	data.Set("grant_type", "authorization_code")
	res, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	res.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return res, err
}

func (gj *GoogleJson) AuthControl(_ string, _ any) *controls.Control {
	var googleAuthControl controls.Control
	googleAuthControl = gj
	return &googleAuthControl
}
