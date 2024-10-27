package modules

import (
	"backend/data"
	"backend/module"
	"backend/storage"
	"backend/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
)

func HandleInput(database *mongo.Database, storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 50); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		jsonPart := r.FormValue("json")
		var entry data.Entry
		if err := json.Unmarshal([]byte(jsonPart), &entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entryUuid, _ := uuid.NewRandom()
		entry.Uuid = hex.EncodeToString(entryUuid[:])

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		var fileBuffer bytes.Buffer
		_, err = io.Copy(&fileBuffer, file)

		contentType := http.DetectContentType(fileBuffer.Bytes())
		for _, contentType = range header.Header.Values("Content-Type") {
		}

		_, err = storage.UploadFile(fileBuffer, entry.Uuid, header.Filename, contentType)
		if err != nil {
			http.Error(w, "Error uploading the file to minio", http.StatusBadRequest)
			return
		}

		data.AddEntry(database, entry)

		http.Post(fmt.Sprintf("%s/queue", util.EmbedderUrl), "application/json", bytes.NewBuffer(util.UnwrapError(json.Marshal(entry))))

		w.WriteHeader(http.StatusOK)
	}
}

func HandleProxyInput(manager *module.ModuleManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moduleName := r.PathValue("module")
		module := manager.FindByName(moduleName)

		if module == nil {
			http.Error(w, "Module not found", http.StatusNotFound)
			return
		}

		defer r.Body.Close()
		response, err := http.Post(fmt.Sprintf("%s/input", module.URL()), r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(response.StatusCode)
		defer response.Body.Close()
		io.Copy(w, response.Body)
	}
}
