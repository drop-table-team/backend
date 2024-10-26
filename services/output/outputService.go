package output

import (
	"backend/models"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterResponse struct {
	MongoDatabase   string `json:"mongoDatabase"`
	MongoCollection string `json:"mongoCollection"`
}

func HandleRegister(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var module models.OutputModule
		if err := json.NewDecoder(r.Body).Decode(&module); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		models.AddModule(client, module)

		response := RegisterResponse{
			MongoDatabase:   "mongo_data",
			MongoCollection: "entries",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func HandleUnregister(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var request models.UnregisterName
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		models.RemoveModule(client, request)
	}
}
