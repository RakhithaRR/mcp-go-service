package mcp

import (
	"encoding/json"
)

func processJsonResponse(inputString string) (string, error) {
	var data map[string]any

	err := json.Unmarshal([]byte(inputString), &data)
	if err != nil {
		logger.Error("Failed to unmarshal JSON", "error", err)
		return "", err
	}

	compactJSONBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error marshalling data back to JSON", "error", err)
		return "", err
	}

	return string(compactJSONBytes), nil
}
