package handlers

import (
	"backend/database"
	"backend/trackers"
	"fmt"
	"io"
	"net/http"
	"os"
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

	prowlarrIds := getProwlarrTrackerIdsFromDB()

	indexer := jsonBody.Get("indexer").Str
	trackerId := ""
	for key, value := range prowlarrIds {
		if strings.Contains(value, indexer) {
			trackerId = key
			break
		}
	}
	username := jsonBody.Get("username").Str
	password := jsonBody.Get("password").Str
	cookies := loginAndGetCookies(indexer, username, password)

	insertSQL := `INSERT OR REPLACE INTO credentials (
	tracker_id, username, password, cookies, api_key
	) VALUES (?, ?, ?, ?, ?);
	`

	args := []interface{}{trackerId, username, password, cookies, jsonBody.Get("api_key").Str}

	database.ExecuteQuery(insertSQL, args)

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func loginAndGetCookies(indexer string, username string, password string) string {
	data, _ := os.ReadFile("config/config.json")
	siteInfo := gjson.Get(string(data), fmt.Sprintf(`#[site_name=="%s"]`, indexer))
	if !siteInfo.Exists() {
		return ""
	}
	trackerType := trackers.DetermineTrackerType(indexer)

	loginURL := siteInfo.Get("login.url").String()

	if trackerType == "unit3d" {
		return trackers.LoginAndGetCookiesUnit3d(username, password, loginURL, siteInfo.Get("domain").Str)
	} else if trackerType == "anthelion" {
		return trackers.LoginAndGetCookiesAnthelion(username, password, loginURL, siteInfo)
	}

	return ""

}

// returns the indexers that have their credentials registered in Dasharr's db by the user
func SavedCredentials(c echo.Context) error {
	sql := `SELECT tracker_id from credentials`
	results := database.ExecuteQuery(sql, []interface{}{})
	indexerNames := getProwlarrTrackerIdsFromDB()

	for _, obj := range results {
		if id, ok := obj["tracker_id"].(int64); ok {
			obj["indexer_name"] = strings.TrimSuffix(indexerNames[fmt.Sprint(id)], " (API)")
		}
	}
	return c.JSON(http.StatusOK, results)
}
