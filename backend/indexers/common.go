package indexers

import (
	"backend/database"
	"backend/helpers"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func ConstructIndexerRequest(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) *http.Request {
	req := &http.Request{}
	indexerType := DetermineIndexerType(indexerName)

	if indexerType == "gazelleApi" {
		req = ConstructRequestGazelleApi(prowlarrIndexerConfig, indexerName)
	} else if indexerType == "unit3d" {
		req = ConstructRequestUnit3d(indexerName, indexerId)
	} else if indexerType == "gazelleScrape" {
		req = ConstructRequestGazelleScrape(prowlarrIndexerConfig, indexerName, indexerId)
	} else if indexerType == "MAM" {
		req = ConstructRequestMAM(prowlarrIndexerConfig)
	} else if indexerType == "TL" {
		req = ConstructRequestTL(prowlarrIndexerConfig, indexerName, indexerId)
	}

	return req
}

// todo : cookie refresh
func ProcessIndexerResponse(response *http.Response, indexerName string) (map[string]interface{}, error) {
	indexerInfo := helpers.GetIndexerInfo(indexerName)
	body, _ := io.ReadAll(response.Body)
	results := map[string]interface{}{}
	indexerType := DetermineIndexerType(indexerName)

	if indexerType == "gazelleApi" {
		results = ProcessIndexerResponseGazelleApi(gjson.Parse(string(body)), indexerInfo)
	} else if indexerType == "unit3d" {
		results = ProcessIndexerResponseUnit3d(string(body), indexerInfo)
	} else if indexerType == "gazelleScrape" {
		results = ProcessIndexerResponseGazelleScrape(string(body), indexerInfo)
	} else if indexerType == "MAM" {
		results = ProcessIndexerResponseMAM(gjson.Parse(string(body)), indexerInfo)
	} else if indexerType == "TL" {
		results = ProcessIndexerResponseTL(string(body), indexerInfo)
	}

	var err error = nil
	if results == nil || len(results) == 0 {
		err = errors.New(fmt.Sprintf("An error occured wile parsing %s's results", indexerName))
	}

	return results, err
}

func DetermineIndexerType(indexerName string) string {
	contains := func(s string, list []string) bool {
		for _, v := range list {
			if v == s {
				return true
			}
		}
		return false
	}

	if contains(indexerName, []string{"Orpheus", "Redacted", "GazelleGames", "BroadcasTheNet"}) {
		return "gazelleApi"
	} else if contains(indexerName, []string{"Blutopia", "Aither", "ItaTorrents", "Oldtoons", "LST", "seedpool", "FearNoPeer"}) {
		return "unit3d"
	} else if contains(indexerName, []string{"Anthelion", "AlphaRatio"}) {
		return "gazelleScrape"
	} else if contains(indexerName, []string{"MyAnonamouse"}) {
		return "MAM"
	} else if contains(indexerName, []string{"TorrentLeech"}) {
		return "TL"
	}
	return "unknown"
}

func addCookiesToRequest(request *http.Request, cookieStr string) *http.Request {
	cookies := strings.Split(cookieStr, ";")
	for _, cookie := range cookies {
		kv := strings.SplitN(cookie, "=", 2)
		if len(kv) == 2 {
			request.AddCookie(&http.Cookie{
				Name:  kv[0],
				Value: kv[1],
			})
		}
	}
	return request

}

func LoginAndSaveCookies(indexer string, username string, password string, twoFaCode string, apiKey string, indexerId interface{}) error {
	indexerInfo := helpers.GetIndexerInfo(indexer)
	if indexerInfo.Get("credentials.method").Str == "prowlarr" {
		credentials := database.GetProwlarrCredentials(indexerId)
		username = credentials["username"]
		password = credentials["password"]
	} else if username == "" {
		// in case of cookie refresh, the credentials are not given by the user once again
		username = database.GetIndexerUsername(indexerId)
		password = database.GetIndexerPassword(indexerId)
	}
	if !indexerInfo.Exists() {
		return errors.New("indexer not found in config")
	}
	indexerType := DetermineIndexerType(indexer)

	loginURL := indexerInfo.Get("login.url").String()

	var cookies string
	if indexerType == "unit3d" {
		cookies = LoginAndGetCookiesUnit3d(username, password, twoFaCode, loginURL, indexerInfo.Get("domain").Str)
	} else if indexerType == "gazelleScrape" {
		cookies = LoginAndGetCookiesGazelleScrape(username, password, twoFaCode, loginURL, indexerInfo)
	} else if indexerType == "TL" {
		cookies = LoginAndGetCookiesTL(username, password, twoFaCode, loginURL, indexerInfo)
	}

	if cookies != "" {
		insertSQL := `INSERT OR REPLACE INTO credentials (
		indexer_id, username, password, cookies, api_key
		) VALUES (?, ?, ?, ?, ?);`
		args := []interface{}{indexerId, username, password, cookies, apiKey}
		database.ExecuteQuery(insertSQL, args)
		// return c.JSON(http.StatusOK, map[string]string{"status": "success"})
		return nil
	} else {
		// return c.JSON(http.StatusUnauthorized, map[string]string{"error": "login_failed"})
	}

	return errors.New("an error occured getting cookies for indexer " + indexer)

}
