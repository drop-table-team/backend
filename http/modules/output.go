package modules

import (
	"backend/module"
	"backend/util"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	MongoAddress     string `json:"mongoAddress"`
	MongoDatabase    string `json:"mongoDatabase"`
	MongoCollection  string `json:"mongoCollection"`
	QdrantAddress    string `json:"qdrantAddress"`
	QdrantCollection string `json:"qdrantCollection"`
}

func HandleOutputRegister(manager *module.ModuleManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		module := manager.FindByName(req.Name)
		if module == nil {
			http.Error(w, "module not found", http.StatusNotFound)
			return
		}

		_ = util.SendJson(w, RegisterResponse{
			MongoAddress:     util.MongoUri,
			MongoDatabase:    util.MongoDatabase,
			MongoCollection:  "entries",
			QdrantAddress:    util.QdrantUrl,
			QdrantCollection: "data",
		})
	}
}

func HandleOutputUnregister(manager *module.ModuleManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		module := manager.FindByName(request.Name)
		if module == nil {
			http.Error(w, "module not found", http.StatusNotFound)
			return
		}
	}
}
