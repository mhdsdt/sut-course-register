package main

import (
	"fmt"
	"flag"
)

func init() {
    flag.IntVar(&delaySeconds, "d", 5, "Delay in seconds between registration attempts")
    flag.IntVar(&maxRetries, "r", 5, "Maximum number of registration retries")
    flag.BoolVar(&infiniteRequests, "i", false, "Request indefinitely until successful")
    flag.BoolVar(&onTimeRegistration, "on-time", false, "Enable on-time registration")
    flag.StringVar(&configFileName, "config", "config.json", "Path to the configuration file")
    flag.IntVar(&offset, "o", 300, "Offset in milliseconds before the first registration request")
    flag.Parse()
}

func main() {
	err := readConfig()
    if err != nil {
        fmt.Printf("Error reading auth token and favorites from file: %v\n", err)
        return
    }
	
    establishWebsocket()
}
