package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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
		result[id] = strings.TrimSuffix(result[id], "2FA")
	}
	return result

}
