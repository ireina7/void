package http

import (
	"errors"
	"net/http"

	"github.com/ireina7/void/conf"
	localConf "github.com/ireina7/void/conf/local"
	"github.com/ireina7/void/context"
	"github.com/ireina7/void/db"
	postgres "github.com/ireina7/void/db/postgres"
	"github.com/ireina7/void/logger"
	fileLogger "github.com/ireina7/void/logger/file"
)

type Conf = conf.Effect
type Logger = logger.Effect
type Database = db.Effect
type Context = context.Effect

// The main app structure
// Having 4 main effects:
// - Context: Context abstraction
// - Conf: Configuration
// - Logger: Logging
// - Database: Database CRUD
type HttpApp struct {
	Context
	Conf
	Logger
	Database
}

func Instance() (HttpApp, error) {
	var app HttpApp = HttpApp{}

	// Configuration
	conf := localConf.Instance()
	app.Conf = &conf

	// Logging
	logger, err := fileLogger.Instance()
	if err != nil {
		return app, err
	}
	app.Logger = &logger

	// Database
	db, err := postgres.Instance(&logger)
	if err != nil {
		return app, err
	}
	app.Database = &db

	// ctx := context.LocalContext{}
	// app.Context = &ctx

	return app, nil
}

func (app *HttpApp) Error() error {
	return nil
}

// The main App logic: Including http routes for all functions
// No error should be returned as they should be handled inside.
func (app *HttpApp) Run() {
	// Handle closing
	defer func() {
		app.Database.Close()
		if err := app.Database.Error(); err != nil {
			app.Logger.Fatal(
				errors.New("Database closing error: " + err.Error()),
			)
		}
	}()

	// Database new creation?
	dbRebuild := false
	if dbRebuild {
		app.DropTable("blog_headers")
		app.DropTable("blogs")
		app.CreateBlogTables()
		if err := app.Database.Error(); err != nil {
			app.Fatal(err)
		}
	}

	// HTTP routes
	http.HandleFunc("/blogs", corsHandler(app.queryBlogs))
	http.HandleFunc("/blog/", corsHandler(app.queryBlog))
	http.HandleFunc("/submit", corsHandler(app.submitBlog))
	http.HandleFunc("/delete/", corsHandler(app.deleteBlog))
	http.HandleFunc("/login", corsHandler(app.checkAccount))
	http.ListenAndServe(conf.Addr(app), nil)
}
