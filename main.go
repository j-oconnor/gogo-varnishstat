package main

import (
    "log"
    "time"
)

func main() {

        log.Println("Started ticker")

        // Tick on the minute
        t := minuteTicker()

        for {
                // Wait for ticker to send
                <-t.C

                // Update the ticker
                t = minuteTicker()

                go printStats()
        }
}

func minuteTicker() *time.Ticker {
        // Return new ticker that triggers on the minute
        return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
