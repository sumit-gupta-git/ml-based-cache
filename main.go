package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ml-based-cache/internal/cache"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/generator/random"
	"ml-based-cache/internal/utils"
)

type Results struct {
	Results []cache.Res
}

func main() {
	Execute()
}

func Training() {
	// LFU Bias
	for i := 0; i < 10000; i++ {
		random.GenerateSplitBiasedRandom(10000, 0, 100, 5, 200)
	}

	// LRU Bias
	for i := 0; i <= 9; i++ {
		for j := 0; j < 1000; j++ {
			random.GenerateRecencyBias(10000, 0, 100, float64(i*10))
		}
	}

	// Random
	for k := 0; k < 10000; k++ {
		random.GenerateRandomArray(10000, 0, 100)
	}

	files, err := utils.ReadJSONArraysFromDir("./data")
	if err != nil {
		fmt.Println(err.Error())
	}

	r := Results{
		Results: make([]cache.Res, 0),
	}
	for _, v := range files {
		c := cache.NewCache(15)
		res := c.SimulateTraining(generator.Converter(v))
		r.Results = append(r.Results, res)
	}

	f, _ := os.Create("./py/test_data2.json")

	data, _ := json.Marshal(&r.Results)

	fmt.Fprint(f, string(data))
}
