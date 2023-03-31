package db

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

type DbParam struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type DbConnection struct {
	DbParam
	raw *sql.DB
}

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
	self.raw.Close()
}

func (self *DbConnection) Exec(expr string, args ...any) (sql.Result, error) {
	return self.raw.Exec(expr, args...)
}

func (self *DbConnection) Query(expr string, args ...any) (*sql.Rows, error) {
	return self.raw.Query(expr, args...)
}
