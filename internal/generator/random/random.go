package random

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GenerateSplitBiasedRandom(size int, min int, max int, biasPoolSize int, splits int) {
	// Create a new random source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Validate inputs
	if biasPoolSize <= 0 || biasPoolSize > max-min+1 {
		fmt.Println("Invalid biasPoolSize")
		return
	}
	if size < 1 || min > max || splits < 1 || size%splits != 0 {
		fmt.Println("Invalid size, range, or splits")
		return
	}

	// Define the bias pool as a fixed sequence
	biasPool := make([]int, biasPoolSize)
	biasPoolSet := make(map[int]bool)
	for i := 0; i < biasPoolSize; i++ {
		biasPool[i] = min + i
		biasPoolSet[min+i] = true
	}

	// Calculate sizes for one segment
	segmentSize := size / splits
	firstSize := segmentSize / 10
	middleSize := (segmentSize * 4) / 10
	lastSize := segmentSize - firstSize - middleSize

	// Generate unique random parts, excluding bias pool
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

	// Generate one segment
	// First part (10% unique random)
	first := generateUniqueRandom(firstSize)
	if first == nil {
		return
	}

	// Middle part (50% from bias pool)
	middle := make([]int, middleSize)
	for i := 0; i < middleSize; i++ {
		middle[i] = biasPool[i%len(biasPool)]
	}

	// Last part (40% unique random)
	last := generateUniqueRandom(lastSize)
	if last == nil {
		return
	}

	// Concatenate parts to form one segment
	segment := append(first, middle...)
	segment = append(segment, last...)

	// Repeat segment splits times
	arr := make([]int, 0, size)
	for i := 0; i < splits; i++ {
		arr = append(arr, segment...)
	}

	// Write the array to a JSON file
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

func contains(arr []int, x int) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// func GenerateBiasedRandom(size, min, max int, biasPercent, biasRangePercent float64) {
// 	src := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(src)
//
// 	result := make([]int, size)
//
// 	// Calculate range limits
// 	totalRange := max - min
// 	biasMax := min + int(float64(totalRange)*biasRangePercent/100.0)
// 	biasedCount := int(float64(size) * biasPercent / 100.0)
//
// 	for i := 0; i < size; i++ {
// 		if i < biasedCount {
// 			// Use [min, biasMax] inclusive
// 			result[i] = r.Intn(biasMax-min+1) + min
// 		} else {
// 			// Use [min, max] inclusive
// 			result[i] = r.Intn(max-min+1) + min
// 		}
// 	}
//
// 	r.Shuffle(size, func(i, j int) { result[i], result[j] = result[j], result[i] })
//
// 	f, err := os.Create(fmt.Sprintf("./data/%d", time.Now().UnixNano()))
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
//
// 	data, _ := json.Marshal(result)
// 	fmt.Fprint(f, string(data))
// }

//	func GenerateRepeatBias(size, min, max int, biasPercent float64) {
//		if size <= 0 || min > max {
//			return
//		}
//
//		src := rand.NewSource(time.Now().UnixNano())
//		r := rand.New(src)
//
//		result := make([]int, 0, size)
//
//		for i := 0; i < size; i++ {
//			useExisting := r.Float64()*100 < biasPercent && len(result) > 0
//
//			if useExisting {
//				// Choose a random existing element
//				existing := result[r.Intn(len(result))]
//				result = append(result, existing)
//			} else {
//				// Generate a new value
//				newVal := r.Intn(max-min+1) + min
//				result = append(result, newVal)
//			}
//		}
//
//		f, err := os.Create(fmt.Sprintf("./data/%d", time.Now().UnixNano()))
//		if err != nil {
//			fmt.Println(err.Error())
//			return
//		}
//
//		data, _ := json.Marshal(result)
//		fmt.Fprint(f, string(data))
//	}
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
			// Determine recent window size (min of 10 or current length)
			windowSize := 10
			if len(result) < windowSize {
				windowSize = len(result)
			}
			start := len(result) - windowSize
			recent := result[start:]

			// Pick from recent elements
			choice := recent[r.Intn(len(recent))]
			result = append(result, choice)
		} else {
			// Generate a fresh random value
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
		z, b = b, z // ensure correct range
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
