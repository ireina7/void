// Supportiung routes:
// - /blogs
// - /blog/:id
// - /submit
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ireina7/void/db"
	"github.com/ireina7/void/utils"
)

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set(
		"Access-Control-Allow-Headers",
		"Authorization,Access-Control-Allow-Headers,Origin,Cookie,Set-Cookie,Accept,X-Requested-With,Content-Type,Access-Control-Request-Method,Access-Control-Request-Headers",
	)
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

func (app *HttpApp) queryBlogs(w http.ResponseWriter, req *http.Request) {
	app.Info("Querying blogs...")
	setupCORS(&w, req)
	blogs := app.QueryBlogHeaders("1 = 1")
	if err := app.Database.Error(); err != nil {
		app.LogError(err)
	}
	fmt.Println(blogs)
	blogJson, err := json.Marshal(blogs)
	if err != nil {
		app.LogError(err)
	}
	w.Write(blogJson)
}

func (app *HttpApp) submitBlog(w http.ResponseWriter, req *http.Request) {
	app.Info("Submitting blog...")
	setupCORS(&w, req)
	var blog db.Blog
	// fmt.Println(req.Body)
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.AddBlog(blog)
	if err = app.Database.Error(); err != nil {
		app.LogError(err)
	}
}

func (app *HttpApp) queryBlog(w http.ResponseWriter, req *http.Request) {
	setupCORS(&w, req)
	var id int
	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/blog/"))
	app.Info(fmt.Sprintf("Querying blog { id: %d }", id))
	if err != nil {
		app.LogError(err)
		return
	}
	blog := app.QueryBlog(id)
	if err = app.Database.Error(); err != nil {
		app.LogError(err)
	}
	blogJson, err := json.Marshal(blog)
	w.Write(blogJson)
}

func (app *HttpApp) deleteBlog(w http.ResponseWriter, req *http.Request) {
	setupCORS(&w, req)
	var id int
	id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/delete/"))
	app.Info(fmt.Sprintf("Deleting blog { id: %d }", id))
	if err != nil {
		app.LogError(err)
		return
	}
	app.DeleteBlog(id)
	if err := app.Database.Error(); err != nil {
		app.LogError(err)
	}
}

type AccountQuery struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

func (app *HttpApp) checkAccount(w http.ResponseWriter, req *http.Request) {
	setupCORS(&w, req)
	var query AccountQuery
	err := json.NewDecoder(req.Body).Decode(&query)
	// var accountName string
	// accountName = strings.TrimPrefix(req.URL.Path, "/login/")
	app.Info(fmt.Sprintf("Checking account { name: %s }", query.Name))
	account := app.QueryAccount(query.Name)
	if err = app.Database.Error(); err != nil {
		app.LogError(err)
		return
	}
	utils.Use(account)
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
