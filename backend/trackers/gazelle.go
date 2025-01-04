package trackers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

func ConstructRequestGazelle(trackerConfig gjson.Result, trackerName string) *http.Request {
	baseUrl := trackerConfig.Get("baseUrl").Str
	if trackerName == "GazelleGames" {
		baseUrl += "api.php?request="
	} else {
		baseUrl += "ajax.php?action="
	}
	apiKey := trackerConfig.Get("apikey").Str
	userId := getUserId(baseUrl, apiKey, trackerName)
	req, _ := http.NewRequest("GET", baseUrl+"user&id="+strconv.Itoa(int(userId)), nil)
	if trackerName == "GazelleGames" {
		req.Header.Add("X-API-Key", apiKey)
	} else {
		req.Header.Add("Authorization", apiKey)
	}
	return req
}

func ProcessTrackerResponseGazelle(results gjson.Result) gjson.Result {
	return results.Get("response")
}

func getUserId(baseUrl string, apiKey string, trackerName string) int64 {
	req, _ := http.NewRequest("", "", nil)
	if trackerName == "GazelleGames" {
		req, _ = http.NewRequest("GET", baseUrl+"quick_user", nil)
		req.Header.Add("X-API-Key", apiKey)
	} else {
		req, _ = http.NewRequest("GET", baseUrl+"index", nil)
		req.Header.Add("Authorization", apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.Status == "200 OK" {
		// fmt.Println(string(body))
		return gjson.Get(string(body), "response.id").Int()
	} else {
		return -1
	}
}
