package main

import (
	httpModules "backend/http/modules"
	"backend/storage"
	"context"
	"errors"
	"net/http"
	"os"

	"backend/module"
	"backend/util"
	"encoding/json"
	"log"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	moduleConfigPath := util.MaybeEnv("MODULE_CONFIG_PATH")

	minioUrl := util.MustEnv("MINIO_HOST")
	minioBucket := util.MustEnv("MINIO_BUCKET")

	minioAccessKey := util.MustEnv("MINIO_ACCESS_KEY")
	minioSecretKey := util.MustEnv("MINIO_SECRET_KEY")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// register interrupt handler to handle keyboard interrupts (ctrl-c)
	initInterruptHandler()

	moduleManager, err := initModules(moduleConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	stopFunctions = append(stopFunctions, func() { _ = moduleManager.StopAll() })

	mongoClient, err := initMongo(util.MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	mongoDatabase := mongoClient.Database(util.MongoDatabase)
	stopFunctions = append(stopFunctions, func() { _ = mongoClient.Disconnect(context.Background()) })

	storageClient, err := initStorage(storage.Config{
		Endpoint: minioUrl,
		Bucket:   minioBucket,
		Credentials: storage.Credentials{
			AccessKey: minioAccessKey,
			SecretKey: minioSecretKey,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// output modules
	// mux.HandleFunc("GET /modules/output", httpModules.HandleGetOutput(moduleManager, mongoDatabaseName))
	mux.HandleFunc("POST /modules/output", httpModules.HandleOutputRegister(moduleManager))

	mux.HandleFunc("GET /modules/input", httpModules.HandleGetInput(moduleManager))
	mux.HandleFunc("POST /modules/input", httpModules.HandleInput(mongoDatabase, storageClient))
	mux.HandleFunc("POST /modules/input/{module}", httpModules.HandleProxyInput(moduleManager))

	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			mux.ServeHTTP(w, r)
		}),
	}
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	for _, stopFunction := range stopFunctions {
		stopFunction()
	}
}

var stopFunctions []func()

func initInterruptHandler() {
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Print("received interrupt signal, stopping gracefully")
		go func() {
			for _, stopFunction := range stopFunctions {
				stopFunction()
			}
			os.Exit(0)
		}()
		<-c
		<-c
		log.Fatal("received 3 interrupt signals, aborting immediately")
	}()
}

func initMongo(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return client, client.Ping(context.Background(), nil)
}

func initStorage(config storage.Config) (*storage.Storage, error) {
	return storage.NewStorage(config)
}

func initModules(moduleConfigPath *string) (*module.ModuleManager, error) {
	var config module.ModuleConfig
	if moduleConfigPath == nil {
		log.Print("no config file given")
	} else {
		// check if module config file is valid
		if _, err := os.Stat(*moduleConfigPath); err != nil && errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("module config file doesn't exist")
		} else if err != nil {
			return nil, err
		}

		configFileContent, err := os.ReadFile(*moduleConfigPath)
		if err != nil {
			return nil, err
		}
		config, err = module.ParseServiceConfig(configFileContent)
		if err != nil {
			return nil, err
		}
		log.Printf("parsed module config: %v", string(util.UnwrapError(json.Marshal(config))))
	}

	if config.Modules == nil {
		config.Modules = make([]string, 0)
	}
	if config.ModuleDefinitions == nil {
		config.ModuleDefinitions = make([]module.ModuleDefinition, 0)
	}

	moduleManager, err := module.NewModuleManager(config)
	if err != nil {
		return nil, err
	}

	return &moduleManager, moduleManager.StartAll()
}
