package cache

import (
	"time"

	"ml-based-cache/internal/cache/algos"
	"ml-based-cache/internal/generator"
	"ml-based-cache/internal/models"
	"ml-based-cache/internal/utils"
)

func (m *Cache) SimulateTraining(items *[]models.CacheItem) Res {
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
