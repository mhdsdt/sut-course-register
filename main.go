package main

import (
    "course_register/src"
	"fmt"
	"flag"
)

func init() {
    flag.IntVar(&src.DelaySeconds, "d", 5, "Delay in seconds between registration attempts")
    flag.IntVar(&src.MaxRetries, "r", 5, "Maximum number of registration retries")
    flag.BoolVar(&src.InfiniteRequests, "i", false, "Request indefinitely until successful")
    flag.BoolVar(&src.OnTimeRegistration, "on-time", false, "Enable on-time registration")
    flag.StringVar(&src.ConfigFileName, "config", "config.json", "Path to the configuration file")
    flag.IntVar(&src.Offset, "o", 300, "Offset in milliseconds before the first registration request")
    flag.Parse()
}

func main() {
	err := src.ReadConfig()
    if err != nil {
        fmt.Printf("Error reading auth token and favorites from file: %v\n", err)
        return
    }
	
    src.EstablishWebsocket()
}
