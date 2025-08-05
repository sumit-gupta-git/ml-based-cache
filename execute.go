package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ml-based-cache/internal/cache"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/generator/random"
	"ml-based-cache/internal/models"
	"ml-based-cache/internal/utils"
)

type AlgoType int

const (
	NA = iota
	LRU
	LFU
)

type Prediction struct {
	Res int `json:"prediction_encoded"`
}

type Result struct {
	Hits          int
	Miss          int
	Throughput    int
	ResponseTimes int
}

func Execute() {
	arr := RandomCache(10_000, 200, 10)
	c := cache.NewCache(15)

	result := Simulate(c, arr)

	data, _ := json.Marshal(result)

	fmt.Println("\n\nRESULTS\n\n")
	fmt.Println(string(data))
}

func Simulate(globalCache *cache.Cache, arr *[]models.CacheItem) Result {
	itemsBatch := make([]models.CacheItem, 0, 1000)
	var currPolicy AlgoType = LRU
	currCache := cache.NewCache(15)
	processed := 0

	for _, v := range *arr {
		if currPolicy == LRU {
			currCache.AccessLRU(&v)
		} else if currPolicy == LFU {
			currCache.AccessLFU(&v)
		}

		itemsBatch = append(itemsBatch, v)
		processed++

		// After every 1000 items, query model
		if processed%1000 == 0 {
			batchRes := cache.Res{Items: generator.ReConverter(&itemsBatch)}
			if currPolicy == LRU {
				batchRes = cache.Res{
					Items:          generator.ReConverter(&itemsBatch),
					LRUHits:        currCache.Hits,
					LRUMiss:        currCache.Misses,
					LRUAvgReaccess: utils.Average(currCache.AvgAccessTime),
				}
			} else if currPolicy == LFU {
				batchRes = cache.Res{
					Items:          generator.ReConverter(&itemsBatch),
					LFUHits:        currCache.Hits,
					LFUMiss:        currCache.Misses,
					LFUAvgReaccess: utils.Average(currCache.AvgAccessTime),
				}
			}

			data, _ := json.Marshal(batchRes)

			currPolicy = QueryModel(string(data), currPolicy)
			fmt.Printf("Using %d\n", currPolicy)

			// Update global cache stats
			globalCache.Hits += currCache.Hits
			globalCache.Misses += currCache.Misses
			globalCache.AvgAccessTime = append(globalCache.AvgAccessTime, currCache.AvgAccessTime...)

			currCache.Misses = 0
			currCache.Hits = 0
			currCache.AvgAccessTime = make([]int64, 0)

			itemsBatch = itemsBatch[:0]
		}
	}

	return Result{
		Hits: globalCache.Hits,
		Miss: globalCache.Misses,
	}
}

func QueryModel(data string, curr AlgoType) AlgoType {
	req, err := http.NewRequest("POST", "http://localhost:5000/predict", bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error Querying Server " + err.Error())
		return curr
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Querying Server " + err.Error())
		return curr
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return curr
	}

	var p Prediction
	json.Unmarshal(body, &p)

	if p.Res == 0 {
		return curr
	} else if p.Res == 1 {
		return LRU
	} else {
		return LFU
	}
}

func RandomCache(size int, splits int, n int) *[]models.CacheItem {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	err := filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "./data" {
			return nil
		}

		if info.IsDir() {
			return os.RemoveAll(path)
		} else {
			return os.Remove(path)
		}
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("All files deleted successfully.")
	}

	for i := 0; i < n; i++ {
		x := r.Intn(3)

		if x == 0 {
			fmt.Println("Generated Random")
			random.GenerateRandomArray(size, 0, 100)
		} else if x == 1 {
			fmt.Println("Generate Split Bias")
			random.GenerateSplitBiasedRandom(size, 0, 100, 5, splits)
		} else if x == 2 {
			fmt.Println("Generated Recency Bias")
			random.GenerateRecencyBias(size, 0, 100, 90)
		}
	}

	total, err := utils.ReadJSONArraysFromDir("./data")
	if err != nil {
		log.Fatal("Error Reading Director ./data")
		return &[]models.CacheItem{}
	}

	arr := make([]int, 0, size*n)
	for _, v := range total {
		arr = append(arr, v...)
	}

	return generator.Converter(arr)
}
