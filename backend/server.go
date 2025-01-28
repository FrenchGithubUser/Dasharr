package main

import (
	"backend/handlers"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	api := e.Group("/api")

	apiKey := os.Getenv("API_KEY")

	apiKeyMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			providedKey := c.Request().Header.Get("X-API-Key")
			if providedKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid_api_key")
			}
			return next(c)
		}
	}

	api.Use(middleware.Logger())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	api.GET("/collectStats", handlers.CollectStats, apiKeyMiddleware)
	api.GET("/stats", handlers.GetStats, apiKeyMiddleware)
	api.GET("/config", handlers.GetConfig, apiKeyMiddleware)
	api.POST("/saveCredentials", handlers.SaveCredentials, apiKeyMiddleware)
	api.GET("/savedCredentials", handlers.SavedCredentials, apiKeyMiddleware)
	api.GET("/prowlarrConfig", handlers.GetProwlarrIndexerIds, apiKeyMiddleware)

	// these routes are not meant to be called by the user
	api.GET("/initdb", handlers.InitDB, apiKeyMiddleware)
	api.GET("/update", handlers.Update, apiKeyMiddleware)
	e.Logger.Fatal(e.Start(":1323"))
}
