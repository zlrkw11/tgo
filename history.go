package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const maxHistoryEntries = 20

type DurationRecord struct {
	DurationMs int64     `json:"duration_ms"`
	Timestamp  time.Time `json:"timestamp"`
}

type TestHistory map[string][]DurationRecord

type Trend struct {
	Direction string        // "faster", "slower", "stable", ""
	Delta     time.Duration // current - average
	Average   time.Duration
}

func historyPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tgo", "history.json")
}

func LoadHistory() (TestHistory, error) {
	data, err := os.ReadFile(historyPath())
	if err != nil {
		return make(TestHistory), nil
	}
	var h TestHistory
	if err := json.Unmarshal(data, &h); err != nil {
		return make(TestHistory), nil
	}
	return h, nil
}

func SaveHistory(h TestHistory) error {
	path := historyPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func RecordRun(h TestHistory, packages []Package) TestHistory {
	now := time.Now()
	for _, pkg := range packages {
		for _, t := range pkg.Tests {
			if t.Status == "" {
				continue
			}
			key := pkg.Name + "/" + t.Name
			record := DurationRecord{
				DurationMs: t.Duration.Milliseconds(),
				Timestamp:  now,
			}
			h[key] = append(h[key], record)
			// trim to last N entries
			if len(h[key]) > maxHistoryEntries {
				h[key] = h[key][len(h[key])-maxHistoryEntries:]
			}
		}
	}
	return h
}

func GetTrend(h TestHistory, pkgName, testName string, current time.Duration) Trend {
	key := pkgName + "/" + testName
	records := h[key]
	if len(records) < 2 {
		return Trend{}
	}

	// average of previous runs (exclude the most recent one we just added)
	var total int64
	count := len(records) - 1 // exclude latest
	for i := 0; i < count; i++ {
		total += records[i].DurationMs
	}
	avg := time.Duration(total/int64(count)) * time.Millisecond
	delta := current - avg

	var direction string
	// only show trend if avg is meaningful (>1ms) and change is significant
	if avg > time.Millisecond {
		ratio := float64(current) / float64(avg)
		if ratio > 1.5 {
			direction = "slower"
		} else if ratio < 0.7 {
			direction = "faster"
		} else {
			direction = "stable"
		}
	}

	return Trend{
		Direction: direction,
		Delta:     delta,
		Average:   avg,
	}
}
