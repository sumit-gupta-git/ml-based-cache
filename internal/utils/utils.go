package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Average(nums []int64) float64 {
	if len(nums) == 0 {
		return 0
	}

	sum := int64(0)
	for _, v := range nums {
		sum += v
	}
	return float64(sum) / float64(len(nums))
}

func ReadJSONArraysFromDir(dir string) ([][]int, error) {
	var result [][]int

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("error reading file %s: %v\n", path, err)
			continue
		}

		var arr []int
		if err := json.Unmarshal(data, &arr); err != nil {
			fmt.Printf("error parsing %s: %v\n", path, err)
			fmt.Println("Raw data:", string(data))
			continue
		}

		result = append(result, arr)
	}

	return result, nil
}
