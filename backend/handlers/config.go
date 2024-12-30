package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func GetUserConfig(c echo.Context) error {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to read config file"})
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to parse JSON data"})
	}

	for _, item := range items {
		if fillable, exists := item["fillable"].(map[string]interface{}); exists {
			for k := range fillable {
				fillable[k] = "*****"
			}
		}
	}

	return c.JSON(http.StatusOK, items)
}
