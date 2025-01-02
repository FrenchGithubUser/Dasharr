package helpers

import (
	"fmt"
	"math"
)

func ConvertBitsToGiB(bits any) float64 {
	switch b := bits.(type) {
	case float64:
		return b / (8 * math.Pow(2, 30))
	case int64:
		return float64(b) / (8 * math.Pow(2, 30))
	default:
		return 0
	}
}

// takes a whole database query result and converts the relevant items
func ProcessDataTransferUnits(results interface{}) []map[string]interface{} {
	var processMap func(map[string]any) map[string]interface{}
	processMap = func(m map[string]any) map[string]interface{} {
		updated := make(map[string]any)
		for k, v := range m {
			switch v := v.(type) {
			case map[string]interface{}:
				updated[k] = processMap(v)
			case []map[string]interface{}:
				updated[k] = ProcessDataTransferUnits(v)
			default:
				if k == "uploaded" || k == "downloaded" {
					updated[k] = ConvertBitsToGiB(v)
				} else {
					updated[k] = v
				}
			}
		}
		return updated
	}

	if dataSlice, ok := results.([]map[string]any); ok {
		fmt.Println("ok")
		newData := make([]map[string]any, len(dataSlice))
		for i, item := range dataSlice {
			newData[i] = processMap(item)
		}
		return newData
	}
	// fmt.Println(results.([]map[string]interface{}))

	// Handle other types or return an error if unexpected
	return nil
}
