package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ml-based-cache/internal/cache"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/utils"
)

type Results struct {
	Results []cache.Res
}

func main() {
	// for i := 0; i < 10; i++ {
	// 	for j := 0; j < 10; j++ {
	// 		for k := 0; k < 10; k++ {
	// 			random.GenerateBiasedRandom(500, 0, 60, float64(i*10), float64(j*10))
	// 		}
	// 	}
	// }
	// for k := 0; k < 100; k++ {
	// 	random.GenerateBiasedRandom(500, 0, 60, float64(0), float64(0))
	// }

	files, err := utils.ReadJSONArraysFromDir("./data")
	if err != nil {
		fmt.Println(err.Error())
	}

	r := Results{
		Results: make([]cache.Res, 0),
	}
	for _, v := range files {
		c := cache.NewCache(25)
		res := c.Simulate(generator.Converter(v), nil)
		r.Results = append(r.Results, res)
	}

	f, _ := os.Create("./data/dataset/testdata")

	data, _ := json.Marshal(&r.Results)

	fmt.Fprint(f, string(data))
}
