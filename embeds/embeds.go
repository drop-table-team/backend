package embeds

import (
	"encoding/json"
	"fmt"
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
	embedQueue       = make(chan Chunk, 100)
	workerRunning    = false
	mu               sync.Mutex
	postURL          string
	embedInitialized = false
)

const (
	chunkSize = 100
)

type Chunk struct {
	index   int
	text    string
	start   int
	end     int
	docUUID string
}

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

func createAndStoreEmbed(text string, docUUID string) error {
	mu.Lock()
	defer mu.Unlock()

	if len(embedQueue) == cap(embedQueue) {
		return fmt.Errorf("queue is full")
	}

	start := 0
	index := 0
	for len(text) > chunkSize {
		chunk := Chunk{
			text:    text[:chunkSize],
			start:   start,
			end:     start + chunkSize,
			docUUID: docUUID,
			index:   index,
		}
		embedQueue <- chunk
		text = text[chunkSize:]
		start += chunkSize
		index++
	}

	chunk := Chunk{
		text:    text,
		start:   start,
		end:     start + len(text),
		docUUID: docUUID,
		index:   index,
	}

	log.Println("Chunk:", chunk)

	embedQueue <- chunk
	if !workerRunning {
		workerRunning = true
		go worker()
	}

	return nil
}

func worker() {
	for chunk := range embedQueue {
		createEmbed(chunk)
	}
	mu.Lock()
	workerRunning = false
	mu.Unlock()
}

func createEmbed(chunk Chunk) {
	payload := fmt.Sprintf(`{"model": "mxbai-embed-large","prompt": "%s"}`, chunk.text)
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

	embeddingChunk := EmbeddingChunk{
		index:   chunk.index,
		vector:  embedResponse.Embedding,
		start:   chunk.start,
		end:     chunk.end,
		docUUID: chunk.docUUID,
	}
	err = upsertVector(embeddingChunk)
	if err != nil {
		fmt.Println("Error upserting vector:", err)
		return
	}

	fmt.Println("Vector upserted successfully")
}
