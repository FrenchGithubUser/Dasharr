package indexers

import (
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
	} else if indexerType == "anthelion" {
		req = ConstructRequestGazelleScrape(prowlarrIndexerConfig, indexerName, indexerId)
	} else if indexerType == "MAM" {
		req = ConstructRequestMAM(prowlarrIndexerConfig)
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
	} else if indexerType == "anthelion" {
		results = ProcessIndexerResponseGazelleScrape(string(body), indexerInfo)
	} else if indexerType == "MAM" {
		results = ProcessIndexerResponseMAM(gjson.Parse(string(body)), indexerInfo)
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
	} else if contains(indexerName, []string{"Blutopia", "Aither", "ItaTorrents", "Oldtoons", "LST", "seedpool"}) {
		return "unit3d"
	} else if contains(indexerName, []string{"Anthelion"}) {
		return "anthelion"
	} else if contains(indexerName, []string{"MyAnonamouse"}) {
		return "MAM"
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
