package trackers

import (
	"backend/helpers"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func ConstructRequestUnit3d(trackerConfig gjson.Result, trackerName string) *http.Request {
	configFile, _ := os.ReadFile(fmt.Sprintf("config/trackers/%s.json", trackerName))
	configFileJson := gjson.Parse(string(configFile))
	baseUrl := configFileJson.Get("base_url").Str + "user"
	apiKey := trackerConfig.Get("extraFieldData.apikey").Str
	req, _ := http.NewRequest("GET", baseUrl, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	// fmt.Println(req)
	return req
}

func ProcessTrackerResponseUnit3d(results gjson.Result, bodyString string) gjson.Result {
	re := regexp.MustCompile(`^([\d\.]+)\s?(GiB|MiB|TiB)$`)

	uploadRegexResult := re.FindStringSubmatch(results.Get("uploaded").Str)
	cleanUpload, _ := strconv.ParseFloat(uploadRegexResult[1], 64)
	edited_results, _ := sjson.Set(bodyString, "uploaded", helpers.AnyUnitToBytes(cleanUpload, uploadRegexResult[2]))
	downloadRegexResult := re.FindStringSubmatch(results.Get("downloaded").Str)
	cleanDownload, _ := strconv.ParseFloat(downloadRegexResult[1], 64)
	edited_results, _ = sjson.Set(edited_results, "downloaded", helpers.AnyUnitToBytes(cleanDownload, downloadRegexResult[2]))
	bufferRegexResult := re.FindStringSubmatch(results.Get("buffer").Str)
	cleanBuffer, _ := strconv.ParseFloat(bufferRegexResult[1], 64)
	edited_results, _ = sjson.Set(edited_results, "buffer", helpers.AnyUnitToBytes(cleanBuffer, downloadRegexResult[2]))

	return gjson.Parse(edited_results)
}
