// I don't know how to make it more transparent and generic at the same time
// QUESTION: Need interface?
//
//	type APIRequestMaker interface {
//		Name() string
//		MakeRequest(method, endpoint string, headers, query_params map[string]string, json map[string]any) (any, error)
//	}
package main

import (
	"time"

	"github.com/imroc/req/v3"
)

type ProjectsAPI struct {
	Client *req.Client
}

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

func NewClient() *req.Client {
	return req.C()
}

func NewProjectsAPI(s SettingsHttp, c *req.Client) *ProjectsAPI {
	baseUrl := s.ProjectsBaseUrl
	apiTimeout := time.Duration(s.ProjectsApiTimeout * 1000 * 1000 * 1000)
	c.SetBaseURL(baseUrl).SetTimeout(apiTimeout)
	return &ProjectsAPI{Client: c}
}

func (api *ProjectsAPI) GetProjects() (resp ProjectsResponse, err error) {
	err = api.Client.Get("/projects").Do().UnmarshalJson(&resp)
	return
}
