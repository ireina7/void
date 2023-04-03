package db

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
	// Query blog header (without content)
	QueryBlogHeaders(expr string) ([]BlogHeader, error)

	// Query the whole blog (with content)
	QueryBlog(id int) (Blog, error)

	// Create blog header and blog content tables
	CreateBlogTables() error

	// Add new blog
	AddBlog(blog Blog) error

	// Delete blog by id
	DeleteBlog(blogId int) error

	// Create account table
	CreateAccountTable() error

	// Query account info by name
	QueryAccount(name string) (Account, error)

	// Add account
	AddAccount(name string, passwd string) error

	// Delete account by id
	DeleteAccount(id int) error

	// Drop table by name
	DropTable(tableName string) error

	// Close database connection
	Close() error
}
