package algos

import (
	"sort"

	"ml-based-cache/internal/models"
)

func LFUReplace(slots *[]*models.CacheItem, item *models.CacheItem) {
	for k, v := range *slots {
		if v == nil {
			(*slots)[k] = item
			return
		}
	}

	// Filter nil
	var nonNil []*models.CacheItem
	for _, v := range *slots {
		if v != nil {
			nonNil = append(nonNil, v)
		}
	}

	// Sort on Frequency
	sort.Slice(nonNil, func(i, j int) bool {
		return nonNil[i].Frequency < nonNil[j].Frequency
	})

	// Replace the least recently used
	lru := nonNil[0]
	for i, v := range *slots {
		if v == lru {
			(*slots)[i] = item
			return
		}
	}
}
