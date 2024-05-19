package main

import (
	"log"

	"github.com/proj-go-5/accounts/internal/api"
)

func main() {
	app, err := api.NewApp()
	if err != nil {
		log.Fatal(err)
		return
	}
	app.Run()
}
