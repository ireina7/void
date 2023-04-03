package http

import (
	"errors"
	"net/http"

	"github.com/ireina7/void/conf"
	localConf "github.com/ireina7/void/conf/local"
	"github.com/ireina7/void/db"
	postgres "github.com/ireina7/void/db/postgres"
	"github.com/ireina7/void/logger"
	fileLogger "github.com/ireina7/void/logger/file"
)

type Conf = conf.Effect
type Logger = logger.Effect
type Database = db.Effect

// The main app structure
// Having 3 main effects:
// - Conf: Configuration
// - Logger: Logging
// - Database: Database CRUD
type App struct {
	Conf
	Logger
	Database
}

func Instance() (App, error) {
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

	// Database
	db, err := postgres.Instance(&logger)
	if err != nil {
		return app, err
	}
	app.Database = &db

	return app, nil
}

// The main App logic: Including http routes for all functions
// No error should be returned as they should be handled inside.
func (app *App) Run() {
	// Handle closing
	defer func() {
		err := app.Database.Close()
		if err != nil {
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
		err := app.CreateBlogTables()
		if err != nil {
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
