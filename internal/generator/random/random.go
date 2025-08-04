package random

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GenerateSplitBiasedRandom(size int, min int, max int, biasPoolSize int, splits int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if biasPoolSize <= 0 || biasPoolSize > max-min+1 {
		fmt.Println("Invalid biasPoolSize")
		return
	}
	if size < 1 || min > max || splits < 1 || size%splits != 0 {
		fmt.Println("Invalid size, range, or splits")
		return
	}

	biasPool := make([]int, biasPoolSize)
	biasPoolSet := make(map[int]bool)
	for i := 0; i < biasPoolSize; i++ {
		biasPool[i] = min + i
		biasPoolSet[min+i] = true
	}

	segmentSize := size / splits
	firstSize := segmentSize / 10
	middleSize := (segmentSize * 4) / 10
	lastSize := segmentSize - firstSize - middleSize

	generateUniqueRandom := func(length int) []int {
		available := make([]int, 0, max-min+1-biasPoolSize)
		for i := min; i <= max; i++ {
			if !biasPoolSet[i] {
				available = append(available, i)
			}
		}
		if len(available) < length {
			fmt.Println("Not enough unique items for random parts")
			return nil
		}
		rng.Shuffle(len(available), func(i, j int) { available[i], available[j] = available[j], available[i] })
		return available[:length]
	}

	first := generateUniqueRandom(firstSize)
	if first == nil {
		return
	}

	middle := make([]int, middleSize)
	for i := 0; i < middleSize; i++ {
		middle[i] = biasPool[i%len(biasPool)]
	}

	// Last part (40% unique random)
	last := generateUniqueRandom(lastSize)
	if last == nil {
		return
	}

	segment := append(first, middle...)
	segment = append(segment, last...)

	// Repeat segment splits times
	arr := make([]int, 0, size)
	for i := 0; i < splits; i++ {
		arr = append(arr, segment...)
	}

	fname := fmt.Sprintf("./data/%d", time.Now().UnixNano())
	f, err := os.Create(fname)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	data, err := json.Marshal(arr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = f.WriteString(string(data))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GenerateRecencyBias(size, min, max int, biasPercent float64) {
	if size <= 0 || min > max {
		return
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	result := make([]int, 0, size)

	for i := 0; i < size; i++ {
		useRecent := r.Float64()*100 < biasPercent && len(result) > 0

		if useRecent {
			windowSize := 10
			if len(result) < windowSize {
				windowSize = len(result)
			}
			start := len(result) - windowSize
			recent := result[start:]

			choice := recent[r.Intn(len(recent))]
			result = append(result, choice)
		} else {
			newVal := r.Intn(max-min+1) + min
			result = append(result, newVal)
		}
	}

	f, err := os.Create(fmt.Sprintf("./data/%d", time.Now().UnixNano()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, _ := json.Marshal(result)
	fmt.Fprint(f, string(data))
}

func GenerateRandomArray(size, z, b int) {
	if z > b {
		z, b = b, z
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = r.Intn(b-z+1) + z
	}

	f, err := os.Create(fmt.Sprintf("./data/%d", time.Now().UnixNano()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, _ := json.Marshal(arr)
	fmt.Fprint(f, string(data))
}
