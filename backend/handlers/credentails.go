package handlers

import (
	"backend/database"
	"backend/trackers"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	body := siteInfo.Get("login.body").String()
	fields := siteInfo.Get("login.fields").Map()

	if trackerType == "unit3d" {
		return trackers.LoginAndGetCookiesUnit3d(username, password, loginURL, siteInfo.Get("domain").Str)
	} else {
		if body == "form_data" {
			formData := url.Values{}
			formData.Add(fields["username"].String(), username)
			formData.Add(fields["password"].String(), password)

			extraFields := fields["extra"].Map()
			for key, val := range extraFields {
				formData.Add(key, val.String())
			}

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse // Prevents redirect
				},
			}
			resp, err := client.PostForm(loginURL, formData)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()

			cookies := resp.Cookies()
			cookieName := siteInfo.Get("login.cookie_name").Str
			for _, cookie := range cookies {
				if cookie.Name == cookieName {
					// todo: return an array of cookies instead
					return fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
				}
			}
		}

	}

	return ""

}
