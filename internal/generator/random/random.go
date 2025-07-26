package random

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GenerateBiasedRandom(size, min, max int, biasPercent, biasRangePercent float64) {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	result := make([]int, size)

	// Calculate range limits
	totalRange := max - min
	biasMax := min + int(float64(totalRange)*biasRangePercent/100.0)
	biasedCount := int(float64(size) * biasPercent / 100.0)

	for i := 0; i < size; i++ {
		if i < biasedCount {
			// Use [min, biasMax] inclusive
			result[i] = r.Intn(biasMax-min+1) + min
		} else {
			// Use [min, max] inclusive
			result[i] = r.Intn(max-min+1) + min
		}
	}

	r.Shuffle(size, func(i, j int) { result[i], result[j] = result[j], result[i] })

	f, err := os.Create(fmt.Sprintf("./data/%d", time.Now().UnixNano()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, _ := json.Marshal(result)
	fmt.Fprint(f, string(data))
}
