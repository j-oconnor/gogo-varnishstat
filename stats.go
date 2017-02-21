package main

import (
    "fmt"

    "github.com/phenomenes/vago"
)

func printStats() {
    // Open the default Varnish Shared Memory file
    v, err := vago.Open("")
    if err != nil {
        fmt.Println(err)
        return
    }
    statnames := []string{"MAIN.cache_hit", "MAIN.cache_miss", "MAIN.cache_hitpass"}
    stats := v.Stats()
    for field, value := range stats {
      for _,statname := range statnames {
        if field == statname {
          fmt.Printf("%s\t%d\n", field, value)
        }
      }
    }
    v.Close()
}
