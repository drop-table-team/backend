package input

import (
	"backend/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func HandleInput(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var entry models.Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		models.AddEntry(client, entry)
	}
}
