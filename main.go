package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

func mergeYAMLFiles(filePaths []string) (map[interface{}]interface{}, error) {
	mergedData := make(map[interface{}]interface{})

	for _, filePath := range filePaths {
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var yamlData map[interface{}]interface{}
		err = yaml.Unmarshal(fileData, &yamlData)
		if err != nil {
			return nil, err
		}

		// Merge YAML data
		mergedData = mergeMaps(mergedData, yamlData)
	}

	return mergedData, nil
}

func mergeMaps(dest, src map[interface{}]interface{}) map[interface{}]interface{} {
	for key, srcValue := range src {
		if destValue, ok := dest[key]; ok {
			// If the key exists in both maps, merge recursively
			if srcMap, ok := srcValue.(map[interface{}]interface{}); ok {
				if destMap, ok := destValue.(map[interface{}]interface{}); ok {
					dest[key] = mergeMaps(destMap, srcMap)
				} else {
					// Conflict: key exists but values are not both maps
					fmt.Printf("Conflict for key %v. Values are not both maps.\n", key)
				}
			} else {
				// Convert values to slices and merge them
				dest[key] = mergeSlices(destValue, srcValue)
			}
		} else {
			// Key doesn't exist in dest map, add it
			dest[key] = srcValue
		}
	}

	return dest
}

func mergeSlices(destValue, srcValue interface{}) interface{} {
	// Convert values to slices and merge them
	destSlice, destIsSlice := convertToSlice(destValue)
	srcSlice, srcIsSlice := convertToSlice(srcValue)

	if destIsSlice && srcIsSlice {
		return append(destSlice, srcSlice...)
	} else if srcIsSlice {
		return srcSlice
	} else {
		return destValue
	}
}

func convertToSlice(value interface{}) ([]interface{}, bool) {
	// Convert a value to a slice, if possible
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Slice {
		slice := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Interface()
		}
		return slice, true
	}
	return nil, false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go file1.yaml file2.yaml ...")
		os.Exit(1)
	}

	filePaths := os.Args[1:]
	mergedData, err := mergeYAMLFiles(filePaths)
	if err != nil {
		fmt.Printf("Error merging YAML files: %v\n", err)
		os.Exit(1)
	}

	mergedYAML, err := yaml.Marshal(mergedData)
	if err != nil {
		fmt.Printf("Error marshaling merged YAML data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(mergedYAML))
}
