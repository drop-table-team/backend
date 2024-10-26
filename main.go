package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"backend/module"
	"backend/util"
	"encoding/json"
	"log"
	"os/signal"

	"backend/services/input"
	"backend/services/output"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Env struct {
	client *mongo.Client
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, HTTP!\n")
}

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://mongo:27017"

func main() {
	moduleConfigPath := util.MaybeEnv("MODULE_CONFIG_PATH")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// register interrupt handler to handle keyboard interrupts (ctrl-c)
	initInterruptHandler()

	moduleManager, err := initModules(moduleConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	stopFunctions = append(stopFunctions, func() { _ = moduleManager.StopAll() })

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	env := Env{client}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	// output
	http.HandleFunc("/modules/output/register", output.HandleRegister(env.client))
	http.HandleFunc("/modules/output/unregister", output.HandleUnregister(env.client))

	// input
	http.HandleFunc("/modules/input", input.HandleInput(env.client))

	err = http.ListenAndServe(":8080", nil)

	fmt.Println("Running")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
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
		config, err := module.ParseServiceConfig(configFileContent)
		if err != nil {
			return nil, err
		}
		log.Printf("parsed module config: %v", string(util.UnwrapError(json.Marshal(config))))
	}

	module, err := module.NewModuleManager(config)
	if err != nil {
		return nil, err
	}

	return &module, module.StartAll()
}
