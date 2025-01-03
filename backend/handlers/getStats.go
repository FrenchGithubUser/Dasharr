package handlers

import (
	"backend/database"
	"backend/helpers"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetStats(c echo.Context) error {

	date_from := c.QueryParam("date_from")
	date_to := c.QueryParam("date_to")
	trackerIds := strings.Split(c.QueryParam("tracker_ids"), ",")

	allResults := make(map[string]interface{})

	// allStatsQuery := `
	// SELECT *, DATE(collected_at) AS day, MAX(collected_at) AS latest_collected_at
	// FROM user_stats
	// WHERE collected_at BETWEEN ? AND ?
	// AND tracker_id IN (?` + strings.Repeat(", ?", len(trackerIds)-1) + `)
	// GROUP BY tracker_id, day;
	// `
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
	allStats := helpers.RemoveNilEntries(database.ExecuteQuery(allStatsQuery, args))

	allResults["all"] = allStats

	summaryQuery := `
		SELECT
	    tracker_id,
		(SELECT downloaded_amount FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) - 
		(SELECT downloaded_amount FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS downloaded_amount,

	    (SELECT uploaded_amount FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) - 
		(SELECT uploaded_amount FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS uploaded_amount, 

		(SELECT snatched FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) - 
		(SELECT snatched FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS snatched, 

		(SELECT seeding FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) -
		(SELECT seeding FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS seeding, 

		(SELECT ratio FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) - 
		(SELECT ratio FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS ratio, 

		(SELECT torrent_comments FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) -
		(SELECT torrent_comments FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS torrent_comments, 

		(SELECT forum_posts FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at DESC LIMIT 1) - 
		(SELECT forum_posts FROM user_stats WHERE tracker_id = u.tracker_id AND collected_at BETWEEN ? AND ? ORDER BY collected_at ASC LIMIT 1) AS forum_posts

		FROM user_stats u
		WHERE collected_at BETWEEN ? AND ?
		AND tracker_id IN (?` + strings.Repeat(", ?", len(trackerIds)-1) + `)
		GROUP BY tracker_id;
	`
	amount_dates_required := 15
	args = []interface{}{}
	for i := 0; i < amount_dates_required; i++ {
		args = append(args, date_from, date_to)
	}
	for _, id := range trackerIds {
		args = append(args, id)
	}

	perTrackerSummaryStats := helpers.RemoveNilEntries(database.ExecuteQuery(summaryQuery, args))

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
