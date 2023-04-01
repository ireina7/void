package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ireina7/void/conf"
	localConf "github.com/ireina7/void/conf/local"
	"github.com/ireina7/void/db"
	postgres "github.com/ireina7/void/db/postgres"
	"github.com/ireina7/void/logger"
	fileLogger "github.com/ireina7/void/logger/file"
)

type Runnable interface {
	Run()
}

type Conf = conf.Effect
type Logger = logger.Effect
type Database = db.Effect
type App struct {
	Conf
	Logger
	Database
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

	// Database
	db, err := postgres.Instance(&logger)
	if err != nil {
		return app, err
	}
	app.Database = &db

	return app, nil
}

// The main App logic
func (app *App) Run() {
	defer func() {
		err := app.Database.Close()
		if err != nil {
			app.Logger.Fatal(
				errors.New("Database closing error: " + err.Error()),
			)
		}
	}()

	dbRebuild := false
	if dbRebuild {
		app.DropTable("blog_headers")
		app.DropTable("blogs")
		err := app.CreateTables()
		if err != nil {
			app.Fatal(err)
		}
	}

	http.HandleFunc("/blogs", corsHandler(app.queryBlogs))
	http.HandleFunc("/blog/", corsHandler(app.queryBlog))
	http.HandleFunc("/submit", corsHandler(app.submitBlog))
	http.ListenAndServe(conf.Addr(app), nil)
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization,Access-Control-Allow-Headers,Origin,Cookie,Set-Cookie,Accept,X-Requested-With,Content-Type,Access-Control-Request-Method,Access-Control-Request-Headers")
}

type Handler = func(w http.ResponseWriter, req *http.Request)

func corsHandler(h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//handle preflight in here
			setupCORS(&w, r)
		} else {
			h(w, r)
		}
	}
}

func (app *App) queryBlogs(w http.ResponseWriter, req *http.Request) {
	app.Info("Querying blogs...")
	setupCORS(&w, req)
	blogs, err := app.QueryBlogHeaders("1 = 1")
	if err != nil {
		app.Error(err)
	}
	fmt.Println(blogs)
	blogJson, err := json.Marshal(blogs)
	w.Write(blogJson)
}

func (app *App) submitBlog(w http.ResponseWriter, req *http.Request) {
	app.Info("Submitting blog...")
	setupCORS(&w, req)
	var blog db.Blog
	// fmt.Println(req.Body)
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.AddBlog(blog)
	if err != nil {
		app.Error(err)
	}
}

func (app *App) queryBlog(w http.ResponseWriter, req *http.Request) {
	app.Info("Querying blog...")
	setupCORS(&w, req)
	var id int
	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/blog/"))
	if err != nil {
		app.Error(err)
	}
	blog, err := app.QueryBlog(id)
	if err != nil {
		app.Error(err)
	}
	blogJson, err := json.Marshal(blog)
	w.Write(blogJson)
}

func debugRequest(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	var err error

	if req.Body != nil {
		bodyBytes, err = ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer req.Body.Close()
	}

	fmt.Printf("Headers: %+v\n", req.Header)

	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}
}
