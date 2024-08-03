package main

import (
	"chat_app/config"
	"chat_app/internal/router"
	"chat_app/pkg/db"
	"chat_app/pkg/utils"
	"log"
	"net/http"
)

func main() {
	// load config
	config, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error In Loading Config %v", err)
	}
	runner, err := db.RunDB(config)
	if err != nil {
		log.Fatal("Error in database connection")
	}
	err = runner.Migrate()
	if err != nil {
		log.Fatal("Error in Migrate Tables")
	}

	r := router.NewRouter(runner.DB, runner.Client, config)

	utils.Print("Server Started At....", config.ServerAddress)
	if err := http.ListenAndServe(config.ServerAddress, r); err != nil {
		log.Fatalf("Error In Serving %v", err)

	}
}
