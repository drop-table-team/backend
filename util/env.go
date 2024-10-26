package util

import (
	"log"
	"os"
)

func MustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s not set", key)
	}
	return v
}
