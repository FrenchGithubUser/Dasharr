package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

// ExtractID extracts a number from a URL after "indexer/"
func GetUserData(trackerConfig gjson.Result) map[string]interface{} {

	trackerInfo, _ := os.ReadFile(fmt.Sprintf("config/trackers/%s.json", trackerConfig.Get("tracker_id")))
	trackerInfoJson := gjson.Parse(string(trackerInfo))
	authType := trackerConfig.Get("auth.type").Str
	var results gjson.Result
	if authType == "api_key" {
		req, _ := http.NewRequest("GET", trackerInfoJson.Get("base_url").Str, nil)
		trackerType := trackerInfoJson.Get("tracker_type").Str
		if trackerType == "gazelle" {
			req.Header.Add(trackerInfoJson.Get("auth_header").Str, trackerConfig.Get("fillable.api_key").Str)
			updatedUrl, _ := url.Parse(req.URL.String() + trackerConfig.Get("user_id").String())
			req.URL = updatedUrl
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		jsonResponse := gjson.Parse(string(body))
		results = jsonResponse.Get("response")
	}

	mappedResults := make(map[string]interface{})
	trackerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})

	return mappedResults
}
