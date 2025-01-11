package trackers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

func GetUserData(prowlarrIndexerConfig gjson.Result, trackerName string, indexerId int64) (map[string]interface{}, error) {

	var results map[string]interface{}

	req := ConstructTrackerRequest(prowlarrIndexerConfig, trackerName, indexerId)
	if req.URL == nil {
		// fmt.Printf("Tracker %s unsupported\n", trackerName)
		return map[string]interface{}{}, errors.New("Tracker not supported")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		results = ProcessTrackerResponse(resp, trackerName)
		return results, nil
	} else {
		fmt.Println(resp.Status)
		fmt.Printf("Tracker %s did not reply with status 200 OK, skipping.", trackerName)
		// fmt.Println(resp)
		return map[string]interface{}{}, errors.New("An error occured when getting tracker's data")
	}
}
