package json2excel

import "encoding/json"

type request struct {
	// The data field only accepts objects and arrays of objects as input.
	// In cases where the object or array is empty, it will not be possible to infer column names,
	// resulting in an empty Excel document.
	// Each JSON object must be consistent and contain the same number of fields,
	// ensuring data types are consistent across them.
	Data exceldata `json:"data"`
}

func NewRequest(jsonData []byte) (*request, error) {
	var data exceldata
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	return &request{Data: data}, nil
}
