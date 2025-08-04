package cache

import (
	"ml-based-cache/internal/models"
)

type Res struct {
	Items          []int
	LRUHits        int
	LRUMiss        int
	LFUHits        int
	LFUMiss        int
	LFUAvgReaccess float64
	LRUAvgReaccess float64
}

type Cache struct {
	Slots         *[]*models.CacheItem
	Size          int
	Misses        int
	Hits          int
	AvgAccessTime []int64
}

func NewCache(size int) *Cache {
	slots := make([]*models.CacheItem, size)
	return &Cache{
		Slots:         &slots,
		Size:          size,
		Misses:        0,
		Hits:          0,
		AvgAccessTime: make([]int64, 0),
	}
}
