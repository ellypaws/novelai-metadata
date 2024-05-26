package main

import (
	"bufio"
	"encoding/json"
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
	save  save
}

func processPath(p string, save save) (map[string]*meta.Metadata, error) {
	var out = make(map[string]*meta.Metadata)
	data := data{
		out:   out,
		mutex: new(sync.Mutex),
		wait:  new(sync.WaitGroup),
		save:  save,
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
			if metadata == nil {
				return
			}
			if data.out != nil {
				data.mutex.Lock()
				data.out[path] = metadata
				data.mutex.Unlock()
			}
			if data.save != nil {
				err := data.save(path, metadata)
				if err != nil {
					log.Printf("Failed to save metadata for file %s: %v", path, err)
					return
				}
			}
			log.Printf("Done: %v", time.Since(now))
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

type save func(path string, metadata *meta.Metadata) error

func saveJSON(path string, metadata *meta.Metadata) error {
	jsonName := filepath.Join(filepath.Dir(path), fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))))
	jsonFile, err := os.Create(jsonName)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", jsonName, err)
		return err
	}
	defer jsonFile.Close()

	bin, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = jsonFile.Write(bin)
	if err != nil {
		return fmt.Errorf("failed to write metadata to file: %w", err)
	}

	log.Printf("Saved %s", jsonName)
	return nil
}
