package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

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
			updatedUrl, _ := url.Parse(req.URL.String() + trackerConfig.Get("fillable.user_id").String())
			req.URL = updatedUrl
		} else if trackerType == "unit3d" {
			req.Header.Add("Authorization", "Bearer "+trackerConfig.Get("fillable.api_key").Str)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		results = gjson.Parse(string(body))
		if trackerType == "gazelle" {
			results = results.Get("response")
		} else if trackerType == "unit3d" {
			re := regexp.MustCompile(`^([\d\.]+)\s?(GiB|MiB|TiB)$`)

			uploadRegexResult := re.FindStringSubmatch(results.Get("uploaded").Str)
			cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
			edited_results, _ := sjson.Set(string(body), "uploaded", AnyUnitToBits(cleanUpload, uploadRegexResult[2]))
			downloadRegexResult := re.FindStringSubmatch(results.Get("downloaded").Str)
			cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
			edited_results, _ = sjson.Set(edited_results, "downloaded", AnyUnitToBits(cleanDownload, downloadRegexResult[2]))
			bufferRegexResult := re.FindStringSubmatch(results.Get("buffer").Str)
			cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
			edited_results, _ = sjson.Set(edited_results, "buffer", AnyUnitToBits(cleanBuffer, downloadRegexResult[2]))

			results = gjson.Parse(edited_results)
		}
	}

	mappedResults := make(map[string]interface{})
	trackerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})

	return mappedResults
}
