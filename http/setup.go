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
	srv := &http.Server{Addr: settings.Addr, Handler: router}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
			go srv.Serve(ln)
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
func CreateDefaultApp() *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewSettingsHttp,
			NewServeMux,
			NewHttpServer,
			fx.Annotate(NewApiV1Router, fx.ResultTags(`name:"ApiV1Router"`)),
			AsRoute(NewTodosHandlerGet, `group:"todoRoutes"`),
			AsRoute(NewTodosHandlerPost, `group:"todoRoutes"`),
			zap.NewExample,
		),
		fx.Invoke(RegisterSentryMiddleware),
		fx.Invoke(RegisterGinZapLogger),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(fx.Annotate(RegisterTodoApi, fx.ParamTags(`name:"ApiV1Router"`, `group:"todoRoutes"`))),
		// TODO: provide way to register redis client/cache warmup
		// TODO: provide way to register postgres client/use pool
	)
}

func RunApp(app *fx.App) {
	app.Run()
}
