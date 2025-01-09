package trackers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func ConstructTrackerRequest(trackerConfig gjson.Result, trackerName string, indexerId int64) *http.Request {
	req := &http.Request{}
	trackerType := DetermineTrackerType(trackerName)
	if trackerType == "gazelle" {
		req = ConstructRequestGazelle(trackerConfig, trackerName)
	} else if trackerType == "unit3d" {
		req = ConstructRequestUnit3d(trackerConfig, trackerName, indexerId)
	}

	return req
}

func ProcessTrackerResponse(response *http.Response, trackerConfig gjson.Result, trackerName string) map[string]interface{} {
	trackerInfo, _ := os.ReadFile(fmt.Sprintf("config/trackers/%s.json", trackerName))
	trackerInfoJson := gjson.Parse(string(trackerInfo))
	body, _ := io.ReadAll(response.Body)
	results := map[string]interface{}{}
	trackerType := DetermineTrackerType(trackerName)

	if trackerType == "gazelle" {
		results = ProcessTrackerResponseGazelle(gjson.Parse(string(body)), trackerInfoJson)
	} else if trackerType == "unit3d" {
		results = ProcessTrackerResponseUnit3d(string(body), trackerInfoJson)
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
	} else if contains(trackerName, []string{"Blutopia", "Aither", "ItaTorrents"}) {
		return "unit3d"
	}
	return "unknown"
}
