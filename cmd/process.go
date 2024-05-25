package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"nai-metadata/pkg/meta"
	"os"
	"path"
	"strings"
)

func getPathsFromArgsOrPrompt() []string {
	if len(os.Args) < 2 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter image files or directories separated by space: ")
		input, _ := reader.ReadString('\n')
		args := strings.Fields(input)
		if len(args) == 0 {
			fmt.Println("No input provided. Exiting.")
			os.Exit(1)
		}
		return args
	}
	return os.Args[1:]
}

func processPath(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		log.Printf("Error stating path %s: %v", p, err)
		return err
	}

	if info.IsDir() {
		return processDirectory(p)
	} else if path.Ext(info.Name()) == ".png" {
		return processFile(p)
	}

	log.Printf("Skipping non-PNG file: %s", p)
	return nil
}

func processDirectory(dirPath string) error {
	log.Printf("Entering directory: %s", dirPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Failed to read directory %s: %v", dirPath, err)
		return err
	}

	for _, entry := range entries {
		entryPath := path.Join(dirPath, entry.Name())
		err := processPath(entryPath)
		if err != nil {
			log.Printf("Failed to process entry %s: %v", entryPath, err)
		}
	}

	return nil
}

func processFile(filePath string) error {
	log.Printf("Opening file: %s", filePath)
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return err
	}
	defer imgFile.Close()

	data, err := meta.ExtractMetadata(imgFile)
	if err != nil {
		log.Printf("Failed to extract metadata from %s: %v", filePath, err)
		return err
	}

	bin, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal json for file %s: %v", filePath, err)
		return err
	}

	jsonName := path.Join(path.Dir(filePath), fmt.Sprintf("%s.json", strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))))
	return saveFile(jsonName, bin)
}

func saveFile(filePath string, data []byte) error {
	log.Printf("Writing metadata to file: %s", filePath)
	jsonFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
		return err
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(data)
	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", filePath, err)
		return err
	}

	log.Printf("Successfully wrote to file: %s", filePath)
	return nil
}
