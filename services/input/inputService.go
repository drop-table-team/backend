package input

import (
	"backend/models"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HandleInput(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		jsonPart := r.FormValue("json")
		var entry models.Entry
		if err := json.Unmarshal([]byte(jsonPart), &entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		var fileBuffer bytes.Buffer
		fileSize, err := io.Copy(&fileBuffer, file)

		hash, err := generateFileHash(&fileBuffer)

		fmt.Println("Hash:" + hash + ", size: " + strconv.Itoa(int(fileSize)))

		models.AddEntry(client, entry)
	}
}

// generateFileHash creates a unique hash from the file data and timestamp.
func generateFileHash(file *bytes.Buffer) (string, error) {
	// Get the current timestamp as a string
	timestamp := time.Now().String()

	// Create a new SHA-256 hash instance
	hasher := sha256.New()

	// Write file contents to the hasher
	if _, err := hasher.Write(file.Bytes()); err != nil {
		return "", fmt.Errorf("failed to write file data to hash: %v", err)
	}

	// Write timestamp to the hasher
	if _, err := hasher.Write([]byte(timestamp)); err != nil {
		return "", fmt.Errorf("failed to write timestamp to hash: %v", err)
	}

	// Compute the hash and encode it as a hexadecimal string
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}
