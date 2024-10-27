package modules

import (
	"backend/module"
	"backend/util"
	"net/http"
)

type ModuleGetResponse struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

func HandleGetInput(manager *module.ModuleManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := make([]ModuleGetResponse, 0)
		for _, definition := range manager.Config().ModuleDefinitions {
			response = append(response, ModuleGetResponse{
				Name:  definition.Name,
				Types: definition.Types,
			})
		}

		_ = util.SendJson(w, response)
	}
}
