package main

import (
	"log"

	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/database"
)

func main() {
	cfg := config.InitConfig()
	db := database.InitDatabase(cfg)

	_, err := db.DB()
	if err != nil {
		log.Fatal("error while connect to db ", err)
	}

	log.Println("success connected to db")
}
