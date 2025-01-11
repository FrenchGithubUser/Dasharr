package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetProwlarrTrackerIds(c echo.Context) error {
	result := getProwlarrTrackerIdsFromDB()

	return c.JSON(http.StatusOK, result)
}

func getProwlarrTrackerIdsFromDB() map[string]string {
	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Id, Name FROM Indexers`
	trackers, _ := prowlarrDb.Query(prowlarrReq)

	defer trackers.Close()

	result := make(map[string]string)

	for trackers.Next() {
		var id, name string
		trackers.Scan(&id, &name)
		result[id] = strings.TrimSuffix(name, " (API)")
	}
	return result

}
