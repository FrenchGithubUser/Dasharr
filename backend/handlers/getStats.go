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

	allResults["all"] = allStats

	summaryQuery := `
		SELECT
	    tracker_id,
	    MIN(collected_at) AS earliest_date,
	    MAX(collected_at) AS latest_date,
		(SELECT downloaded_amount FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT downloaded_amount FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS downloaded,
	    (SELECT uploaded_amount FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT uploaded_amount FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS uploaded, 
		(SELECT snatched FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT snatched FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS snatched, 
		(SELECT seeding FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT seeding FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS seeding, 
		(SELECT ratio FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT ratio FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS ratio, 
		(SELECT torrent_comments FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT torrent_comments FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS torrent_comments, 
		(SELECT forum_posts FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at DESC LIMIT 1) - (SELECT forum_posts FROM user_stats WHERE tracker_id = u.tracker_id ORDER BY collected_at ASC LIMIT 1) AS forum_posts
		FROM user_stats u
		WHERE collected_at BETWEEN ? AND ?
		GROUP BY tracker_id;
	`

	perTrackerSummaryStats := database.ExecuteQuery(summaryQuery, []interface{}{date_from, date_to})

	allResults["per_tracker_summary"] = perTrackerSummaryStats

	totalSummaryStats := map[string]float64{}
	for _, item := range perTrackerSummaryStats {
		for k, v := range item {
			if k != "tracker_id" {
				switch num := v.(type) {
				case float64:
					totalSummaryStats[k] += num
				case int64:
					totalSummaryStats[k] += float64(num)
				}
			}
		}
	}
	allResults["total_summary"] = totalSummaryStats

	return c.JSON(http.StatusOK, allResults)

}
