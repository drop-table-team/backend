package util

import (
	"os"
)

func MaybeEnv(key string) *string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}
	return &v
}
