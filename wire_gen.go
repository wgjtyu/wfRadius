// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"wfRadius/src/root"
	"wfRadius/src/root/startup"
	"wfRadius/ws"
)

// Injectors from wire.go:

func BuildApp() (*root.App, error) {
	mConfig := startup.InitCfg()
	logger := startup.NewZap(mConfig)
	db, err := startup.InitGorm(logger)
	if err != nil {
		return nil, err
	}
	worker := ws.NewWorker(mConfig, logger)
	app := &root.App{
		Config: mConfig,
		DB:     db,
		Logger: logger,
		Worker: worker,
	}
	return app, nil
}
