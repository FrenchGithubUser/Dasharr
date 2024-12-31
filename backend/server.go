package main

import (
	"backend/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.GET("/initdb", handlers.InitDB)
	e.GET("/update", handlers.Update)
	e.GET("/stats", handlers.GetStats)
	e.GET("/userConfig", handlers.GetUserConfig)
	e.Logger.Fatal(e.Start(":1323"))
}
