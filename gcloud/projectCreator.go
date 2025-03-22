package gcloud

import (
	"context"
	"errors"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/controls"
	"fold/oauth"
	"fold/openapi"
	"fold/router"
	"fold/security"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
	"net/http"
)

type ProjectCreator string

func (pc ProjectCreator) Do(data map[string]any, w http.ResponseWriter, r *http.Request) {
	principle := security.FromRequest(r)
	if principle == nil || principle.Token() == "" {
		router.Unauthorized(errors.New("token not found"), w)
		return
	}
	if data == nil {
		console.RedPrintln("request is empty")
		router.BadRequest(errors.New("request is empty"), w)
		return
	}
	fmt.Println(data)
	service_account := data["credentials"].(string)
	tokenSource := oauth.SourceToken(principle.Token())

	scopes := []string{"https://www.googleapis.com/auth/cloud-platform"}

	// Create a context.
	ctx := context.Background()

	// Create an IAM Credentials client.

	ts, err := tokenSource.Build(ctx)
	if err != nil {
		router.ServerError(err, w)
		return
	}

	client, err := iamcredentials.NewService(ctx, option.WithTokenSource(*ts))
	if err != nil {
		router.ServerError(err, w)
		return
	}
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", service_account)
	request := &iamcredentials.GenerateAccessTokenRequest{
		Scope:    scopes,
		Lifetime: "3600s", // Token lifetime (max 1 hour).
	}
	/*
	   This part doesn't work any way. It seems that user authenticated in main project can't perform action
	   in his own project. User-friendly setup implementation is postponed
	*/
	response, err := client.Projects.ServiceAccounts.GenerateAccessToken(resource, request).Do()
	if err != nil {
		router.ServerError(err, w)
		return
	}

	fmt.Println("Access Token:", response.AccessToken)
	x := api.Ok()
	op := &x
	router.WriteResponse(*op, w)
}

func (pc ProjectCreator) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return []openapi.Parameter{
			{
				Name: "credentials",
				Schema: openapi.Schema{
					Type:        "string",
					Description: "Service account email address.",
				},
			},
		}, map[string]openapi.Response{
			"200": {
				Description: "Google cloud response",
				Content: map[string]openapi.Content{
					openapi.ApplicationJson: {
						Schema: openapi.AnObject,
					},
				},
			},
		}
}

func ConfigureProjectCreator(dataPath string) controls.ControlFactory {
	return func(s string, a any) *controls.Control {
		var ctr controls.Control
		ctr = ProjectCreator(dataPath)
		return &ctr
	}
}
