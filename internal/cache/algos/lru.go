package algos

import (
	"sort"
	"time"

	"ml-based-cache/internal/models"
)

func LRUReplace(slots *[]*models.CacheItem, item *models.CacheItem) {
	item.EntryTime = time.Now()

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

	// Sort on lastUsed
	sort.Slice(nonNil, func(i, j int) bool {
		return nonNil[i].LastUsed.Before(nonNil[j].LastUsed)
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
