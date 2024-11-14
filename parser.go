package jsonref

import (
	"encoding/json"
	"os"
)

func ParseJSON(data []byte, filename string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return parseReference(v, filename), nil
}

func parseReference(v interface{}, filename string) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		if ref, ok := v["$ref"].(string); ok {
			refJSON, err := os.ReadFile(ref)
			if err != nil {
				panic(err)
			}
			refData, err := ParseJSON(refJSON, ref)
			if err != nil {
				panic(err)
			}
			return refData
		}

		for key, value := range v {
			v[key] = parseReference(value, filename)
		}

	case []interface{}:
		for i, item := range v {
			v[i] = parseReference(item, filename)
		}
	}

	return v
}
