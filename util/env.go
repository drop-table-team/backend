package util

import (
	"log"
	"os"
)

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
