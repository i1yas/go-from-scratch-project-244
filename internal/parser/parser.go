package parser

import "encoding/json"

func Parse(fileContent []byte) (map[string]any, error) {
	var result map[string]any

	err := json.Unmarshal(fileContent, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
