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
    fmt.Println(err)
    return
  }

  stats := varnish.Stats()
  statFound := false
  for _, statname := range statnames {
    found = false
    for field, _ := range stats {
      if statname == field {
        statFound = true
      }
    }
    if statFound != true {
      err := statname + " is not valid varnishstat name"
      return err
    }
  }
  varnish.Close()
  return nil
}

func getStats(statnames []string) stats map[string]uint64 {
  // Open the default Varnish Shared Memory file
  varnish, err := vago.Open("")
  if err != nil {
    fmt.Println(err)
    return
  }

  stats := varnish.Stats()
  for field, value := range stats {
    for _, statname := range statnames {
      if field == statname {
        fmt.Printf("%s\t%d\n", field, value)
      }
    }
  }
  varnish.Close()
}
