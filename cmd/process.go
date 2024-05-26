package main

import (
	"bufio"
	"fmt"
	"log"
	"nai-metadata/pkg/meta"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

type data struct {
	out   map[string]*meta.Metadata
	mutex *sync.Mutex
	wait  *sync.WaitGroup
}

func processPath(p string) (map[string]*meta.Metadata, error) {
	var out = make(map[string]*meta.Metadata)
	data := data{
		out:   out,
		mutex: new(sync.Mutex),
		wait:  new(sync.WaitGroup),
	}
	_ = filepath.Walk(p, processWalk(data))
	data.wait.Wait()

	return out, nil
}

func processWalk(data data) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".png" {
			return nil
		}
		data.wait.Add(1)
		go func() {
			defer data.wait.Done()
			now := time.Now()
			metadata, err := processFile(path)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
				return
			}
			if metadata != nil {
				data.mutex.Lock()
				data.out[path] = metadata
				data.mutex.Unlock()
				log.Printf("Done: %v", time.Since(now))
			}
		}()
		return nil
	}
}

func processFile(filePath string) (*meta.Metadata, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer imgFile.Close()

	data, err := meta.ExtractMetadata(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from file: %w", err)
	}

	valid, err := meta.IsNovelAI(*data)
	if err != nil {
		return nil, fmt.Errorf("failed to verify metadata for file: %w", err)
	}
	if !valid {
		log.Printf("Warning: Invalid signature for %s", filePath)
		return nil, nil
	}

	return data, nil
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
