package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/compute/metadata"
)

func main() {
	// hardcoded for now - will take from user input later
	statnames := []string{"MAIN.sess_conn", "MAIN.sess_drop", "MAIN.sess_fail", "MAIN.cache_hit", 
			"MAIN.cache_hitpass", "MAIN.cache_miss", "MAIN.client_req", "MAIN.backend_conn", 
			"MAIN.backend_busy", "MAIN.backend_reuse", "MAIN.threads", "MAIN.n_object", 
			"MAIN.n_lru_nuked", "MAIN.bans_obj", "MAIN.bans_req"}

	applicationName,err := metadata.InstanceAttributeValue("artifact-id")
	if err != nil {
		log.Fatalf("Can't get application name from 'artifact-id' instance metedata attr: %v", err)
	} 

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
		
		// stats is the filtered set of varnishstats
		stats, err := getStats(statnames)
		if err != nil {
			log.Fatalf("Error getting stats: %v", err)
		}
		
		// TODO: add checks for stats < oldstats (in varnishd restart scenario)
		// don't do delta calc on first ticker
		if len(oldStats) != 0 {
			// do delta calculations
			for stat := range stats {
				delta := stats[stat] - oldStats[stat]
				if delta >= 0 {
					if err := writeTimeSeriesValue(s, projectID, stat, applicationName, delta); err != nil {
						log.Println(err)
					}
				}	
			}
		} else { log.Println("No metrics pushed on first minute loop") }
		// overwrite old stats
		oldStats = stats
	}
}

func minuteTicker() *time.Ticker {
	// Return new ticker that triggers on the minute
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
