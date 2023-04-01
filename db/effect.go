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

type Effect interface {
	QueryBlogHeaders(expr string) ([]BlogHeader, error)
	QueryBlog(id int) (Blog, error)
	CreateTables() error
	AddBlog(blog Blog) error
	DeleteBlog(blogId int) error
	DropTable(tableName string) error
	Close() error
}
