package cache

import (
	"fmt"
	"time"

	"ml-based-cache/internal/cache/algos"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/models"
	"ml-based-cache/internal/utils"
)

type Res struct {
	Items       []int
	LRUHits     int
	LRUMiss     int
	LFUHits     int
	LFUMiss     int
	FIFOHits    int
	FIFOMiss    int
	AvgReaccess float64
}

type Cache struct {
	Slots         *[]*models.CacheItem
	Size          int
	Misses        int
	Hits          int
	ResTimes      []int64
	ReplaceTimes  []int64
	LatencyTimes  []int64
	AvgAccessTime []int64
}

func NewCache(size int) *Cache {
	slots := make([]*models.CacheItem, size)
	return &Cache{
		Slots:        &slots,
		Size:         size,
		Misses:       0,
		Hits:         0,
		ResTimes:     make([]int64, 0),
		ReplaceTimes: make([]int64, 0),
		LatencyTimes: make([]int64, 0),
	}
}

func (m *Cache) Simulate(items *[]models.CacheItem, sl *models.MLSelection) Res {
	r := Res{
		Items: generator.ReConverter(items),
	}

	// LRU
	for _, v := range *items {
		v.EntryTime = time.Now()

		now := time.Now()
		m.AccessLRU(&v)
		m.ResTimes = append(m.ResTimes, time.Since(now).Microseconds())
	}
	r.LRUHits = m.Hits
	r.LRUMiss = m.Misses
	m = NewCache(m.Size)

	// LFU
	for _, v := range *items {
		v.EntryTime = time.Now()

		now := time.Now()
		m.AccessLFU(&v)
		m.ResTimes = append(m.ResTimes, time.Since(now).Microseconds())
	}
	r.LFUHits = m.Hits
	r.LFUMiss = m.Misses
	m = NewCache(m.Size)

	// FIFO
	for _, v := range *items {
		v.EntryTime = time.Now()

		now := time.Now()
		m.AccessFIFO(&v)
		m.ResTimes = append(m.ResTimes, time.Since(now).Microseconds())
	}
	r.FIFOHits = m.Hits
	r.FIFOMiss = m.Misses

	r.AvgReaccess = utils.Average(m.AvgAccessTime)

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

	now := time.Now()
	algos.LRUReplace(m.Slots, val)
	m.ReplaceTimes = append(m.ReplaceTimes, time.Since(now).Microseconds())
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

	now := time.Now()
	algos.LFUReplace(m.Slots, val)
	m.ReplaceTimes = append(m.ReplaceTimes, time.Since(now).Microseconds())
}

func (m *Cache) AccessFIFO(val *models.CacheItem) {
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

	now := time.Now()
	algos.FIFOReplace(m.Slots, val)
	m.ReplaceTimes = append(m.ReplaceTimes, time.Since(now).Microseconds())
}

func (m Res) Print() {
	fmt.Printf(`
-------------------LRU--------------------
Misses: %d
Hits: %d

-------------------LFU--------------------
Misses: %d
Hits: %d

-------------------FIFO--------------------
Misses: %d
Hits: %d
		`, m.LRUMiss, m.LRUHits, m.LFUHits, m.LFUMiss, m.FIFOHits, m.FIFOMiss)
}
