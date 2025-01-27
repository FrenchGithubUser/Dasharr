package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// apply the necessary modifications to update an already existing dasharr install
// it will run every time the container is updated or restarted
func Update(c echo.Context) error {
	existingConfig, _ := os.ReadFile("config/config.json")
	updatedConfig, _ := os.ReadFile("config_sample/config_sample.json")

	result := string(existingConfig)
	gjson.ParseBytes(updatedConfig).ForEach(func(_, value gjson.Result) bool {
		if !gjson.Get(string(existingConfig), `#(indexer_name=="`+value.Get("indexer_name").String()+`")`).Exists() {
			result, _ = sjson.SetRaw(result, "-1", value.Raw)
		}
		return true
	})

	os.WriteFile("config/config.json", []byte(result), 0644)
	fmt.Println("Updated dasharr's config")
	return c.JSON(http.StatusOK, "")
}
