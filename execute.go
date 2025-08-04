package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"ml-based-cache/internal/cache"
	"ml-based-cache/internal/generator"
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
	arr := RandomCache(10_000, 0, 100)
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

func RandomCache(size, z, b int) *[]models.CacheItem {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = r.Intn(b-z+1) + z
	}

	return generator.Converter(arr)
}
