package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func newProjectsRouteTestOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *RedisContainer { return nil },
			func() *PostgresContainer { return nil },
			NewTestSettingsHttp,
			NewServeMux,
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
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/projects/", nil)
	t.Logf("Req %v", req)
	if err != nil {
		t.Errorf("Error building request ", err)
	}
	srv.ServeHTTP(w, req)

	t.Logf("Recorder %v", w)

	if w.Code != http.StatusOK {
		t.Errorf("Status is not OK")
	}
}

func TestAPI(t *testing.T) {
	// TODO: mock http call to projects_base_url

	t.Run("Get Projects", func(t *testing.T) {
		NewTestApp(t, newProjectsRouteTestOptions(), fx.Invoke(testGetProjects)).Stop()
	})
}
