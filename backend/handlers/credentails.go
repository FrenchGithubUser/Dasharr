package handlers

import (
	"backend/database"
	"backend/indexers"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func SaveCredentials(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	jsonBody := gjson.Parse(string(body))

	prowlarrIds := getProwlarrIndexerIdsFromDB()

	indexer := jsonBody.Get("indexer").Str
	indexerId := ""
	for key, value := range prowlarrIds {
		if strings.Contains(value, indexer) {
			indexerId = key
			break
		}
	}
	username := jsonBody.Get("username").Str
	password := jsonBody.Get("password").Str
	twoFaCode := jsonBody.Get("twoFaCode").Str
	loginError := indexers.LoginAndSaveCookies(indexer, username, password, twoFaCode, jsonBody.Get("api_key").Str, indexerId)

	if loginError == nil {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	} else {
		fmt.Printf(loginError.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "login_failed"})
	}
}

// returns the indexers that have their credentials registered in Dasharr's db by the user
func SavedCredentials(c echo.Context) error {
	sql := `SELECT indexer_id from credentials`
	results := database.ExecuteQuery(sql, []interface{}{})
	if results == nil {
		results = []map[string]interface{}{}
	}
	indexerNames := getProwlarrIndexerIdsFromDB()

	for _, obj := range results {
		if id, ok := obj["indexer_id"].(int64); ok {
			obj["indexer_name"] = strings.TrimSuffix(indexerNames[fmt.Sprint(id)], " (API)")
			obj["indexer_name"] = strings.TrimSuffix(indexerNames[fmt.Sprint(id)], "2FA")
		}
	}
	return c.JSON(http.StatusOK, results)
}
