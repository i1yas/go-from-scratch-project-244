package parser

import "encoding/json"

// Parse parses file content
func Parse(fileContent []byte) (map[string]any, error) {
	var result map[string]any

	err := json.Unmarshal(fileContent, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
