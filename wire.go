//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"wfRadius/src"
	"wfRadius/src/root"
)

func BuildApp() (*root.App, error) {
	wire.Build(
		src.Set,
	)
	return new(root.App), nil
}
