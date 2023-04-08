package app

import (
	"github.com/ireina7/void/conf"
	"github.com/ireina7/void/context"
	"github.com/ireina7/void/db"
	"github.com/ireina7/void/logger"
)

type Context = context.Effect
type Conf = conf.Effect
type Logger = logger.Effect
type Database = db.Effect

// Runnable interface
type Runnable interface {
	Run()
}

type App interface {
	Context
	Conf
	Logger
	Database
	Runnable
}
