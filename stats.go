package qml

import (
	"sync"
)

var stats *Stats
var statsMutex sync.Mutex

func SetStats(enabled bool) {
	statsMutex.Lock()
	if enabled {
		if stats == nil {
			stats = &Stats{}
		}
	} else {
		stats = nil
	}
	statsMutex.Unlock()
}

func GetStats() (snapshot Stats) {
	statsMutex.Lock()
	snapshot = *stats
	statsMutex.Unlock()
	return
}

func ResetStats() {
	statsMutex.Lock()
	old := stats
	stats = &Stats{}
	// These are absolute values:
	stats.EnginesAlive = old.EnginesAlive
	stats.ValuesAlive = old.ValuesAlive
	statsMutex.Unlock()
	return
}

type Stats struct {
	EnginesAlive int
	ValuesAlive int
}

func (stats *Stats) enginesAlive(delta int) {
	if stats != nil {
		statsMutex.Lock()
		stats.EnginesAlive += delta
		statsMutex.Unlock()
	}
}

func (stats *Stats) valuesAlive(delta int) {
	if stats != nil {
		statsMutex.Lock()
		stats.ValuesAlive += delta
		statsMutex.Unlock()
	}
}
