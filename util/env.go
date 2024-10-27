package util

import (
	"log"
	"os"
)

var OllamaUrl = MustEnv("OLLAMA_URL")
var QdrantUrl = MustEnv("QDRANT_URL")
var NetworkName = MustEnv("NETWORK_NAME")
var MongoUri = MustEnv("MONGO_URI")
var MongoDatabase = MustEnv("MONGO_DATABASE")

var EmbedderUrl = MustEnv("BACKEND_EMBEDDEDER_URL")

func MaybeEnv(key string) *string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}
	return &v
}

func MustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required environment variable not set: %s", key)
	}
	return v
}
