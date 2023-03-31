package main

import (
	"github.com/ireina7/void/conf"
	localConf "github.com/ireina7/void/conf/local"
	"github.com/ireina7/void/db"
	"github.com/ireina7/void/logger"
	fileLogger "github.com/ireina7/void/logger/file"
)

type App struct {
	Conf   conf.Effect
	Logger logger.Effect
	Db     db.Effect
}

func Build() (App, error) {
	var app App = App{}

	// Configuration
	conf := localConf.Instance()
	app.Conf = &conf

	// Logging
	logger, err := fileLogger.Instance()
	if err != nil {
		return app, err
	}
	app.Logger = &logger

	return app, nil
}

func (app *App) Run() {}
