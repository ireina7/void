package main

import (
	"fmt"
	"log"

	appapi "github.com/ireina7/void/app"
	httpApp "github.com/ireina7/void/app/http"
	"github.com/ireina7/void/utils"
	_ "github.com/lib/pq"
)

func main() {
	utils.LoadEnv()
	var app appapi.App
	httpApp, err := httpApp.Instance()
	app = &httpApp
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Void version %s\n", app.Version())
	fmt.Println("Void link start!")
	app.Run()
}
