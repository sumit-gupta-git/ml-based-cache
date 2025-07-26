package models

import "time"

type CacheItem struct {
	Val       int
	LastUsed  time.Time
	Frequency int
	EntryTime time.Time
}

// If !Bit1 && !Bit2 => LRU
// If Bit1 || Bit2 => LFU
// If Bit1 && Bit2 => FIFO
type MLSelection struct {
	Bit1 bool
	Bit2 bool
}
