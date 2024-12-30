package handlers

import (
	"backend/helpers"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetStats(c echo.Context) error {

	trackerIds := strings.Split(c.QueryParam("tracker_ids"), ",")

	query := `
    SELECT *
    FROM user_stats
    WHERE collected_at BETWEEN ? AND ?
      AND tracker_id IN (?` + strings.Repeat(", ?", len(trackerIds)-1) + `);
	`

	args := []interface{}{c.QueryParam("date_from"), c.QueryParam("date_to")}
	for _, id := range trackerIds {
		args = append(args, id)
	}

	db, err := sql.Open("sqlite3", "database.db")
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Query failed: %s", err)
	}
	defer rows.Close()
	jsonRows, _ := helpers.RowsToJSON(rows)

	return c.JSON(http.StatusOK, jsonRows)

}
