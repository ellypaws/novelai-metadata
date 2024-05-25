package main

import (
	"log"
)

func main() {
	paths := getPathsFromArgsOrPrompt()

	for _, p := range paths {
		log.Printf("Processing path: %s", p)
		err := processPath(p)
		if err != nil {
			log.Printf("Failed to process path %s: %v", p, err)
		}
	}
}
