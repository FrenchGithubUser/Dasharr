package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func GetProwlarrIndexerIds(c echo.Context) error {
	result := getProwlarrIndexerIdsFromDB()

	return c.JSON(http.StatusOK, result)
}

func getProwlarrIndexerIdsFromDB() map[string]string {
	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Id, Name FROM Indexers`
	indexers, _ := prowlarrDb.Query(prowlarrReq)

	defer indexers.Close()

	result := make(map[string]string)

	for indexers.Next() {
		var id, name string
		indexers.Scan(&id, &name)
		result[id] = strings.TrimSuffix(name, " (API)")
	}
	return result

}

func getProwlarrCredentials(indexerId interface{}) map[string]string {
	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Settings FROM Indexers WHERE Id = ?`
	indexers, _ := prowlarrDb.Query(prowlarrReq, indexerId)

	defer indexers.Close()

	var result string

	for indexers.Next() {
		var settings string
		indexers.Scan(&settings)
		result = settings
	}
	jsonResult := gjson.Parse(result)

	credentials := map[string]string{
		"username": jsonResult.Get("username").String(),
		"password": jsonResult.Get("password").String(),
	}

	return credentials
}
