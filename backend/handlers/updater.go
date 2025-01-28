package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// apply the necessary modifications to update an already existing dasharr install
// it will run every time the container is updated or restarted
func Update(c echo.Context) error {
	existingData, _ := os.ReadFile("config/config.json")
	updatedData, _ := os.ReadFile("config_sample/config_sample.json")

	var existingConfig, updatedConfig []map[string]interface{}
	_ = json.Unmarshal(existingData, &existingConfig)
	_ = json.Unmarshal(updatedData, &updatedConfig)

	existingMap := make(map[string]map[string]interface{})
	for _, obj := range existingConfig {
		indexerName := obj["indexer_name"].(string)
		existingMap[indexerName] = obj
	}

	for _, updatedObj := range updatedConfig {
		indexerName := updatedObj["indexer_name"].(string)
		if existingObj, found := existingMap[indexerName]; found {
			for k, v := range updatedObj {
				if k == "enabled" {
					continue
				}
				existingObj[k] = v
			}
		} else {
			existingConfig = append(existingConfig, updatedObj)
		}
	}

	result, _ := json.MarshalIndent(existingConfig, "", "  ")
	os.WriteFile("config/config.json", result, 0644)

	return c.JSON(http.StatusOK, "")
}
