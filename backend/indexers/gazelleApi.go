package indexers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

func ConstructRequestGazelleApi(prowlarrIndexerConfig gjson.Result, indexerName string) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str
	var apiUrl string
	if indexerName == "GazelleGames" {
		apiUrl = baseUrl + "api.php?request="
	} else {
		apiUrl = baseUrl + "ajax.php?action="
	}

	var req *http.Request
	if indexerName == "BroadcasTheNet" {
		apiKey := prowlarrIndexerConfig.Get("apiKey").Str
		payload := map[string]interface{}{
			"method":  "userInfo",
			"params":  []string{apiKey},
			"jsonrpc": "2.0",
			"id":      1,
		}
		body, _ := json.Marshal(payload)
		req, _ = http.NewRequest("POST", baseUrl, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		apiKey := prowlarrIndexerConfig.Get("apikey").Str
		userId := getUserIdGazelleApi(apiUrl, apiKey, indexerName)
		req, _ = http.NewRequest("GET", apiUrl+"user&id="+strconv.Itoa(int(userId)), nil)
		if indexerName == "GazelleGames" {
			req.Header.Add("X-API-Key", apiKey)
		} else {
			req.Header.Add("Authorization", apiKey)
		}
	}
	return req
}

func ProcessIndexerResponseGazelleApi(results gjson.Result, indexerInfoJson gjson.Result) map[string]interface{} {
	if indexerInfoJson.Get("indexer_name").Str == "BroadcasTheNet" {
		results = results.Get("result")
	} else {
		results = results.Get("response")
	}
	mappedResults := make(map[string]interface{})
	indexerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})

	return mappedResults
}

func getUserIdGazelleApi(baseUrl string, apiKey string, indexerName string) int64 {
	var req *http.Request
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
		return gjson.Get(string(body), "response.id").Int()
	} else {
		return -1
	}
}
