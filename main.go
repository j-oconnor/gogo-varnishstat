package main

import (
	"log"
	"time"
)

func main() {
	// hardcoded for now - will take from user input later
	statnames := []string{"MAIN.cache_hit", "MAIN.cache_miss", "MAIN.cache_hitpass"}
	applicationName := "myapp"

	if err := validateStats(statnames); err != nil {
		log.Fatalf("Stat validation failed: %v", err)
	}

	log.Println("Start collection loop")

	// Tick on the minute
	t := minuteTicker()

	for {
		// Wait for ticker to send
		<-t.C

		// Update the ticker
		t = minuteTicker()

		stats := getStats(statnames)
	}
}

func minuteTicker() *time.Ticker {
	// Return new ticker that triggers on the minute
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
