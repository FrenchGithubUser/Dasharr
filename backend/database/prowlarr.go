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

	usernamePath := "username"
	passwordPath := "password"
	cookiePath := "cookie"

	// sometimes the username/password dont have the same key in the prowlarr db
	if !jsonResult.Get(usernamePath).Exists() {
		usernamePath = "extraFieldData." + usernamePath
		passwordPath = "extraFieldData." + passwordPath
		cookiePath = "extraFieldData." + cookiePath
	}

	credentials := map[string]string{
		"username": jsonResult.Get(usernamePath).String(),
		"password": jsonResult.Get(passwordPath).String(),
		"cookie":   jsonResult.Get(cookiePath).String(), // sometimes empty
	}

	return credentials
}
