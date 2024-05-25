package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nai-metadata/pkg/meta"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image files>")
		return
	}

	for _, fn := range os.Args[1:] {
		imgFile, err := os.Open(fn)
		if err != nil {
			log.Printf("failed to open file: %v", err)
			continue
		}
		defer imgFile.Close()

		metaData, err := meta.ExtractMetadata(imgFile)
		if err != nil {
			log.Printf("failed to extract metadata: %v", err)
			continue
		}

		bin, err := json.MarshalIndent(metaData, "", "  ")
		if err != nil {
			log.Printf("failed to marshal json: %v", err)
			continue
		}

		// Write the metadata to a file
		jsonFileName := fmt.Sprintf("%s.json", path.Base(fn[:len(fn)-len(path.Ext(fn))]))
		jsonFile, err := os.Create(jsonFileName)
		if err != nil {
			log.Printf("failed to create file: %v", err)
			continue
		}
		defer jsonFile.Close()

		written, err := jsonFile.Write(bin)
		if err != nil {
			log.Printf("failed to write json: %v", err)
			continue
		}

		fmt.Printf("Wrote %d bytes to %s\n", written, jsonFileName)
	}
}
