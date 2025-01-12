package helpers

import (
	"fmt"
	"math"
	"os"

	"github.com/tidwall/gjson"
)

func BytesToGiB(bits int64) float64 {
	return float64(bits) / (math.Pow(2, 30))
}

func AnyUnitToBytes(value float64, unit string) int64 {
	if unit == "GiB" {
		return int64(value * math.Pow(2, 30))
	} else if unit == "TiB" {
		return int64(value * math.Pow(2, 40))
	} else if unit == "B" {
		return int64(value)
	} else {
		return 0
	}
}

func RemoveNilEntries(data []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, m := range data {
		filtered := make(map[string]interface{})
		for k, v := range m {
			if v != nil {
				filtered[k] = v
			}
		}
		if len(filtered) > 0 {
			result = append(result, filtered)
		}
	}
	return result
}

func GetIndexerInfo(indexerName string) gjson.Result {

	indexersInfo, _ := os.ReadFile("config/config.json")
	indexerInfo := gjson.Get(string(indexersInfo), fmt.Sprintf(`#[indexer_name=="%s"]`, indexerName))
	return indexerInfo
	// result := gjson.Get(string(indexersInfo), "#")

	// result.ForEach(func(key, value gjson.Result) bool {
	// 	indexerName := value.Get("indexer_name").String()
	// 	return !strings.Contains(indexerName, indexerName)
	// })

	// return result
}

// takes a whole database query result and converts the relevant items
// func ProcessQueryResults(results interface{}) []map[string]interface{} {
// 	var processMap func(map[string]any) map[string]interface{}
// 	processMap = func(m map[string]any) map[string]interface{} {
// 		updated := make(map[string]any)
// 		for k, v := range m {
// 			switch v := v.(type) {
// 			case map[string]interface{}:
// 				updated[k] = processMap(v)
// 			case []map[string]interface{}:
// 				updated[k] = ProcessQueryResults(v)
// 			default:
// 				if k == "uploaded" || k == "downloaded" {
// 					updated[k] = BytesToGiB(v.(int64))
// 				} else {
// 					updated[k] = v
// 				}
// 			}
// 		}
// 		return updated
// 	}

// 	if dataSlice, ok := results.([]map[string]any); ok {
// 		fmt.Println("ok")
// 		newData := make([]map[string]any, len(dataSlice))
// 		for i, item := range dataSlice {
// 			newData[i] = processMap(item)
// 		}
// 		return newData
// 	}

// 	// Handle other types or return an error if unexpected
// 	return nil
// }
