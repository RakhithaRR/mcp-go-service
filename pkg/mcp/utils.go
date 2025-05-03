package mcp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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

func mapToXMLElements(m map[string]any) []XMLElement {
	elements := []XMLElement{}
	for k, v := range m {
		elem := XMLElement{XMLName: xml.Name{Local: k}}

		switch val := v.(type) {
		case string:
			elem.Content = val
		case map[string]any:
			elem.Children = mapToXMLElements(val)
		case []any:
			for _, item := range val {
				if itemMap, ok := item.(map[string]any); ok {
					elem.Children = append(elem.Children, mapToXMLElements(itemMap)...)
				} else {
					elem.Children = append(elem.Children, XMLElement{XMLName: xml.Name{Local: "item"}, Content: fmt.Sprint(item)})
				}
			}
		default:
			elem.Content = fmt.Sprint(val)
		}

		elements = append(elements, elem)
	}
	return elements
}
