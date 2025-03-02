package gcloud

import (
	"context"
	"errors"
	"fmt"
	"fold/console"
	"fold/controls"
	"fold/oauth"
	"fold/openapi"
	"fold/router"
	"fold/security"
	"fold/util"
	"golang.org/x/oauth2"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/option"
	"net/http"
)

type CreateProjectRequest struct {
	ProjectId string `json:"projectId"`
	Name      string `json:"name"`
}

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

	var req CreateProjectRequest
	err := util.Restructure(data, &req)
	if err != nil {
		router.ServerError(err, w)
		return
	}
	tokenSource := oauth.SourceToken(principle.Token())
	token, err := tokenSource.Token()
	if err != nil {
		router.ServerError(err, w)
		return
	}
	fmt.Println(token)
	manager, err := cloudresourcemanager.NewService(context.Background(),
		option.WithTokenSource(oauth2.StaticTokenSource(token)))

	if err != nil {

		router.ServerError(err, w)
		return
	}

	op, err := manager.Projects.Create(&cloudresourcemanager.Project{
		ProjectId:   req.ProjectId,
		DisplayName: req.Name,
		Tags: map[string]string{
			"managed": "fold",
		},
	}).Do()

	if err != nil {
		var m map[string]any
		e := util.Restructure(err, &m)
		if e == nil {
			fmt.Println(m)
			router.WriteServerResponse(m, int(m["code"].(float64)), w)
			return
		}
		router.ServerError(err, w)
		return
	}

	router.WriteResponse(*op, w)
}

func (pc ProjectCreator) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return []openapi.Parameter{
			{
				Name: "projectId",
				Schema: openapi.Schema{
					Type:        "string",
					Description: "Unique Project identifier. It must be 6 to 30 lowercase ASCII letters, digits, or hyphens",
				},
			},
			{
				Name: "name",
				Schema: openapi.Schema{
					Type:        "string",
					Description: "Optional display name, 4 to 30 characters: lowercase and uppercase letters, numbers, hyphen, single-quotes",
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

func GetProjectCreator(id string, _ any) *controls.Control {
	var ctr controls.Control
	ctr = ProjectCreator(id)
	return &ctr
}
