package cache

import (
	"time"

	"ml-based-cache/internal/cache/algos"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/models"
	"ml-based-cache/internal/utils"
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

func (m *Cache) Simulate(items *[]models.CacheItem, sl *models.MLSelection) Res {
	r := Res{
		Items: generator.ReConverter(items),
	}

	// LRU
	for _, v := range *items {
		m.AccessLRU(&v)
	}
	r.LRUHits = m.Hits
	r.LRUMiss = m.Misses
	r.LRUAvgReaccess = utils.Average(m.AvgAccessTime)
	s := NewCache(m.Size)

	// LFU
	for _, v := range *items {
		s.AccessLFU(&v)
	}
	r.LFUHits = s.Hits
	r.LFUMiss = s.Misses
	r.LFUAvgReaccess = utils.Average(s.AvgAccessTime)

	// // FIFO
	// for _, v := range *items {
	// 	v.EntryTime = time.Now()
	//
	// 	now := time.Now()
	// 	m.AccessFIFO(&v)
	// 	m.ResTimes = append(m.ResTimes, time.Since(now).Microseconds())
	// }
	// r.FIFOHits = m.Hits
	// r.FIFOMiss = m.Misses

	return r
}

func (m *Cache) AccessLRU(val *models.CacheItem) {
	for _, v := range *m.Slots {
		if v == nil {
			continue
		}

		if v.Val == val.Val {
			m.Hits += 1
			v.Frequency += 1
			m.AvgAccessTime = append(m.AvgAccessTime, time.Since(v.EntryTime).Microseconds())
			v.LastUsed = time.Now()

			return
		}
	}
	m.Misses += 1

	algos.LRUReplace(m.Slots, val)
}

func (m *Cache) AccessLFU(val *models.CacheItem) {
	for _, v := range *m.Slots {
		if v == nil {
			continue
		}

		if v.Val == val.Val {
			m.Hits += 1
			v.Frequency += 1
			m.AvgAccessTime = append(m.AvgAccessTime, time.Since(v.EntryTime).Microseconds())
			v.LastUsed = time.Now()

			return
		}
	}
	m.Misses += 1

	algos.LFUReplace(m.Slots, val)
}

// func (m *Cache) AccessFIFO(val *models.CacheItem) {
// 	for _, v := range *m.Slots {
// 		if v == nil {
// 			continue
// 		}
//
// 		if v.Val == val.Val {
// 			m.Hits += 1
// 			v.Frequency += 1
// 			m.AvgAccessTime = append(m.AvgAccessTime, time.Since(v.EntryTime).Microseconds())
// 			v.LastUsed = time.Now()
//
// 			return
// 		}
// 	}
// 	m.Misses += 1
//
// 	now := time.Now()
// 	algos.FIFOReplace(m.Slots, val)
// 	m.ReplaceTimes = append(m.ReplaceTimes, time.Since(now).Microseconds())
// }
//
// func (m Res) Print() {
// 	fmt.Printf(`
// -------------------LRU--------------------
// Misses: %d
// Hits: %d
//
// -------------------LFU--------------------
// Misses: %d
// Hits: %d
//
// -------------------FIFO--------------------
// Misses: %d
// Hits: %d
// 		`, m.LRUMiss, m.LRUHits, m.LFUHits, m.LFUMiss, m.FIFOHits, m.FIFOMiss)
// }
