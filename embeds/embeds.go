package main

import (
	"encoding/json"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

var (
	embedQueue       = make(chan string, 100)
	workerRunning    = false
	mu               sync.Mutex
	postURL          string
	embedInitialized = false
)

func initializeEmbeds() error {
	if embedInitialized {
		return nil
	}

	ollamaUrl, exists := os.LookupEnv("BACKEND_OLLAMA_URL")
	if !exists {
		return fmt.Errorf("environment variable OLLAMA_URL is not set")
	}

	postURL = ollamaUrl + "/api/embeddings"

	err := initializeQdrant()
	if err != nil {
		return err
	}

	embedInitialized = true

	return nil
}

func createAndStoreEmbed(text string) error {
	mu.Lock()
	defer mu.Unlock()

	if len(embedQueue) == cap(embedQueue) {
		return fmt.Errorf("queue is full")
	}

	embedQueue <- text
	if !workerRunning {
		workerRunning = true
		go worker()
	}

	return nil
}

func worker() {
	for text := range embedQueue {
		createEmbed(text)
	}
	mu.Lock()
	workerRunning = false
	mu.Unlock()
}

func createEmbed(text string) {
	payload := fmt.Sprintf(`{"model": "mxbai-embed-large","prompt": "%s"}`, text)
	resp, err := http.Post(postURL, "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var embedResponse EmbedResponse
	err = json.Unmarshal(body, &embedResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	err = upsertVector(embedResponse.Embedding)
	if err != nil {
		fmt.Println("Error upserting vector:", err)
		return
	}

	var points []*qdrant.ScoredPoint
	points, err = queryVector(embedResponse.Embedding)

	if err != nil {
		log.Fatalf("Could not search points: %v", err)
	}
	log.Printf("Found points: %s", points)

	fmt.Println("Vector upserted successfully")
}
