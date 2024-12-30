package helpers

// func RowsToJSON(rows *sql.Rows) ([]map[string]interface{}, error) {
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var results []map[string]interface{}
// 	for rows.Next() {
// 		// Create a slice of `interface{}` to hold column values
// 		values := make([]interface{}, len(columns))
// 		valuePtrs := make([]interface{}, len(columns))
// 		for i := range values {
// 			valuePtrs[i] = &values[i]
// 		}

// 		// Scan the row into the value pointers
// 		err := rows.Scan(valuePtrs...)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Create a map for the row and populate it with column data
// 		rowMap := make(map[string]interface{})
// 		for i, col := range columns {
// 			val := values[i]
// 			if b, ok := val.([]byte); ok {
// 				rowMap[col] = string(b) // Convert byte slices to string
// 			} else {
// 				rowMap[col] = val
// 			}
// 		}
// 		results = append(results, rowMap)
// 	}
// 	return results, nil
// }
