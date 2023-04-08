package postgres

import (
	"database/sql"
	"fmt"
	"os"

	database "github.com/ireina7/void/db"
	"github.com/ireina7/void/logger"
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
	err error
}

type BlogHeader = database.BlogHeader
type Blog = database.Blog
type Account = database.Account

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
func (self *DbConnection) Close() {
	err := self.raw.Close()
	if err != nil {
		self.err = err
	}
}

func (self *DbConnection) Exec(expr string, args ...any) (sql.Result, error) {
	return self.raw.Exec(expr, args...)
}

func (self *DbConnection) Query(expr string, args ...any) (*sql.Rows, error) {
	return self.raw.Query(expr, args...)
}

func (conn *DbConnection) QueryBlogHeaders(cond string) []BlogHeader {
	if conn.err != nil {
		return nil
	}
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM blog_headers WHERE %s`, cond))
	if err != nil {
		conn.err = err
		return nil
	}
	defer rows.Close()

	var blogs []BlogHeader
	for rows.Next() {
		var blog BlogHeader
		if err := rows.Scan(&blog.Id, &blog.Title, &blog.Author, &blog.Date, &blog.Preview); err != nil {
			conn.err = fmt.Errorf("Blogerror %v", err)
			return nil
		}
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		conn.err = fmt.Errorf("Blogerror %v", err)
		return nil
	}
	return blogs
}

func (conn *DbConnection) QueryBlog(id int) Blog {
	if conn.err != nil {
		return Blog{}
	}
	blog := Blog{}

	// Query blog headers
	blogHeader := conn.QueryBlogHeaders(fmt.Sprintf(`id = %d`, id))
	if conn.Error() != nil {
		return blog
	}
	blog.BlogHeader = blogHeader[0]

	// Query blogs table for content
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM blogs WHERE id = %d`, id))
	if err != nil {
		return blog
	}
	if rows.Next() {
		if err := rows.Scan(&blog.Id, &blog.Content); err != nil {
			conn.err = err
			return blog
		}
	}
	return blog
}

func (conn *DbConnection) QueryAccount(name string) Account {
	var account Account
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM accounts WHERE name = "%s"`, name))
	if err != nil {
		conn.err = err
		return account
	}
	if rows.Next() {
		if err := rows.Scan(&account.Id, &account.Name, &account.Passwd); err != nil {
			conn.err = err
			return account
		}
	} else {
		conn.err = fmt.Errorf("Account name %s not found", name)
		return account
	}
	return account
}

func (conn *DbConnection) AddAccount(name string, passwd string) {
	if conn.err != nil {
		return
	}
	cmd := fmt.Sprintf(`
		INSERT INTO accounts (name, passwd) 
		VALUES ('%s', '%s');`,
		name, passwd,
	)
	_, err := conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) DeleteAccount(id int) {
	if conn.err != nil {
		return
	}
	cmd := fmt.Sprintf(`DELETE FROM accounts WHERE id = %d`, id)
	_, err := conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) CreateBlogTables() {
	if conn.err != nil {
		return
	}
	conn.CreateBlogHeaderTable()
	conn.CreateBlogTable()
}

func (conn *DbConnection) CreateBlogHeaderTable() {
	if conn.err != nil {
		return
	}
	_, err := conn.Exec(
		`CREATE TABLE blog_headers (
			id      SERIAL PRIMARY KEY,
			title   varchar(50),
			author  varchar(30),
			date    date,
			preview varchar(100)
		);`,
	)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) CreateBlogTable() {
	if conn.err != nil {
		return
	}
	_, err := conn.Exec(
		`CREATE TABLE blogs (
			id      SERIAL PRIMARY KEY,
			content text
		);`,
	)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) CreateAccountTable() {
	if conn.err != nil {
		return
	}
	_, err := conn.Exec(
		`CREATE TABLE accounts (
			id      SERIAL PRIMARY KEY,
			name 	varchar(20),
			passwd  varchar(20)
		);`,
	)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) DropTable(tableName string) {
	if conn.err != nil {
		return
	}
	_, err := conn.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName))
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) AddBlog(blog Blog) {
	if conn.err != nil {
		return
	}
	// Insert into blog headers
	cmd := fmt.Sprintf(`
		INSERT INTO blog_headers (title, author, date, preview) 
		VALUES ('%s', '%s', '%s', '%s');`,
		blog.Title, blog.Author, blog.Date, blog.Preview,
	)
	_, err := conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}

	// Insert into blogs(content)
	cmd = fmt.Sprintf(`
		INSERT INTO blogs (content)
		VALUES ('%s');`,
		blog.Content,
	)
	_, err = conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) UpdateTable(
	tableName string,
	condition string,
	field string,
	value string,
) error {
	if conn.err != nil {
		return nil
	}
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

func (conn *DbConnection) DeleteBlog(blogId int) {
	if conn.err != nil {
		return
	}
	cmd := fmt.Sprintf(`DELETE FROM blog_headers WHERE id = %d`, blogId)
	_, err := conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}
	cmd = fmt.Sprintf(`DELETE FROM blogs WHERE id = %d`, blogId)
	_, err = conn.Exec(cmd)
	if err != nil {
		conn.err = err
		return
	}
}

func (conn *DbConnection) Error() error {
	return conn.err
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
