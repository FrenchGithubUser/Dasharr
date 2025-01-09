package handlers

import (
	"backend/trackers"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func Update(c echo.Context) error {
	db, _ := sql.Open("sqlite3", "config/database.db")

	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Id, Name, Settings FROM Indexers`
	trackers, _ := prowlarrDb.Query(prowlarrReq)
	defer trackers.Close()

	cols, _ := trackers.Columns()
	trackerConfig := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range ptrs {
		ptrs[i] = &trackerConfig[i]
	}

	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for trackers.Next() {
		trackers.Scan(ptrs...)
		// Make a copy of trackerConfig to avoid race conditions
		configCopy := make([]interface{}, len(trackerConfig))
		copy(configCopy, trackerConfig)

		wg.Add(1)
		semaphore <- struct{}{}
		go func(configCopy []interface{}) {
			defer wg.Done()
			defer func() { <-semaphore }()

			//temp only do blu
			// if configCopy[1] == "Blutopia (API)" {
			processTrackerProwlarr(configCopy, db)
			// }
		}(configCopy)
	}

	wg.Wait()

	return c.String(http.StatusOK, "Data inserted successfully!")
}

func processTrackerProwlarr(trackerConfig []interface{}, db *sql.DB) bool {
	trackerName := trackerConfig[1].(string)
	trackerName = strings.TrimSuffix(trackerName, " (API)")
	// enabled := trackerConfig.Get("fillable.enabled").Bool()

	// if enabled {
	fmt.Printf("Updating %s's stats\n", trackerName)

	trackerStats, error := trackers.GetUserData(gjson.Parse(trackerConfig[2].(string)), trackerName, trackerConfig[0].(int64))
	if error != nil {
		return false
	}

	insertSQL := `INSERT INTO user_stats (
			tracker_id, uploaded_torrents, uploaded_amount, downloaded_amount, snatched, seeding, leeching,
			ratio, required_ratio, last_access, torrent_comments, invited, forum_posts, warned, class,
			donor, uploaded_rank, downloaded_rank, uploads_rank, requests_rank, bounty_rank, posts_rank,
			artists_rank, overall_rank, buffer, bonus_points, seeding_size, freeleech_tokens
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.Exec(insertSQL,
		trackerConfig[0],
		trackerStats["uploaded_torrents"],
		trackerStats["uploaded_amount"],
		trackerStats["downloaded_amount"],
		trackerStats["snatched"],
		trackerStats["seeding"],
		trackerStats["leeching"],
		trackerStats["ratio"],
		trackerStats["required_ratio"],
		trackerStats["last_access"],
		trackerStats["torrent_comments"],
		trackerStats["invited"],
		trackerStats["forum_posts"],
		trackerStats["warned"],
		trackerStats["class"],
		trackerStats["donor"],
		trackerStats["uploaded_rank"],
		trackerStats["downloaded_rank"],
		trackerStats["uploads_rank"],
		trackerStats["requests_rank"],
		trackerStats["bounty_rank"],
		trackerStats["posts_rank"],
		trackerStats["artists_rank"],
		trackerStats["overall_rank"],
		trackerStats["buffer"],
		trackerStats["bonus_points"],
		trackerStats["seeding_size"],
		trackerStats["freeleech_tokens"],
	)

	if err != nil {
		fmt.Println("Failed to insert data:", err.Error())
	}
	// }

	return true
}
