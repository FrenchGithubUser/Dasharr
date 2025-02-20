package handlers

import (
	"backend/helpers"
	"backend/indexers"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func CollectStats(c echo.Context) error {
	db, _ := sql.Open("sqlite3", "config/database.db")

	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Id, Name, Settings FROM Indexers`
	indexers, _ := prowlarrDb.Query(prowlarrReq)
	defer indexers.Close()

	cols, _ := indexers.Columns()
	prowlarrIndexerConfig := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range ptrs {
		ptrs[i] = &prowlarrIndexerConfig[i]
	}

	// quite arbitrary value, can be changed if needed
	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for indexers.Next() {
		indexers.Scan(ptrs...)
		// Make a copy of prowlarrIndexerConfig to avoid race conditions
		prowlarrIndexerConfigCopy := make([]interface{}, len(prowlarrIndexerConfig))
		copy(prowlarrIndexerConfigCopy, prowlarrIndexerConfig)

		wg.Add(1)
		semaphore <- struct{}{}
		go func(prowlarrIndexerConfigCopy []interface{}) {
			defer wg.Done()
			defer func() { <-semaphore }()

			processIndexerProwlarr(prowlarrIndexerConfigCopy, db)
		}(prowlarrIndexerConfigCopy)
	}

	wg.Wait()

	return c.String(http.StatusOK, "Data inserted successfully!")
}

func processIndexerProwlarr(prowlarrIndexerConfig []interface{}, db *sql.DB) bool {
	indexerName := prowlarrIndexerConfig[1].(string)
	indexerName = strings.TrimSuffix(indexerName, " (API)")
	indexerName = strings.TrimSuffix(indexerName, "2FA")

	if helpers.GetIndexerInfo(indexerName).Get("enabled").Bool() {

		fmt.Printf("Updating %s's stats\n", indexerName)

		indexerStats, error := indexers.GetUserData(gjson.Parse(prowlarrIndexerConfig[2].(string)), indexerName, prowlarrIndexerConfig[0].(int64))
		if error != nil {
			return false
		}

		insertSQL := `INSERT INTO user_stats (
			indexer_id, uploaded_torrents, uploaded_amount, downloaded_amount, snatched, seeding, leeching,
			ratio, required_ratio, last_access, torrent_comments, invited, forum_posts, warned, class,
			donor, uploaded_rank, downloaded_rank, uploads_rank, requests_rank, bounty_rank, posts_rank,
			artists_rank, overall_rank, buffer, bonus_points, seeding_size, freeleech_tokens
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

		_, err := db.Exec(insertSQL,
			prowlarrIndexerConfig[0],
			indexerStats["uploaded_torrents"],
			indexerStats["uploaded_amount"],
			indexerStats["downloaded_amount"],
			indexerStats["snatched"],
			indexerStats["seeding"],
			indexerStats["leeching"],
			indexerStats["ratio"],
			indexerStats["required_ratio"],
			indexerStats["last_access"],
			indexerStats["torrent_comments"],
			indexerStats["invited"],
			indexerStats["forum_posts"],
			indexerStats["warned"],
			indexerStats["class"],
			indexerStats["donor"],
			indexerStats["uploaded_rank"],
			indexerStats["downloaded_rank"],
			indexerStats["uploads_rank"],
			indexerStats["requests_rank"],
			indexerStats["bounty_rank"],
			indexerStats["posts_rank"],
			indexerStats["artists_rank"],
			indexerStats["overall_rank"],
			indexerStats["buffer"],
			indexerStats["bonus_points"],
			indexerStats["seeding_size"],
			indexerStats["freeleech_tokens"],
		)

		if err != nil {
			fmt.Println("Failed to insert data:", err.Error())
		}
	}

	return true
}
