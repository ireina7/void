package main

import (
	"fmt"
	"log"

	"github.com/ireina7/void/utils"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Void link start!")
	utils.LoadEnv()
	app, err := Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()

	// err = db.CreateTables()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
