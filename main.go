package main

import (
	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/database"
	"github.com/playground-pro-project/playground-pro-api/app/router"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := database.InitDatabase(cfg)
	router.InitRouter(db, e)
	e.Logger.Fatal(e.Start(":8080"))
}
