package main

import (
	"backend/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	const apiKey = "your-secure-api-key"

	apiKeyMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			providedKey := c.Request().Header.Get("X-API-Key")
			if providedKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
			}
			return next(c)
		}
	}

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.GET("/initdb", handlers.InitDB, apiKeyMiddleware)
	// e.GET("/updateOld", handlers.UpdateOld)
	e.GET("/update", handlers.Update, apiKeyMiddleware)
	e.GET("/stats", handlers.GetStats, apiKeyMiddleware)
	e.GET("/userConfig", handlers.GetUserConfig, apiKeyMiddleware)
	e.GET("/prowlarrConfig", handlers.GetProwlarrTrackerIds, apiKeyMiddleware)
	e.Logger.Fatal(e.Start(":1323"))
}
