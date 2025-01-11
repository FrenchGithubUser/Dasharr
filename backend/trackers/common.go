package trackers

import (
	"backend/helpers"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func ConstructTrackerRequest(prowlarrIndexerConfig gjson.Result, trackerName string, indexerId int64) *http.Request {
	req := &http.Request{}
	trackerType := DetermineTrackerType(trackerName)

	if trackerType == "gazelle" {
		req = ConstructRequestGazelle(prowlarrIndexerConfig, trackerName)
	} else if trackerType == "unit3d" {
		// LoginAndGetCookiesUnit3d(username, password, loginURL, indexerInfo.Get("domain").Str)
		req = ConstructRequestUnit3d(trackerName, indexerId)
	} else if trackerType == "anthelion" {
		// LoginAndGetCookiesAnthelion(username, password, loginURL, indexerInfo)
		req = ConstructRequestAnthelion(prowlarrIndexerConfig, trackerName, indexerId)
	}

	return req
}

// todo : cookie refresh
func ProcessTrackerResponse(response *http.Response, trackerName string) map[string]interface{} {
	indexerInfo := helpers.GetIndexerInfo(trackerName)
	body, _ := io.ReadAll(response.Body)
	results := map[string]interface{}{}
	trackerType := DetermineTrackerType(trackerName)

	if trackerType == "gazelle" {
		results = ProcessTrackerResponseGazelle(gjson.Parse(string(body)), indexerInfo)
	} else if trackerType == "unit3d" {
		results = ProcessTrackerResponseUnit3d(string(body), indexerInfo)
	} else if trackerType == "anthelion" {
		results = ProcessTrackerResponseAnthelion(string(body), indexerInfo)
	}

	return results
}

func DetermineTrackerType(trackerName string) string {
	contains := func(s string, list []string) bool {
		for _, v := range list {
			if v == s {
				return true
			}
		}
		return false
	}

	if contains(trackerName, []string{"Orpheus", "Redacted", "GazelleGames"}) {
		return "gazelle"
	} else if contains(trackerName, []string{"Blutopia", "Aither", "ItaTorrents", "Oldtoons"}) {
		return "unit3d"
	} else if contains(trackerName, []string{"Anthelion"}) {
		return "anthelion"
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
