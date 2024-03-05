package main

import (
	"go.uber.org/fx"
	"testing"
)

func TestValidateApp(t *testing.T) {
	if err := fx.ValidateApp(CreateDefaultApp()); err != nil {
		t.Fatal(err)
	}
}
