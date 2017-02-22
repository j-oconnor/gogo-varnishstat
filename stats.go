package main

import (
	"fmt"

	"github.com/phenomenes/vago"
)

func printStats(statnames []string) {
	// Open the default Varnish Shared Memory file
	v, err := vago.Open("")
	if err != nil {
		fmt.Println(err)
		return
	}
	// statnames := []string{"MAIN.cache_hit", "MAIN.cache_miss", "MAIN.cache_hitpass"}
	stats := v.Stats()
	for field, value := range stats {
		for _, statname := range statnames {
			if field == statname {
				fmt.Printf("%s\t%d\n", field, value)
			}
		}
	}
	v.Close()
}

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
    for field, _ := range stats {
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

func getStats(statnames []string) (map[string]uint64, error) {
  // Open the default Varnish Shared Memory file
  varnish, err := vago.Open("")
  if err != nil {
    return nil, fmt.Errorf("Could not open VSM Connection: %v", err)
  }

  stats := varnish.Stats()
  //for field, value := range stats {
  //  for _, statname := range statnames {
  //    if field == statname {
  //      fmt.Printf("%s\t%d\n", field, value)
  //    }
  //  }
  //}
  varnish.Close()
  return stats, nil
}
