package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imroc/req/v3"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var SuccessResponseProjects string = `{"projects":[{"id":1,"name":"Project 1"},{"id":2,"name":"Project 2"},{"id":3,"name":"Project 3"}]}`

// NOTE: Do not format!
var SuccessResponseProjectsAPI string = `{"data":[{"id":1,"name":"Project 1"},{"id":2,"name":"Project 2"},{"id":3,"name":"Project 3"}],"error":null}`

func NewTestClient() *req.Client {
	client := req.C()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder(http.MethodGet, "http://projects_base_url/projects",
		httpmock.NewStringResponder(http.StatusOK, SuccessResponseProjects))
	return client
}

func newProjectsRouteTestOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *RedisContainer { return nil },
			func() *PostgresContainer { return nil },
			NewTestSettingsHttp,
			NewServeMux,
			NewTestClient,
			// NewHttpServer,
			NewProjectsAPI,
			NewProjectsServiceGet,
			fx.Annotate(NewApiV1Router, fx.ResultTags(`name:"ApiV1Router"`)),
			AsRoute(NewProjectsHandlerGet, `group:"projectsRoutes"`),
		),
		fx.Invoke(
			fx.Annotate(RegisterProjectsApi, fx.ParamTags(`name:"ApiV1Router"`, `group:"projectsRoutes"`)),
		),
	)
}

func testGetProjects(srv *gin.Engine, t fxtest.TB) {
	defer httpmock.DeactivateAndReset()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/projects/", nil)
	t.Logf("Req %v", req)
	if err != nil {
		t.Errorf("Error building request ", err)
	}
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status is not OK")
	}
	if w.Body.String() != SuccessResponseProjectsAPI {
		t.Errorf("Wrong body! Get another!")
	}
}

func TestAPI(t *testing.T) {
	t.Run("Get Projects", func(t *testing.T) {
		NewTestApp(t, newProjectsRouteTestOptions(), fx.Invoke(testGetProjects)).Stop()
	})
}
