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
	MongoDatabase   string `json:"mongo_database"`
	MongoCollection string `json:"mongo_collection"`
}

func HandleOutputRegister(manager *module.ModuleManager, mongoDatabase string) http.HandlerFunc {
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

		module.Register()

		_ = util.SendJson(w, RegisterResponse{
			MongoDatabase:   mongoDatabase,
			MongoCollection: "entries",
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

		module.Unregister()
	}
}
