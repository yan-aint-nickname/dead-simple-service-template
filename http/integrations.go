package main

import (
	"github.com/dghubble/sling"
)


type ProjectsClientHttp struct {
	Client *sling.Sling
}

type HttpClient[T any] interface {
	MakeRequest(method, endpoint string, requestData any) (T, error)
}

func NewProjectsClient(s SettingsHttp) *ProjectsClientHttp {
	c := sling.New()
	c.Base(s.ProjectsBaseUrl)

	return &ProjectsClientHttp{
		Client: c,
	}
}
