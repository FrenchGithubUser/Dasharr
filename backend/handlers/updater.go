package handlers

import (
	"backend/helpers"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func Update(c echo.Context) error {
	db, _ := sql.Open("sqlite3", "database.db")

	configFile, _ := os.ReadFile("config/config.json")
	trackerConfigs := gjson.Parse(string(configFile))

	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	trackerConfigs.ForEach(func(_, trackerConfig gjson.Result) bool {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(trackerConfig gjson.Result) {
			defer wg.Done()
			defer func() { <-semaphore }()

			processTracker(trackerConfig, db)
		}(trackerConfig)

		return true
	})

	wg.Wait()

	return c.String(http.StatusOK, "Data inserted successfully!")
}

func processTracker(trackerConfig gjson.Result, db *sql.DB) bool {
	trackerName := trackerConfig.Get("tracker_name").String()
	enabled := trackerConfig.Get("fillable.enabled").Bool()

	if enabled {
		fmt.Printf("Updating %s's stats\n", trackerName)

		trackerStats := helpers.GetUserData(trackerConfig)

		insertSQL := `INSERT INTO user_stats (
			tracker_id, uploaded_torrents, uploaded_amount, downloaded_amount, snatched, seeding, leeching, 
			ratio, required_ratio, last_access, torrent_comments, invited, forum_posts, warned, class, 
			donor, uploaded_rank, downloaded_rank, uploads_rank, requests_rank, bounty_rank, posts_rank, 
			artists_rank, overall_rank, buffer, bonus_points
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

		_, err := db.Exec(insertSQL,
			trackerConfig.Get("tracker_id").Int(),
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
		)

		if err != nil {
			fmt.Println("Failed to insert data:", err.Error())
		}
	}

	return true
}
