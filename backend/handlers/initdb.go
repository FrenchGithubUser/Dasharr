package handlers

import (
	"backend/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitDB(c echo.Context) error {
	err := database.InitDB()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to initialize database: "+err.Error())
	}
	return c.String(http.StatusOK, "Database initialized successfully!")
}
