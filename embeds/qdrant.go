package embeds

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	qdrantInitialized    bool
	defaultSegmentNumber uint64 = 2
	client               *qdrant.Client
	ctx                  context.Context
)

const (
	collectionName        = "mxbai-embed-large"
	vectorSize     uint64 = 1024
	distance              = qdrant.Distance_Cosine
)

type EmbeddingChunk struct {
	index   int
	vector  []float32
	start   int
	end     int
	docUUID string
}

func initializeQdrant() error {
	if qdrantInitialized {
		return nil
	}

	var err error
	ctx = context.Background()

	qdrantUrl, exists := os.LookupEnv("QDRANT_URL")
	if !exists {
		return fmt.Errorf("environment variable OLLAMA_URL is not set")
	}

	parts := strings.Split(qdrantUrl, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid QDRANT_URL format, expected host:port")
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port number: %w", err)
	}

	client, err = qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	healthCheckResult, err := client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("error checking health: %w", err)
	}
	log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())

	// List collections
	collections, err := client.ListCollections(ctx)
	if err != nil {
		return fmt.Errorf("error listing collections: %w", err)
	} else {
		log.Printf("List of collections: %s", &collections)
	}

	collectionExists := false
	for _, collection := range collections {
		if collection == collectionName {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     vectorSize,
				Distance: distance,
			}),
			OptimizersConfig: &qdrant.OptimizersConfigDiff{
				DefaultSegmentNumber: &defaultSegmentNumber,
			},
		})
		if err != nil {
			return fmt.Errorf("error creating collection: %w", err)
		}
	}

	qdrantInitialized = true

	return nil
}

func upsertVector(chunk EmbeddingChunk) error {
	if !qdrantInitialized {
		return fmt.Errorf("qdrant not initialized")
	}

	if len(chunk.vector) != int(vectorSize) {
		return fmt.Errorf("vector size is not %d", vectorSize)
	}

	waitUpsert := false
	vectorUUID := uuid.New().String()
	upsertPoints := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDUUID(vectorUUID),
			Vectors: qdrant.NewVectors(chunk.vector...),
			Payload: qdrant.NewValueMap(map[string]any{
				"docUUID": chunk.docUUID,
				"start":   chunk.start,
				"end":     chunk.end,
				"index":   chunk.index,
			}),
		},
	}

	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})

	if err != nil {
		return fmt.Errorf("error upserting vector: %w", err)
	}

	return nil
}
