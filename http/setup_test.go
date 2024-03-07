package main

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type FxTestApp struct {
	*fxtest.App
}

func NewTestApp(t *testing.T, opts ...fx.Option) *FxTestApp {
	opts = append(
		opts,
		fx.Supply(fx.Annotate(t, fx.As(new(fxtest.TB)))),
		fx.Provide(fxtest.NewLifecycle),
	)
	app := fxtest.New(t, opts...)

	if err := app.Err(); err != nil {
		t.Fatal("App initialization failed", err)
	}
	app.RequireStart()
	return &FxTestApp{app}
}

func (t *FxTestApp) Stop() {
	t.RequireStop()
}

func TestValidateApp(t *testing.T) {
	if err := fx.ValidateApp(CreateDefaultApp()); err != nil {
		t.Fatal(err)
	}
}
