package trackers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

func ConstructRequestGazelle(trackerConfig gjson.Result, trackerName string) *http.Request {
	baseUrl := trackerConfig.Get("baseUrl").Str + "ajax.php?action="
	apiKey := trackerConfig.Get("apikey").Str
	userId := getUserId(baseUrl, apiKey)
	req, _ := http.NewRequest("GET", baseUrl+"user&id="+strconv.Itoa(int(userId)), nil)
	req.Header.Add("Authorization", apiKey)
	fmt.Println(req)
	return req
}

func getUserId(baseUrl string, apiKey string) int64 {
	req, _ := http.NewRequest("GET", baseUrl+"index", nil)
	req.Header.Add("Authorization", apiKey)

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
