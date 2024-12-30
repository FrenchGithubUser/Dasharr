package handlers

import (
	"backend/database"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetStats(c echo.Context) error {

	date_from := c.QueryParam("date_from")
	date_to := c.QueryParam("date_to")
	trackerIds := strings.Split(c.QueryParam("tracker_ids"), ",")

	allResults := make(map[string]interface{})

	allStatsQuery := `
    SELECT *
    FROM user_stats
    WHERE collected_at BETWEEN ? AND ?
      AND tracker_id IN (?` + strings.Repeat(", ?", len(trackerIds)-1) + `);
	`

	args := []interface{}{date_from, date_to}
	for _, id := range trackerIds {
		args = append(args, id)
	}
	allStats := database.ExecuteQuery(allStatsQuery, args)

	allResults["all_stats"] = allStats

	summaryQuery := `
		SELECT
	    tracker_id,
	    MIN(collected_at) AS earliest_date,
	    MAX(collected_at) AS latest_date,
	    MAX(downloaded_amount) - MIN(downloaded_amount) AS downloaded,
	    MAX(uploaded_amount) - MIN(uploaded_amount) AS uploaded,
	    MAX(snatched) - MIN(snatched) AS snatched,
	    MAX(seeding) - MIN(seeding) AS seeding,
	    MAX(ratio) - MIN(ratio) AS ratio,
	    MAX(torrent_comments) - MIN(torrent_comments) AS torrent_comments,
	    MAX(forum_posts) - MIN(forum_posts) AS forum_posts
		FROM user_stats
		WHERE collected_at BETWEEN ? AND ?
		GROUP BY tracker_id;
	`

	summaryStats := database.ExecuteQuery(summaryQuery, []interface{}{date_from, date_to})

	allResults["summary_stats"] = summaryStats

	return c.JSON(http.StatusOK, allResults)

}
