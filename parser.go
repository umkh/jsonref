package main

import (
	"encoding/json"
	"os"
	// "path/filepath"
)
/*
type RefResolver struct {
	baseDir string
	data    map[string]any
}

func (r *RefResolver) Decode(data []byte) (any, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return r.resolveRefs(m), nil
}

func (r *RefResolver) resolveRefs(m map[string]any) any {
	for k, v := range m {
		if ref, ok := v.(map[string]any); ok && ref["$ref"] != nil {
			refPath := filepath.Join(r.baseDir, ref["$ref"].(string))
			refData, err := os.ReadFile(refPath)
			if err != nil {
				panic(err)
			}

			refResolver := &RefResolver{
				baseDir: filepath.Dir(refPath),
				data:    make(map[string]any),
			}

			m[k], _ = refResolver.Decode(refData)
		} else if nested, ok := v.(map[string]any); ok {
			m[k] = r.resolveRefs(nested)
		} else if arr, ok := v.([]any); ok {
			m[k] = r.resolveArrayRefs(arr)
		}
	}

	return m
}

func (r *RefResolver) resolveArrayRefs(arr []any) []any {
	fmt.Println(arr)
	for i, v := range arr {
		if ref, ok := v.(map[string]any); ok && ref["$ref"] != nil {
			refPath := filepath.Join(r.baseDir, ref["$ref"].(string))
			refData, err := os.ReadFile(refPath)
			if err != nil {
				panic(err)
			}

			refResolver := &RefResolver{
				baseDir: filepath.Dir(refPath),
				data:    make(map[string]any),
			}
			arr[i], _ = refResolver.Decode(refData)
		} else if nested, ok := v.(map[string]any); ok {
			arr[i] = r.resolveRefs(nested)
		}
	}
	return arr
}

func main() {
	// Read the input JSON file
	inputData, err := os.ReadFile("./config/krakend.json")
	if err != nil {
		panic(err)
	}

	// Create a RefResolver and decode the input JSON
	refResolver := &RefResolver{
		baseDir: ".",
		data:    make(map[string]any),
	}
	result, err := refResolver.Decode(inputData)
	if err != nil {
		panic(err)
	}

	// Marshal the resolved JSON data to a new file
	outputData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	os.WriteFile("output.json", outputData, 0644)

	fmt.Println("JSON file with $ref pointers resolved saved as output.json")
}
*/ 


func main() {
	// Read the main JSON file
	mainJSON, err := os.ReadFile("./config/krakend.json")
	if err != nil {
		panic(err)
	}

	// Parse the main JSON file
	config, err := parseJSON(mainJSON, ".")
	if err != nil {
		panic(err)
	}

	// Marshal the final result
	result, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("output.json", result, 0644)
}

func parseJSON(data []byte, filename string) (interface{}, error) {
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
			refData, err := parseJSON(refJSON, ref)
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

