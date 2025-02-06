package indexers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

func GetUserData(prowlarrIndexerConfig gjson.Result, indexerName string, indexerId int64) (map[string]interface{}, error) {

	// var results map[string]interface{}

	req := ConstructIndexerRequest(prowlarrIndexerConfig, indexerName, indexerId)
	if req.URL == nil {
		// fmt.Printf("Indexer %s unsupported\n", indexerName)
		return map[string]interface{}{}, errors.New("Indexer not supported")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		results, processErr := ProcessIndexerResponse(resp, indexerName)
		return results, processErr
	} else {

		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Printf("Indexer %s did not reply with status 200 OK, skipping.", indexerName)
		// fmt.Println(resp)
		return map[string]interface{}{}, errors.New("An error occured when getting indexer's data")
	}
}
