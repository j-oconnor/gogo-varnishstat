package main

import (
	"fmt"

	"github.com/phenomenes/vago"
)

func validateStats(statnames []string) error {
	// Open the default Varnish Shared Memory file
	varnish, err := vago.Open("")
	if err != nil {
		return fmt.Errorf("Could not open VSM Connection: %v", err)
	}

	stats := varnish.Stats()
	statFound := false
	for _, statname := range statnames {
		statFound = false
		for field := range stats {
			if statname == field {
				statFound = true
			}
		}
		if statFound != true {
			return fmt.Errorf("%s is not a valid varnishstat counter", statname)
		}
	}
	varnish.Close()
	return nil
}

// return only stats from selected stats
func getStats(statnames []string) (map[string]int64, error) {
	// Open the default Varnish Shared Memory file
	varnish, err := vago.Open("")
	if err != nil {
		return nil, fmt.Errorf("Could not open VSM Connection: %v", err)
	}
	responseStats := make(map[string]int64)
	stats := varnish.Stats()
	for _, statname := range statnames {
		if value, ok := stats[statname]; ok {
			responseStats[statname] = int64(value)
		}
	}
	varnish.Close()
	return responseStats, nil
}
