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
	// LFU Bias
	for i := 0; i < 10; i++ {
		random.GenerateSplitBiasedRandom(10000, 0, 100, 5, 200)
	}

	// // FIFO Bias
	// for i := 0; i < 10; i++ {
	// 	for j := 0; j < 900; j++ {
	// 		random.GenerateRepeatBias(500, 0, 100, float64(i*10))
	// 	}
	// }
	//
	// // LRU Bias
	// for i := 0; i <= 9; i++ {
	// 	for j := 0; j < 100; j++ {
	// 		random.GenerateRecencyBias(10000, 0, 100, float64(i*10))
	// 	}
	// }
	//
	// // Random
	// for k := 0; k < 1000; k++ {
	// 	random.GenerateRandomArray(10000, 0, 100)
	// }

	files, err := utils.ReadJSONArraysFromDir("./data")
	if err != nil {
		fmt.Println(err.Error())
	}

	r := Results{
		Results: make([]cache.Res, 0),
	}
	for _, v := range files {
		c := cache.NewCache(15)
		res := c.Simulate(generator.Converter(v), nil)
		r.Results = append(r.Results, res)
	}

	f, _ := os.Create("./py/train_data.json")

	data, _ := json.Marshal(&r.Results)

	fmt.Fprint(f, string(data))
}
