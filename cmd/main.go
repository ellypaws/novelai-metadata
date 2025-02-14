package main

import (
	"log"
	"time"
)

func main() {
	args, processor := getPathsFromArgsOrPrompt()

	now := time.Now()
	for _, p := range args {
		log.Printf("Processing path: %s", p)
		_, err := processPath(p, processor)
		if err != nil {
			log.Printf("Failed to process path %s: %v", p, err)
		}
	}

	log.Printf("Finished processing in %v", time.Since(now))
}
