package helpers

import (
	"backend/trackers"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func ConstructTrackerRequest(trackerConfig gjson.Result, trackerName string) *http.Request {
	req := &http.Request{}
	if determineTrackerType(trackerName) == "gazelle" {
		req = trackers.ConstructRequestGazelle(trackerConfig, trackerName)
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
		results = results.Get("response")
	} else if trackerType == "unit3d" {
		re := regexp.MustCompile(`^([\d\.]+)\s?(GiB|MiB|TiB)$`)

		uploadRegexResult := re.FindStringSubmatch(results.Get("uploaded").Str)
		cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
		edited_results, _ := sjson.Set(string(body), "uploaded", AnyUnitToBytes(cleanUpload, uploadRegexResult[2]))
		downloadRegexResult := re.FindStringSubmatch(results.Get("downloaded").Str)
		cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
		edited_results, _ = sjson.Set(edited_results, "downloaded", AnyUnitToBytes(cleanDownload, downloadRegexResult[2]))
		bufferRegexResult := re.FindStringSubmatch(results.Get("buffer").Str)
		cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
		edited_results, _ = sjson.Set(edited_results, "buffer", AnyUnitToBytes(cleanBuffer, downloadRegexResult[2]))

		results = gjson.Parse(edited_results)
	}

	mappedResults := make(map[string]interface{})
	trackerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})
	return mappedResults
}

func determineTrackerType(trackerName string) string {
	if func(s string, list []string) bool {
		for _, v := range list {
			if v == s {
				return true
			}
		}
		return false
	}(trackerName, []string{"Orpheus", "Redacted"}) {
		return "gazelle"
	}
	return "unkown"
}
