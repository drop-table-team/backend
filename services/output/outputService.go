package output

import (
	"encoding/json"
	"net/http"
	"os"

	"backend/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterResponse struct {
	MongoAddress    string `json:"mongoAddress"`
	MongoDatabase   string `json:"mongoDatabase"`
	MongoCollection string `json:"mongoCollection"`
	QdrantAddress   string `json:"qdrantAddress"`
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
			MongoAddress:    os.Getenv("MONGO_ADDRESS"),
			MongoDatabase:   "mongo_data",
			MongoCollection: "entries",
			QdrantAddress:   os.Getenv("QDRANT_ADDRESS"),
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
