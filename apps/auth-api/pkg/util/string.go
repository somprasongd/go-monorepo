package util

import "encoding/json"

func StringToMap(encoded string, dest *map[string]any) error {
	return json.Unmarshal([]byte(encoded), dest)
}

func MapToString(source map[string]any) (string, error) {
	jsonStr, err := json.Marshal(source)

	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}

func StructToString(source interface{}) (string, error) {
	jsonStr, err := json.Marshal(source)

	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}
