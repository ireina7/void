package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ireina7/void/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := db.DbConnection{
		DbParam: db.DbParam{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
		},
	}
	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// err = db.CreateTables()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(`"Hello, void"`)
	app, err := Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Logger.Info("Database connected")
	app.Run()
}
