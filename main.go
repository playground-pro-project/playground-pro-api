package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/database"
	"github.com/playground-pro-project/playground-pro-api/app/router"
)

func main() {
	cfg := config.InitConfig()
	db := database.InitDatabase(cfg)

	_, err := db.DB()
	if err != nil {
		log.Fatal("error while connect to db ", err)
	}

	log.Println("success connected to db")

	e := echo.New()
	router.InitRouter(db, e)

	e.Logger.Fatal(e.Start(":8080"))
}
