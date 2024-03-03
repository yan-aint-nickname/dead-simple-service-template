package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func RegisterSentryMiddleware(lc fx.Lifecycle, app *gin.Engine, settings SettingsHttp) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := sentry.Init(sentry.ClientOptions{
				Dsn:              settings.SentryDsn,
				EnableTracing:    true,
				Environment:      settings.Env,
				AttachStacktrace: true,
			})
			if err != nil {
				return err
			}
			return nil
		},
	})
	app.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
}

func RegisterGinZapLogger(app *gin.Engine, logger *zap.Logger) {
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	app.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	app.Use(ginzap.RecoveryWithZap(logger, true))
}

func RegisterGinRecovery(app *gin.Engine) {
	app.Use(gin.Recovery())
}

func AsRoute(f any, group string) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(group),
	)
}

func NewServeMux() *gin.Engine {
	return gin.New()
}

func NewHttpServer(lc fx.Lifecycle, log *zap.Logger, settings SettingsHttp, router *gin.Engine) *http.Server {
	// NOTE: I need to specify timeouts because of gosec G112
	srv := &http.Server{
		Addr:              settings.Addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))

			// NOTE: better approach?
			// I don't know better way to _not_ handle this error
			// because I need to recover from panics not shutdown
			go srv.Serve(ln) // nolint: errcheck
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

// NOTE: Invoke for funcs what have no return, they *provide nothing*
// Returs app state
// NOTE: Tags works for specifing producers and consumers of values
func CreateDefaultApp() *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewSettingsHttp,
			NewServeMux,
			NewHttpServer,
			NewRedisClient,
			NewPostgresPool,
			NewTodosServiceGet,
			NewTodosServicePost,
			// TODO: rewrite with routers module?
			fx.Annotate(NewApiV1Router, fx.ResultTags(`name:"ApiV1Router"`)),
			AsRoute(NewTodosHandlerGet, `group:"todoRoutes"`),
			AsRoute(NewTodosHandlerPost, `group:"todoRoutes"`),
			zap.NewExample,
		),
		fx.Invoke(RegisterSentryMiddleware),
		fx.Invoke(RegisterGinRecovery),
		fx.Invoke(RegisterGinZapLogger),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(fx.Annotate(RegisterTodosApi, fx.ParamTags(`name:"ApiV1Router"`, `group:"todoRoutes"`))),
	)
}

func RunApp(app *fx.App) {
	app.Run()
}
