package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/compute/metadata"
)

func main() {
	// hardcoded for now - will take from user input later
	statnames := []string{"MAIN.cache_hit", "MAIN.cache_miss", "MAIN.cache_hitpass"}
	applicationName := "myapp"
	interval := 60

	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatalf("Can't get project ID from metadata: %v", err)
	}

	if err := validateStats(statnames); err != nil {
		log.Fatalf("Stat validation failed: %v", err)
	}

	// get monitoring service handle
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, stat := range statnames {
		createCustomMetric(s, projectID, stat)
	}
	log.Println("Start collection loop")
	oldStats := make(map[string]int64)
	// Tick on the minute
	t := minuteTicker()

	for {
		// Wait for ticker to send
		<-t.C

		// Update the ticker
		t = minuteTicker()
		// stats is the entire varnishstat output
		stats, err := getStats(statnames)
		if err != nil {
			log.Fatalf("Error getting stats: %v", err)
		}
		if len(oldStats) != 0 {
			for stat := range stats {
				delta := stats[stat] - oldStats[stat]
				tps := delta / int64(interval)
				writeTimeSeriesValue(s, projectID, stat, applicationName, tps)
			}
		}
		oldStats = stats
	}
}

func minuteTicker() *time.Ticker {
	// Return new ticker that triggers on the minute
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
