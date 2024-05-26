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
		return fmt.Errorf("failed to stat path %s: %v", p, err)
	}

	if info.IsDir() {
		return processDirectory(p)
	} else if path.Ext(info.Name()) == ".png" {
		return processFile(p)
	}

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
	imgFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer imgFile.Close()

	data, err := meta.ExtractMetadata(imgFile)
	if err != nil {
		return fmt.Errorf("failed to extract metadata from file: %w", err)
	}

	valid, err := meta.IsNovelAI(*data)
	if err != nil {
		return fmt.Errorf("failed to verify metadata for file: %w", err)
	}
	if !valid {
		log.Printf("Warning: Invalid signature for %s", filePath)
		return nil
	}

	bin, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	jsonName := path.Join(path.Dir(filePath), fmt.Sprintf("%s.json", strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))))
	return saveFile(jsonName, bin)
}

func saveFile(filePath string, data []byte) error {
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
