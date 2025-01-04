package trackers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func ConstructTrackerRequest(trackerConfig gjson.Result, trackerName string) *http.Request {
	req := &http.Request{}
	trackerType := determineTrackerType(trackerName)
	if trackerType == "gazelle" {
		req = ConstructRequestGazelle(trackerConfig, trackerName)
	} else if trackerType == "unit3d" {
		req = ConstructRequestUnit3d(trackerConfig, trackerName)
	}

	return req
}

func ProcessTrackerResponse(response *http.Response, trackerConfig gjson.Result, trackerName string) map[string]interface{} {
	trackerInfo, _ := os.ReadFile(fmt.Sprintf("config/trackers/%s.json", trackerName))
	trackerInfoJson := gjson.Parse(string(trackerInfo))
	body, _ := io.ReadAll(response.Body)
	results := gjson.Parse(string(body))
	trackerType := determineTrackerType(trackerName)

	if trackerType == "gazelle" {
		results = ProcessTrackerResponseGazelle(results)
	} else if trackerType == "unit3d" {
		results = ProcessTrackerResponseUnit3d(results, string(body))
	}

	mappedResults := make(map[string]interface{})
	trackerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})
	return mappedResults
}

func determineTrackerType(trackerName string) string {
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
	} else if contains(trackerName, []string{"Blutopia (API)", "Aither (API)", "ItaTorrents"}) {
		return "unit3d"
	}
	return "unknown"
}
