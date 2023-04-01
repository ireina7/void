package postgres

import (
	"database/sql"
	"fmt"
	"os"

	database "github.com/ireina7/void/db"
	"github.com/ireina7/void/logger"
	"github.com/ireina7/void/utils"
)

// // BlogHeader table
// // | id | author | date | preview |
// type BlogHeader struct {
// 	id      int
// 	title   string
// 	author  string
// 	date    string
// 	preview string
// }

// // Blog table
// // | id | content |
// type Blog struct {
// 	BlogHeader
// 	content string
// }

var db *sql.DB

type DbParam struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type Logger = logger.Effect
type DbConnection struct {
	DbParam
	Logger
	raw *sql.DB
}

type BlogHeader = database.BlogHeader
type Blog = database.Blog

// Connect to database
func (self *DbConnection) Init() error {
	psqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		self.Host, self.Port, self.User, self.Password, self.DbName,
	)
	// fmt.Println(psqlconn)
	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	self.raw = db
	return nil
}

// Close database connection
func (self *DbConnection) Close() error {
	return self.raw.Close()
}

func (self *DbConnection) Exec(expr string, args ...any) (sql.Result, error) {
	return self.raw.Exec(expr, args...)
}

func (self *DbConnection) Query(expr string, args ...any) (*sql.Rows, error) {
	return self.raw.Query(expr, args...)
}

func (conn *DbConnection) QueryBlogHeaders(cond string) ([]BlogHeader, error) {
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM blog_headers WHERE %s`, cond))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []BlogHeader
	for rows.Next() {
		var blog BlogHeader
		if err := rows.Scan(&blog.Id, &blog.Title, &blog.Author, &blog.Date, &blog.Preview); err != nil {
			return nil, fmt.Errorf("Blogerror %v", err)
		}
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Blogerror %v", err)
	}
	return blogs, nil
}

func (conn *DbConnection) QueryBlog(id int) (Blog, error) {
	blog := Blog{}

	// Query blog headers
	blogHeader, err := conn.QueryBlogHeaders(fmt.Sprintf(`id = %d`, id))
	if err != nil {
		return blog, err
	}
	blog.BlogHeader = blogHeader[0]

	// Query blogs table for content
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM blogs WHERE id = %d`, id))
	if err != nil {
		return blog, nil
	}
	if rows.Next() {
		if err := rows.Scan(&blog.Id, &blog.Content); err != nil {
			return blog, err
		}
	}
	return blog, nil
}

func (conn *DbConnection) CreateTables() error {
	err := conn.CreateBlogHeaderTable()
	if err != nil {
		return err
	}
	err = conn.CreateBlogTable()
	if err != nil {
		return err
	}
	return nil
}

func (conn *DbConnection) CreateBlogHeaderTable() error {
	res, err := conn.Exec(
		`CREATE TABLE blog_headers (
			id      SERIAL PRIMARY KEY,
			title   varchar(50),
			author  varchar(30),
			date    date,
			preview varchar(100)
		);`,
	)
	if err != nil {
		return err
	}
	utils.Use(res)
	return nil
}

func (conn *DbConnection) CreateBlogTable() error {
	_, err := conn.Exec(
		`CREATE TABLE blogs (
			id      SERIAL PRIMARY KEY,
			content text
		);`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (conn *DbConnection) DropTable(tableName string) error {
	res, err := conn.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName))
	if err != nil {
		return err
	}
	utils.Use(res)
	return nil
}

func (conn *DbConnection) AddBlog(blog Blog) error {
	// Insert into blog headers
	cmd := fmt.Sprintf(`
		INSERT INTO blog_headers (title, author, date, preview) 
		VALUES ('%s', '%s', '%s', '%s');`,
		blog.Title, blog.Author, blog.Date, blog.Preview,
	)
	_, err := conn.Exec(cmd)
	if err != nil {
		return err
	}

	// Insert into blogs(content)
	cmd = fmt.Sprintf(`
		INSERT INTO blogs (content)
		VALUES ('%s');`,
		blog.Content,
	)
	_, err = conn.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (conn *DbConnection) UpdateTable(
	tableName string,
	condition string,
	field string,
	value string,
) error {
	cmd := fmt.Sprintf(`
		UPDATE %s SET %s = %s WHERE %s;`,
		tableName, field, value, condition,
	)
	_, err := conn.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (conn *DbConnection) DeleteBlog(blogId int) error {
	cmd := fmt.Sprintf(`DELETE FROM blog_headers WHERE id = %d`, blogId)
	_, err := conn.Exec(cmd)
	if err != nil {
		return err
	}
	cmd = fmt.Sprintf(`DELETE FROM blogs WHERE id = %d`, blogId)
	_, err = conn.Exec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func Instance(logger logger.Effect) (DbConnection, error) {
	db := DbConnection{
		DbParam: DbParam{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
		},
		Logger: logger,
	}
	err := db.Init()
	if err != nil {
		return db, err
	}
	db.Info("Database connected")
	return db, nil
}
