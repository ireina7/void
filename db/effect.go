package db

import "github.com/ireina7/void/context"

// BlogHeader table
// | id | author | date | preview |
type BlogHeader struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Preview string `json:"preview"`
}

// Blog table
// | id | content |
type Blog struct {
	BlogHeader
	Content string `json:"content"`
}

// Account info
type Account struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

type Effect interface {
	context.Effect
	// Query blog header (without content)
	QueryBlogHeaders(expr string) []BlogHeader

	// Query the whole blog (with content)
	QueryBlog(id int) Blog

	// Create blog header and blog content tables
	CreateBlogTables()

	// Add new blog
	AddBlog(blog Blog)

	// Delete blog by id
	DeleteBlog(blogId int)

	// Create account table
	CreateAccountTable()

	// Query account info by name
	QueryAccount(name string) Account

	// Add account
	AddAccount(name string, passwd string)

	// Delete account by id
	DeleteAccount(id int)

	// Drop table by name
	DropTable(tableName string)

	// Close database connection
	Close()
}
