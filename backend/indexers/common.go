package indexers

import (
	"backend/helpers"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func ConstructIndexerRequest(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) *http.Request {
	req := &http.Request{}
	indexerType := DetermineIndexerType(indexerName)

	if indexerType == "gazelle" {
		req = ConstructRequestGazelle(prowlarrIndexerConfig, indexerName)
	} else if indexerType == "unit3d" {
		// LoginAndGetCookiesUnit3d(username, password, loginURL, indexerInfo.Get("domain").Str)
		req = ConstructRequestUnit3d(indexerName, indexerId)
	} else if indexerType == "anthelion" {
		// LoginAndGetCookiesAnthelion(username, password, loginURL, indexerInfo)
		req = ConstructRequestAnthelion(prowlarrIndexerConfig, indexerName, indexerId)
	}

	return req
}

// todo : cookie refresh
func ProcessIndexerResponse(response *http.Response, indexerName string) map[string]interface{} {
	indexerInfo := helpers.GetIndexerInfo(indexerName)
	body, _ := io.ReadAll(response.Body)
	results := map[string]interface{}{}
	indexerType := DetermineIndexerType(indexerName)

	if indexerType == "gazelle" {
		results = ProcessIndexerResponseGazelle(gjson.Parse(string(body)), indexerInfo)
	} else if indexerType == "unit3d" {
		results = ProcessIndexerResponseUnit3d(string(body), indexerInfo)
	} else if indexerType == "anthelion" {
		results = ProcessIndexerResponseAnthelion(string(body), indexerInfo)
	}

	return results
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

	if contains(indexerName, []string{"Orpheus", "Redacted", "GazelleGames"}) {
		return "gazelle"
	} else if contains(indexerName, []string{"Blutopia", "Aither", "ItaTorrents", "Oldtoons", "LST"}) {
		return "unit3d"
	} else if contains(indexerName, []string{"Anthelion"}) {
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
