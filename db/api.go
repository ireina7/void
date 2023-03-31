package db

import (
	"fmt"

	"github.com/ireina7/void/utils"
)

// BlogHeader table
// | id | author | date | preview |
type BlogHeader struct {
	id      int
	title   string
	author  string
	date    string
	preview string
}

// Blog table
// | id | content |
type Blog struct {
	BlogHeader
	content string
}

type Effect interface {
	QueryBlogHeaders(expr string) ([]BlogHeader, error)
	QueryBlog(id int) (Blog, error)
	CreateTables() error
	AddBlog(blog Blog) error
	DeleteBlog(blogId int) error
	DropTable(tableName string) error
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
		if err := rows.Scan(&blog.id, &blog.author, &blog.date, &blog.preview); err != nil {
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
		if err := rows.Scan(&blog.id, &blog.content); err != nil {
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
			id      integer PRIMARY KEY,
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
			id      integer PRIMARY KEY,
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
	cmd := fmt.Sprintf(`
		INSERT INTO blogs (id, title, author, date, preview) 
		VALUES (%d, '%s', '%s', '%s', '%s');`,
		blog.id, blog.title, blog.author, blog.date, blog.preview,
	)
	_, err := conn.Exec(cmd)
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
