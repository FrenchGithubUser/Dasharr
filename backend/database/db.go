package database

import (
	"backend/helpers"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tracker_id INTEGER,
    uploaded_torrents INTEGER DEFAULT 0,
    uploaded_amount INTEGER DEFAULT 0,
    downloaded_amount INTEGER DEFAULT 0,
    snatched INTEGER DEFAULT 0,
    seeding INTEGER DEFAULT 0,
    leeching INTEGER DEFAULT 0,
    ratio REAL DEFAULT 0.0,
    required_ratio REAL DEFAULT 0.0,
    last_access DATETIME DEFAULT NULL,
    torrent_comments INTEGER DEFAULT 0,
    invited INTEGER DEFAULT 0,
    forum_posts INTEGER DEFAULT 0,
    warned BOOLEAN DEFAULT 0,
    class TEXT DEFAULT NULL,
    donor BOOLEAN DEFAULT 0,
    uploaded_rank INTEGER DEFAULT 0,
    downloaded_rank INTEGER DEFAULT 0,
    uploads_rank INTEGER DEFAULT 0,
    requests_rank INTEGER DEFAULT 0,
    bounty_rank INTEGER DEFAULT 0,
    posts_rank INTEGER DEFAULT 0,
    artists_rank INTEGER DEFAULT 0,
    overall_rank INTEGER DEFAULT 0,
    collected_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	ExecuteQuery(createTableSQL, []interface{}{})

	log.Println("Database initialized successfully")
	return nil
}

func ExecuteQuery(query string, args []interface{}) []map[string]interface{} {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			if col == "uploaded_amount" || col == "downloaded_amount" {
				switch values[i].(type) {
				case int64:
					row[col] = helpers.BytesToGiB(values[i].(int64))
				case float64:
					row[col] = helpers.BytesToGiB(int64(values[i].(float64)))
				}
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	return results
}
