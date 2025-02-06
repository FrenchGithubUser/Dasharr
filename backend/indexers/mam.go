package indexers

import (
	"net/http"

	"github.com/tidwall/gjson"
)

func ConstructRequestMAM(prowlarrIndexerConfig gjson.Result) *http.Request {
	baseUrl := prowlarrIndexerConfig.Get("baseUrl").Str

	var req *http.Request
	apiKey := prowlarrIndexerConfig.Get("mamId").Str
	req, _ = http.NewRequest("GET", baseUrl+"jsonLoad.php", nil)
	req.Header.Set("Cookie", "mam_id="+apiKey)

	return req
}

func ProcessIndexerResponseMAM(results gjson.Result, indexerInfoJson gjson.Result) map[string]interface{} {
	mappedResults := make(map[string]interface{})
	indexerInfoJson.Get("stats_keys").ForEach(func(key, value gjson.Result) bool {
		mappedResults[value.String()] = results.Get(key.String()).Value()
		return true
	})

	return mappedResults
}
