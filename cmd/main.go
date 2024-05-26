package main

import (
	"log"
	"time"
)

func main() {
	paths := getPathsFromArgsOrPrompt()

	now := time.Now()
	for _, p := range paths {
		log.Printf("Processing path: %s", p)
		_, err := processPath(p, saveCaption)
		if err != nil {
			log.Printf("Failed to process path %s: %v", p, err)
		}
	}
	log.Printf("Finished processing in %v", time.Since(now))
}
