package database

import (
	"database/sql"

	"github.com/tidwall/gjson"
)

func GetProwlarrCredentials(indexerId interface{}) map[string]string {
	prowlarrDb, _ := sql.Open("sqlite3", "prowlarr/prowlarr.db")
	prowlarrReq := `SELECT Settings FROM Indexers WHERE Id = ?`
	indexers, _ := prowlarrDb.Query(prowlarrReq, indexerId)

	defer indexers.Close()

	var result string

	for indexers.Next() {
		var settings string
		indexers.Scan(&settings)
		result = settings
	}
	jsonResult := gjson.Parse(result)

	credentials := map[string]string{
		"username": jsonResult.Get("username").String(),
		"password": jsonResult.Get("password").String(),
	}

	return credentials
}
