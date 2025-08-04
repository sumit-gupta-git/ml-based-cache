package models

import "time"

type CacheItem struct {
	Val       int
	LastUsed  time.Time
	Frequency int
	EntryTime time.Time
}
