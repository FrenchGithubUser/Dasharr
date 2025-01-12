package indexers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

func ConstructRequestGazelle(prowlarrIndexerConfig gjson.Result, indexerName string) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str
	if indexerName == "GazelleGames" {
		baseUrl += "api.php?request="
	} else {
		baseUrl += "ajax.php?action="
	}
	apiKey := prowlarrIndexerConfig.Get("apikey").Str
	userId := getUserIdGazelle(baseUrl, apiKey, indexerName)
	req, _ := http.NewRequest("GET", baseUrl+"user&id="+strconv.Itoa(int(userId)), nil)
	if indexerName == "GazelleGames" {
		req.Header.Add("X-API-Key", apiKey)
	} else {
		req.Header.Add("Authorization", apiKey)
	}
	return req
}

func ProcessIndexerResponseGazelle(results gjson.Result, indexerInfoJson gjson.Result) map[string]interface{} {

	results = results.Get("response")
	mappedResults := make(map[string]interface{})
	indexerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})
	// fmt.Println(mappedResults)

	return mappedResults
}

func getUserIdGazelle(baseUrl string, apiKey string, indexerName string) int64 {
	req, _ := http.NewRequest("", "", nil)
	if indexerName == "GazelleGames" {
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
