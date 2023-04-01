package main

import (
	"fmt"
	"log"

	app "github.com/ireina7/void/app/http"
	"github.com/ireina7/void/utils"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Void link start!")
	utils.LoadEnv()
	app, err := app.Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
